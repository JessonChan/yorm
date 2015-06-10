package yorm

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

var tableSqlMap = map[string]string{}

type sqlScanner interface {
	Scan(dest ...interface{}) error
}

//Select query records(s) from a table.
func Select(i interface{}, condition string, args ...interface{}) error {
	return defaultExecutor.Select(i, condition, args...)
}

//SelectByPK select by pk.
func SelectByPK(i interface{}, tableName ...string) error {
	return defaultExecutor.SelectByPK(i, tableName...)
}

// select by the primary key,the table name param means you can select from other tables
func (ex *executor) SelectByPK(i interface{}, tableName ...string) error {
	iv := reflect.ValueOf(i)

	if !iv.IsValid() {
		return ErrNotSupported
	}
	q, err := newTableSetter(iv)
	if q == nil {
		return err
	}
	table := q.table
	if len(tableName) > 0 {
		table = tableName[0]
	}
	clause := tableSqlMap[table]
	if clause == "" {
		queryClause := buildSelectSql(q, tableName...)
		queryClause.WriteString("WHERE ")
		queryClause.WriteString(q.pkColumn.name)
		queryClause.WriteString("=? LIMIT 1")
		clause = queryClause.String()
		tableSqlMap[table] = clause
	}
	return ex.query(i, clause, iv.Elem().FieldByName(q.pkColumn.fieldName).Int())
}

func Count(i interface{}, where ...interface{}) int64 {
	// i maybe not ptr
	q, _ := newTableSetter(reflect.ValueOf(i))
	if q == nil {
		return 0
	}
	var count int64
	if len(where) > 0 {
		Select(&count, fmt.Sprintf("select count(0) from %s %s", q.table, where[0].(string)), where[1:]...)
	} else {
		Select(&count, fmt.Sprintf("select count(0) from %s", q.table))
	}
	return count
}

func R(i interface{}, cond ...interface{}) error {
	return defaultExecutor.R(i, cond...)
}

// the s method is a super method,easy but not so fast
func (ex *executor) R(i interface{}, cond ...interface{}) error {
	var typ reflect.Type
	var q *tableSetter
	var err error

	q, err = newTableSetter(reflect.ValueOf(i))
	if q == nil {
		typ, q, err = newTableSetterBySlice(i)
		if typ == nil {
			return err
		}
	}

	switch len(cond) {
	case 0:
		if typ == nil {
			return ex.SelectByPK(i)
		} else {
			return ex.Select(i, buildSelectSql(q).String())
		}
	default:
		if reflect.TypeOf(cond[0]).Kind() != reflect.String {
			return ErrIllegalParams
		}
		return ex.Select(i, cond[0].(string), cond[1:]...)
	}
}
func (ex *tranExecutor) Select(i interface{}, condition string, args ...interface{}) error {
	if ex == nil {
		return ErrNilMethodReceiver
	}

	if strings.HasPrefix(strings.ToUpper(condition), "SELECT") {
		return ex.query(i, condition, args...)
	}

	q, _ := newTableSetter(reflect.ValueOf(i))
	if q == nil {
		return ErrNotSupported
	}
	queryClause := buildSelectSql(q)

	if !strings.HasPrefix(strings.ToUpper(condition), "WHERE") {
		queryClause.WriteString("WHERE ")
	}
	queryClause.WriteString(condition)

	return ex.query(i, queryClause.String(), args...)
}

//Query do a select operation.
// if the is a struct ,you need not write select x,y,z,you need only write the where condition ...
func (ex *executor) Select(i interface{}, condition string, args ...interface{}) error {
	if ex == nil {
		return ErrNilMethodReceiver
	}

	if strings.HasPrefix(strings.ToUpper(condition), "SELECT") {
		return ex.query(i, condition, args...)
	}

	q, _ := newTableSetter(reflect.ValueOf(i))
	if q == nil {
		return ErrNotSupported
	}
	queryClause := buildSelectSql(q)

	if !strings.HasPrefix(strings.ToUpper(condition), "WHERE") {
		queryClause.WriteString("WHERE ")
	}
	queryClause.WriteString(condition)

	return ex.query(i, queryClause.String(), args...)
}

func buildSelectSql(q *tableSetter, tableName ...string) *bytes.Buffer {
	queryClause := bytes.NewBufferString("SELECT ")
	queryClause.WriteString(buildFullColumnSql(q))
	queryClause.WriteString("FROM ")
	queryClause.WriteString(append(tableName, q.table)[0])
	queryClause.WriteString(" ")
	return queryClause
}
func buildFullColumnSql(q *tableSetter) string {
	if q == nil || len(q.columns) == 0 {
		return ""
	}
	queryClause := bytes.NewBuffer([]byte{})
	for loop := 0; loop < len(q.columns); loop++ {
		queryClause.WriteByte(',')
		queryClause.WriteString(q.columns[loop].name)
	}
	queryClause.WriteByte(' ')
	//1 means a ","
	return string(queryClause.Bytes()[1:])
}

//Query do a query operation.
func (ex *tranExecutor) query(i interface{}, query string, args ...interface{}) error {
	typ := reflect.TypeOf(i)
	if typ.Kind() != reflect.Ptr {
		return ErrNonPtr
	}
	typ = typ.Elem()

	if strings.Contains(query, "*") {
		vl := reflect.ValueOf(i)
		if typ.Kind() == reflect.Slice {
			vl = reflect.New(typ.Elem())
		}
		q, _ := newTableSetter(vl)
		query = strings.Replace(query, "*", buildFullColumnSql(q), -1)
	}

	yogger.Debug("%s;%v", query, args)
	if typ.Kind() == reflect.Slice {
		rows, err := ex.Query(query, args...)
		if rows == nil {
			return err
		}
		return queryList(i, rows)
	}
	return queryOne(i, ex.QueryRow(query, args...))

}

//Query do a query operation.
func (ex *executor) query(i interface{}, query string, args ...interface{}) error {
	typ := reflect.TypeOf(i)
	if typ.Kind() != reflect.Ptr {
		return ErrNonPtr
	}
	typ = typ.Elem()

	if strings.Contains(query, "*") {
		vl := reflect.ValueOf(i)
		if typ.Kind() == reflect.Slice {
			vl = reflect.New(typ.Elem())
		}
		q, _ := newTableSetter(vl)
		query = strings.Replace(query, "*", buildFullColumnSql(q), -1)
	}

	var err error
	var stmt *sql.Stmt
	yogger.Debug("%s;%v", query, args)
	stmt, err = ex.getStmt(query)
	if stmt == nil {
		return err
	}
	if typ.Kind() == reflect.Slice {
		rows, err := stmt.Query(args...)
		if rows == nil {
			return err
		}
		return queryList(i, rows)
	}
	return queryOne(i, stmt.QueryRow(args...))

}

func queryOne(i interface{}, row *sql.Row) error {
	if row == nil {
		return ErrIllegalParams
	}
	return convertAssignRow(i, row)
}
func queryList(i interface{}, rows *sql.Rows) error {
	if rows == nil {
		return ErrIllegalParams
	}
	return convertAssignRows(i, rows)
}

func newTableSetterBySlice(i interface{}) (reflect.Type, *tableSetter, error) {
	typ := reflect.TypeOf(i)
	if typ.Kind() != reflect.Ptr {
		return nil, nil, ErrNonPtr
	}
	typ = typ.Elem()
	if typ.Kind() != reflect.Slice {
		return nil, nil, ErrNonSlice
	}
	typ = typ.Elem()
	if typ.Kind() == reflect.Struct {
		q, e := newTableSetter(reflect.New(typ))
		return typ, q, e
	} else {
		return typ, nil, nil
	}
}

func convertAssignRows(i interface{}, rows *sql.Rows) error {
	typ, q, err := newTableSetterBySlice(i)
	if typ == nil {
		return err
	}
	size := 0
	v := reflect.Indirect(reflect.ValueOf(i))
	for rows.Next() {
		if size >= v.Cap() {
			newCap := v.Cap()
			if newCap == 0 {
				newCap = 1
			} else {
				if newCap < 1024 {
					newCap += newCap
				} else {
					newCap += newCap / 4
				}
			}
			nv := reflect.MakeSlice(v.Type(), v.Len(), newCap)
			reflect.Copy(nv, v)
			v.Set(nv)
		}
		if size >= v.Len() {
			v.SetLen(size + 1)
		}
		st := reflect.New(typ)
		st = st.Elem()
		if q != nil {
			err := scanValue(rows, q, st)
			if err != nil {
				return err
			}
		} else {
			ti := newPtrInterface(typ)
			rows.Scan(ti)
			setValue(st, ti)
		}
		v.Index(size).Set(st)
		size++
	}
	return nil
}
func convertAssignRow(i interface{}, row *sql.Row) error {
	typ := reflect.TypeOf(i)

	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() != reflect.Struct {
		return row.Scan(i)
	}

	q, err := newTableSetter(reflect.ValueOf(i))
	if q == nil {
		return err
	}
	st := reflect.ValueOf(i).Elem()
	return scanValue(row, q, st)
}

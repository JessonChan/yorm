package yorm

import (
	"bytes"
	"database/sql"
	"reflect"
	"strings"
)

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
	if !reflect.ValueOf(i).IsValid() {
		return ErrNotSupported
	}
	q, err := newTableSetter(reflect.ValueOf(i))
	if q == nil {
		return err
	}
	queryClause := buildSelectSql(q, tableName...)
	queryClause.WriteString("WHERE ")
	queryClause.WriteString(q.pkColumn.name)
	queryClause.WriteString("=? LIMIT 1")
	return ex.query(i, queryClause.String(), reflect.ValueOf(i).Elem().FieldByName(q.pkColumn.fieldName).Int())
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
func (ex *executor) query(i interface{}, query string, args ...interface{}) error {
	typ := reflect.TypeOf(i)
	if typ.Kind() != reflect.Ptr {
		return ErrNonPtr
	}
	typ = typ.Elem()

	if strings.Contains(query, "*") {
		q, _ := newTableSetter(reflect.ValueOf(i))
		query = strings.Replace(query, "*", buildFullColumnSql(q), -1)
	}

	var err error
	var stmt *sql.Stmt
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

func convertAssignRows(i interface{}, rows *sql.Rows) error {
	typ := reflect.TypeOf(i)
	if typ.Kind() != reflect.Ptr {
		return ErrNonPtr
	}
	typ = typ.Elem()
	if typ.Kind() != reflect.Slice {
		return ErrNonSlice
	}
	typ = typ.Elem()
	var q *tableSetter
	var err error
	if typ.Kind() == reflect.Struct {
		q, err = newTableSetter(reflect.New(typ))
		if q == nil {
			return err
		}
	}
	size := 0
	v := reflect.Indirect(reflect.ValueOf(i))
	ti := newPtrInterface(typ)
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
			rows.Scan(ti)
			st.Set(reflect.ValueOf(ti).Elem())
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

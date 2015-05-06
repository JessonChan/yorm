package yorm

import (
	"bytes"
	"database/sql"
	"reflect"
	"strings"
	"time"
)

type sqlScanner interface {
	Scan(dest ...interface{}) error
}

func Select(i interface{}, condition string, args ...interface{}) error {
	return defaultExecutor.Select(i, condition, args...)
}

func SelectByPk(i interface{}, tableName ...string) error {
	return defaultExecutor.SelectByPk(i, tableName...)
}

// 这个设计是否合理？
func (this *executor) SelectByPk(i interface{}, tableName ...string) error {
	if !reflect.ValueOf(i).IsValid() {
		return ErrNotSupported
	}
	q, err := newTableSetter(reflect.ValueOf(i))
	if q == nil {
		return err
	}
	queryClause := buildSelectSql(q, append(tableName, q.table)[0])
	queryClause.WriteString("WHERE ")
	queryClause.WriteString(q.pkColumn.name)
	queryClause.WriteString("=?")
	return this.query(i, queryClause.String(), reflect.ValueOf(i).Elem().FieldByName(q.pkColumn.fieldName).Int())
}

//Query do a select operation.
// if the is a struct ,you need not write select x,y,z,you need only write the where condition ...
func (this *executor) Select(i interface{}, condition string, args ...interface{}) error {
	if this == nil {
		return ErrNilMethodReceiver
	}

	if strings.HasPrefix(strings.ToUpper(condition), "SELECT") {
		return this.query(i, condition, args...)
	}
	q, _ := newTableSetter(reflect.ValueOf(i))
	if q == nil {
		return this.query(i, condition, args...)
	}

	queryClause := buildSelectSql(q)

	if !strings.HasPrefix(strings.ToUpper(condition), "WHERE") {
		queryClause.WriteString("WHERE ")
	}
	queryClause.WriteString(condition)

	return this.query(i, queryClause.String(), args...)
}

func buildSelectSql(q *tableSetter, tabelName ...string) *bytes.Buffer {
	queryClause := bytes.NewBufferString("SELECT ")
	splitDot := ","
	for loop := 0; loop < len(q.columns); loop++ {
		if loop == len(q.columns)-1 {
			splitDot = " "
		}
		queryClause.WriteString(q.columns[loop].name)
		queryClause.WriteString(splitDot)
	}

	queryClause.WriteString("FROM ")
	queryClause.WriteString(q.table)
	queryClause.WriteString(" ")
	return queryClause
}

//Query do a query operation.
func (this *executor) query(i interface{}, query string, args ...interface{}) error {
	typ := reflect.TypeOf(i)
	if typ.Kind() != reflect.Ptr {
		return ErrNonPtr
	}
	typ = typ.Elem()
	var err error
	var stmt *sql.Stmt
	stmt, err = this.getStmt(query)
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

func scanValue(sc sqlScanner, q *tableSetter, st reflect.Value) error {
	err := sc.Scan(q.dests...)
	if err != nil {
		return err
	}
	for idx, c := range q.columns {
		// different assign func here
		switch c.typ.Kind() {
		case reflect.Int:
			st.Field(c.fieldNum).SetInt(int64(*(q.dests[idx].(*int))))
		case reflect.Int64:
			st.Field(c.fieldNum).SetInt(int64(*(q.dests[idx].(*int64))))
		case reflect.String:
			st.Field(c.fieldNum).SetString(string(*(q.dests[idx].(*string))))
		case reflect.Struct:
			switch c.typ {
			case TIME_TYPE:
				timeStr := string(*(q.dests[idx].(*string)))
				var layout string
				if len(timeStr) == 10 {
					layout = "2006-01-02"
				}
				if len(timeStr) == 19 {
					layout = "2006-01-02 15:04:05"
				}
				timeTime, err := time.ParseInLocation(layout, timeStr, time.Local)
				if timeTime.IsZero() {
					return err
				}
				st.Field(c.fieldNum).Set(reflect.ValueOf(timeTime))
			}
		}
	}
	return nil
}

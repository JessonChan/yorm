package yorm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"
)

type querySetter struct {
	table   string
	dests   []interface{}
	columns []column
}

type sqlScanner interface {
	Scan(dest ...interface{}) error
}

//Query do a query operation.
func Query(i interface{}, query string, args ...interface{}) error {
	typ := reflect.TypeOf(i)
	if typ.Kind() != reflect.Ptr {
		return ErrNonPtr
	}
	typ = typ.Elem()
	var err error
	var stmt *sql.Stmt
	stmt, err = getStmt(query)
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

func newQuery(ri reflect.Value) *querySetter {
	if q, ok := tableMap[ri.Kind()]; ok {
		return q
	}
	if ri.Kind() != reflect.Ptr || ri.IsNil() {
		return nil
	}
	q := new(querySetter)
	defer func() {
		tableMap[ri.Kind()] = q
	}()
	table, cs := structToTable(reflect.Indirect(ri).Interface())
	q.table = table
	q.columns = cs
	q.dests = make([]interface{}, len(cs))
	for k, v := range cs {
		q.dests[k] = newPtrInterface(v.typ)
	}
	return q
}

func newPtrInterface(t reflect.Type) interface{} {
	k := t.Kind()
	var ti interface{}
	switch k {
	case reflect.Int:
		ti = new(int)
	case reflect.Int64:
		ti = new(int64)
	case reflect.String:
		ti = new(string)
	case reflect.Struct:
		switch t {
		case reflect.TypeOf(time.Time{}):
			ti = new(string)
		}
	}
	return ti
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
	var q *querySetter
	if typ.Kind() == reflect.Struct {
		q = newQuery(reflect.New(typ))
		if q == nil {
			return errors.New("q is not support")
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

	q := newQuery(reflect.ValueOf(i))
	if q == nil {
		return errors.New("nil struct")
	}
	st := reflect.ValueOf(i).Elem()
	return scanValue(row, q, st)
}

func scanValue(sc sqlScanner, q *querySetter, st reflect.Value) error {
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
			case reflect.TypeOf(time.Time{}):
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

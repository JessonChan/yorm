package yorm

import (
	"database/sql"
	"reflect"
)

type Query struct {
	query string
	// todo when return array[] not a single value
	dests   []interface{}
	where   string
	table   string
	columns []column
}
type Where string

func (q *Query) String() string {
	if q.where == "" {
		return q.query
	}
	return q.query + " where " + q.where
}

func (q *Query) Where(w string) *Query {
	q.where = w
	return q
}

func newQuery(i interface{}) *Query {
	ri := reflect.ValueOf(i)
	if ri.Kind() != reflect.Ptr || ri.IsNil() {
		return nil
	}
	q := new(Query)
	table, cs := structToTable(reflect.Indirect(ri).Interface())
	q.table = table
	q.columns = cs
	q.query = query(table, cs)
	q.dests = make([]interface{}, len(cs))
	for k, v := range cs {
		var ti interface{}
		switch v.typ.Kind() {
		case reflect.Int:
			ti = new(int)
		case reflect.Int64:
			ti = new(int64)
		}
		q.dests[k] = ti
	}
	return q
}

// 对返回值进行赋值
func convertAssign(i interface{}, rows *sql.Rows, q *Query) error {
	err := berforeAssign(rows, q)
	if err != nil {
		return err
	}
	st := reflect.ValueOf(i).Elem()
	for idx, c := range q.columns {
		// different assign func here
		switch c.typ.Kind() {
		case reflect.Int:
			st.Field(c.fieldNum).SetInt(int64(*(q.dests[idx].(*int))))
		case reflect.Int64:
			st.Field(c.fieldNum).SetInt(int64(*(q.dests[idx].(*int64))))
		}
	}
	return nil
}
func berforeAssign(rows *sql.Rows, q *Query) error {
	for rows.Next() {
		err := rows.Scan(q.dests...)
		if err != nil {
			return err
		}
	}
	return nil
}

func query(table string, cs []column) string {
	s := "select "
	for k, v := range cs {
		s += v.name
		if k != len(cs)-1 {
			s += ","
		}
	}
	s += " from " + table
	return s
}

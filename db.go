package yorm

import (
	"database/sql"
	"reflect"
)

type Query struct {
	query   string
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
		if v.typ.Kind() == reflect.Int {
			q.dests[k] = new(int)
		}
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
		st.Field(c.fieldNum).SetInt(int64(*(q.dests[idx].(*int))))
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

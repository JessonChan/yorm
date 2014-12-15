package yorm

import "reflect"

type Query struct {
	query string
	dests []interface{}
	where string
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
	q.query = query(table, cs)
	q.dests = make([]interface{}, len(cs))
	for k, v := range cs {
		if v.typ.Kind() == reflect.Int {
			q.dests[k] = new(int)
		}
	}
	return q
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

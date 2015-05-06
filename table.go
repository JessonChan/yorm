package yorm

import "reflect"

type tableSetter struct {
	table   string
	dests   []interface{}
	columns []column
}

func newTableSetter(ri reflect.Value) *tableSetter {
	if q, ok := tableMap[ri.Kind()]; ok {
		return q
	}
	if ri.Kind() != reflect.Ptr || ri.IsNil() {
		return nil
	}
	q := new(tableSetter)
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

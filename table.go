package yorm

import (
	"reflect"
	"strings"
	"time"
)

type tableSetter struct {
	table    string
	dests    []interface{}
	columns  []*column
	pkColumn *column
}

func newTableSetter(ri reflect.Value) (*tableSetter, error) {
	if q, ok := tableMap[ri.Kind()]; ok {
		return q, nil
	}
	if ri.Kind() != reflect.Ptr || ri.IsNil() {
		return nil, ErrNotSupported
	}
	q := new(tableSetter)
	defer func() {
		tableMap[ri.Kind()] = q
	}()
	table, cs := structToTable(reflect.Indirect(ri).Interface())
	var err error
	q.pkColumn, err = findPkColumn(cs)
	if q.pkColumn == nil {
		return nil, err
	}
	q.table = table
	q.columns = cs
	q.dests = make([]interface{}, len(cs))
	for k, v := range cs {
		q.dests[k] = newPtrInterface(v.typ)
	}
	return q, nil
}

func findPkColumn(cs []*column) (*column, error) {
	var c *column
	var idColumn *column
	isPk := false

	for _, v := range cs {
		if strings.ToLower(v.name) == "id" {
			idColumn = v
		}
		if v.isPk {
			if isPk {
				return c, ErrDuplicatePkColumn
			}
			isPk = true
			c = v
		}
	}
	if c == nil {
		c = idColumn
	}
	if c == nil {
		return nil, ErrNonePkColumn
	}
	return c, nil
}

var (
	TIME_TYPE = reflect.TypeOf(time.Time{})
)

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
		case TIME_TYPE:
			ti = new(string)
		}
	}
	return ti
}

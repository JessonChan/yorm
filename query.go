package yorm

import (
	"database/sql"
	"errors"
	"reflect"
)

type QuerySetter struct {
	table   string
	dests   []interface{}
	columns []column
}

func QueryOne(i interface{}, row *sql.Row) error {
	if row == nil {
		return errors.New("nil row")
	}
	return convertAssignRow(i, row)
}
func QueryList(i interface{}, rows *sql.Rows) error {
	if rows == nil {
		return errors.New("rows nil")
	}
	return convertAssignRows(i, rows)
}

type sqlScanner interface {
	Scan(dest ...interface{}) error
}


var tableMap map[reflect.Kind]*QuerySetter = make(map[reflect.Kind]*QuerySetter)

func newQuery(ri reflect.Value) *QuerySetter {
	if q, ok := tableMap[ri.Kind()]; ok {
		return q
	}
	if ri.Kind() != reflect.Ptr || ri.IsNil() {
		return nil
	}
	q := new(QuerySetter)
	defer func() {
		tableMap[ri.Kind()] = q
	}()
	table, cs := structToTable(reflect.Indirect(ri).Interface())
	q.table = table
	q.columns = cs
	q.dests = make([]interface{}, len(cs))
	for k, v := range cs {
		q.dests[k] = newInterface(v.typ.Kind())
	}
	return q
}

func newInterface(k reflect.Kind) interface{} {
	var ti interface{}
	switch k {
	case reflect.Int:
		ti = new(int)
	case reflect.Int64:
		ti = new(int64)
	case reflect.String:
		ti = new(string)
	}
	return ti
}

func convertAssignRows(i interface{}, rows *sql.Rows) error {
	typ := reflect.TypeOf(i)
	if typ.Kind() != reflect.Ptr {
		return errors.New("not ptr")
	}
	typ = typ.Elem()
	if typ.Kind() != reflect.Slice {
		return errors.New("need a slice container")
	}
	typ = typ.Elem()
	var q *QuerySetter
	if typ.Kind() == reflect.Struct {
		q = newQuery(reflect.New(typ))
		if q == nil {
			return errors.New("q is not support")
		}
	}
	size := 0
	v := reflect.Indirect(reflect.ValueOf(i))
	ti := newInterface(typ.Kind())
	for rows.Next() {
		if size >= v.Cap() {
			newCap := v.Cap()
			if newCap == 0 {
				newCap = 1
			} else {
				newCap *= 2
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
			scanValue(rows, q, st)
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

	if typ.Kind() == reflect.Ptr && typ.Kind() != reflect.Struct {
		return row.Scan(i)
	}

	q := newQuery(reflect.ValueOf(i))
	if q == nil {
		return errors.New("nil struct")
	}
	st := reflect.ValueOf(i).Elem()
	return scanValue(row, q, st)
}

func scanValue(sc sqlScanner, q *QuerySetter, st reflect.Value) error {
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
		}
	}
	return nil
}

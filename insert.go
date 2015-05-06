package yorm

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

//Insert  return lastInsertId and error if has
func (this *executor) Insert(i interface{}, args ...string) (int64, error) {
	if this == nil {
		return 0, ErrNilMethodReceiver
	}

	q, err := newTableSetter(reflect.ValueOf(i))
	if q == nil {
		return 0, err
	}

	clause := &bytes.Buffer{}
	clause.WriteString("INSERT INTO ")
	clause.WriteString(q.table)
	clause.WriteString(" SET ")

	fs := ""
	e := reflect.ValueOf(i).Elem()
	var pk reflect.Value
	dests := []interface{}{}
	var columns []column
	if len(args) == 0 {
		columns = q.columns
	} else {
		for _, arg := range args {
			arg = strings.ToLower(arg)
			for _, c := range q.columns {
				if strings.ToLower(c.fieldName) == arg || strings.ToLower(c.name) == arg {
					columns = append(columns, c)
				}
			}
		}
	}
	for _, c := range columns {
		//todo this is auto increase field
		v := filedByName(e, c.fieldName)
		if strings.ToLower(c.fieldName) == "id" {
			pk = v
			continue
		}
		fs += fmt.Sprintf(",%v=?", c.name)
		dests = append(dests, fmt.Sprintf("%v", v.Interface()))
	}
	if fs == "" {
		return 0, errors.New("no filed to insert")
	}

	clause.WriteString(fs[1:])

	stmt, err := getStmt(clause.String())
	if err != nil {
		return 0, err
	}
	r, err := stmt.Exec(dests...)
	if err != nil {
		return 0, err
	}
	id, err := r.LastInsertId()
	if id > 0 && pk.IsValid() {
		pk.SetInt(id)
	}
	return id, err
}

func filedByName(e reflect.Value, names ...string) reflect.Value {
	var f reflect.Value
	for _, name := range names {
		f = e.FieldByName(name)
		if f.IsValid() {
			return f
		}
	}
	return f
}

func Insert(i interface{}, args ...string) (int64, error) {
	return defaultExecutor.Insert(i, args...)
}

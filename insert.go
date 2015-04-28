package yorm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// return lastInsertId and error if has
func Insert(i interface{}, args ...interface{}) (int64, error) {
	q := newQuery(reflect.ValueOf(i))
	if q == nil {
		return 0, errors.New("your object not support")
	}
	clause := "insert into " + q.table + " ("

	fs := ""
	vs := ""
	e := reflect.ValueOf(i).Elem()
	var pk reflect.Value
	for _, c := range q.columns {
		//todo this is auto increase field
		v := filedByName(e, c)
		if strings.ToLower(c.name) == "id" || strings.ToLower(c.fieldName) == "id" {
			pk = v
			continue
		}
		fs += fmt.Sprintf(",%v", c.name)
		vs += fmt.Sprintf(",%v", v.Interface())
	}
	if fs == "" || vs == "" {
		return 0, errors.New("no filed to insert")
	}
	clause += fs[1:] + ") values ("
	clause += vs[1:] + ")"
	stmt, err := getStmt(clause)
	if err != nil {
		return 0, err
	}
	r, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	id, err := r.LastInsertId()
	if id > 0 {
		pk.SetInt(id)
	}
	return id, err
}

func filedByName(e reflect.Value, c column) reflect.Value {
	var f reflect.Value
	for _, name := range []string{c.fieldName, c.name} {
		f = e.FieldByName(name)
		if f.IsValid() {
			return f
		}
	}
	return f
}

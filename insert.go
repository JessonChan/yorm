package yorm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// return lastInsertId and error if has
func Insert(i interface{}, args ...string) (int64, error) {
	q := newQuery(reflect.ValueOf(i))
	if q == nil {
		return 0, errors.New("your object not support")
	}
	clause := "insert into " + q.table + " set "

	fs := ""
	e := reflect.ValueOf(i).Elem()
	var pk reflect.Value
	dests := []interface{}{}

	if len(args) > 0 {
		for _, name := range args {
			v := filedByName(e, name)
			if strings.ToLower(name) == "id" {
				pk = v
				continue
			}
			fs += fmt.Sprintf(",%v=?", name)
			dests = append(dests, v.Interface())
		}
		if !pk.IsValid() {
			pk = filedByName(e, "id")
		}
	} else {
		for _, c := range q.columns {
			//todo this is auto increase field
			v := filedByName(e, c.fieldName, c.name)
			if strings.ToLower(c.name) == "id" || strings.ToLower(c.fieldName) == "id" {
				pk = v
				continue
			}
			fs += fmt.Sprintf(",%v=?", c.name)
			dests = append(dests, v.Interface())
		}
	}
	if fs == "" {
		return 0, errors.New("no filed to insert")
	}
	clause += fs[1:]
	fmt.Println(clause)
	stmt, err := getStmt(clause)
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

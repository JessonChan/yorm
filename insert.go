package yorm

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

//Insert  return lastInsertId and error if has
func (ex *executor) Insert(i interface{}, args ...string) (int64, error) {
	return insertExec(ex.exec, i, args...)
}

//Insert  return lastInsertId and error if has
func (ex *tranExecutor) Insert(i interface{}, args ...string) (int64, error) {
	return insertExec(ex.exec, i, args...)
}

//Insert insert a record.
func Insert(i interface{}, args ...string) (int64, error) {
	return defaultExecutor.Insert(i, args...)
}

//Insert  return lastInsertId and error if has
func insertExec(exec func(clause string, args ...interface{}) (sql.Result, error), i interface{}, args ...string) (int64, error) {

	q, err := newTableSetter(reflect.ValueOf(i))
	if q == nil {
		return 0, err
	}

	clause := &bytes.Buffer{}
	clause.WriteString("INSERT INTO ")
	clause.WriteString(q.table)
	clause.WriteString(" SET ")

	fs := &bytes.Buffer{}
	dests := []interface{}{}

	e := reflect.ValueOf(i).Elem()

	var pk reflect.Value
	var columns []*column

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
		v := e.FieldByName(c.fieldName)
		//todo how to handle a zero value ???
		//		if !v.IsValid() {
		//			continue
		//		}

		if c.isAuto {
			pk = v
			continue
		}
		vi := v.Interface()
		switch v.Type() {

		case TimeType:
			//zero time ,skip insert
			if vi.(time.Time).IsZero() {
				continue
			}
			vi = vi.(time.Time).Format(longSimpleTimeFormat)

		case BoolType:
			if vi.(bool) {
				vi = 1
			} else {
				vi = 0
			}
		}

		fs.WriteString("," + c.name + "=?")
		dests = append(dests, fmt.Sprintf("%v", vi))
	}
	if fs.Len() == 0 {
		return 0, errors.New("no filed to insert")
	}

	clause.Write(fs.Bytes()[1:])

	r, err := exec(clause.String(), dests...)
	if err != nil {
		return 0, err
	}
	id, err := r.LastInsertId()
	if id > 0 && pk.IsValid() {
		pk.SetInt(id)
	}
	return id, err
}

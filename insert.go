package yorm

import (
	"bytes"
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
func insertExec(exec ExecHandler, i interface{}, args ...string) (int64, error) {

	q, err := newTableSetter(reflect.ValueOf(i))

	var typ reflect.Type
	if q == nil {
		typ, q, err = newTableSetterBySlice(i)
		if q == nil {
			return 0, err
		}
	}

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

	fs := &bytes.Buffer{}
	clause := &bytes.Buffer{}
	dests := []interface{}{}
	if typ != nil {
		clause.WriteString("INSERT INTO ")
		clause.WriteString(q.table)
		clause.WriteString("(")
		for _, c := range columns {
			if c.isAuto {
				continue
			}
			fs.WriteString(",`" + c.name + "`")
		}
		if fs.Len() == 0 {
			return 0, errors.New("no filed to insert")
		}
		clause.Write(fs.Bytes()[1:])
		clause.WriteString(")")
		clause.WriteString("VALUES")
		is := reflect.ValueOf(i).Elem()
		for l := 0; l < is.Len(); l++ {
			fs.Reset()
			fs.WriteString("(")
			for _, c := range columns {
				v := is.Index(l).FieldByName(c.fieldName)
				if c.isAuto {
					continue
				}
				vi := v.Interface()
				switch v.Type() {
				case StringType:
					vi = strings.Replace(vi.(string), `'`, `\'`, -1)
					vi = strings.Replace(vi.(string), `"`, `\"`, -1)
				case TimeType:
					vi = vi.(time.Time).Format(longSimpleTimeFormat)

				case BoolType:
					if vi.(bool) {
						vi = 1
					} else {
						vi = 0
					}
				}
				fs.WriteString(fmt.Sprintf("'%v',", vi))
			}
			fs.Truncate(fs.Len() - 1)
			fs.WriteString(")")
			clause.WriteString(fs.String())
			if l < is.Len()-1 {
				clause.WriteString(",")
			}
		}
	} else {
		clause.WriteString("INSERT INTO ")
		clause.WriteString(q.table)
		clause.WriteString(" SET ")

		e := reflect.ValueOf(i).Elem()

		for _, c := range columns {
			v := e.FieldByName(c.fieldName)
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
	}

	r, err := exec(clause.String(), dests...)
	if err != nil {
		return 0, err
	}
	id, err := r.LastInsertId()
	if id > 0 && pk.IsValid() && typ == nil {
		pk.SetInt(id)
	}
	return id, err
}

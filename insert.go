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
func insertExec(exec ExecHandler, anyModel interface{}, args ...string) (int64, error) {

	q, err := newTableSetter(reflect.ValueOf(anyModel))

	var typ reflect.Type
	if q == nil {
		typ, q, err = newTableSetterBySlice(anyModel)
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

	fields := &bytes.Buffer{}
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
			fields.WriteString(",`" + c.name + "`")
		}
		if fields.Len() == 0 {
			return 0, errors.New("no row is inserted.")
		}
		clause.Write(fields.Bytes()[1:]) //skip the first ","
		clause.WriteString(")")
		clause.WriteString("VALUES")

		model := reflect.ValueOf(anyModel).Elem()

		for i := 0; i < model.Len(); i++ {
			fields.Reset()
			fields.WriteString("(")

			for _, c := range columns {
				v := model.Index(i).FieldByName(c.fieldName)
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

				fields.WriteString(fmt.Sprintf("'%v',", vi))
			}
			fields.Truncate(fields.Len() - 1)

			fields.WriteString(")")
			clause.WriteString(fields.String())

			if i < model.Len()-1 {
				clause.WriteString(",")
			}
		}

	} else {
		clause.WriteString("INSERT INTO ")
		clause.WriteString(q.table)
		clause.WriteString(" SET ")

		e := reflect.ValueOf(anyModel).Elem()

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

			fields.WriteString("," + c.name + "=?")
			dests = append(dests, fmt.Sprintf("%v", vi))
		}
		if fields.Len() == 0 {
			return 0, errors.New("no row is inserted.")
		}

		clause.Write(fields.Bytes()[1:])
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

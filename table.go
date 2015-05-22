package yorm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

type tableSetter struct {
	table    string
	dests    []interface{}
	columns  []*column
	pkColumn *column
}

var (
	//TimeType time's reflect type.
	TimeType = reflect.TypeOf(time.Time{})

	// one struct reflect to a table query setter
	tableMap = map[reflect.Value]*tableSetter{}
	//table lock
	tableRWLock sync.RWMutex
)

func newTableSetter(ri reflect.Value) (*tableSetter, error) {
	if q, ok := tableMap[ri]; ok {
		return q, nil
	}
	tableRWLock.Lock()
	defer tableRWLock.Unlock()
	if q, ok := tableMap[ri]; ok {
		return q, nil
	}
	if ri.Kind() != reflect.Ptr {
		return nil, ErrNonPtr
	}
	if ri.IsNil() {
		return nil, ErrNotSupported
	}
	q := new(tableSetter)
	table, cs := structToTable(reflect.Indirect(ri).Interface())
	var err error
	q.pkColumn, err = findPkColumn(cs)
	if q.pkColumn == nil {
		tableMap[ri] = nil
		return nil, err
	}
	q.table = table
	q.columns = cs
	q.dests = make([]interface{}, len(cs))
	for k, v := range cs {
		q.dests[k] = newPtrInterface(v.typ)
	}
	tableMap[ri] = q
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
		if v.isPK {
			if isPk {
				return c, ErrDuplicatePkColumn
			}
			isPk = true
			c = v
		}
	}
	if c == nil && idColumn != nil {
		idColumn.isPK = true
		idColumn.isAuto = true
		c = idColumn
	}
	if c == nil {
		return nil, ErrNonePkColumn
	}
	return c, nil
}

func newPtrInterface(t reflect.Type) interface{} {
	k := t.Kind()
	var ti interface{}
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ti = new(sql.NullInt64)
	case reflect.String:
		ti = new(sql.NullString)
	case reflect.Float32, reflect.Float64:
		ti = new(sql.NullFloat64)
	case reflect.Struct:
		switch t {
		case TimeType:
			ti = new(sql.NullString)
		}
	}
	return ti
}

func scanValue(sc sqlScanner, q *tableSetter, st reflect.Value) error {
	err := sc.Scan(q.dests...)
	if err != nil {
		return err
	}
	for idx, c := range q.columns {
		// different assign func here
		fv := st.Field(c.fieldNum)
		fi := q.dests[idx]
		err := setValue(fv, fi)
		if err != nil {
			fmt.Println(q.table, c.name, fv, fi, err)
			continue
		}
	}
	return nil
}

func setValue(fv reflect.Value, fi interface{}) error {
	switch fv.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		sqlValue := sql.NullInt64(*(fi.(*sql.NullInt64)))
		if !sqlValue.Valid {
			return errors.New("sqlValue is invalid")
		}
		fv.SetInt(sqlValue.Int64)
	case reflect.String:
		sqlValue := sql.NullString(*(fi.(*sql.NullString)))
		if !sqlValue.Valid {
			return errors.New("sqlValue is invalid")
		}
		fv.SetString(sqlValue.String)
	case reflect.Float32, reflect.Float64:
		sqlValue := sql.NullFloat64(*(fi.(*sql.NullFloat64)))
		if !sqlValue.Valid {
			return errors.New("sqlValue is invalid")
		}
		fv.SetFloat(sqlValue.Float64)
	case reflect.Struct:
		switch fv.Type() {
		case TimeType:
			sqlValue := sql.NullString(*(fi.(*sql.NullString)))
			if !sqlValue.Valid {
				return errors.New("sqlValue is invalid")
			}
			timeStr := sqlValue.String
			var layout string
			if len(timeStr) == 10 {
				layout = shortSimpleTimeFormat
			}
			if len(timeStr) == 19 {
				layout = longSimpleTimeFormat
			}
			timeTime, err := time.ParseInLocation(layout, timeStr, time.Local)
			if timeTime.IsZero() {
				return err
			}
			fv.Set(reflect.ValueOf(timeTime))
		}
	}
	return nil
}

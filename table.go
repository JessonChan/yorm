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

var (
	//TimeType time's reflect type.
	TimeType = reflect.TypeOf(time.Time{})

	// one struct reflect to a table query setter
	tableMap = map[reflect.Value]*tableSetter{}
)

func newTableSetter(ri reflect.Value) (*tableSetter, error) {
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
	defer func() {
		tableMap[ri] = q
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
		if v.isPK {
			if isPk {
				return c, ErrDuplicatePkColumn
			}
			isPk = true
			c = v
		}
	}
	if c == nil {
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
	case reflect.Int:
		ti = new(int)
	case reflect.Int64:
		ti = new(int64)
	case reflect.String:
		ti = new(string)
	case reflect.Float32:
		ti = new(float32)
	case reflect.Float64:
		ti = new(float64)
	case reflect.Struct:
		switch t {
		case TimeType:
			ti = new(string)
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
		switch c.typ.Kind() {
		case reflect.Int:
			st.Field(c.fieldNum).SetInt(int64(*(q.dests[idx].(*int))))
		case reflect.Int64:
			st.Field(c.fieldNum).SetInt(int64(*(q.dests[idx].(*int64))))
		case reflect.String:
			st.Field(c.fieldNum).SetString(string(*(q.dests[idx].(*string))))
		case reflect.Float32:
			st.Field(c.fieldNum).SetFloat(float64(*(q.dests[idx].(*float32))))
		case reflect.Float64:
			st.Field(c.fieldNum).SetFloat(float64(*(q.dests[idx].(*float64))))
		case reflect.Struct:
			switch c.typ {
			case TimeType:
				timeStr := string(*(q.dests[idx].(*string)))
				var layout string
				if len(timeStr) == 10 {
					layout = "2006-01-02"
				}
				if len(timeStr) == 19 {
					layout = "2006-01-02 15:04:05"
				}
				timeTime, err := time.ParseInLocation(layout, timeStr, time.Local)
				if timeTime.IsZero() {
					return err
				}
				st.Field(c.fieldNum).Set(reflect.ValueOf(timeTime))
			}
		}
	}
	return nil
}

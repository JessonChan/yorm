package yorm

import (
	"reflect"
	"time"
)

//field name
const (
	CamelToUnderscore = iota
	FieldName
)

// A column  represents a single column on a db record
type column struct {
	fieldNum  int
	fieldName string
	name      string
	typ       reflect.Type
	isInner   bool //inner struct ?
	isPK      bool
}

var structColumnCache = map[reflect.Type][]*column{}

func structToTable(i interface{}) (tableName string, columns []*column) {
	typ := reflect.TypeOf(i)
	if typ.Kind() != reflect.Struct {
		return
	}
	return camelToUnderscore(typ.Name()), structColumns(typ)
}

func structColumns(t reflect.Type) (columns []*column) {
	if t.Kind() != reflect.Struct {
		return
	}
	if cs, ok := structColumnCache[t]; ok {
		return cs
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		//unexpected struct type,ommit
		if field.PkgPath != "" {
			continue
		}

		fieldType := field.Type
		tag := parseTag(field.Tag.Get("yorm"))
		if tag.skip {
			continue
		}
		//todo if ft is ptr'ptr or three deep ptr?
		if fieldType.Name() == "" && fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		name := camelToUnderscore(field.Name)
		var isInner bool
		if tag.columnIsSet {
			if tag.columnName != "" {
				name = tag.columnName
			}
		} else {
			if fieldType.Kind() == reflect.Struct {
				if fieldType.Kind() == reflect.TypeOf(time.Time{}).Kind() {
				} else {
					isInner = true
				}

			}
		}
		c := &column{
			fieldNum:  i,
			fieldName: field.Name,
			name:      name,
			typ:       fieldType,
			isInner:   isInner,
			isPK:      tag.pkIsSet,
		}
		if c.isInner {

			// recursive unwind  inner struct
			columns = append(columns, structColumns(c.typ)...)
		} else {
			columns = append(columns, c)
		}
	}
	if len(columns) > 0 {
		structColumnCache[t] = columns
	}
	return
}

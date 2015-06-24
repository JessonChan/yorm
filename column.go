package yorm

import "reflect"

//field name
const (
	CamelToUnderscore = iota
	FieldName
)

var structColumnCache = map[reflect.Type][]*column{}

// A column  represents a single column on a db record
type column struct {
	fieldNum  int
	fieldName string
	name      string
	typ       reflect.Type
	isPK      bool
	isAuto    bool
}

type YormTableStruct interface {
	YormTableName() string
}

func structToTable(i interface{}) (tableName string, columns []*column) {
	typ := reflect.TypeOf(i)
	if typ.Kind() != reflect.Struct {
		return
	}
	if yt, ok := i.(YormTableStruct); ok {
		tableName = yt.YormTableName()
	} else {
		tableName = camelToUnderscore(typ.Name())
		tableName = tableFunc(tableName)
	}
	return tableName, structColumns(typ)
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
		//unexpected struct type,omit
		if field.PkgPath != "" {
			continue
		}
		fieldType := field.Type
		if fieldType.Kind() == reflect.Struct {
			if fieldType != TimeType {
				continue
			}
		}
		tag := parseTag(field.Tag.Get("yorm"))
		if tag.skip {
			continue
		}
		name := camelToUnderscore(field.Name)
		if tag.columnIsSet {
			if tag.columnName != "" {
				name = tag.columnName
			}
		}
		c := &column{
			fieldNum:  i,
			fieldName: field.Name,
			name:      name,
			typ:       fieldType,
			isPK:      tag.pkIsSet,
			isAuto:    tag.autoIsSet,
		}
		columns = append(columns, c)
	}
	if len(columns) > 0 {
		structColumnCache[t] = columns
	}
	return
}

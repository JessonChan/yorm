package yorm

import "reflect"

// A column  represents a single column on a db record
type column struct {
	name   string
	typ    reflect.Type
	follow bool
}

var structColumnCache map[reflect.Type][]column

func structColumns(t reflect.Type) (columns []column) {
	if t.Kind() != reflect.Struct {
		return
	}
	if cs, ok := structColumnCache[t]; ok {
		return cs
	}
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		//unexpected struct type,ommit
		if tf.PkgPath != "" {
			continue
		}
		ft := tf.Type
		tag := parseTag(tf.Tag.Get("yorm"))
		if tag.skip {
			continue
		}
		//todo if ft is ptr'ptr or three deep ptr?
		if ft.Name() == "" && ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		name := tf.Name
		var follow bool
		if tag.columnIsSet {
			if tag.columnName != "" {
				name = tag.columnName
			}
		} else {
			if ft.Kind() == reflect.Struct {
				follow = true
			}
		}
		c := column{name: name, typ: ft, follow: follow}
		if c.follow {
			columns = append(columns, structColumns(c.typ)...)
		} else {
			columns = append(columns, c)
		}
	}
	if len(columns) > 0 {
		if structColumnCache == nil {
			structColumnCache = make(map[reflect.Type][]column)
		}
		structColumnCache[t] = columns
	}
	return
}

package yorm

import "reflect"

// A column  represents a single column on a db record
type column struct {
	name string
	typ  reflect.Type
}

var structColumnCache map[reflect.Type][]column

func structColumns(t reflect.Type) (columns []column) {
	if cs, ok := structColumnCache[t]; ok {
		return cs
	}
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		//unexpected struct fields,ommit
		if tf.PkgPath != "" {
			continue
		}
		ft := tf.Type
		if ft.Name() == "" && ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		columns = append(columns, column{name: tf.Name, typ: ft})
	}
	if len(columns) > 0 {
		if structColumnCache == nil {
			structColumnCache = make(map[reflect.Type][]column)
		}
		structColumnCache[t] = columns
	}
	return
}

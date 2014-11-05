package yorm

import (
	"fmt"
	"reflect"
	"testing"
)

func dumpColumns(cs []column, t *testing.T) {
	for k, v := range cs {
		t.Log(fmt.Sprintf("%d,%s,%v", k, v.name, v.typ))
	}
}
func TestStructColumns_2(t *testing.T) {
	type A struct {
		UserName string
		Uid      string
		name     string
	}
	cs := structColumns(reflect.TypeOf(A{}))
	if len(cs) != 2 {
		t.Error("structColumns A struct ")
		t.Fail()
	}
	dumpColumns(structColumns(reflect.TypeOf(A{})), t)
}
func TestStructColumns_1(t *testing.T) {
	cs := structColumns(
		reflect.TypeOf(struct {
			Id   int
			Name string
		}{}))
	if len(cs) != 2 {
		t.Fail()
	}
	dumpColumns(cs, t)
}

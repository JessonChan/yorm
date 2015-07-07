package yorm

import (
	"fmt"
	"reflect"
	"testing"
)

func TestStructToTable_1(t *testing.T) {
	i := 1
	n, cs := structToTable(i)
	assert(t, n == "")
	assert(t, len(cs) == 0)
}

func TestStructToTable_2(t *testing.T) {
	a := struct {
		Name string
	}{}
	testStructToTable(t, a, "", 1)
}
func TestStructToTable_3(t *testing.T) {
	type A struct {
		Id   int64
		Name string
	}
	testStructToTable(t, A{}, "a", 2)
}
//func TestStructToTable_4(t *testing.T) {
//	type A struct {
//		Id   int64
//		Name string
//	}
//	type B struct {
//		A
//	}
//	type C struct {
//		A `yorm:"-"`
//		B `yorm:"column(test)"`
//	}
//	type D struct {
//		B
//		Log string `yorm:"column(log2)"`
//	}
//	testStructToTable(t, B{}, "b", 2)
//	testStructToTable(t, C{}, "c", 1)
//	testStructToTable(t, D{}, "d", 3)
//	_, cs := structToTable(D{})
//	for _, v := range cs {
//		t.Log(v.name)
//		for _, v1 := range []string{"id", "name", "log2"} {
//			if v.name == v1 {
//				goto S
//			}
//		}
//		t.Fail()
//	S:
//	}
//}

func testStructToTable(t *testing.T, i interface{}, name string, numField int) {
	n, cs := structToTable(i)
	assert(t, n == name)
	assert(t, len(cs) == numField)
}
func TestStructColumns_2(t *testing.T) {
	type A struct {
		UserName string
		Uid      string
		name     string
	}
	testStructColumns(A{}, 2, t)
}
func TestStructColumns_1(t *testing.T) {
	testStructColumns(struct {
		Id   int
		Name string
	}{}, 2, t)
}

func TestStructColumns_3(t *testing.T) {
	testStructColumns(struct{ Name *string }{}, 1, t)
}
//func TestStructColumns_4(t *testing.T) {
//	type A struct {
//		Id int
//	}
//	type B struct {
//		A
//	}
//	testStructColumns(B{}, 1, t)
//}

func TestStructColumns_5(t *testing.T) {
	type A struct {
		Id int
	}
	type B struct {
		A `yorm:"-"`
	}
	testStructColumns(B{}, 0, t)
}

//func TestStructColumns_6(t *testing.T) {
//	type A struct {
//		Id int
//	}
//	type B struct {
//		A `yorm:"column(user_name)"`
//	}
//	testStructColumns(B{}, 1, t)
//}
func testStructColumns(itf interface{}, numField int, t *testing.T) {
	cs := structColumns(reflect.TypeOf(itf))
	for k, v := range cs {
		t.Log(fmt.Sprintf("%d,%s,%v", k, v.name, v.typ))
	}
	if len(cs) != numField {
		t.Fail()
	}
}

package yorm

import (
	"fmt"
	"reflect"
	"testing"
)

type Table struct {
	Id int64
}

func (t *Table) YormDynamicTableName() string {
	return fmt.Sprintf("table_%v", t.Id)
}

func TestNewTableSetter_1(t *testing.T) {
	t1 := &Table{Id: 1}
	q, _ := newTableSetter(reflect.ValueOf(t1))
	t.Log(q.table, t1.YormDynamicTableName())
	if q.table != t1.YormDynamicTableName() {
		t.FailNow()
	}
	t2 := &Table{Id: 2}
	q, _ = newTableSetter(reflect.ValueOf(t2))
	if q.table != t2.YormDynamicTableName() {
		t.FailNow()
	}
}

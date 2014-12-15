package yorm

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type Movie struct {
	Id  int
	Id1 int `yorm:"column(mid)"`
	Id2 int `yorm:"column(did)"`
}

func TestQuery_1(t *testing.T) {
	s := query(structToTable(Movie{}))
	t.Log(s)
}
func TestQuery_2(t *testing.T) {
	dbpath := "root:@tcp(127.0.0.1:3306)/yorm_test?charset=utf8"
	db, err := sql.Open("mysql", dbpath)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	mapping := Movie{Id: 1}
	q := newQuery(&Movie{})
	t.Log(q)
	q.Where("id = 1")
	rows, err := db.Query(q.String())
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer rows.Close()
	//	berforeAssign(rows, q)
	convertAssign(&mapping, rows, q)
	t.Log(mapping)
	//  mapping.Id1 = *(q.dests[1].(*int))
	//  mapping.Id2 = *(q.dests[2].(*int))
	//  t.Log(mapping)
}

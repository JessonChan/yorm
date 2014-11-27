package yorm

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type MtMovieMapping struct {
	Id       int
	MtimeId  int `yorm:"column(mtimeId)"`
	DoubanId int `yorm:"column(DoubanId)"`
}

func TestQurey_1(t *testing.T) {
	s := query(structToTable(MtMovieMapping{}))
	t.Log(s)
}
func TestQurey_2(t *testing.T) {
	dbpath := "q3boy:123@tcp(192.168.2.218:3306)/movie_crawler?charset=utf8"
	db, err := sql.Open("mysql", dbpath)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	mapping := MtMovieMapping{Id: 1}
	q := newQuery(MtMovieMapping{})
	t.Log(q)
	q.Where("id = 1")
	rows, err := db.Query(q.String())
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(q.dests...)
		if err != nil {
			return
		}
	}
	mapping.MtimeId = *(q.dests[1].(*int))
	mapping.DoubanId = *(q.dests[2].(*int))
	t.Log(mapping)
}

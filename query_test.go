package yorm

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func Test_DB_Query(t *testing.T) {
	dbpath := "root:@tcp(127.0.0.1:3306)/yorm_test?charset=utf8"
	m := &Movie{}
	SetDb(dbpath)
	err := Q(m)
	t.Log(err)
	t.Log(*m)
}

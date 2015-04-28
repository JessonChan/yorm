package yorm

import (
	"database/sql"
	"reflect"
)

var sqlDb *sql.DB
var tableMap map[reflect.Kind]*querySetter = make(map[reflect.Kind]*querySetter)
var stmtMap map[string]*sql.Stmt = make(map[string]*sql.Stmt)

const (
	MYSQL = "mysql"
)

func Register(dbPath string) error {
	var err error
	sqlDb, err = sql.Open(MYSQL, dbPath)
	if sqlDb == nil {
		return err
	}
	return sqlDb.Ping()
}

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

func getStmt(clause string) (*sql.Stmt, error) {
	stmt := stmtMap[clause]
	var err error
	if stmt == nil {
		stmt, err = sqlDb.Prepare(clause)
		if stmt != nil {
			stmtMap[clause] = stmt
		}
	}
	return stmt, err
}

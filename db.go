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
	var err error
	clause, err = validClause(clause)
	if err != nil {
		return err
	}
	stmt := stmtMap[clause]
	if stmt == nil {
		stmt, err = sqlDb.Prepare(clause)
		if stmt != nil {
			stmtMap[clause] = stmt
		}
	}
	return stmt, err
}

func validClause(clause string) (string, error) {
	return clause, nil
}

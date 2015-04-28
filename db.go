package yorm

import (
	"database/sql"
	"reflect"
)

const (
	DriverMySQL = "mysql"

	//DriverDefault is  the default driver(mysql)
	DriverDefault = DriverMySQL
)

var sqlDb *sql.DB
var tableMap = map[reflect.Kind]*querySetter{}
var stmtMap = map[string]*sql.Stmt{}

// Register register a database driver.
func Register(dsn string, driver ...string) error {
	var err error

	if len(driver) != 0 {
		sqlDb, err = sql.Open(driver[0], dsn)

	} else {
		sqlDb, err = sql.Open(DriverDefault, dsn)
	}

	if sqlDb == nil {
		return err
	}
	return sqlDb.Ping()
}

func getStmt(clause string) (*sql.Stmt, error) {
	var err error
	clause, err = validClause(clause)
	if err != nil {
		return nil, err
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

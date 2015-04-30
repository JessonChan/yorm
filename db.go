package yorm

import (
	"database/sql"
	"reflect"
	"sync"
)

const (
	DriverMySQL = "mysql"
)

// main db to operate ,maybe will support multi dbs(read/write ...)
var sqlDb *sql.DB
var dbMutex sync.RWMutex

// one struct reflect to a table query setter
var tableMap = map[reflect.Kind]*querySetter{}

// stmt to prepare db conn
var stmtMap = map[string]*sql.Stmt{}

// Register register a database driver.
func Register(dsn string, driver ...string) error {

	dbMutex.Lock()
	defer dbMutex.Unlock()

	var err error
	if sqlDb != nil {
		return nil
	}
	sqlDb, err = sql.Open(append(driver, DriverMySQL)[0], dsn)

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

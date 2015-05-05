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

type sqlExecutor interface {
	Select(i interface{}, clause string, args ...interface{}) error
	Insert(i interface{}, args ...string) (int64, error)
	Update(clause string, args ...interface{}) (int64, error)
	Delete(clause string, args ...interface{}) (int64, error)
}

type executor struct {
	*sql.DB
}

var defaultExecutor *executor

var executorMap = map[string]*executor{}

const defaultExecutorName = "default"

func RegisterWithName(dsn, name string, driver ...string) error {
	if executorMap[name] != nil {
		return nil
	}
	dbMutex.Lock()
	defer dbMutex.Unlock()
	sqlDb, err := sql.Open(append(driver, DriverMySQL)[0], dsn)
	if sqlDb == nil {
		return err
	}
	err = sqlDb.Ping()
	if err != nil {
		return err
	}
	executorMap[name] = &executor{sqlDb}
	return nil
}

// Register register a database driver.
func Register(dsn string, driver ...string) error {
	err := RegisterWithName(dsn, defaultExecutorName, driver...)
	if err == nil {
		defaultExecutor = executorMap[defaultExecutorName]
	}
	return err
}

func Using(name string) sqlExecutor {
	return executorMap[name]
}

func (this *executor) getStmt(clause string) (*sql.Stmt, error) {
	if this == nil {
		return nil, ErrNotInitDefaultExecutor
	}
	var err error
	clause, err = validClause(clause)
	if err != nil {
		return nil, err
	}
	stmt := stmtMap[clause]
	if stmt == nil {
		stmt, err = this.Prepare(clause)
		if stmt != nil {
			stmtMap[clause] = stmt
		}
	}
	return stmt, err
}

func getStmt(clause string) (*sql.Stmt, error) {
	return defaultExecutor.getStmt(clause)
}

func validClause(clause string) (string, error) {
	return clause, nil
}

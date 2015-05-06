package yorm

import (
	"database/sql"
	"reflect"
	"sync"
)

const (
	DriverMySQL = "mysql"
)

var dbMutex sync.RWMutex

// one struct reflect to a table query setter
var tableMap = map[reflect.Kind]*tableSetter{}

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

type nilSqlExecutor struct {
}

func (n nilSqlExecutor) Select(i interface{}, clause string, args ...interface{}) error {
	return ErrNilSqlExecutor
}
func (n nilSqlExecutor) Insert(i interface{}, args ...string) (int64, error) {
	return 0, ErrNilSqlExecutor
}
func (n nilSqlExecutor) Update(clause string, args ...interface{}) (int64, error) {
	return 0, ErrNilSqlExecutor
}
func (n nilSqlExecutor) Delete(clause string, args ...interface{}) (int64, error) {
	return 0, ErrNilSqlExecutor
}

func Using(name string) sqlExecutor {
	if e, ok := executorMap[name]; ok {
		return e
	}
	return nilSqlExecutor{}
}

func (this *executor) getStmt(clause string) (*sql.Stmt, error) {

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
	if defaultExecutor == nil {
		return nil, ErrNilSqlExecutor
	}
	return defaultExecutor.getStmt(clause)
}

func validClause(clause string) (string, error) {
	return clause, nil
}

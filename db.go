package yorm

import (
	"database/sql"
	"sync"
)

const (

	//DriverMySQL driver for mysql
	DriverMySQL         = "mysql"
	defaultExecutorName = "default"

	shortSimpleTimeFormat = "2006-01-02"
	longSimpleTimeFormat  = "2006-01-02 15:04:05"
)

var (
	dbMutex         sync.RWMutex
	defaultExecutor *executor
	// stmt to prepare db conn
	stmtMap     = map[string]*sql.Stmt{}
	executorMap = map[string]*executor{}
)

type sqlExecutor interface {
	SelectByPK(i interface{}, tableName ...string) error
	Select(i interface{}, clause string, args ...interface{}) error
	Insert(i interface{}, args ...string) (int64, error)
	Update(clause string, args ...interface{}) (int64, error)
	Delete(clause string, args ...interface{}) (int64, error)
}

type executor struct {
	*sql.DB
}

//RegisterWithName register a database dirver with specific name.
func RegisterWithName(dsn, name string, driver ...string) (err error) {
	if executorMap[name] != nil {
		return nil
	}
	dbMutex.Lock()
	defer dbMutex.Unlock()
	var sqlDb *sql.DB
	sqlDb, err = sql.Open(append(driver, DriverMySQL)[0], dsn)
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

//Register register a database driver.
func Register(dsn string, driver ...string) error {
	err := RegisterWithName(dsn, defaultExecutorName, driver...)
	if err == nil {
		defaultExecutor = executorMap[defaultExecutorName]
	}
	return err
}

type nilSqlExecutor struct {
}

func (n nilSqlExecutor) SelectByPK(i interface{}, args ...string) error {
	return ErrNilSqlExecutor
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

//Using using the executor with name, if not exist  nilSqlExecutor instead.
func Using(name string) sqlExecutor {
	if e, ok := executorMap[name]; ok {
		return e
	}
	return nilSqlExecutor{}
}

func (ex *executor) getStmt(clause string) (*sql.Stmt, error) {

	var err error
	clause, err = validClause(clause)
	if err != nil {
		return nil, err
	}
	stmt := stmtMap[clause]
	if stmt == nil {
		stmt, err = ex.Prepare(clause)
		if stmt != nil {
			stmtMap[clause] = stmt
		}
	}
	return stmt, err
}

func (ex *executor) exec(clause string, args ...interface{}) (sql.Result, error) {
	stmt, err := ex.getStmt(clause)
	if err != nil {
		return nil, err
	}
	yogger.Debug("%s;%v", clause, args)
	return stmt.Exec(args...)
}


func validClause(clause string) (string, error) {
	return clause, nil
}

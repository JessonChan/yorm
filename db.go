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
	Select(i interface{}, clause string, args ...interface{}) error
	Insert(i interface{}, args ...string) (int64, error)
	Update(clause string, args ...interface{}) (int64, error)
	Delete(clause string, args ...interface{}) (int64, error)
}

type ExecHandler func(clause string, args ...interface{}) (sql.Result, error)

type executor struct {
	*sql.DB
}

type tranExecutor struct {
	*sql.Tx
}

var tableFunc = func(table string) string {
	return table
}

/*
  table func 可以根据项目需要定制不同的model对应表名，比如添加前缀、后缀；删除model中的前缀等
  对于实现了YormTableStruct (	YormTableName() string)则不会去处理
*/
func RegisterTableFunc(fn func(string) string) {
	tableFunc = fn
}

//RegisterWithName register a database driver with specific name.
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
func (n nilSqlExecutor) SetMaxOpenConns(int) {
}
func (n nilSqlExecutor) SetMaxIdleConns(int) {
}

//Using using the executor with name, if not exist  nilSqlExecutor instead.
func Using(name string) sqlExecutor {
	if e, ok := executorMap[name]; ok {
		return e
	}
	return nilSqlExecutor{}
}

func Begin(name ...string) (sqlExecutor, error) {
	var err error
	name = append(name, defaultExecutorName)
	if e, ok := executorMap[name[0]]; ok {
		var tx *sql.Tx
		tx, err = e.Begin()
		if err == nil {
			return &tranExecutor{tx}, err
		}
	}
	return nilSqlExecutor{}, err
}
func Commit(e sqlExecutor) error {
	return e.(*tranExecutor).Commit()
}
func RollBack(e sqlExecutor) error {
	return e.(*tranExecutor).Rollback()
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
	// 如果是单句执行，不需要再使用stmt，防止过多的prepare
	if len(args) == 0 {
		log.Debug("%s", clause)
		return ex.Exec(clause)
	}
	log.Debug("%s;%v", clause, args)
	stmt, err := ex.getStmt(clause)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(args...)
}

func (ex *tranExecutor) exec(clause string, args ...interface{}) (sql.Result, error) {
	log.Debug("%s;%v", clause, args)
	return ex.Exec(clause, args...)
}

func validClause(clause string) (string, error) {
	return clause, nil
}

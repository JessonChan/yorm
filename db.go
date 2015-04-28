package yorm

import (
	"database/sql"
	"reflect"
)

//DBType standard for database's type such as: mysql, oralce, db2 and so on.
type DBType string

const (
	MySQL DBType = "mysql"

	//DBDefault the default database(mysql)
	DBDefault = MySQL
)

var sqlDb *sql.DB

var tableMap = map[reflect.Kind]*querySetter{}

var stmtMap = map[string]*sql.Stmt{}

// Register register a database driver.
func Register(dsn string, db ...DBType) error {
	var err error

	if len(db) != 0 {
		sqlDb, err = sql.Open(string(db[0]), dsn)

	} else {
		sqlDb, err = sql.Open(string(DBDefault), dsn)
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

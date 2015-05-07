package yorm

import (
	"errors"
	"strings"
)

func (ex *executor) Delete(clause string, args ...interface{}) (int64, error) {

	if !strings.HasPrefix(strings.ToUpper(clause), "DELETE") {
		return 0, errors.New("must be begin with delete keyword")
	}
	stmt, err := ex.getStmt(clause)
	if err != nil {
		return 0, err
	}
	r, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	id, err := r.RowsAffected()
	return id, err
}

//Delete delete record(s) from a table
func Delete(clause string, args ...interface{}) (int64, error) {
	return defaultExecutor.Delete(clause, args...)
}

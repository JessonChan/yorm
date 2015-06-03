package yorm

import (
	"errors"
	"strings"
)

func (ex *executor) Delete(clause string, args ...interface{}) (int64, error) {

	if !strings.HasPrefix(strings.ToUpper(clause), "DELETE") {
		return 0, errors.New("must be begin with delete keyword")
	}
	r, err := ex.exec(clause, args...)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

//Delete delete record(s) from a table
func Delete(clause string, args ...interface{}) (int64, error) {
	return defaultExecutor.Delete(clause, args...)
}

func (ex *tranExecutor) Delete(clause string, args ...interface{}) (int64, error) {

	if !strings.HasPrefix(strings.ToUpper(clause), "DELETE") {
		return 0, errors.New("must be begin with delete keyword")
	}
	r, err := ex.exec(clause, args...)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

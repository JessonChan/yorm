package yorm

import (
	"strings"
)

func (ex *executor) Update(clause string, args ...interface{}) (int64, error) {
	return update(ex.exec, clause, args...)
}
func (ex *tranExecutor) Update(clause string, args ...interface{}) (int64, error) {
	return update(ex.exec, clause, args...)
}

func update(exec ExecHandler, clause string, args ...interface{}) (int64, error) {
	if !strings.HasPrefix(strings.ToUpper(clause), "UPDATE") {
		return 0, ErrUpdateBadSql
	}
	r, err := exec(clause, args...)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

//Update update record(s)
func Update(clause string, args ...interface{}) (int64, error) {
	return defaultExecutor.Update(clause, args...)
}

package yorm

import "strings"

func (ex *executor) Update(clause string, args ...interface{}) (int64, error) {
	if !strings.HasPrefix(strings.ToUpper(clause), "UPDATE") {
		return 0, ErrUpdateBadSql
	}
	r, err := ex.exec(clause, args...)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

//Update update record(s)
func Update(clause string, args ...interface{}) (int64, error) {
	return defaultExecutor.Update(clause, args...)
}

package yorm

import "strings"

func (ex *executor) Update(clause string, args ...interface{}) (int64, error) {
	if !strings.HasPrefix(strings.ToUpper(clause), "UPDATE") {
		return 0, ErrUpdateBadSql
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

//Update update record(s)
func Update(clause string, args ...interface{}) (int64, error) {
	return defaultExecutor.Update(clause, args...)
}

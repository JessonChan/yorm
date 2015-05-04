package yorm

import (
	"errors"
	"strings"
)

func (this *executor) Delete(clause string, args ...interface{}) (int64, error) {
	if !strings.HasPrefix(clause, "delete ") {
		return 0, errors.New("delete clause must be start with delete keyword")
	}
	stmt, err := this.getStmt(clause)
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

func Delete(clause string, args ...interface{}) (int64, error) {
	return defaultExecutor.Delete(clause, args...)
}

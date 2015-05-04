package yorm

import (
	"errors"
	"strings"
)

func (this *executor) Update(clause string, args ...interface{}) (int64, error) {
	if !strings.HasPrefix(clause, "update ") {
		return 0, errors.New("update clause must be start with update keyword")
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

func Update(clause string, args ...interface{}) (int64, error) {
	return defaultExecutor.Update(clause, args...)
}

package yorm

import (
	"database/sql"
	"errors"
)

func Q(i interface{}, rows *sql.Rows) error {
//	q := newQuery(i)
	// todo
//	return convertAssign(i, rows, q)
	return nil
}

func QueryOne(i interface{}, row *sql.Row) error {
	if row == nil {
		return errors.New("nil row")
	}
	return convertAssignRow(i, row)
}
func QueryList(i interface{}, rows *sql.Rows) error {
	if rows == nil {
		return errors.New("rows nil")
	}
	return convertAssignRows(i, rows)
}

package yorm

import "database/sql"

func Q(i interface{}, rows *sql.Rows) error {
	q := newQuery(i)
	// todo
	return convertAssign(i, rows, q)
}

func QueryOne(i interface{}, row *sql.Row) error {
	q := newQuery(i)
	q.String()
	return nil
}

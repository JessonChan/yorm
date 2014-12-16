package yorm

import "database/sql"

var dbpath string

func SetDb(dp string) {
	dbpath = dp
}

func Q(i interface{}) error {
	db, err := sql.Open("mysql", dbpath)
	if err != nil {
		return err
	}
	defer db.Close()
	q := newQuery(i)
	// todo
	q.Where("id = 1")
	rows, err := db.Query(q.String())
	if err != nil {
		return err
	}
	defer rows.Close()
	return convertAssign(i, rows, q)
}

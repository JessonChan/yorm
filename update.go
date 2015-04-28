package yorm

func Update(i interface{}, clause string, args ...interface{}) (int64, error) {
	stmt, err := getStmt(clause)
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

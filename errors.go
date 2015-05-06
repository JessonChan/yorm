package yorm

import (
	"errors"
)

var (
	ErrIllegalParams          = errors.New("illegal parameter(s)")
	ErrNonPtr                 = errors.New("must be the pointer to the modified object")
	ErrNonSlice               = errors.New("must be a slice")
	ErrNotSupported           = errors.New("not supported now")
	ErrNotInitDefaultExecutor = errors.New("not init default sql executor")
	ErrNilMethodReceiver      = errors.New("the method receiver is nil.")
	ErrNilSqlExecutor         = errors.New("yorm not register the db config")
	ErrUpdateBadSql           = errors.New("must be begin with update keyword")
	ErrDuplicatePkColumn      = errors.New("duplicate pk column find")
	ErrNonePkColumn           = errors.New("none pk column find")
)

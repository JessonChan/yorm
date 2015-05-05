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
)

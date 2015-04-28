package yorm

import "errors"

var (
	ErrNotSupported  = errors.New("ErrNotSupported")
	ErrNonPtr        = errors.New("ErrNonPtr")
	ErrIllegalParams = errors.New("ErrIllegalParams")
	ErrNonSlice      = errors.New("ErrNonSlice")
)

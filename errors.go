package yorm

import (
	"errors"
	
)

var (
	ErrIllegalParams = errors.New("illegal parameter(s)")
	ErrNonPtr        = errors.New("must be the pointer to the modified object")
	ErrNonSlice      = errors.New("must be a slice")
	ErrNotSupported  = errors.New("Not supported now")
)

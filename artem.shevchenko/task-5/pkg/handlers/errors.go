package handlers

import "errors"

var (
	ErrCantBeDecorated    = errors.New("can't be decorated")
	ErrHandlersGroupWait  = errors.New("handlers group wait error")
)

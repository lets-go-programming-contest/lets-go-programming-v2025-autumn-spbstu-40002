package conveyer

import "errors"

var (
	ErrChanNotFound      = errors.New("chan not found")
	ErrHandlersGroupWait = errors.New("handlers group wait error")
)

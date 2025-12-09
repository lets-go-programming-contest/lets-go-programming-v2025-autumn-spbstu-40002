package conveyer

import "errors"

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrChanFull     = errors.New("channel is full")
	ErrNoData       = errors.New("no data available")
)

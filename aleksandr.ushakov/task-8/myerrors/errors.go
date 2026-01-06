package myerrors

import "errors"

var (
	ErrConfigRead  = errors.New("config read error")
	ErrConfigParse = errors.New("did not find expected key")
)

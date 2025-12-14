package myerrors

import "errors"

var (
	ErrNoChannel  = errors.New("chan not found")
	ErrNoDecorate = errors.New("can't be decorated")
)

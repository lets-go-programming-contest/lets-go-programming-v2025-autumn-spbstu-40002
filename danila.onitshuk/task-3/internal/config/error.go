package config

import "errors"

var (
	ErrPath = errors.New("invalid path to the config")
	ErrCfg  = errors.New("incorrect config")
)

package utils

import "errors"

var (
	ErrConfigRead  = errors.New("failed to read config")
	ErrConfigParse = errors.New("failed to parse config")
	ErrXMLRead     = errors.New("failed to read xml")
	ErrXMLParse    = errors.New("failed to parse xml")
	ErrJSONWrite   = errors.New("failed to write json")
	ErrDirCreate   = errors.New("failed to create directory")
)

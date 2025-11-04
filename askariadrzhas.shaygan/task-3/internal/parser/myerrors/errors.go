package myerrors

import "errors"

var (
	ErrConfigPath  = errors.New("config path not provided")
	ErrConfigRead  = errors.New("failed to read config file")
	ErrConfigParse = errors.New("invalid config: missing input/output fields")
	ErrParseXML    = errors.New("failed to parse XML file")
)

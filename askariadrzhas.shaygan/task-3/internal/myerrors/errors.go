package myerrors

import "errors"

var (
	ErrConfigPath   = errors.New("config path not provided")
	ErrConfigRead   = errors.New("config read error")
	ErrConfigParse  = errors.New("invalid config: missing input/output fields")
	ErrFileNotFound = errors.New("no such file or directory")
	ErrXMLRead      = errors.New("xml read error")
	ErrXMLDecode    = errors.New("xml decode error")
	ErrDirCreate    = errors.New("directory create error")
	ErrOutOpen      = errors.New("output file open error")
	ErrOutEncode    = errors.New("output file encode error")
)

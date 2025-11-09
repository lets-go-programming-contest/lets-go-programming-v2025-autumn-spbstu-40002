package myerrors

import "errors"

var (
	ErrConfigPath   = errors.New("config path error")
	ErrConfigRead   = errors.New("config read error")
	ErrConfigParse  = errors.New("did not find expected key")
	ErrFileNotFound = errors.New("no such file or directory")
	ErrXMLDecode    = errors.New("xml decode error")
	ErrDirCreate    = errors.New("directory create error")
	ErrOutOpen      = errors.New("output file open error")
	ErrOutEncode    = errors.New("output file encode error")
	ErrCloseFile    = errors.New("file close error")
)

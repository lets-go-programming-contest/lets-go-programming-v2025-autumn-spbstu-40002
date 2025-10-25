package parser

import "errors"

var (
	ErrFileNotFound = errors.New("no such file or directory")
	ErrCloseFile    = errors.New("failed to close file")
	ErrDecodeXML    = errors.New("failed to decode data from XML")
	ErrCreatDir     = errors.New("could not create a directory")
	ErrCreatJSON    = errors.New("could not create a json")
	ErrWriteJSON    = errors.New("could not write a json")
)

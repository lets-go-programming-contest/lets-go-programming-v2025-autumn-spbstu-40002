package parser

import "errors"

var (
	ErrFileNotFound = errors.New("file not found")
	ErrCloseFile    = errors.New("failed to close file")
	ErrDecodeXML    = errors.New("failed to decode data from XML")
	ErrCreatDir     = errors.New("")
	ErrCreatJson    = errors.New("")
	ErrWriteJson    = errors.New("")

	Err = errors.New("")
)

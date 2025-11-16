package myerrors

import (
	"errors"
)

var (
	ErrNoConfigFileFound         = errors.New("unable to open configuration file")
	ErrFailedToDeserializeConfig = errors.New("did not find expected key")
	ErrFailedToOpenXML           = errors.New("no such file or directory")
	ErrFailedToDecodeXML         = errors.New("failed to decode data from XML")
	ErrFailedToOpenOutputFile    = errors.New("failed to open output file")
	ErrFailedToSerializeJSON     = errors.New("failed to create json from data")
	ErrFailedToWriteData         = errors.New("failed to write output file")
	ErrFailedToCreateDir         = errors.New("failed to create directory")
	ErrFailedToCloseFile         = errors.New("failed to close file")
)

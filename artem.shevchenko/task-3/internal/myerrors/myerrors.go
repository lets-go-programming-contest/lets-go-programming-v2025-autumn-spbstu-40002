package myerrors

import (
	"errors"
)

var (
	ErrNoConfigFileFound         = errors.New("unable to open configuration file")
	ErrFailedToDeserializeConfig = errors.New("failed to read configuration file")
	ErrNoConfigFileProvided      = errors.New("the configuration file was not provided")
	ErrFailedToOpenXML           = errors.New("failed to open XML file")
	ErrFailedToDecodeXML         = errors.New("failed to decode data from XML")
	ErrFailedToOpenOutputFile    = errors.New("failed to open output file")
	ErrFailedToSerializeJSON     = errors.New("failed to create json from data")
	ErrFailedToWriteData         = errors.New("failed to write output file")
	ErrFailedToCreateDir         = errors.New("failde to create directory")
	ErrNumCodeIsNotIneger        = errors.New("numcode is not integer")
	ErrValueIsNotFloat           = errors.New("value is not float")
)

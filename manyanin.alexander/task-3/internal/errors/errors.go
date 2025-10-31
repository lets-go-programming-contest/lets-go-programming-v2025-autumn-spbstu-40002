package errors

import (
	"errors"
)

var (
	ErrConfigRead      = errors.New("failed to read config")
	ErrConfigParse     = errors.New("failed to parse config")
	ErrConfigInvalid   = errors.New("invalid config")
	ErrConfigPathEmpty = errors.New("config path is empty")

	ErrFileRead          = errors.New("failed to read file")
	ErrFileWrite         = errors.New("failed to write file")
	ErrFileNotFound      = errors.New("file not found")
	ErrDirCreate         = errors.New("failed to create directory")
	ErrInputFileNotExist = errors.New("input file does not exist")

	ErrXMLRead   = errors.New("failed to read xml")
	ErrXMLParse  = errors.New("failed to parse xml")
	ErrXMLDecode = errors.New("failed to decode xml")
	ErrXMLEmpty  = errors.New("xml file is empty")

	ErrJSONWrite   = errors.New("failed to write json")
	ErrJSONMarshal = errors.New("failed to marshal json")

	ErrDataConvert           = errors.New("failed to convert data")
	ErrDataInvalid           = errors.New("invalid data")
	ErrDataEmpty             = errors.New("no data found")
	ErrNoCurrenciesExtracted = errors.New("no currencies extracted from xml")
)

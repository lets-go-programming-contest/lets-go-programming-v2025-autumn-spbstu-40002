package errors

import (
	"errors"
)

var (
	ErrConfigPathEmpty       = errors.New("config path is empty")
	ErrConfigRead            = errors.New("failed to read config")
	ErrConfigParse           = errors.New("failed to parse config")
	ErrConfigInvalid         = errors.New("invalid config")
	ErrInputFileNotExist     = errors.New("no such file or directory")
	ErrXMLRead               = errors.New("failed to read xml")
	ErrXMLDecode             = errors.New("failed to decode xml")
	ErrNoCurrenciesExtracted = errors.New("no currencies extracted from xml")
	ErrJSONWrite             = errors.New("failed to write json")
	ErrJSONMarshal           = errors.New("failed to marshal json")
	ErrDirCreate             = errors.New("failed to create directory")
	ErrDataEmpty             = errors.New("no data found")
	ErrXMLEmpty              = errors.New("xml file is empty")
	ErrXMLFileRead           = errors.New("error reading XML file")
	ErrOutputFileCreate      = errors.New("error creating output file")
)

package errors

import "errors"

var (
	ErrConfigPathNotSpecified = errors.New("config path not specified")
	ErrConfigFileRead         = errors.New("error reading config file")
	ErrConfigFieldsMissing    = errors.New("not all required fields specified in config file")
	ErrOutputDirCreate        = errors.New("error creating output directory")
	ErrInputFileNotExist      = errors.New("no such file or directory")
	ErrXMLFileRead            = errors.New("error reading XML file")
	ErrXMLDecode              = errors.New("error decoding XML")
	ErrNoCurrenciesExtracted  = errors.New("failed to extract currency data from XML file")
	ErrOutputFileCreate       = errors.New("error creating output file")
	ErrJSONEncode             = errors.New("error encoding JSON")
)

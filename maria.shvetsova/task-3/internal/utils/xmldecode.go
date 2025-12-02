package utils

import (
	"encoding/xml"
	"os"

	"golang.org/x/net/html/charset"
)

func GetCurrencyData(inputFile string) (*ValCurs, error) {
	fileInfo, err := os.Stat(inputFile)
	if os.IsNotExist(err) {
		return nil, errFileDoesntExist
	}

	if fileInfo.Size() == 0 {
		return nil, errEmptyFile
	}

	data, err := os.Open(inputFile)
	if err != nil {
		return nil, errOpeningFile
	}

	defer func() {
		if err = data.Close(); err != nil {
			panic(err)
		}
	}()

	decoder := xml.NewDecoder(data)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, errDecoding
	}

	return &valCurs, nil
}

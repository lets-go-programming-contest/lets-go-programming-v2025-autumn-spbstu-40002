package core

import (
	"encoding/xml"
	"os"

	"golang.org/x/net/html/charset"
)

func GetCurrencyData(inputFile string) (*ValCurs, error) {
	info, err := os.Stat(inputFile)
	if os.IsNotExist(err) {
		return nil, errFileDoesntExist
	}
	if info.Size() == 0 {
		return nil, errEmptyFile
	}

	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var data ValCurs
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

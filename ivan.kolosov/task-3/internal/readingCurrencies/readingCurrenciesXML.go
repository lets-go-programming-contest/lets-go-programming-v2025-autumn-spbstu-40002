package readingcurrencies

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"

	cur "github.com/InsomniaDemon/task-3/internal/currenciesTypes"
	"golang.org/x/net/html/charset"
)

var (
	errParsingXML     = errors.New("error occurred while parsing xml file")
	errOpeningXMLFile = errors.New("error occurred while opening xml file")
	errClosingXMLFile = errors.New("error occurred while closing xml file")
)

func GetCurrencies(path string) (_ cur.Currencies, returnError error) {
	var currencies cur.Currencies

	file, err := os.Open(path)
	if err != nil {
		return currencies, fmt.Errorf("%w: %w", errOpeningXMLFile, err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			if returnError != nil {
				returnError = fmt.Errorf("%w: %w; %w", returnError, err, errClosingXMLFile)
			} else {
				returnError = fmt.Errorf("%w: %w", errClosingXMLFile, err)
			}
		}
	}()

	dc := xml.NewDecoder(file)
	dc.CharsetReader = charset.NewReaderLabel

	err = dc.Decode(&currencies)
	if err != nil {
		return currencies, fmt.Errorf("%w: %w", errParsingXML, err)
	}

	return currencies, nil
}

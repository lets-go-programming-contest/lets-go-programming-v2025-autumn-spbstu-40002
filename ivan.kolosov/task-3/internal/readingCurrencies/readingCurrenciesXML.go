package readingCurrencies

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

var (
	errParsingXML     = errors.New("error occurred while parsing xml file")
	errParsingFloat   = errors.New("error occurred while parsing float")
	errOpeningXMLFile = errors.New("error occurred while opening xml file")
	errClosingXMLFile = errors.New("error occurred while closing xml file")
)

func (cur *CurrencyXML) UnmarshalXML(dc *xml.Decoder, start xml.StartElement) error {
	err := dc.DecodeElement(&cur, &start)
	if err != nil {
		return errParsingXML
	}

	s := strings.ReplaceAll(cur.Value, ",", ".")

	cur.ValueFloat, err = strconv.ParseFloat(s, 64)
	if err != nil {
		return errParsingFloat
	}

	s = strings.ReplaceAll(cur.VunitRate, ",", ".")

	cur.VunitRateFloat, err = strconv.ParseFloat(s, 64)
	if err != nil {
		return errParsingFloat
	}

	return nil
}

func GetCurrencies(path string) (cur CurrenciesXML, returnError error) {
	file, err := os.Open(path)
	if err != nil {
		return cur, errOpeningXMLFile
	}

	defer func() {
		err := file.Close()
		if err != nil {
			if returnError != nil {
				returnError = fmt.Errorf("%w; %v", returnError, errClosingXMLFile)
			} else {
				returnError = errClosingXMLFile
			}
		}
	}()

	dc := xml.NewDecoder(file)
	dc.CharsetReader = charset.NewReaderLabel
	err = dc.Decode(&cur)
	if err != nil {
		return cur, errParsingXML
	}

	return cur, nil
}

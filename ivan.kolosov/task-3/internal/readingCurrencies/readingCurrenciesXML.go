package readingcurrencies

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

func (cur *CurrencyXML) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var temp struct {
		ID        string `xml:"ID,attr"`
		NumCode   int    `xml:"NumCode"`
		CharCode  string `xml:"CharCode"`
		Nominal   int    `xml:"Nominal"`
		Name      string `xml:"Name"`
		Value     string `xml:"Value"`
		VunitRate string `xml:"VunitRate"`
	}

	err := dec.DecodeElement(&temp, &start)
	if err != nil {
		return fmt.Errorf("%w: %w", errParsingXML, err)
	}

	str := strings.ReplaceAll(temp.Value, ",", ".")
	if str == "" {
		cur.ValueFloat = 0.0
	} else {
		cur.ValueFloat, err = strconv.ParseFloat(str, 64)
		if err != nil {
			return fmt.Errorf("%w: %w", errParsingFloat, err)
		}
	}

	str = strings.ReplaceAll(temp.VunitRate, ",", ".")
	if str == "" {
		cur.VunitRateFloat = 0.0
	} else {
		cur.VunitRateFloat, err = strconv.ParseFloat(str, 64)
		if err != nil {
			return fmt.Errorf("%w: %w", errParsingFloat, err)
		}
	}

	cur.ID = temp.ID
	cur.NumCode = temp.NumCode
	cur.CharCode = temp.CharCode
	cur.Nominal = temp.Nominal
	cur.Name = temp.Name
	cur.Value = temp.Value
	cur.VunitRate = temp.VunitRate

	return nil
}

func GetCurrencies(path string) (_ CurrenciesXML, returnError error) {
	var cur CurrenciesXML

	file, err := os.Open(path)
	if err != nil {
		return cur, fmt.Errorf("%w: %w", errOpeningXMLFile, err)
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

	err = dc.Decode(&cur)
	if err != nil {
		return cur, fmt.Errorf("%w: %w", errParsingXML, err)
	}

	return cur, nil
}

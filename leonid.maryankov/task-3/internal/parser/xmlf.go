package parser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID       string  `xml:"ID,attr"`
	NumCode  int     `xml:"NumCode"`
	CharCode string  `xml:"CharCode"`
	Nominal  int     `xml:"Nominal"`
	ValueStr string  `xml:"Value"`
	Value    float64 `xml:"-"`
}

func ParseXML(path string) ([]Valute, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("xml reading error: %s: %w", path, err)
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = charset.NewReaderLabel

	var curs ValCurs
	if err := decoder.Decode(&curs); err != nil {
		return nil, fmt.Errorf("xml parsing error: %s: %w", path, err)
	}

	for ind := range curs.Valutes {
		value := &curs.Valutes[ind]

		str := strings.TrimSpace(strings.ReplaceAll(value.ValueStr, ",", "."))
		if str == "" {
			return nil, fmt.Errorf("empty value for currency %s", value.CharCode)
		}

		fl, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, fmt.Errorf("value parsing error for %s: %w", value.CharCode, err)
		}

		value.Value = fl
	}

	return curs.Valutes, nil
}

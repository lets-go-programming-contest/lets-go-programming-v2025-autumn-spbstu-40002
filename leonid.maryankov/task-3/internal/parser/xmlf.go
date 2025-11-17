package parser

import (
	"bytes"
	"encoding/xml"
	"errors"
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
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	ValueStr string  `json:"-"         xml:"Value"`
	Value    float64 `json:"value"     xml:"-"`
}

var errEmptyValue = errors.New("empty value for currency")

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
			return nil, fmt.Errorf("%w %s", errEmptyValue, value.CharCode)
		}

		fl, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, fmt.Errorf("value parsing error for %s: %w", value.CharCode, err)
		}

		value.Value = fl
	}

	return curs.Valutes, nil
}

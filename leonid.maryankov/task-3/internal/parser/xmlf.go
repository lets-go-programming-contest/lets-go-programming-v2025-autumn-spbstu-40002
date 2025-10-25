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

	for i := range curs.Valutes {
		v := &curs.Valutes[i]

		s := strings.TrimSpace(strings.ReplaceAll(v.ValueStr, ",", "."))
		if s == "" {
			return nil, fmt.Errorf("empty value for currency %s", v.CharCode)
		}

		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("value parsing error for %s: %w", v.CharCode, err)
		}

		denom := 1
		if v.Nominal > 0 {
			denom = v.Nominal
		}

		v.Value = f / float64(denom)
	}

	return curs.Valutes, nil
}

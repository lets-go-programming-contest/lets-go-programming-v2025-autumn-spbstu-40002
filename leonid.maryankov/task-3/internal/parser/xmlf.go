package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
)

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int     `xml:"NumCode"`
	CharCode string  `xml:"CharCode"`
	ValueStr string  `xml:"Value"`
	Value    float64 `xml:"-"`
}

func ParseXML(path string) ([]Valute, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("XML reading error: " + path + err.Error())
	}

	var curs ValCurs
	if err := xml.Unmarshal(data, &curs); err != nil {
		return nil, fmt.Errorf("XML parsing error: " + path + err.Error())
	}

	for i := range curs.Valutes {
		v := &curs.Valutes[i]
		v.ValueStr = replaceComma(v.ValueStr)
		v.Value, _ = strconv.ParseFloat(v.ValueStr, 64)
	}

	return curs.Valutes, nil
}

func replaceComma(s string) string {
	for i := range s {
		if s[i] == ',' {
			s = s[:i] + "." + s[i+1:]
		}
	}
	return s
}

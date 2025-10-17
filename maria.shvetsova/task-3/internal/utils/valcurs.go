package utils

import (
	"encoding/xml"
	"sort"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

func (v *ValCurs) SortByValue() {
	sort.Sort(Valutes(v.Valutes))
}

func (v *ValCurs) ConvertToOutput() ([]Output, error) {
	var output []Output

	for _, valute := range v.Valutes {
		value, err := valute.GetFloatValue()
		if err != nil {
			return nil, errInvalidFormat
		}

		output = append(output, Output{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	return output, nil
}

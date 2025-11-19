package utils

import (
	"sort"
)

type ValCurs struct {
	Date    string   `xml:"Date,attr" json:"date"`
	Name    string   `xml:"name,attr" json:"name"`
	Valutes []Valute `xml:"Valute" json:"valutes"`
}

func (v *ValCurs) SortByValue() {
	sort.Sort(Valutes(v.Valutes))
}

func (v *ValCurs) ConvertToOutput() ([]Valute, error) {
	output := make([]Valute, 0, len(v.Valutes))

	for _, valute := range v.Valutes {

		output = append(output, Valute{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    valute.Value,
		})
	}

	return output, nil
}

package utils

import (
	"strconv"
	"strings"
)

type Valute struct {
	ID        string `xml:"ID,attr"`
	NumCode   int    `xml:"NumCode"`
	CharCode  string `xml:"CharCode"`
	Nominal   int    `xml:"Nominal"`
	Name      string `xml:"Name"`
	Value     string `xml:"Value"`
	VunitRate string `xml:"VunitRate"`
}

func (v *Valute) GetFloatValue() (float64, error) {
	normalized := strings.Replace(v.Value, ",", ".", -1)

	return strconv.ParseFloat(normalized, 64)
}

type Valutes []Valute

func (v Valutes) Len() int { return len(v) }

func (v Valutes) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

func (v Valutes) Less(i, j int) bool {
	valueI, err := v[i].GetFloatValue()
	if err != nil {
		panic(err)
	}

	valueJ, err := v[j].GetFloatValue()
	if err != nil {
		panic(err)
	}

	return valueI > valueJ
}

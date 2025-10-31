package utils

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID        string `xml:"ID,attr"`
	NumCode   int `xml:"NumCode"`
	CharCode  string `xml:"CharCode"`
	Nominal   int    `xml:"Nominal"`
	Name      string `xml:"Name"`
	Value     string `xml:"Value"`
	VunitRate string `xml:"VunitRate"`
}

func (v *Valute) ConvertValue() (float64, error) {
	strValue := strings.ReplaceAll(v.Value, ",", ".")

	value, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return 0, fmt.Errorf("parse float value: %w", err)
	}

	return value, nil
}

type Valutes []Valute

func (v Valutes) Len() int { return len(v) }

func (v Valutes) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

func (v Valutes) Less(indI, indJ int) bool {
	valueI, err := v[indI].ConvertValue()
	if err != nil {
		panic(err)
	}

	valueJ, err := v[indJ].ConvertValue()
	if err != nil {
		panic(err)
	}

	return valueI > valueJ
}

func (v *ValCurs) SortByValue() {
	sort.Sort(Valutes(v.Valutes))
}

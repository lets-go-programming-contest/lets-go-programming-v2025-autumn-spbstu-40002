package utils

import (
	"encoding/xml"
	"sort"
)

type ValCurs struct {
	Date    string   `json:"date"    xml:"Date,attr"`
	Name    string   `json:"name"    xml:"name,attr"`
	Valutes []Valute `json:"valutes" xml:"Valute"`
}

type Valute struct {
	ID        string  `xml:"ID,attr"`
	NumCode   int     `json:"num_code"  xml:"NumCode"`
	CharCode  string  `json:"char_code" xml:"CharCode"`
	Nominal   int     `xml:"Nominal"`
	Name      string  `xml:"Name"`
	Value     float64 `json:"value"     xml:"Value"`
	VunitRate string  `xml:"VunitRate"`
}

type Valutes []Valute

func (v Valutes) Len() int { return len(v) }

func (v Valutes) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

func (v Valutes) Less(indI, indJ int) bool {
	return v[indI].Value > v[indJ].Value
}

func (v *ValCurs) SortByValue() {
	sort.Sort(Valutes(v.Valutes))
}

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
	Date    string   `xml:"Date,attr" json:"date"`
	Name    string   `xml:"name,attr" json:"name"`
	Valutes []Valute `xml:"Valute"    json:"valutes"`
}

type Valute struct {
	ID        string  `xml:"ID,attr"`
	NumCode   int     `xml:"NumCode"  json:"num_code"`
	CharCode  string  `xml:"CharCode" json:"char_code"`
	Nominal   int     `xml:"Nominal"`
	Name      string  `xml:"Name"`
	Value     float64 `xml:"Value"    json:"value"`
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

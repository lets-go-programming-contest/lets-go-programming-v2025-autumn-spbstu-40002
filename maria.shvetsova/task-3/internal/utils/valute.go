package utils

import (
	"fmt"
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
	normalized := strings.ReplaceAll(v.Value, ",", ".")

	value, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return 0, fmt.Errorf("parse float value: %w", err)
	}

	return value, nil
}

type Valutes []Valute

func (v Valutes) Len() int { return len(v) }

func (v Valutes) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

func (v Valutes) Less(indI, indJ int) bool {
	valueI, err := v[indI].GetFloatValue()
	if err != nil {
		panic(err)
	}

	valueJ, err := v[indJ].GetFloatValue()
	if err != nil {
		panic(err)
	}

	return valueI > valueJ
}

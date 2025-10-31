package xml

import (
	"fmt"
	"strconv"
	"strings"
)

type Currency struct {
	ID       string `xml:"ID,attr"`
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

func (c Currency) ToFloat() (float64, error) {
	normalized := strings.ReplaceAll(c.Value, ",", ".")
	value, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return 0, fmt.Errorf("parse float: %w", err)
	}
	return value, nil
}

type ByExchangeRate []Currency

func (b ByExchangeRate) Len() int {
	return len(b)
}

func (b ByExchangeRate) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByExchangeRate) Less(firstIndex, secondIndex int) bool {
	rateI, err := b[firstIndex].ToFloat()
	if err != nil {
		return false
	}
	rateJ, err := b[secondIndex].ToFloat()
	if err != nil {
		return false
	}
	return rateI > rateJ
}

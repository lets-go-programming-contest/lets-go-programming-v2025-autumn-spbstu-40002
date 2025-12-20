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

type CurrencyRecord struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func (c Currency) ToRecord() (CurrencyRecord, error) {
	value, err := c.ToFloat()
	if err != nil {
		return CurrencyRecord{}, fmt.Errorf("convert to record: %w", err)
	}

	return CurrencyRecord{
		NumCode:  c.NumCode,
		CharCode: c.CharCode,
		Value:    value,
	}, nil
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

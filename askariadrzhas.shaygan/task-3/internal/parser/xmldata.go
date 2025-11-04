package parser

import (
	"encoding/xml"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type ExchangeData struct {
	Items []CurrencyItem `xml:"Valute"`
}

type CurrencyItem struct {
	NumberCode string `xml:"NumCode"`
	Symbol     string `xml:"CharCode"`
	Rate       string `xml:"Value"`
}

type ProcessedCurrency struct {
	Code   string  `json:"numCode"`
	Symbol string  `json:"charCode"`
	Rate   float64 `json:"value"`
}

func ExtractCurrencyData(filePath string) []ProcessedCurrency {
	file, err := os.Open(filePath)
	if err != nil {
		panic("cannot open source file: " + err.Error())
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			panic("cannot close file: " + cerr.Error())
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var raw ExchangeData
	if err = decoder.Decode(&raw); err != nil {
		panic("invalid XML format: " + err.Error())
	}

	var out []ProcessedCurrency
	for _, item := range raw.Items {
		if p := convertCurrencyItem(item); p != nil {
			out = append(out, *p)
		}
	}

	if len(out) == 0 {
		panic("no valid currency data found")
	}

	return out
}

func convertCurrencyItem(item CurrencyItem) *ProcessedCurrency {
	clean := strings.Replace(item.Rate, ",", ".", 1)

	val, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		return nil
	}

	return &ProcessedCurrency{
		Code:   item.NumberCode,
		Symbol: item.Symbol,
		Rate:   val,
	}
}

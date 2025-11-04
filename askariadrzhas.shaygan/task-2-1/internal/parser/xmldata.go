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
	Code   string  `json:"num_code"`
	Symbol string  `json:"char_code"`
	Rate   float64 `json:"value"`
}

func ExtractCurrencyData(filePath string) []ProcessedCurrency {
	file, err := os.Open(filePath)
	if err != nil {
		panic("cannot open source file: " + err.Error())
	}
	defer file.Close()

	xmlDecoder := xml.NewDecoder(file)
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	var rawData ExchangeData
	err = xmlDecoder.Decode(&rawData)
	if err != nil {
		panic("invalid XML format: " + err.Error())
	}

	var processed []ProcessedCurrency
	for _, item := range rawData.Items {
		converted := convertCurrencyItem(item)
		if converted != nil {
			processed = append(processed, *converted)
		}
	}

	if len(processed) == 0 {
		panic("no valid currency data found")
	}

	return processed
}

func convertCurrencyItem(item CurrencyItem) *ProcessedCurrency {
	cleanedRate := strings.Replace(item.Rate, ",", ".", 1)
	rateValue, err := strconv.ParseFloat(cleanedRate, 64)
	if err != nil {
		return nil
	}

	return &ProcessedCurrency{
		Code:   item.NumberCode,
		Symbol: item.Symbol,
		Rate:   rateValue,
	}
}

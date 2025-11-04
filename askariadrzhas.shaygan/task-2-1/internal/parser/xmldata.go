// internal/parser/xmldata.go
package parser

import (
	"encoding/xml"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"

	"github.com/XShaygaND/task-3/internal/types"
)

type ExchangeData struct {
	Items []CurrencyItem `xml:"Valute"`
}

type CurrencyItem struct {
	NumberCode string `xml:"NumCode"`
	Symbol     string `xml:"CharCode"`
	Rate       string `xml:"Value"`
}

func ExtractCurrencyData(filePath string) []types.ProcessedCurrency {
	file, err := os.Open(filePath)
	if err != nil {
		panic("cannot open source file: " + err.Error())
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic("cannot close file: " + closeErr.Error())
		}
	}()

	xmlDecoder := xml.NewDecoder(file)
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	var rawData ExchangeData
	err = xmlDecoder.Decode(&rawData)
	if err != nil {
		panic("invalid XML format: " + err.Error())
	}

	var processed []types.ProcessedCurrency
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

func convertCurrencyItem(item CurrencyItem) *types.ProcessedCurrency {
	cleanedRate := strings.Replace(item.Rate, ",", ".", 1)
	rateValue, err := strconv.ParseFloat(cleanedRate, 64)
	if err != nil {
		return nil
	}

	return &types.ProcessedCurrency{
		Code:   item.NumberCode,
		Symbol: item.Symbol,
		Rate:   rateValue,
	}
}

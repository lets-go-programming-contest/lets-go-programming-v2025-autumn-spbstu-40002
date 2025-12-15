package xml

import (
	"fmt"
	"strconv"
	"strings"
)

type Currency struct {
	NumCode  int    `json:"num_code"  xml:"NumCode"`
	CharCode string `json:"char_code" xml:"CharCode"`
	Value    string `json:"value"     xml:"Value"`
}

func (currency Currency) GetFloat() (float64, error) {
	commaReplacement := strings.ReplaceAll(currency.Value, ",", ".")

	value, err := strconv.ParseFloat(commaReplacement, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot convert value %s to float: %w", commaReplacement, err)
	}

	return value, nil
}

type ByValue []Currency

func (currency ByValue) Len() int {
	return len(currency)
}

func (currency ByValue) Swap(i, j int) {
	currency[i], currency[j] = currency[j], currency[i]
}

func (currency ByValue) Less(iCurr, jCurr int) bool {
	currencyI, err := currency[iCurr].GetFloat()
	if err != nil {
		panic(err)
	}

	currencyJ, err := currency[jCurr].GetFloat()
	if err != nil {
		panic(err)
	}

	return currencyI > currencyJ
}

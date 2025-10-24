package utils

import (
	curxml "github.com/InsomniaDemon/task-3/internal/readingCurrencies"
	curjson "github.com/InsomniaDemon/task-3/internal/writingCurrencies"
)

func GetCurrenciesJSON(curIn curxml.CurrenciesXML) (curOut curjson.CurrenciesJSON) {
	for i, currency := range curIn.Currencies {
		curOut = append(curOut, curjson.CurrencyJSON{})
		curOut[i].NumCode = currency.NumCode
		curOut[i].CharCode = currency.CharCode
		curOut[i].Value = currency.ValueFloat
	}
	return curOut
}

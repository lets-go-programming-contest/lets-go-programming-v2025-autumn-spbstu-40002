package utils

import (
	curxml "github.com/InsomniaDemon/task-3/internal/readingCurrencies"
	curjson "github.com/InsomniaDemon/task-3/internal/writingCurrencies"
)

func GetCurrenciesJSON(curIn curxml.CurrenciesXML) curjson.CurrenciesJSON {
	var curOut curjson.CurrenciesJSON
	for _, currency := range curIn.Currencies {
		curOut = append(curOut, curjson.CurrencyJSON{
			NumCode:  currency.NumCode,
			CharCode: currency.CharCode,
			Value:    currency.ValueFloat,
		})
	}

	return curOut
}

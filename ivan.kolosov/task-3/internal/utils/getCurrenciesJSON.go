package utils

import (
	cur "github.com/InsomniaDemon/task-3/internal/currenciesTypes"
)

func GetCurrenciesJSON(curIn cur.Currencies) cur.Currencies {
	var curOut cur.Currencies
	for _, currency := range curIn.CurrencyArray {
		curOut.CurrencyArray = append(curOut.CurrencyArray, cur.Currency{
			NumCode:  currency.NumCode,
			CharCode: currency.CharCode,
			Value:    currency.Value,
		})
	}

	return curOut
}

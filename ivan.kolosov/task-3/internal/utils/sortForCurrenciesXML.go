package utils

import (
	"sort"

	cur "github.com/InsomniaDemon/task-3/internal/currenciesTypes"
)

func SortCurrenciesXML(cur *cur.Currencies) {
	sort.Slice(cur.CurrencyArray, func(i, j int) bool {
		return cur.CurrencyArray[i].Value > cur.CurrencyArray[j].Value
	})
}

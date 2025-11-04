package utils

import (
	"sort"

	"github.com/XShaygaND/task-3/internal/currency"
)

func SortValutesByValue(valutes []currency.Valute) {
	sort.Slice(valutes, func(i, j int) bool {
		return ParseCurrencyValue(valutes[i].Value) > ParseCurrencyValue(valutes[j].Value)
	})
}

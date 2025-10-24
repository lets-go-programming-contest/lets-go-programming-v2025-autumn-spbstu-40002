package utils

import (
	"sort"

	xmlcur "github.com/InsomniaDemon/task-3/internal/readingCurrencies"
)

func SortCurrenciesXML(cur *xmlcur.CurrenciesXML) {
	sort.Slice(cur.Currencies, func(i, j int) bool {
		return cur.Currencies[i].ValueFloat > cur.Currencies[j].ValueFloat
	})
}

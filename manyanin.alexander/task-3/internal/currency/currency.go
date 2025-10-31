package converter

import (
	"sort"
	"strconv"
	"strings"

	"github.com/manyanin.alexander/task-3/internal/errors"
	parser "github.com/manyanin.alexander/task-3/internal/xml_parser"
)

type Currency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func Convert(valCurs *parser.ValCurs) []Currency {
	var currencies []Currency

	for _, valute := range valCurs.Valutes {
		numCode, err := strconv.Atoi(valute.NumCode)
		if err != nil {
			continue
		}

		cleanValue := strings.ReplaceAll(valute.Value, ",", ".")
		value, err := strconv.ParseFloat(cleanValue, 64)
		if err != nil {
			continue
		}

		currencies = append(currencies, Currency{
			NumCode:  numCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	if len(currencies) == 0 {
		panic(errors.ErrNoCurrenciesExtracted)
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	return currencies
}

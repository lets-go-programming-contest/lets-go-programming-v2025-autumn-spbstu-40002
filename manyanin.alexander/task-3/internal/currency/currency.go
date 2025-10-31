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

type ByValueDesc []Currency

func (a ByValueDesc) Len() int           { return len(a) }
func (a ByValueDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByValueDesc) Less(i, j int) bool { return a[i].Value > a[j].Value }

func Convert(valCurs *parser.ValCurs) []Currency {
	var currencies []Currency

	for _, v := range valCurs.Valutes {
		numCode, err := strconv.Atoi(v.NumCode)
		if err != nil {
			continue
		}

		valueStr := strings.Replace(v.Value, ",", ".", -1)
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue
		}

		currencies = append(currencies, Currency{
			NumCode:  numCode,
			CharCode: v.CharCode,
			Value:    value,
		})
	}

	if len(currencies) == 0 {
		panic(errors.ErrNoCurrenciesExtracted)
	}

	sort.Sort(ByValueDesc(currencies))

	return currencies
}

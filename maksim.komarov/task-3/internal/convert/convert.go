package convert

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/megurumacabre/task-3/internal/cbr"
)

var ErrParseValue = errors.New("parse currency value")

type CurrencyOut struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func MapAndSort(doc cbr.Document) ([]CurrencyOut, error) {
	out := make([]CurrencyOut, 0, len(doc.Valutes))

	for _, val := range doc.Valutes {
		num := strings.ReplaceAll(strings.TrimSpace(val.Value), ",", ".")
		parsed, err := strconv.ParseFloat(num, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %q", ErrParseValue, val.Value)
		}

		out = append(out, CurrencyOut{
			NumCode:  val.NumCode,
			CharCode: val.CharCode,
			Value:    parsed,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Value > out[j].Value
	})

	return out, nil
}

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

func Sort(doc cbr.Document) ([]cbr.Valute, error) {
	out := make([]cbr.Valute, 0, len(doc.Valutes))

	for _, val := range doc.Valutes {
		num := strings.ReplaceAll(val.ValueRaw, ",", ".")
		parsed, err := strconv.ParseFloat(num, 64)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", ErrParseValue, err)
		}

		val.Value = parsed
		out = append(out, val)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Value > out[j].Value
	})

	return out, nil
}

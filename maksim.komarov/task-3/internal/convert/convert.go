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

func ToUnifiedSorted(doc cbr.Document) ([]cbr.Currency, error) {
	vals := make([]cbr.Currency, 0, len(doc.Valutes))

	for _, cur := range doc.Valutes {
		num := strings.ReplaceAll(cur.ValueRaw, ",", ".")

		parsed, err := strconv.ParseFloat(num, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrParseValue, err.Error())
		}

		cur.Value = parsed
		vals = append(vals, cur)
	}

	sort.Slice(vals, func(i, j int) bool { return vals[i].Value > vals[j].Value })

	return vals, nil
}

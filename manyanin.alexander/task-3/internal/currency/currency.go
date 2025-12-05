package currency

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/manyanin.alexander/task-3/internal/errors"
)

type Currency struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

func (c Currency) MarshalJSON() ([]byte, error) {
	cleanValue := strings.ReplaceAll(c.Value, ",", ".")

	value, _ := strconv.ParseFloat(cleanValue, 64)

	numCode, _ := strconv.Atoi(c.NumCode)

	data, err := json.Marshal(struct {
		NumCode  int     `json:"num_code"`
		CharCode string  `json:"char_code"`
		Value    float64 `json:"value"`
	}{
		NumCode:  numCode,
		CharCode: c.CharCode,
		Value:    value,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrJSONMarshal, err)
	}

	return data, nil
}

func parseValue(value string) float64 {
	cleanValue := strings.ReplaceAll(value, ",", ".")

	val, _ := strconv.ParseFloat(cleanValue, 64)

	return val
}

func SortByValue(currencies []Currency) []Currency {
	sorted := make([]Currency, len(currencies))

	copy(sorted, currencies)

	sort.Slice(sorted, func(i, j int) bool {
		return parseValue(sorted[i].Value) > parseValue(sorted[j].Value)
	})

	return sorted
}

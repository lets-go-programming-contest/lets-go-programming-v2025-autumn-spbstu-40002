package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func convertValutes(root valCurs) ([]OutputCurrency, error) {
	out := make([]OutputCurrency, 0, len(root.Valutes))
	for _, v := range root.Valutes {
		nc := strings.TrimSpace(v.NumCode)
		cc := strings.TrimSpace(v.CharCode)
		valStr := strings.TrimSpace(v.Value)
		valStr = strings.ReplaceAll(valStr, ",", ".")
		numCode, err := strconv.Atoi(nc)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid num code %q: %v", ErrXMLParse, nc, err)
		}
		value, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid value %q: %v", ErrXMLParse, valStr, err)
		}
		nominalStr := strings.TrimSpace(v.Nominal)
		if nominalStr == "" {
			nominalStr = "1"
		}
		nominal, err := strconv.Atoi(nominalStr)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid nominal %q: %v", ErrXMLParse, nominalStr, err)
		}
		normalized := value / float64(nominal)
		out = append(out, OutputCurrency{
			NumCode:  numCode,
			CharCode: cc,
			Value:    normalized,
		})
	}
	return out, nil
}

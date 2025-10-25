package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func SortValute(vs []Valute) {
	sort.Slice(vs, func(i, j int) bool {
		return vs[i].Value > vs[j].Value
	})
}

func SaveToJSON(path string, valutes []Valute) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("error creating directory for %s: %w", path, err)
	}

	type Currency struct {
		NumCode  int     `json:"num_code"`
		CharCode string  `json:"char_code"`
		Value    float64 `json:"value"`
	}

	output := make([]Currency, 0, len(valutes))
	for _, v := range valutes {
		output = append(output, Currency{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    v.Value,
		})
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("error in json serialization: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("file write error: %s: %w", path, err)
	}

	return nil
}

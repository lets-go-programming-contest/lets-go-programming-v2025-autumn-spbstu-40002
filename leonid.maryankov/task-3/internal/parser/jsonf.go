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

func SaveToJson(path string, valutes []Valute) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("Error creating a directory for " + path + err.Error())
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
		return fmt.Errorf("Error in JSON serialization: " + err.Error())
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("File recording error: " + path + err.Error())
	}
	return nil
}

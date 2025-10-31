package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xanaoxx/task-3/internal/xml"
)

const defaultDirPerm = 0750

type CurrencyRecord struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func SaveCurrencies(data *xml.Currencies, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, defaultDirPerm); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("close error: %v\n", closeErr)
		}
	}()

	records := make([]CurrencyRecord, len(data.Currencies))

	for i, currency := range data.Currencies {
		value, err := currency.ToFloat()
		if err != nil {
			return fmt.Errorf("convert: %w", err)
		}

		records[i] = CurrencyRecord{
			NumCode:  currency.NumCode,
			CharCode: currency.CharCode,
			Value:    value,
		}
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(records); err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}

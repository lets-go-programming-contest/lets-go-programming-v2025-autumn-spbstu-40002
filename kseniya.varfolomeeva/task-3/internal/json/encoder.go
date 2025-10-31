package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xanaoxx/task-3/internal/xml"
)

type CurrencyRecord struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func SaveCurrencies(currencyData *xml.Currencies, outputPath string) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create directory: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	defer file.Close()

	records := make([]CurrencyRecord, len(currencyData.Currencies))
	for i, currency := range currencyData.Currencies {
		value, err := currency.ToFloat()
		if err != nil {
			return err
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
		return fmt.Errorf("JSON encoding failed: %w", err)
	}

	return nil
}


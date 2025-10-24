package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Exam-Play/task-3/internal/xml"
)

const dirPermission = 0o755

type CurrencySimple struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func EncodeJSON(currencies *xml.CurrenciesXML, outputFile string) error {
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, dirPermission); err != nil {
		return fmt.Errorf("unable to create directory: %w", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	simpleCurrencies := make([]CurrencySimple, 0, len(currencies.Currencies))

	for _, currency := range currencies.Currencies {
		value, err := currency.GetFloat()
		if err != nil {
			return fmt.Errorf("invalid format float: %w", err)
		}

		simpleCurrencies = append(simpleCurrencies, CurrencySimple{
			NumCode:  currency.NumCode,
			CharCode: currency.CharCode,
			Value:    value,
		})
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(simpleCurrencies); err != nil {
		return fmt.Errorf("unable to encode json: %w", err)
	}

	return nil
}

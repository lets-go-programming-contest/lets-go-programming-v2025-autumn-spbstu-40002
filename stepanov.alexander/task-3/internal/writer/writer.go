package writer

import (
	"encoding/json"
	"os"

	"stepanov.alexander/task-3/internal/processor"
)

func WriteJSON(filepath string, rates []processor.CurrencyRate) error {
	data, err := json.MarshalIndent(rates, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}

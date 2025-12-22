package writer

import (
	"encoding/json"
	"os"
	"path/filepath"

	"stepanov.alexander/task-3/internal/processor"
)

func WriteJSON(path string, rates []processor.CurrencyRate) error {
	// создаём директорию, если её нет (.output, .output/subdir/one_more и т.п.)
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	data, err := json.MarshalIndent(rates, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

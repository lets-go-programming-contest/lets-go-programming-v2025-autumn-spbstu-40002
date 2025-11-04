package writer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/XShaygaND/task-3/internal/parser"
)

const dirPerm = 0o755

func WriteJSON(outputPath string, data []parser.CurrencyItem) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), dirPerm); err != nil {
		return fmt.Errorf("failed to create output directories: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create json file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}

	return nil
}

package jsonwriter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ControlShiftEscape/task-3/internal/models"
)

const dirPerm = 0o755

func WriteSortedReducedJSON(curs *models.ValCurs, outputPath string) error {
	if curs == nil {
		return nil
	}

	if err := models.SortByValueDesc(curs); err != nil {
		return fmt.Errorf("failed to sort currencies by value: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), dirPerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Errorf("failed to close output file %s: %w", outputPath, err))
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(curs.Valutes)
}

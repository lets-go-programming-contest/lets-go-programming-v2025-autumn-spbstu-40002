package jsonwriter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ControlShiftEscape/task-3/internal/models"
)

type ReducedValute struct {
	NumCode  string `json:"num_code"`
	CharCode string `json:"char_code"`
	Value    string `json:"value"`
}

func WriteSortedReducedJSON(curs *models.ValCurs, outputPath string) error {
	if curs == nil {
		return nil
	}

	if err := models.SortByValueDesc(curs); err != nil {
		return err
	}

	reduced := make([]ReducedValute, len(curs.Valutes))
	for i, v := range curs.Valutes {
		reduced[i] = ReducedValute{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    strings.ReplaceAll(v.Value, ",", "."),
		}
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("Warning: failed to close file %s: %v", outputPath, closeErr)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(reduced); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

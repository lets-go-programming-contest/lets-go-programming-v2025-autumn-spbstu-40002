package jsonwriter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ControlShiftEscape/task-3/internal/models"
)

type ReducedValute struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
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
		num, _ := strconv.Atoi(v.NumCode)
		val, _ := strconv.ParseFloat(strings.ReplaceAll(v.Value, ",", "."), 64)

		reduced[i] = ReducedValute{
			NumCode:  num,
			CharCode: v.CharCode,
			Value:    val,
		}
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", outputPath, err)
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
	return encoder.Encode(reduced)
}

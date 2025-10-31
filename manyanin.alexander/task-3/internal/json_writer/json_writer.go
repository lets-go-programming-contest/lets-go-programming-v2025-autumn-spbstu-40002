package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	converter "github.com/manyanin.alexander/task-3/internal/currency"
	"github.com/manyanin.alexander/task-3/internal/errors"
)

func SaveToJSON(currencies []converter.Currency, outputPath string) {
	if len(currencies) == 0 {
		panic(errors.ErrDataEmpty)
	}

	outputDir := filepath.Dir(outputPath)
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		panic(errors.ErrDirCreate.Error() + ": " + outputDir)
	}

	outputData, err := json.MarshalIndent(currencies, "", "    ")
	if err != nil {
		panic(errors.ErrJSONMarshal.Error() + ": " + err.Error())
	}

	err = os.WriteFile(outputPath, outputData, 0644)
	if err != nil {
		panic(errors.ErrJSONWrite.Error() + ": " + outputPath)
	}
}

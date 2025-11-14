package jsonwrite

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

	file, err := os.Create(outputPath)
	if err != nil {
		panic(errors.ErrJSONWrite.Error() + ": " + err.Error())
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	err = encoder.Encode(currencies)
	if err != nil {
		panic(errors.ErrJSONMarshal.Error() + ": " + err.Error())
	}
}

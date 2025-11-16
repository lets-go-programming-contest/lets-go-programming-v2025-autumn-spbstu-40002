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

	const dirPerm = 0o755

	err := os.MkdirAll(outputDir, dirPerm)
	if err != nil {
		panic(errors.ErrDirCreate.Error() + ": " + outputDir)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		panic(errors.ErrJSONWrite.Error() + ": " + err.Error())
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic(errors.ErrOutputFileCreate.Error() + ": " + closeErr.Error())
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	err = encoder.Encode(currencies)
	if err != nil {
		panic(errors.ErrJSONMarshal.Error() + ": " + err.Error())
	}
}

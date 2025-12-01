package jsonwrite

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/F0LY/task-3/internal/currency"
	"github.com/F0LY/task-3/internal/errors"
)

func WriteCurrenciesToFile(currencies []currency.Currency, filePath string) {
	dir := filepath.Dir(filePath)

	const dirPerm = 0o755

	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		panic(errors.ErrOutputDirCreate.Error() + ": " + err.Error())
	}

	data, err := json.MarshalIndent(currencies, "", "    ")
	if err != nil {
		panic(errors.ErrJSONEncode.Error() + ": " + err.Error())
	}

	file, err := os.Create(filePath)
	if err != nil {
		panic(errors.ErrOutputFileCreate.Error() + ": " + err.Error())
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic(errors.ErrOutputFileCreate.Error() + ": " + closeErr.Error())
		}
	}()

	_, err = file.Write(data)
	if err != nil {
		panic(errors.ErrOutputFileCreate.Error() + ": " + err.Error())
	}
}

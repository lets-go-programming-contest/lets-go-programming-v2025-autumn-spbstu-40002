package jsonwrite

import (
	"encoding/json"
	"os"

	"github.com/F0LY/task-3/internal/currency"
	"github.com/F0LY/task-3/internal/errors"
)

func WriteCurrenciesToFile(currencies []currency.Currency, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(errors.ErrOutputFileCreate.Error() + ": " + err.Error())
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	err = encoder.Encode(currencies)
	if err != nil {
		panic(errors.ErrJSONEncode.Error() + ": " + err.Error())
	}
}

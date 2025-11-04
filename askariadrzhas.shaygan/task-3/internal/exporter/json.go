package exporter

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/XShaygaND/task-3/internal/currency"
	"github.com/XShaygaND/task-3/internal/myerrors"
	"github.com/XShaygaND/task-3/internal/utils"
)

const dirPerm = 0o755

func WriteToJSON(valutes []currency.Valute, path string) {
	jsonValutes := utils.ParseValutesToJSON(valutes)

	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		panic(myerrors.ErrDirCreate)
	}

	file, err := os.Create(path)
	if err != nil {
		panic(myerrors.ErrOutOpen)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	if err := enc.Encode(jsonValutes); err != nil {
		panic(myerrors.ErrOutEncode)
	}
}

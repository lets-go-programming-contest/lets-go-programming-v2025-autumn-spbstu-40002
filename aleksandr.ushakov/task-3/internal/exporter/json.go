package exporter

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rachguta/task-3/internal/currency"
	"github.com/rachguta/task-3/internal/myerrors"
	"github.com/rachguta/task-3/internal/utils"
)

const dirPerm = 0o755

func WriteToJSON(valutes []currency.Valute, path string) {
	jsonValutes := utils.ParseValutesToJSON(valutes)

	err := os.MkdirAll(filepath.Dir(path), dirPerm)
	if err != nil {
		panic(myerrors.ErrDirCreate)
	}

	file, err := os.Create(path)
	if err != nil {
		panic(myerrors.ErrOutOpen)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			panic(myerrors.ErrCloseFile)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(jsonValutes)
	if err != nil {
		panic(myerrors.ErrOutEncode)
	}
}

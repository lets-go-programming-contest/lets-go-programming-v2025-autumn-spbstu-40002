package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const dirPerm = 0o755

func SaveJSON(path string, items []OutputItem) {
	dir := filepath.Dir(path)
	if dir != "." {
		err := os.MkdirAll(dir, dirPerm)
		if err != nil {
			panic(err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	defer func() {
		cerr := file.Close()
		if cerr != nil {
			panic(cerr)
		}
	}()

	data, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		panic(err)
	}

	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}
}

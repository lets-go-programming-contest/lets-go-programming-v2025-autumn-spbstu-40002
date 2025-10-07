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
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			panic(err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			panic(cerr)
		}
	}()

	data, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		panic(err)
	}

	if _, err := file.Write(data); err != nil {
		panic(err)
	}
}

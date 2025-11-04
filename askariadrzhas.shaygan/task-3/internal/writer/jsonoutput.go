package writer

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const defaultDirPerm = 0o755

func SaveAsJSON(data interface{}, outputPath string) {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, defaultDirPerm); err != nil {
		panic("cannot create output directory: " + err.Error())
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("cannot create output file: " + err.Error())
	}

	defer func() {
		if closeErr := outputFile.Close(); closeErr != nil {
			panic("cannot close file: " + closeErr.Error())
		}
	}()

	jsonEncoder := json.NewEncoder(outputFile)
	jsonEncoder.SetIndent("", "    ")

	if err := jsonEncoder.Encode(data); err != nil {
		panic("cannot write JSON data: " + err.Error())
	}
}

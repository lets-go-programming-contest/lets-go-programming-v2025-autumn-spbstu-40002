package writer

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func SaveAsJSON(data interface{}, outputPath string) {
	dir := filepath.Dir(outputPath)
	err := os.MkdirAll(dir, 0o755)

	if err != nil {
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

	err = jsonEncoder.Encode(data)
	if err != nil {
		panic("cannot write JSON data: " + err.Error())
	}
}

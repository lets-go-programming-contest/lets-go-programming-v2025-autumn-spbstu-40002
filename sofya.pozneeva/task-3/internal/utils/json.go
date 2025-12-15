package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	dirPerm  = 0o755
	filePerm = 0o600
)

func GetValutesForJSON(valutes *Valutes) ([]byte, error) {
	outputData := make([]Valute, 0)

	for _, valute := range *valutes {
		output := Valute{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    valute.Value,
		}
		outputData = append(outputData, output)
	}

	jsonData, err := json.MarshalIndent(outputData, "", "    ")
	if err != nil {
		panic(err)
	}

	return jsonData, nil
}

func JSONWrite(config *Config, jsonData []byte) {
	outputDir := filepath.Dir(config.OutputFile)

	err := os.MkdirAll(outputDir, dirPerm)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(config.OutputFile, jsonData, filePerm)
	if err != nil {
		panic(err)
	}
}

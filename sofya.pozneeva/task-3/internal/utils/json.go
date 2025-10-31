package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type CurrencyOutput struct {
	NumCode  string  `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func GetValutesForJSON(valutes *Valutes) ([]byte, error) {
	outputData := make([]CurrencyOutput, 0)

	for _, valute := range *valutes {
		value, err := valute.ConvertValue()
		if err != nil {
			return nil, err
		}

		output := CurrencyOutput{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    value,
		}
		outputData = append(outputData, output)
	}

	jsonData, err := json.MarshalIndent(outputData, "", "    ")
	if err != nil {
		panic(err)
	}

	return jsonData, nil
}

const (
	dirPerm = 0755
	filePerm = 0600
)

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

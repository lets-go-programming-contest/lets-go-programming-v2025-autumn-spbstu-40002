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

func GetValutesForJson(valutes *Valutes) ([]byte, error) {
	var outputData []CurrencyOutput

	for _, valute := range *valutes {
		value, err := valute.ConvertValue()
		if err != nil {
			return nil, err
		}
		outputData = append(outputData, CurrencyOutput{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	jsonData, err := json.MarshalIndent(outputData, "", "    ")
	if err != nil {
		panic(err)
	}
	return jsonData, nil
}

func JsonWrite(config *Config, jsonData []byte) {
	outputDir := filepath.Dir(config.OutputFile)

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(config.OutputFile, jsonData, 0644)
	if err != nil {
		panic(err)
	}
}

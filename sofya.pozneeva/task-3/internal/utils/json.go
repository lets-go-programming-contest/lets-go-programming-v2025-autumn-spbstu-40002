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

	var outputData = make([]CurrencyOutput, 0)

	for _, valute := range *valutes {
		value, err := valute.ConvertValue()
		if err != nil {
			return nil, err
		}
		var output CurrencyOutput = CurrencyOutput{
			NumCode: valute.NumCode,
			CharCode:v alute.CharCode,
			Value: value}
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

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(config.OutputFile, jsonData, 0600)
	if err != nil {
		panic(err)
	}
}

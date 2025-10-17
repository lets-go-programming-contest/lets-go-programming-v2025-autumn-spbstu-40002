package json

import (
	"encoding/json"
	"os"
	"path/filepath"

	merr "github.com/slendycs/go-lab-3/internal/myerrors"
	"github.com/slendycs/go-lab-3/internal/xml"
)

type JsonData struct {
	NumCode  string `json:"num_code"`
	CharCode string `json:"char_code"`
	Value    string `json:"value"`
}

func MakeJsonFromData(path string, data *xml.ValCurs) {
	// Creating output directory
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(merr.ErrFailedToCreateDir)
	}

	// Try to open output file.
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(merr.ErrFailedToOpenOutputFile)
	}
	defer file.Close()

	// Create a slice of valute data.
	output := make([]JsonData, 0)
	for _, valute := range data.Valute {
		// Create a JSON record of data.
		outputData := JsonData{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    valute.Value,
		}

		output = append(output, outputData)
	}

	// Serialize data
	rawData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		panic(merr.ErrFailedToSerializeJSON)
	}

	// Write data
	_, err = file.Write(rawData)
	if err != nil {
		panic(merr.ErrFailedToWriteData)
	}
}

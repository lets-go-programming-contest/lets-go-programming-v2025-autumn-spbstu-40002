package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	merr "github.com/slendycs/go-lab-3/internal/myerrors"
	"github.com/slendycs/go-lab-3/internal/xml"
)

type JsonData struct {
	NumCode  int `json:"num_code"`
	CharCode string `json:"char_code"`
	Value    float64 `json:"value"`
}

func MakeJsonData(data *xml.ValCurs) []JsonData {
	// Create a slice of valute data.
	output := make([]JsonData, 0)
	for _, valute := range data.Valute {
		// Create a JSON record of data.
		intNumCode, err := strconv.Atoi(valute.NumCode)
		if err != nil {
			fmt.Println(merr.ErrNumCodeIsNotIneger)
			return nil
		}
		formatedValue := strings.ReplaceAll(valute.Value, ",", ".")
		floatValue, err := strconv.ParseFloat(formatedValue, 64)
		if err != nil {
			fmt.Println(merr.ErrValueIsNotFloat)
			return nil
		}

		outputData := JsonData{
			NumCode:  intNumCode,
			CharCode: valute.CharCode,
			Value:    floatValue,
		}

		output = append(output, outputData)
	}
	return output
}


func WriteJSONData(path string, data *xml.ValCurs) {
	// Serialize data
	rawData, err := json.MarshalIndent(MakeJsonData(data), "", "  ")
	if err != nil {
		panic(merr.ErrFailedToSerializeJSON)
	}

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

	// Write data
	_, err = file.Write(rawData)
	if err != nil {
		panic(merr.ErrFailedToWriteData)
	}
}

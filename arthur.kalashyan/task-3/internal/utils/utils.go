package utils

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/net/html/charset"
)

const (
	dirPerm  = 0o755
	filePerm = 0o600
)

type Exchange struct {
	Currencies []Currency `json:"currencies" xml:"Valute"`
}

type Currency struct {
	NumCode  string `json:"num_code"  xml:"NumCode"`
	CharCode string `json:"char_code" xml:"CharCode"`
	Value    string `json:"value"     xml:"Value"`
}

func ReadXML(path string) (*Exchange, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open xml file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic(closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var exch Exchange
	if err := decoder.Decode(&exch); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	return &exch, nil
}

func SaveToJSON(data any, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	bytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	err = os.WriteFile(path, bytes, filePerm)
	if err != nil {
		return fmt.Errorf("write json: %w", err)
	}

	return nil
}

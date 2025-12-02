package bank

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/charmap"
)

type Bank struct {
	Currencies []Currency `json:"currencies" xml:"Valute"`
}

type Currency struct {
	NumCode  int    `json:"num_code"  xml:"NumCode"`
	CharCode string `json:"char_code" xml:"CharCode"`
	Value    string `json:"value"     xml:"Value"`
}

func charsetReader(charset string, reader io.Reader) (io.Reader, error) {
	switch charset {
	case "windows-1251":
		return charmap.Windows1251.NewDecoder().Reader(reader), nil
	default:
		return reader, nil
	}
}

func DecodeXML(reader io.Reader) (*Bank, error) {
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charsetReader

	bankData := new(Bank)
	if err := decoder.Decode(bankData); err != nil {
		return nil, fmt.Errorf("failed to decode XML to Bank struct: %w", err)
	}

	return bankData, nil
}

func LoadFromXML(path string) (*Bank, error) {
	fileHandle, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open input XML file: %w", err)
	}

	bankData, decodeErr := DecodeXML(fileHandle)

	if closeErr := fileHandle.Close(); closeErr != nil && decodeErr == nil {
		return nil, fmt.Errorf("failed to close input XML file: %w", closeErr)
	}

	return bankData, decodeErr
}

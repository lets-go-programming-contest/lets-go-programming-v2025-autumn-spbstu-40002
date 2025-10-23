package parsers

import (
	"encoding/xml"
	"os"
	"strconv"
	"strings"

	merr "github.com/slendycs/go-lab-3/internal/myerrors"
	"golang.org/x/net/html/charset"
)

func (value *CommaFloat64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	// Checking if a string is empty.
	if s == "" {
		*value = CommaFloat64(0.0)
		return nil
	}

	// Replace the comma with a dot and decode the value as float64.
	s = strings.Replace(s, ",", ".", 1)
	dotFloat64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*value = CommaFloat64(dotFloat64)
	return nil
}

func ReadXML(path string, data *ValStruct) {
	// Opening XML file with data.
	file, err := os.Open(path)
	if err != nil {
		panic(merr.ErrFailedToOpenXML)
	}
	defer file.Close()

	// Create a new decoder for XML file.
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	// Decode XML file.
	err = decoder.Decode(data)
	if err != nil {
		panic(merr.ErrFailedToDecodeXML)
	}
}

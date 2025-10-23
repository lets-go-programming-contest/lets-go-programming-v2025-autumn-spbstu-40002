package xmlfile

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type ValCurs struct {
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int     `xml:"num_code" json:"num_code" yaml:"num_code"`
	CharCode string  `xml:"char_code" json:"char_code" yaml:"char_code"`
	Value    float64 `xml:"value" json:"value" yaml:"value"`
}

func (v *Valute) UnmarshalXML(decod *xml.Decoder, start xml.StartElement) error {
	var aux struct {
		NumCode  int    `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	if err := decod.DecodeElement(&aux, &start); err != nil {
		return fmt.Errorf("couldn't decode elem: %w", err)
	}

	valStr := strings.ReplaceAll(strings.TrimSpace(aux.Value), ",", ".")

	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return fmt.Errorf("couldn't parse Value: %w", err)
	}

	v.NumCode = aux.NumCode
	v.CharCode = aux.CharCode
	v.Value = val

	return nil
}

func GetValCursStruct(inputPath string) (ValCurs, error) {
	var doc ValCurs

	file, err := os.Open(inputPath)
	if err != nil {
		return doc, fmt.Errorf("couldn't open XML file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("error closing file: %v", err)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&doc); err != nil {
		return doc, fmt.Errorf("xml parsing error: %w", err)
	}

	return doc, nil
}

func SortValCursByValue(doc *ValCurs) {
	sort.Slice(doc.Valutes, func(i, j int) bool {
		return doc.Valutes[i].Value > doc.Valutes[j].Value
	})
}

package currency

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"sort"
)

// Currency represents a currency entry
type Currency struct {
	NumCode  int     `xml:"NumCode" json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Value    float64 `xml:"Value" json:"value"`
}

// ValCurs represents the XML structure from CBR
type ValCurs struct {
	XMLName   xml.Name   `xml:"ValCurs"`
	Date      string     `xml:"Date,attr"`
	Name      string     `xml:"Name,attr"`
	Currencies []Currency `xml:"Valute"`
}

// Processor handles currency data processing
type Processor struct{}

// NewProcessor creates a new currency processor
func NewProcessor() *Processor {
	return &Processor{}
}

// ParseXML parses XML data from file
func (p *Processor) ParseXML(filename string) (*ValCurs, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic("Failed to read input file: " + err.Error())
	}

	var valCurs ValCurs
	err = xml.Unmarshal(data, &valCurs)
	if err != nil {
		panic("Failed to parse XML: " + err.Error())
	}

	return &valCurs, nil
}

// SortByValue sorts currencies by value in descending order
func (p *Processor) SortByValue(currencies []Currency) []Currency {
	sorted := make([]Currency, len(currencies))
	copy(sorted, currencies)
	
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	
	return sorted
}

// SaveToJSON saves currencies to JSON file
func (p *Processor) SaveToJSON(currencies []Currency, filename string) error {
	// Create directory if it doesn't exist
	if len(filename) > 0 {
		lastSlash := -1
		for i := len(filename) - 1; i >= 0; i-- {
			if filename[i] == '/' || filename[i] == '\\' {
				lastSlash = i
				break
			}
		}
		if lastSlash >= 0 {
			dir := filename[:lastSlash]
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				panic("Failed to create directory: " + err.Error())
			}
		}
	}

	// Create output file
	file, err := os.Create(filename)
	if err != nil {
		panic("Failed to create output file: " + err.Error())
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(currencies)
	if err != nil {
		panic("Failed to encode JSON: " + err.Error())
	}

	return nil
}

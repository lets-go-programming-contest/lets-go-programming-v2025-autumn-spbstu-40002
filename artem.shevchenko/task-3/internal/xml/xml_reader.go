package xml

import (
	"encoding/xml"
	"os"
	"sort"

	merr "github.com/slendycs/go-lab-3/internal/myerrors"
	"golang.org/x/net/html/charset"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Text    string   `xml:",chardata"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valute  []struct {
		Text      string `xml:",chardata"`
		ID        string `xml:"ID,attr"`
		NumCode   string `xml:"NumCode"`
		CharCode  string `xml:"CharCode"`
		Nominal   string `xml:"Nominal"`
		Name      string `xml:"Name"`
		Value     string `xml:"Value"`
		VunitRate string `xml:"VunitRate"`
	} `xml:"Valute"`
}

func ReadXML(path string, data *ValCurs) {
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

func SortVal(data *ValCurs) {
	sort.Slice(data.Valute, func(i, j int) bool {
		return data.Valute[i].Value > data.Valute[j].Value
	})
}

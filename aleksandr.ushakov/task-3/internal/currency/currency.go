package currency

import "encoding/xml"

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valute  []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

type JSONValute struct {
	NumCode  string  `json:"num-code"`
	CharCode string  `json:"char-code"`
	Value    float64 `json:"value"`
}

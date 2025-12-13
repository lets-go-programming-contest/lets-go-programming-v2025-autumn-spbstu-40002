package currency

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/manyanin.alexander/task-3/internal/errors"
)

type Currency struct {
	XMLName  xml.Name `json:"-"         xml:"Valute"`
	NumCode  int      `json:"num_code"  xml:"NumCode"`
	CharCode string   `json:"char_code" xml:"CharCode"`
	Value    float64  `json:"value"     xml:"Value"`
}

func (c *Currency) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var raw struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	err := decoder.DecodeElement(&raw, &start)
	if err != nil {
		return fmt.Errorf("%w: %w", errors.ErrXMLDecode, err)
	}

	num, err := strconv.Atoi(raw.NumCode)
	if err != nil {
		c.NumCode = 0
	} else {
		c.NumCode = num
	}

	cleanValue := strings.ReplaceAll(raw.Value, ",", ".")
	val, err := strconv.ParseFloat(cleanValue, 64)

	if err != nil {
		c.Value = 0
	} else {
		c.Value = val
	}

	c.CharCode = raw.CharCode

	return nil
}

func (c Currency) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(struct {
		NumCode  int     `json:"num_code"`
		CharCode string  `json:"char_code"`
		Value    float64 `json:"value"`
	}{
		NumCode:  c.NumCode,
		CharCode: c.CharCode,
		Value:    c.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrJSONMarshal, err)
	}

	return data, nil
}

func SortByValue(currencies []Currency) []Currency {
	sorted := make([]Currency, len(currencies))

	copy(sorted, currencies)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	return sorted
}

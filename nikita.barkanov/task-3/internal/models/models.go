package models

import (
	"encoding/json"
	"encoding/xml"
	"sort"
	"strconv"
	"strings"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute" json:"-"`
}

type Valute struct {
	NumCode  string `xml:"NumCode"  json:"num_code"`
	CharCode string `xml:"CharCode" json:"char_code"`
	Value    string `xml:"Value"    json:"-"`
}

func (v Valute) ValueFloat() (float64, error) {
	cleaned := strings.ReplaceAll(v.Value, ",", ".")
	return strconv.ParseFloat(cleaned, 64)
}

func (v Valute) MarshalJSON() ([]byte, error) {
	value, err := v.ValueFloat()
	if err != nil {
		return nil, err
	}

	type Alias struct {
		NumCode  int     `json:"num_code"`
		CharCode string  `json:"char_code"`
		Value    float64 `json:"value"`
	}

	num, _ := strconv.Atoi(v.NumCode)

	return jsonMarshal(Alias{
		NumCode:  num,
		CharCode: v.CharCode,
		Value:    value,
	})
}

var jsonMarshal = func(v any) ([]byte, error) {
	return json.Marshal(v)
}

func SortByValueDesc(curs *ValCurs) error {
	if curs == nil || len(curs.Valutes) == 0 {
		return nil
	}

	sort.Slice(curs.Valutes, func(i, j int) bool {
		vi, errI := curs.Valutes[i].ValueFloat()
		vj, errJ := curs.Valutes[j].ValueFloat()

		if errI != nil {
			return false
		}
		if errJ != nil {
			return true
		}
		return vi > vj
	})

	return nil
}

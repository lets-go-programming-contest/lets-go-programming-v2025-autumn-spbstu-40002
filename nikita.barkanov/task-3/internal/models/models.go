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
	Valutes []Valute `json:"-"      xml:"Valute"`
}

type Valute struct {
	NumCode  string `json:"num_code"  xml:"NumCode"`
	CharCode string `json:"char_code" xml:"CharCode"`
	Value    string `json:"-"         xml:"Value"`
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

	return json.Marshal(Alias{
		NumCode:  num,
		CharCode: v.CharCode,
		Value:    value,
	})
}

func SortByValueDesc(curs *ValCurs) error {
	if curs == nil || len(curs.Valutes) == 0 {
		return nil
	}

	sort.Slice(curs.Valutes, func(i, j int) bool {
		valuteI, errI := curs.Valutes[i].ValueFloat()
		valuteJ, errJ := curs.Valutes[j].ValueFloat()

		if errI != nil {
			return false
		}

		if errJ != nil {
			return true
		}

		return valuteI > valuteJ
	})

	return nil
}

package models

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
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
		return nil, fmt.Errorf("convert value to float: %w", err)
	}

	type Alias struct {
		NumCode  int     `json:"num_code"`
		CharCode string  `json:"char_code"`
		Value    float64 `json:"value"`
	}

	num, err := strconv.Atoi(v.NumCode)
	if err != nil {
		return nil, fmt.Errorf("convert NumCode to int: %w", err)
	}

	data, err := json.Marshal(Alias{
		NumCode:  num,
		CharCode: v.CharCode,
		Value:    value,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal Alias to JSON: %w", err)
	}

	return data, nil
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

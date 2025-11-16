package models

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID        string `xml:"ID,attr"`
	NumCode   string `xml:"NumCode"`
	CharCode  string `xml:"CharCode"`
	Nominal   int    `xml:"Nominal"`
	Name      string `xml:"Name"`
	Value     string `xml:"Value"`
	VunitRate string `xml:"VunitRate"`
}

func (v *Valute) ValueFloat() (float64, error) {
	cleaned := strings.ReplaceAll(v.Value, ",", ".")

	value, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse float from %q: %w", cleaned, err)
	}

	return value, nil
}

func SortByValueDesc(curs *ValCurs) error {
	return sortValutes(curs, true)
}

func sortValutes(curs *ValCurs, desc bool) error {
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

		if desc {
			return valuteI > valuteJ
		}

		return valuteI < valuteJ
	})

	return nil
}

package structures

import (
	"fmt"
	"strconv"
	"strings"
)

type XMLStruct struct {
	ValCurs []struct {
		NumCode  int    `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	} `xml:"Valute"`
}

type JSONSValute struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

type JSONStruct struct {
	ValCurs []JSONSValute `json:"valute"`
}

func ConvertXMLToJSON(xmlStruct XMLStruct) (JSONStruct, error) {
	jsonStruct := new(JSONStruct)

	for _, valute := range xmlStruct.ValCurs {
		floatValue, err := strconv.ParseFloat(strings.Replace(valute.Value, ",", ".", 1), 64)
		if err != nil {
			return JSONStruct{}, fmt.Errorf("Invalid value: %w", err)
		}

		jsonStruct.ValCurs = append(jsonStruct.ValCurs, JSONSValute{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    floatValue,
		})
	}

	return *jsonStruct, nil
}

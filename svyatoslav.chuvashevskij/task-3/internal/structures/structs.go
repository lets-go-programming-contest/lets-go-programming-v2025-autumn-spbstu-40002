package structures

type ValuteStruct struct {
	ValCurs []struct {
		NumCode  int     `xml:"NumCode" json:"num_code"`
		CharCode string  `xml:"CharCode" json:"char_code"`
		Value    float64 `xml:"Value" json:"value"`
	} `xml:"Valute" json:"valute"`
}

/*type ValuteStruct struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}*/

/*type JSONStruct struct {
	ValCurs []ValuteStruct `json:"valute"`
}*/

/*func ConvertXMLToJSON(xmlStruct ValuteStruct) (JSONStruct, error) {
	jsonStruct := new(JSONStruct)

	for _, valute := range xmlStruct.ValCurs {
		floatValue, err := strconv.ParseFloat(strings.Replace(valute.Value, ",", ".", 1), 64)
		if err != nil {
			return JSONStruct{}, fmt.Errorf("invalid value: %w", err)
		}

		jsonStruct.ValCurs = append(jsonStruct.ValCurs, ValuteStruct{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    floatValue,
		})
	}

	return *jsonStruct, nil
}*/

package structures

type ValuteStruct struct {
	ValCurs []struct {
		NumCode  int     `json:"num_code"  xml:"NumCode"`
		CharCode string  `json:"char_code" xml:"CharCode"`
		Value    float64 `json:"value"     xml:"Value"`
	} `json:"valute" xml:"Valute"`
}

package structures

type ValuteStruct struct {
	ValCurs []struct {
		NumCode  int     `xml:"NumCode" json:"num_code"`
		CharCode string  `xml:"CharCode" json:"char_code"`
		Value    float64 `xml:"Value" json:"value"`
	} `xml:"Valute" json:"valute"`
}

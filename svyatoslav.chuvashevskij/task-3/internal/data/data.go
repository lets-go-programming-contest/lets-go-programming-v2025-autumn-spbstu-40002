package data

type DataStruct struct {
	ValCurs []struct {
		NumCode  int    `xml:"NumCode" json:"num_code"`
		CharCode string `xml:"CharCode" json:"char_code"`
		Value    string `xml:"Value" json:"value"`
	} `xml:"Valute" json:""`
}

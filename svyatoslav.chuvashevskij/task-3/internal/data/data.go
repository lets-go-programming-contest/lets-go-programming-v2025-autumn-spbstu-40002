package data

type DataStruct struct {
	ValCurs []struct {
		NumCode  int    `json:"num_code" xml:"NumCode"`
		CharCode string `json:"char_code" xml:"CharCode"`
		Value    string `json:"value" xml:"Value"`
	} `json:"valute" xml:"Valute"`
}

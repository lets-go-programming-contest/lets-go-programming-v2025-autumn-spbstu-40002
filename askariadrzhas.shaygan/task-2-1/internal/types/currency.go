package types

type ProcessedCurrency struct {
	Code   string  `json:"num_code"`
	Symbol string  `json:"char_code"`
	Rate   float64 `json:"value"`
}

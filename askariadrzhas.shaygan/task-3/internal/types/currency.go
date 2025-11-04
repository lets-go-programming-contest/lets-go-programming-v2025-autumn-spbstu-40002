package types

type ProcessedCurrency struct {
	Code   string  `json:"numCode"`
	Symbol string  `json:"charCode"`
	Rate   float64 `json:"value"`
}

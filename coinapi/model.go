package coinapi

type ExchangeRate struct {
	Rate  float64 `json:"rate"`
	Error string  `json:"error"`
}

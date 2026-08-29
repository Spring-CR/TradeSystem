package options

type PricingResult struct {
	Ytm  float64 `json:"ytm"`
	DirtyPrice float64 `json:"dirtyPrice"`
	CleanPrice float64 `json:"cleanPrice"`
}

type PricingForm struct {
	Symbol string `form:"symbol"`
	Ytm float64 `form:"ytm"`
	TradeDate string `form:"tradeDate"`
}
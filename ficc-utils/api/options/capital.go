package options

// 资金账户信息
type CapitalResult struct {
	Account         int     `json:"account"`
	Currency        string  `json:"currency"`
	TotalBalance   float64 `json:"totalBalance"`
	AvailableAmount float64 `json:"availableAmount"`
}

type InternalCapitalResult struct {
	CapitalResult
	KeyCapAcctId int    `json:"keyCapAcctId"`
	CapAcctCode  string `json:"capAcctCode"`
}

// TITANS资金查询response
type CapitalServiceResponse struct {
	ServiceID string `json:"serviceId"`
	ErrCode   struct {
		Code           int    `json:"code"`
		CHS            string `json:"chs"`
		ENG            string `json:"eng"`
		ExceptionClass string `json:"exceptionClass"`
	} `json:"errCode"`
	Data []struct {
		KeyCtptyId   int     `json:"keyCtptyId"`
		Purpose      string  `json:"purpose"`
		Currency     string  `json:"currency"`
		TotalBalance      float64 `json:"totalBalance"`
		AvailableAmt float64 `json:"availableAmt"`
		FreezeAmt    float64 `json:"freezeAmt"`
		KeyCapAcctId int     `json:"keyCapAcctId"`
		CapAcctCode  string  `json:"capAcctCode"`
	} `json:"data"`
	Timestamp int64 `json:"timestamp"`
}

// 保证金占用
type MarginResponse struct {
	Total 	int 				`json:"total"`
	Data 	[]*MarginOccupancy 	`json:"data"`
	Code    int          		`json:"Code"`
	Message string          	`json:"Message"`
}

type MarginOccupancy struct {
	Account                 int     `json:"account"`
	MaxMarginOccupancy      float64 `json:"maxMarginOccupancy"`
	MaxLongMarginOccupancy  float64 `json:"maxLongMarginOccupancy"`
	MaxShortMarginOccupancy float64 `json:"maxShortMarginOccupancy"`
}

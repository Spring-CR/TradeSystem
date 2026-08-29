package olts_fut

type NewOrderSingle struct {
	MsgType            string                 `json:"msgType"`
	ClOrdID            string                 `json:"secondaryClOrdID"`
	Account            string                 `json:"clientID"` // 对冲账号
	HandlInst          string                 `json:"handlInst"`
	Symbol             string                 `json:"symbol"`
	SecurityID         string                 `json:"securityID"`
	SecurityExchange   string                 `json:"marketCode"`
	Side               string                 `json:"side"`
	TransactTime       string                 `json:"transactTime"`
	OrderQty           int                    `json:"orderQty"`
	OrdType            string                 `json:"ordType"`
	Price              float64                `json:"price"`
	Currency           string                 `json:"-"`
	CashOrderQty       int                    `json:"cashOrderQty"`
	TimeInForce        string                 `json:"timeInForce"`
	AlgoName           string                 `json:"targetStrategy"`
	BusinessType       string                 `json:"businessType"`
	OpenClose          string                 `json:"openClose"`
	Operator           string                 `json:"operator"`
	CounterParty       []string               `json:"counterParty"`
	ContractID         *string                `json:"contractID"`
	StrategyParameters map[string]interface{} `json:"strategyParameters"`
	TradePurpose       string                 `json:"tradePurpose"`
}

type ExecutionReport struct {
	MsgType          string  `json:"msgType"`
	Account          string  `json:"clientID"`
	HandlInst        string  `json:"handlInst"`
	OrdQty           float64 `json:"orderQty"`
	Price            float64 `json:"price"`
	Side             string  `json:"side"`
	Symbol           string  `json:"symbol"`
	SecurityType     string  `json:"-"`
	SecurityExchange string  `json:"marketCode"`
	OrdType          string  `json:"ordType"`
	Currency         string  `json:"-"`
	TimeInForce      string  `json:"timeInForce"`
	ClOrdID          string  `json:"secondaryClOrdID"`
	OrdID            string  `json:"clOrdID"`
	OrigClOrdID      string  `json:"-"`
	ExecID           string  `json:"msgid"`
	ExecType         string  `json:"execType"`
	ExecTransType    string  `json:"execTransType"`
	ExecRefID        string  `json:"-"`
	OrdStatus        string  `json:"ordStatus"`
	LastQty          float64 `json:"lastShares"`
	LastPx           float64 `json:"lastPx"`
	AvgPx            float64 `json:"avgPx"`
	LeavesQty        float64 `json:"leavesQty"`
	CumQty           float64 `json:"cumAmt"`
	TransactTime     string  `json:"transactTime"`
	OrdRejReason     string  `json:"-"`
	Text             string  `json:"text"`
}

type OrderCancelRequest struct {
	MsgType      string `json:"msgType"`
	Account      string `json:"-"`
	Symbol       string `json:"symbol"`
	Side         string `json:"side"`
	ClOrdID      string `json:"-"`
	OrigClOrdID  string `json:"secondaryClOrdID"`
	TransactTime string `json:"transactTime"`
	AlgoName     string `json:"targetStrategy"`
}

type OrderCancelReject struct {
	MsgType      string `json:"msgType"`
	OrdID        string `json:"clOrdID"`
	ClOrdID      string `json:"-"`
	OrigClOrdID  string `json:"secondaryClOrdID"`
	OrdStatus    string `json:"-"`
	TransactTime string `json:"transactTime"`
	CxlRejReason string `json:"text"`
}

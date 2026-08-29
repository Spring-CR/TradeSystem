package ficc

type NewOrderSingle struct {
	MsgType          string
	ClOrdID          string
	SecurityExchange string
	TraderID         string
	Symbol           string
	Side             string
	OrdType          string
	Quantity         float64
	Price            float64
	YTM              float64
	InvestorID       string
	SettlType        string
	TransactTime     string
	CounterpartyID   string
	OrdExecUser      string
	HandlInst        string
	AlgName          string
}

type ExecutionReport struct {
	MsgType          string
	MsgID            string
	SecurityExchange string
	ClOrdID          string
	OrigClOrdID      string `json:"-"`
	OrdID            string
	Symbol           string
	ExecType         string
	ExecRefID        string
	ExecTransType    string
	OrdStatus        string
	LastQty          float64
	LastPx           float64
	Quantity         float64
	Price            float64
	Currency         string
	LeavesQty        float64
	OrdType          string
	ExecID           string
	TransactTime     string
	Side             string
	OrdRejReason     string
	Text             string
	HandlInst        string
	Broker           string
}

type OrderCancelRequest struct {
	MsgType      string
	ClOrdID      string
	OrigClOrdID  string
	TransactTime string
}

type OrderCancelReject struct {
	MsgType      string
	MsgID        string
	OrdID        string
	ClOrdID      string
	OrigClOrdID  string
	OrdStatus    string
	TransactTime string
	CxlRejReason string
}

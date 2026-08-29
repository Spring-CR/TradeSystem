package stars_fut

type NewOrderSingle struct {
	MsgType          string
	ClOrdID          string
	Account          string // 对冲账号
	HandlInst        string
	Symbol           string
	SecurityExchange string
	Side             string
	TransactTime     string
	OrderQty         float64
	OrdType          string
	Price            float64
	Currency         string
}

type ExecutionReport struct {
	MsgType          string
	Account          string
	HandlInst        string
	OrdQty           float64
	Price            float64
	Side             string
	Symbol           string
	SecurityType     string
	SecurityExchange string
	OrdType          string
	Currency         string
	TimeInForce      string
	ClOrdID          string
	OrdID            string
	OrigClOrdID      string
	ExecID           string
	ExecType         string
	ExecTransType    string
	ExecRefID        string
	OrdStatus        string
	LastQty          float64
	LastPx           float64
	AvgPx            float64
	LeavesQty        float64
	CumQty           float64
	TransactTime     string
	OrdRejReason     string
	Text             string
}

type OrderCancelRequest struct {
	MsgType      string
	Account      string
	Symbol       string
	Side         string
	ClOrdID      string
	OrigClOrdID  string
	TransactTime string
}

type OrderCancelReject struct {
	MsgType      string
	OrdID        string
	ClOrdID      string
	OrigClOrdID  string
	OrdStatus    string
	TransactTime string
	CxlRejReason string
}

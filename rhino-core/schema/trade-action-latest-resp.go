package schema

import "strconv"

type TradeActionLatestResp struct {
	ID                     int64
	ActionUser             string `sql:"size: 64"`
	ActionMsgTime          int64
	ActionTime             int64
	ActionType             string `sql:"size: 2"`
	ActionKey              string `sql:"unique: uq_talr_actkey, size: 188"`
	OrdDraftBeforeUpdate   string `sql:"type: MEDIUMTEXT"`
	AppOrdID               string `sql:"size: 188"`
	OrderID                string `sql:"size: 188"`
	RootClOrdID            string `sql:"size: 188"`
	ClOrdID                string `sql:"size: 188"`
	OrigClOrdID            string `sql:"size: 188"`
	ExecID                 string `sql:"size: 188"`
	ExecType               string `sql:"size: 2"`
	OrdStatus              string `sql:"size: 2"`
	OrdRejReason           string
	CxlRejResponseTo       string
	ExecRestatementReason  string
	Account                string `sql:"size: 64"`
	SecurityExchange       string `sql:"size: 8"`
	SecurityExchangeRegion string `sql:"size: 4"`
	Symbol                 string `sql:"size: 64"`
	SymbolSfx              string `sql:"size: 8"`
	SecurityID             string `sql:"size: 64"`
	IDSource               string `sql:"size: 2"`
	SecurityType           string `sql:"size: 2"`
	Side                   string `sql:"size: 2"`
	OpenClose              string `sql:"size: 2"`
	OrderQty               float64
	CashOrderQty           float64
	OrdType                string `sql:"size: 2"`
	Price                  float64
	Currency               string `sql:"size: 4"`
	EffectiveTime          string `sql:"size: 64"`
	ExpireTime             string `sql:"size: 64"`
	LastShares             int64
	LastPx                 float64
	LeavesQty              int64
	CumQty                 int64
	AvgPx                  float64
	TransactTime           int64
	MsgTime                int64
	MsgSeq                 int64
	StreamInputMsgSeq      int64  // 当使用stream api传入订单，这个字段设置传入参数的消息序号
	ChannelCode            string `sql:"size: 32"`
}

func (o *TradeActionLatestResp) GetCacheKey() string {
	return o.AppOrdID + o.ActionUser + strconv.Itoa(int(o.ActionMsgTime)) + o.ActionType + o.ActionKey
}

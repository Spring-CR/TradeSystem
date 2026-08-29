package options

type TradeOrder interface {
	*Order | *ExecReport | map[string]any
}

type Order struct {
	Account                  int     `json:"-"`
	BusinessType             string  `json:"-"`
	CapitalAcctID            int     `json:"-"`
	ClOrdID                  string  `json:"clOrdID"`      //export 报单编号
	Counterparty             string  `json:"counterparty"` //export 账户名称 output:accountName
	Currency                 string  `json:"currency"`     //export 币种
	Customer                 string  `json:"-"`
	DirtyPrice               float64 `json:"dirtyPrice"`        //export 意向全价 成交全价(不含费)
	DirtyPriceWithFee        float64 `json:"dirtyPriceWithFee"` //export 成交全价(含费)
	EstiFrozenAmount         float64 `json:"-"`
	ExecUser                 string  `json:"-"`
	F_alg_params             string  `json:"-"`
	F_app_ord_id             string  `json:"-"`
	F_approve_status         int     `json:"f_approve_status"` //export 审批状态 output:approveStatus
	F_avg_px                 float64 `json:"-"`
	F_cl_ord_id              string  `json:"-"`
	F_cum_qty                int64   `json:"f_cum_qty"` //export 成交面额(万元) output:cumQty
	F_db_insert_time         int64   `json:"-"`
	F_last_px                float64 `json:"-"`
	F_last_shares            int64   `json:"-"`
	F_leaves_qty             int64   `json:"-"`
	F_ord_create_time        int64   `json:"f_ord_create_time"` //export 日期 交易申请时间 output:ordCreateTime
	F_ord_rej_reason         string  `json:"-"`
	F_ord_status             string  `json:"f_ord_status"`  //export 状态 output:ordStatus
	OrdStatusText            string  `json:"ordStatusText"` //状态描述
	F_ord_status_update_time int64   `json:"-"`
	F_reviewer               string  `json:"-"`
	F_trade_date             int     `json:"-"`
	F_transact_time          int64   `json:"-"`
	HandlInst                string  `json:"handlInst"` //export 交易效率
	InitMarginAmount         float64 `json:"-"`
	InitMarginRatio          float64 `json:"-"`
	InvestorID               string  `json:"-"`
	LimitCheckResult         int     `json:"-"`
	Locked                   bool    `json:"-"`
	OrdSource                string  `json:"ordSource"` //export 订单来源
	OrdType                  string  `json:"-"`
	ParValue                 float64 `json:"-"`
	PlanCode                 string  `json:"-"`
	Price                    float64 `json:"price"`    //export 意向净价 成交净价(不含费)
	Quantity                 float64 `json:"quantity"` //export 券面总额(万元) / 10000
	Remark                   string  `json:"-"`
	SecurityExchange         string  `json:"-"`
	SecurityID               int     `json:"-"`
	SecurityType             string  `json:"-"`
	SettlType                string  `json:"settlType"`      //export 清算速度
	Side                     string  `json:"side"`           //export 交易方向
	CommissionRate           float64 `json:"commissionRate"` //export 佣金费率
	Spread                   float64 `json:"spread"`         //export 利差
	Symbol                   string  `json:"symbol"`         //export 标的代码
	SymbolName               string  `json:"symbolName"`     //export 标的名称
	TraderID                 string  `json:"-"`
	TransactTime             string  `json:"-"`
	UltraContractCode        string  `json:"ultraContractCode"`    //export 大合约编号
	Ytm                      float64 `json:"ytm"`                  //export 意向到期收益(%)
	AvgPrice                 float64 `json:"avgPrice"`             // 成交净价
	AvgDirtyPrice            float64 `json:"avgDirtyPrice"`        // 成交全价
	AvgDirtyPriceWithFee     float64 `json:"avgDirtyPriceWithFee"` // 成交全价（含费）
}

type OrderOut struct {
	ClOrdID              string  `json:"clOrdID"`              //export 报单编号
	Counterparty         string  `json:"accountName"`          //export 账户名称 counterparty
	Currency             string  `json:"currency"`             //export 币种
	DirtyPrice           float64 `json:"dirtyPrice"`           //export 意向全价 成交全价(不含费)
	DirtyPriceWithFee    float64 `json:"dirtyPriceWithFee"`    //export 成交全价(含费)
	F_approve_status     int     `json:"approveStatus"`        //export 审批状态 f_approve_status
	F_cum_qty            int64   `json:"cumQty"`               //export 成交面额(万元) f_cum_qty
	F_ord_create_time    int64   `json:"ordCreateTime"`        //export 日期 交易申请时间 f_ord_create_time
	F_ord_status         string  `json:"ordStatus"`            //export 状态 f_ord_status
	OrdStatusText        string  `json:"ordStatusText"`        //状态描述
	HandlInst            string  `json:"handlInst"`            //export 交易效率
	OrdSource            string  `json:"ordSource"`            //export 订单来源
	Price                float64 `json:"price"`                //export 意向净价 成交净价(不含费)
	Quantity             float64 `json:"quantity"`             //export 券面总额(万元) / 10000
	SettlType            string  `json:"settlType"`            //export 清算速度
	Side                 string  `json:"side"`                 //export 交易方向
	CommissionRate       float64 `json:"commissionRate"`       //export 佣金费率
	Spread               float64 `json:"spread"`               //export 利差
	Symbol               string  `json:"symbol"`               //export 标的代码
	SymbolName           string  `json:"symbolName"`           //export 标的名称
	UltraContractCode    string  `json:"ultraContractCode"`    //export 大合约编号
	Ytm                  float64 `json:"ytm"`                  //export 意向到期收益(%)
	AvgPrice             float64 `json:"avgPrice"`             //export 成交净价
	AvgDirtyPrice        float64 `json:"avgDirtyPrice"`        //export 成交全价
	AvgDirtyPriceWithFee float64 `json:"avgDirtyPriceWithFee"` //export 成交全价（含费）
}

func (o *Order) ToOrderOut() *OrderOut {
	if o.HandlInst == "4" {
		o.HandlInst = "3"
	}
	return &OrderOut{
		ClOrdID:              o.ClOrdID,
		Counterparty:         o.Counterparty,
		Currency:             o.Currency,
		DirtyPrice:           o.DirtyPrice,
		DirtyPriceWithFee:    o.DirtyPriceWithFee,
		F_approve_status:     o.F_approve_status,
		F_cum_qty:            o.F_cum_qty,
		F_ord_create_time:    o.F_ord_create_time,
		F_ord_status:         o.F_ord_status,
		OrdStatusText:        o.OrdStatusText,
		HandlInst:            o.HandlInst,
		OrdSource:            o.OrdSource,
		Price:                o.Price,
		Quantity:             o.Quantity,
		SettlType:            o.SettlType,
		Side:                 o.Side,
		CommissionRate:       o.CommissionRate,
		Spread:               o.Spread,
		Symbol:               o.Symbol,
		SymbolName:           o.SymbolName,
		UltraContractCode:    o.UltraContractCode,
		Ytm:                  o.Ytm,
		AvgPrice:             o.AvgPrice,
		AvgDirtyPrice:        o.AvgDirtyPrice,
		AvgDirtyPriceWithFee: o.AvgDirtyPriceWithFee,
	}
}

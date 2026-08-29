package options

type ExecReport struct {
	Account                  int     `json:"account"`
	BusinessType             string  `json:"businessType"` //业务类型
	CapitalAcctID            int     `json:"capitalAcctID"`
	ClOrdID                  string  `json:"clOrdID"`
	CommissionRate           float64 `json:"commissionRate"` // 佣金费率
	Counterparty             string  `json:"counterparty"`   // 账户名称
	Currency                 string  `json:"currency"`       // 币种
	Customer                 string  `json:"customer"`
	DirtyPrice               float64 `json:"dirtyPrice"`        // 意向全价 成交全价(不含费)
	DirtyPriceWithFee        float64 `json:"dirtyPriceWithFee"` // 成交全价(含费)
	EstiFrozenAmount         float64 `json:"estiFrozenAmount"`
	ExecUser                 string  `json:"execUser"` // 执行交易员
	F_alg_params             string  `json:"f_alg_params"`
	F_app_ord_id             string  `json:"f_app_ord_id"`
	F_approve_status         int     `json:"f_approve_status"` // 状态
	F_avg_px                 float64 `json:"f_avg_px"`         // 成交均价
	F_channel_code           string  `json:"f_channel_code"`
	F_cl_ord_id              string  `json:"f_cl_ord_id"`
	F_cum_qty                int64   `json:"f_cum_qty"` // 成交面额(万元)
	F_db_insert_time         int64   `json:"f_db_insert_time"`
	F_exec_id                string  `json:"f_exec_id"`
	F_exec_ref_id            string  `json:"f_exec_ref_id"`
	F_exec_trans_type        string  `json:"f_exec_trans_type"`
	F_exec_type              string  `json:"f_exec_type"`
	F_last_px                float64 `json:"f_last_px"`
	F_last_shares            int64   `json:"f_last_shares"`
	F_leaves_qty             int64   `json:"f_leaves_qty"`
	F_msg_time               int64   `json:"f_msg_time"`
	F_ord_create_time        int64   `json:"f_ord_create_time"` // 日期 交易申请时间
	F_ord_rej_reason         string  `json:"f_ord_rej_reason"`
	F_ord_status             string  `json:"f_ord_status"` // 状态
	F_ord_status_update_time int64   `json:"f_ord_status_update_time"`
	F_orig_cl_ord_id         string  `json:"f_orig_cl_ord_id"`
	F_reviewer               string  `json:"f_reviewer"`
	F_trade_date             int     `json:"f_trade_date"`
	F_transact_time          int64   `json:"f_transact_time"`
	HandlInst                string  `json:"handlInst"` // 交易效率
	InitMarginAmount         float64 `json:"initMarginAmount"`
	InitMarginRatio          float64 `json:"initMarginRatio"`
	InvestorID               string  `json:"investorID"`
	LimitCheckResult         int     `json:"limitCheckResult"`
	Locked                   bool    `json:"locked"`
	OrdSource                string  `json:"ordSource"` // 订单来源
	OrdType                  string  `json:"ordType"`
	ParValue                 float64 `json:"parValue"`
	PlanCode                 string  `json:"planCode"`
	Price                    float64 `json:"price"`    // 意向净价 成交净价(不含费)
	Quantity                 float64 `json:"quantity"` // 券面总额(万元) / 10000
	Remark                   string  `json:"remark"`
	SecurityExchange         string  `json:"securityExchange"`
	SecurityID               int     `json:"securityID"`
	SecurityType             string  `json:"securityType"`
	SettlType                string  `json:"settlType"`  // 清算速度
	Side                     string  `json:"side"`       // 交易方向
	Spread                   float64 `json:"spread"`     // 利差
	Symbol                   string  `json:"symbol"`     // 标的代码
	SymbolName               string  `json:"symbolName"` // 标的名称
	TraderID                 string  `json:"traderID"`
	TransactTime             string  `json:"transactTime"`
	UltraContractCode        string  `json:"ultraContractCode"` // 大合约编号
	Ytm                      float64 `json:"ytm"`               // 意向到期收益(%)
}
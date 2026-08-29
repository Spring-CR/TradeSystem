package schema

type TaskStock struct {
	MarketNo         string  `json:"market_no"`         // 市场
	ReportCode       string  `json:"report_code"`       // 证券代码
	EntrustDirection string  `json:"entrust_direction"` // 委托方向
	OpenClose        string  `json:"open_close"`        // 开平仓标志
	Amount           float64 `json:"amount"`            // 指令数量
	Balance          float64 `json:"balance"`           // 指令金额
	ContractSize     float64 `json:"contract_size"`     // 合约乘数
	Price            float64 `json:"price"`             // 指令价格
	InvestType       string  `json:"invest_type"`       // 投资类型
}

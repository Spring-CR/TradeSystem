package options

import (
	"ficc-utils/common/utils"
)

type ContractPosition struct {
	CounterpartyID             int          `json:"-"`
	Counterparty               string       `json:"Counterparty"`       // 账户名称
	PlanID                     int          `json:"-"`
	PlanCode                   string       `json:"-"`
	UltraContractID            int          `json:"-"`
	UltraContractCode          string       `json:"UltraContractCode"`  // 大合约编号
	ContractID                 int          `json:"-"`
	ContractCode               string       `json:"ContractCode"`       // 合约编号
	SecurityID                 int64        `json:"-"`
	Symbol                     string       `json:"Symbol"`             // 标的代码
	SecurityName               string       `json:"SecurityName"`       // 标的名称
	SecurityExchange           string       `json:"SecurityExchange"`   // 交易市场
	LongShort                  string       `json:"LongShort"`          // 合约方向 LONG SHORT
	TradeDate                  string       `json:"TradeDate"`          // 交易达成日
	StartDate                  string       `json:"StartDate"`          // 期初定价日
	EndDate                    string       `json:"EndDate"`            // 期末定价日
	Notional                   float64      `json:"Notional"`           // 期初名义本金
	OpenTradeAmount            float64      `json:"OpenTradeAmount"`    // 开仓面额(元)
	OpenBondNetPrice           float64      `json:"OpenBondNetPrice"`   // 开仓净价
	OpenBondGrossPrice         float64      `json:"OpenBondGrossPrice"` // 开仓全价
	InitPrice                  float64      `json:"InitPrice"`          // 期初价格
	InitQuantity               int64        `json:"-"`
	DynamicNotional            float64      `json:"DynamicNotional"`    // 剩余名义本金
	Currency                   string       `json:"Currency"`           // 币种
	InterestSettlementDate     string       `json:"-"`
	InterestSettlementSpeed    string       `json:"-"`
	Quantity                   float64      `json:"-"`
	StructureSettlementDate    string       `json:"-"`
	StructureSettlementSpeed   string       `json:"-"`
}

type ContractPositionOut struct {
	CounterpartyID             int          `json:"-"`
	Counterparty               string       `json:"accountName"`        // 账户名称
	PlanID                     int          `json:"-"`
	PlanCode                   string       `json:"-"`
	UltraContractID            int          `json:"-"`
	UltraContractCode          string       `json:"ultraContractCode"`  // 大合约编号
	ContractID                 int          `json:"-"`
	ContractCode               string       `json:"contractCode"`       // 合约编号
	SecurityID                 int64        `json:"-"`
	Symbol                     string       `json:"symbol"`             // 标的代码
	SecurityName               string       `json:"symbolName"`         // 标的名称
	SecurityExchange           string       `json:"securityExchange"`   // 交易市场
	LongShort                  string       `json:"longShort"`          // 合约方向 LONG SHORT
	TradeDate                  string       `json:"tradeDate"`          // 交易达成日
	StartDate                  string       `json:"startDate"`          // 期初定价日
	EndDate                    string       `json:"endDate"`            // 期末定价日
	Notional                   float64      `json:"notional"`           // 期初名义本金
	OpenTradeAmount            float64      `json:"openTradeAmount"`    // 开仓面额(元)
	OpenBondNetPrice           float64      `json:"openBondNetPrice"`   // 开仓净价
	OpenBondGrossPrice         float64      `json:"openBondGrossPrice"` // 开仓全价
	InitPrice                  float64      `json:"initPrice"`          // 期初价格
	InitQuantity               int64        `json:"-"`
	DynamicNotional            float64      `json:"dynamicNotional"`    // 剩余名义本金
	Currency                   string       `json:"currency"`           // 币种
	InterestSettlementDate     string       `json:"-"`
	InterestSettlementSpeed    string       `json:"-"`
	Quantity                   float64      `json:"-"`
	StructureSettlementDate    string       `json:"-"`
	StructureSettlementSpeed   string       `json:"-"`
}

func (c *ContractPosition) ToContractPositionOut() *ContractPositionOut {
	out := ContractPositionOut{}
	utils.CopyStruct(c, &out)
	return &out
}
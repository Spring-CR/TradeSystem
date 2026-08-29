package options

type PositionResult struct {
	Total int 			`json:"total"`
	Data []Position 	`json:"data"`
}

type Position struct {
    // 标的属性
    Account                    int          `json:"account"`
    CounterpartyID             int          `json:"-"`
    Counterparty               string       `json:"counterparty"`
    Symbol                     string       `json:"symbol"`
    SymbolName                 string       `json:"symbolName"`
    Currency                   string       `json:"currency"`
    PlanCode                   string       `json:"planCode"`
    UltraContractCode          string       `json:"ultraContractCode"`
    SecurityExchange           string       `json:"securityExchange"`
    SecurityType               string       `json:"securityType"`
    ParValue                   float64      `json:"parValue"`
    // 统计属性
    NetPositionT0              float64      `json:"netPositionT0"`              // T+0清算速度的净持仓
    NetPositionT1              float64      `json:"netPositionT1"`              // T+1清算速度的净持仓
    LongAvailablePositionT0    float64      `json:"longAvailablePositionT0"`    // T+0多头可用持仓
    LongAvailablePositionT1    float64      `json:"longAvailablePositionT1"`    // T+1多头可用持仓
    ShortAvailablePositionT0   float64      `json:"shortAvailablePositionT0"`   // T+0空头可用持仓
    ShortAvailablePositionT1   float64      `json:"shortAvailablePositionT1"`   // T+1空头可用持仓
    LongCleanPriceCost         float64      `json:"-"`                          // 多头净价持仓成本
    LongDirtyPriceCost         float64      `json:"-"`                          // 多头全价持仓成本
    LongDirtyPriceWithFeeCost  float64      `json:"-"`                          // 多头全价（含费）持仓成本
    ShortCleanPriceCost        float64      `json:"-"`                          // 空头净价持仓成本
    ShortDirtyPriceCost        float64      `json:"-"`                          // 空头全价持仓成本
    ShortDirtyPriceWithFeeCost float64      `json:"-"`                          // 空头全价（含费）持仓成本
    LongAvgCleanPrice          float64      `json:"longAvgCleanPrice"`          // 多头净价持仓均价
    LongAvgDirtyPrice          float64      `json:"-"`                          // 多头全价持仓均价
    LongAvgDirtyPriceWithFee   float64      `json:"longAvgDirtyPriceWithFee"`   // 多头全价（含费）持仓均价
    ShortAvgCleanPrice         float64      `json:"shortAvgCleanPrice"`         // 空头净价持仓均价
    ShortAvgDirtyPrice         float64      `json:"-"`                          // 空头全价持仓均价
    ShortAvgDirtyPriceWithFee  float64      `json:"shortAvgDirtyPriceWithFee"`  // 空头全价（含费）持仓均价
    MaxLongMarginOccupancy     float64      `json:"-"`                          // 最大多头保证金占用
    MaxShortMarginOccupancy    float64      `json:"-"`                          // 最大空头保证金占用
    MaxMarginOccupancy         float64      `json:"-"`                          // 最大单边保证金占用
}
package schema

// 交易台交易指令
type TradeInstr struct {
	// 表主键，自增
	ID int64 `json:"-"`
	// 消息类型，默认设置为 D
	MsgType string `sql:"size: 1" json:"msgType"`
	// 取指令的组合编号
	ClientID string `sql:"size: 16" json:"clientID"`
	// 为了方便与instr-stock设置的key
	ParentKey string `sql:"index: i_ti_parent, size: 64" json:"-"`
	// 自定义母单识别标志
	SecondaryClOrdID string `sql:"unique: pk_trade_instr, size: 128" json:"secondaryClOrdID"`
	// 取指令的证券代码
	SecurityID string `sql:"size: 16" json:"securityID"`
	// 取指令的证券代码
	Symbol string `sql:"size: 16" json:"symbol"`
	// 取指令的买卖方向：3 - 买；1 - 卖
	Side string `sql:"size: 1" json:"side"`
	// 订单创建时间，格式如：20220512-08:54:49
	TransactTime string `sql:"size: 17" json:"transactTime"`
	// 订单数量，默认取指令数量
	OrderQty float64 `json:"orderQty"`
	// 价格类型,  默认："2" ，限价。'H' - 最优五档立即成交剩余撤单（沪市）；'I' - 最优五档立即成交剩余转限价（沪市）；'U' - 对方最优价格（深市）；'V' - 本方最优价格（深市）；'W' - 立即成交否则撤单（深市）；'X' - 最优五档立即成交否则撤销（深市）；'Y' - 全部成交否则撤单（深市）；‘L’ - 五档即成剩撤（中金所五档市价）；‘M’ -  五档即成剩转（中金所五档市价转限价）；‘N’ - 最优一档即成剩撤（中金所最优价）；‘O’ - 最优一档即成剩转（中金所最优价）
	OrdType string `sql:"size: 1" json:"ordType"`
	// 订单有效时间类型，默认："0"，当日有效
	TimeInForce string `sql:"size: 1" json:"timeInForce"`

	// 委托价格，取指令价格
	Price float64 `json:"price"`

	// 策略名，默认：“0”，单笔策略。
	TargetStrategy string `sql:"size: 1" json:"targetStrategy"`
	// 策略参数，json字符串
	StrategyParametersText string                 `sql:"type: MEDIUMTEXT" json:"-"`
	StrategyParameters     map[string]interface{} `sql:"-" json:"strategyParameters"`

	// 取指令的市场代码
	MarketCode string `sql:"size: 8" json:"marketCode"`

	// 执行交易的人员的ERP id，取titans登录用户的信息（erp id）
	UserText string                 `json:"-"`
	User     map[string]interface{} `sql:"-" json:"user"`

	// 开平标志，取指令开平仓标志
	OpenClose string `sql:"size: 1" json:"openClose"`

	// ATP接口用户，测试环境暂时写死“4012”，生产再看ATP柜台账号情况
	ApiOperator string `sql:"size: 32" json:"operator"`

	// 以下两个字段值保留成交时的值
	AvgPx       float64 `json:"-"`
	CumAmt      float64 `json:"-"`
	CumQty      float64 `json:"-"`
	CumTotalFee float64 `json:"-"`

	// 以下4个状态只保留终态的状态
	OrdStatus string `sql:"size: 2" json:"-"`
	// 状态更新时间
	StatusUpdateTime int64 `json:"-"`
	// 状态更新文本
	StatusUpdateText string `sql:"type: MEDIUMTEXT" json:"-"`
	// Kafka偏移量
	StatusKafkaOffset int64 `json:"-"`
	// 消息时间戳
	MessageTime int64 `json:"-"`

	// 从kafka获取消息后更新状态时的消息offset
	// KafkaOffset int64 `json:"-"`
	// GFQG反馈文本
	// RespText string `sql:"type: MEDIUMTEXT" json:"text"`

	// 以下几个属性，是为了测试增加的属性
	/*
		"securityExchange": "SS",
		"cashOrderQty": 0,
		"currency": "CNY",
		"businessType": "AIRBAG",
		"tradePurpose": "期权其他",
		"dynamicNotional": 1000000,
		 "closeOutNotional": 1000000,
		 "contractInsId": 10086,
	*/
	//SecurityExchange string `sql:"-" json:"securityExchange"`
	CashOrderQty int `sql:"-" json:"cashOrderQty"`
	// Currency         string `sql:"-" json:"currency"`
	// BusinessType     string `sql:"-" json:"businessType"`
	// TradePurpose     string `sql:"-" json:"tradePurpose"`
	// DynamicNotional  int64  `sql:"-" json:"dynamicNotional"`
	// CloseOutNotional int64  `sql:"-" json:"closeOutNotional"`
	// ContractInsId    int64  `sql:"-" json:"contractInsId"`
	ClOrdID     string `sql:"size: 128" json:"clOrdID"`
	OrigClOrdID string `sql:"size: 128" json:"origClOrdID"`
}

type TradeInstrResp struct {
	ID                int64  `json:"-"`
	SecondaryClOrdID  string `sql:"unique: pk_tir, index: i_trade_instr_resp, size: 128" json:"secondaryClOrdID"`
	MessageTime       int64
	MsgType           string  `sql:"size: 1" json:"msgType"` // msgType，配合四个状态才是终态
	TransactTime      string  `sql:"size: 17" json:"transactTime"` // 订单创建时间，格式如：20220512-08:54:49
	StatusKafkaOffset int64   `sql:"index: pk_tir;i_tir_o" json:"-"`
	OrdStatus         string  `sql:"size: 2" json:"ordStatus"`
	StatusUpdateText  string  `sql:"type: MEDIUMTEXT" json:"text"`
	AvgPx             float64 `json:"avgPx"`
	CumAmt            float64 `json:"cumAmt"`
	CumQty            float64 `json:"cumQty"`
	CumTotalFee       float64 `json:"cumTotalFee"`
	ClOrdID           string  `sql:"size: 128" json:"clOrdID"`
	OrigClOrdID       string  `sql:"size: 128" json:"origClOrdID"`
}

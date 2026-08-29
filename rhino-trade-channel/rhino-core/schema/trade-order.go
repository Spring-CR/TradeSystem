package schema

type TradeOrder struct {
	ID                     int64
	SystemCode             string `sql:"size: 32"`
	BusinessCode           string `sql:"size: 32"`
	ClGroupOrdID           string `sql:"size: 188"`
	ClOrdID                string `sql:"index: idx_to_clordid, size: 188"`
	Account                string `sql:"size: 64"`
	HandlInst              string `sql:"size: 2"`
	AppOrdID               string `sql:"unique: uq_appordid, size: 188"`
	OrdID                  string `sql:"size: 188"` //`sql:"index: idx_to_ordid, size: 188"`// 去掉这个索引快很多
	ParentClOrdID          string `sql:"size: 188"`
	IsDirectOrd            bool
	IsAlgOrd               bool
	IsSubAlgOrd            bool
	IsInstrOrd             bool
	IsSubInstrOrd          bool
	IsCrossDateOrd         bool
	IsSubCrossDateOrd      bool
	MinQty                 float64
	SecurityExchange       string `sql:"size: 8"`
	SecurityExchangeRegion string `sql:"size: 4"`
	Symbol                 string `sql:"size: 64"`
	Symbol2                string `sql:"-"`
	SymbolSfx              string `sql:"size: 8"`
	SecurityID             string `sql:"size: 64"`
	IDSource               string `sql:"size: 2"`
	SecurityType           string `sql:"size: 16"`
	Side                   string `sql:"size: 2"`
	TransactTime           int64
	TradeDate              int64
	OrderQty               float64
	CashOrderQty           float64
	OrdType                string `sql:"size: 2"`
	Price                  float64
	Currency               string `sql:"size: 4"`
	EffectiveTime          int64
	ExpireDate             int64
	ExpireTime             int64
	OpenClose              string `sql:"size: 2"`
	ContractMultiplier     float64
	OrdChangedCount        int
	OrdCancelCount         int
	ExtendAttr             string `sql:"type: MEDIUMTEXT"`
	AlgParams              string `sql:"type: MEDIUMTEXT"`
	AlgName                string `sql:"size: 32"`
	OrdCreator             string `sql:"size: 64"`
	OrdCreateTime          int64
	OrdDraftUpdateUser     string `sql:"size: 64"`
	OrdDraftUpdateTime     int64
	OrdDraftDelFlag        int
	OrdDraftDelUser        string `sql:"size: 64"`
	OrdDraftDelTime        int64
	OrdExecUserScope       string
	OrdExecUser            string `sql:"size: 64"`
	OrdStatusUpdateTime    int64
	OrdStatusUpdateTime2   int64  `sql:"-"` // 这里不包含状态为pending cancel的状态更新时间
	OrdStatus              string `sql:"size: 2"`
	OrdStatus2             string `sql:"-"`
	OrdFillStatus          string `sql:"-"` //`sql:"size: 2"` 不需要存数据库，根据OrdStatus、剩余单量、订单下单参数的单量，就能判断出成交状态了
	ReviewFlag             string `sql:"size: 2"`
	ReviewerScope          string
	Reviewer               string
	ApproveStatus          int
	ReviewTime             int64
	OrderSubmitFailReason  string
	PushInQueueBeforeTrade bool
	LastShares             int64   `sql:"-"` //
	LastPx                 float64 `sql:"-"` //
	LeavesQty              int64   `sql:"-"` //
	CumQty                 int64   `sql:"-"` //
	AvgPx                  float64 `sql:"-"` //
	OrdRejReason           string  `sql:"-"`
	LatestActionType       string  `sql:"size: 2"`
	ChannelCode            string  `sql:"size: 32"`
	DBInsertTime           int64
	MsgSeq                 int64 //`sql:"index: idx_to_msgseq"` // 去掉这个索引快很多
	WorkerAffinity         int
	ExtendAttrMap          map[string]interface{} `sql:"-"`
	AlgParamsMap           map[string]interface{} `sql:"-"`
	DBInsertOnOrdExec      bool                   `sql:"-"` // 是否发生数据库插入，true代表数据库插入，false代表数据库更新。对于是数据库新纪录的插入还是草稿更新，会影响orderCache AddRootTradeOrder的逻辑
	QuotaValidateTime      int64
}

func (o *TradeOrder) GetCacheKey() string {
	return o.AppOrdID
}

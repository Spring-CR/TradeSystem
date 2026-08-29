package schema

type TaskInstrStock struct {

	// 表主键，自增
	ID int64

	// 业务日期，记录创建的日期(8位，类型万徳)
	Date int `sql:"unique: pk_tis, fk: date@task_instrs, fkGroup: fk_tis_to_ti, index: i_f_tis;i_tis_date"`
	// 指令的ID，每天从1开始，全系统当日唯一
	DailyInstrNo int64 `sql:"fk: daily_instr_no@task_instrs, fkGroup: fk_tis_to_ti, index: pk_tis;i_f_tis"`
	// 指令新下达时，该字段的值为1，当指令有修改时，该字段的值依次递增
	IndexDailyModify int64 `sql:"fk: index_daily_modify@task_instrs, fkGroup: fk_tis_to_ti, index: pk_tis;i_f_tis"`
	// 指令证券序号，批量指令的ID，每天从1开始，全系统当日唯一
	StockSerialNo int64 `sql:"index: pk_tis"`

	// 市场
	MarketNo string `sql:"size: 8"`
	// 证券代码
	ReportCode string `sql:"size: 16"`
	// 委托方向
	EntrustDirection string`sql:"size: 2"`
	// 开平标志，取指令开平仓标志
	OpenClose string `sql:"size: 8"`
	// 投资类型，a:投机； b:套保; c:套利
	InvestType string `sql:"size: 2"`
	// 指令数量
	Amount float64
	// 指令金额
	Balance float64
	// 合约乘数
	ContractSize float64
	// 指令价格
	Price float64
	// 委托完成状态 1：未执行 2：部分执行 3：完成
	StockEntrustExecuteStatus string`sql:"size: 1"`
	// 成交完成状态 1：未执行 2：部分执行 3：完成
	// 空白表示未执行，r-就绪(Ready) t-超时(Timeout) 0-新（New）(接受) 1-部分成交（Partially filled） 8-已拒绝（Rejected）2-已成交（Filled）C-已过期（Expired）4-已撤销（Canceled）5-已修改（Replaced）6-待撤销（Pending Cancel）9-已暂停（结合tag8001判断暂停类型）D-重申（Restated )
	StockDealExecuteStatus string`sql:"size: 1"`
	// 累计成交数量
	TotalDealAmount float64
	// 累计成交金额
	TotalDealBalance float64
	// 累计成交均价
	CumAvgPrice float64
	// 累计委托数量
	TotalEntrustAmount float64
	// 累计委托金额
	TotalEntrustBalance float64
	// 指令成交完成时间
	DealCompleteDateTime int64
	// 预估费用
	EstimateFee float64

	// 执行时间。与执行人一起，在创建交易台交易指令之后后更新
	StockInstrExecutionTime int64
	// 执行人
	StockInstrOperator string `sql:"size: 32"`
	// // 状态更新时间
	// StatusUpdateTime int64
	// // 状态更新文本
	// StatusUpdateText string `sql:"type: MEDIUMTEXT" json:"text"`
	// // Kafka偏移量
	// StatusKafkaOffset int64
}
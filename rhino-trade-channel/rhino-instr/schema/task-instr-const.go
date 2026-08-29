package schema

const (
	TaskInstrInstrTypeStock      = "1" // 股票
	TaskInstrInstrTypeCombi      = "2" // 组合
	TaskInstrInstrTypeStockBatch = "3" // 股票批量
	TaskInstrInstrTypeCombiBatch = "4" // 组合批量

	TaskInstrMarketNoSH = "1" // 上海
	TaskInstrMarketNoSZ = "2" // 深圳
	TaskInstrMarketNoZJ = "3" // 中金所

	TaskInstrEntrustDirectionBuy  = "3" // 买
	TaskInstrEntrustDirectionSell = "1" // 卖

	TaskInstrOpenCloseOpen  = "O"    // 开仓
	TaskInstrOpenCloseClose = "C"    // 平仓
	TaskInstrOpenAuto       = "AUTO" // 自动转开平

	TaskInstrEntrustExecuteStatusNotBegin = "1"
	TaskInstrEntrustExecuteStatusRunning = "2"
	TaskInstrEntrustExecuteStatusFinish = "3"

	TaskInstrDealExecuteStatusNotBegin = "1"
	TaskInstrDealExecuteStatusRunning = "2"
	TaskInstrDealExecuteStatusFinish = "3"
)

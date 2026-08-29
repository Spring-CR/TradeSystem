package schema

/*
空白表示未执行，r-就绪(Ready) t-超时(Timeout) 0-新（New）(接受) 1-部分成交（Partially filled） 8-已拒绝（Rejected）2-已成交（Filled）C-已过期（Expired）
4-已撤销（Canceled）5-已修改（Replaced）6-待撤销（Pending Cancel）9-已暂停（结合tag8001判断暂停类型）
*/
const (
	TradeInstrOrdStatusEmpty           = ""
	TradeInstrOrdStatusReady           = "r"
	TradeInstrOrdStatusTimeout         = "t"
	TradeInstrOrdStatusNew             = "0"
	TradeInstrOrdStatusPartiallyFilled = "1"
	TradeInstrOrdStatusRejected        = "8"
	TradeInstrOrdStatusFilled          = "2"
	TradeInstrOrdStatusExpired         = "C"
	TradeInstrOrdStatusCanceled        = "4"
	TradeInstrOrdStatusReplaced        = "5"
	TradeInstrOrdStatusPendingCancel   = "6"
	TradeInstrOrdStatusStop            = "9"
)


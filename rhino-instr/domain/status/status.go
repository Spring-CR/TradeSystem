package status

const (
	// 1：未分发 2：已分发 3：分发拒绝
	DispenseStatusUnprocessed = "1"
	DispenseStatusProcessed   = "2"
	DispenseStatusRefused     = "3"
)

const (
	// 指令状态：1、有效 2、已修改 3、已撤销 4、已暂停 5、审批拒绝 6、分发拒绝 7、指令录入 8、分仓失败 9、草稿指令 a、临时下达 b、临时修改 c、撤消失败
	InstrStatusEffective            = "1" // 有效
	InstrStatusModified             = "2" // 已修改
	InstrStatusRevoked              = "3" // 已撤销
	InstrStatusSuspended            = "4" // 已暂停
	InstrStatusApprovalRejected     = "5" // 审批拒绝
	InstrStatusDistributionRejected = "6" // 分发拒绝
	InstrStatusEntry                = "7" // 指令录入
	InstrStatusWarehouseFailure     = "8" // 分仓失败
	InstrStatusDraft                = "9" // 草稿指令
	InstrStatusTemporaryIssued      = "a" // 临时下达
	InstrStatusTemporaryModified    = "b" // 临时修改
	InstrStatusRevocationFailed     = "c" // 撤消失败
)
/*
const (
	// 委托完成状态：1：未执行 2：部分执行 3：完成 4：委托失败
	EntrustExecuteStatusNotExecuted       = "1" // 未执行
	EntrustExecuteStatusPartiallyExecuted = "2" // 部分执行
	EntrustExecuteStatusCompleted         = "3" // 完成
	EntrustExecuteStatusFailed            = "4" // 失败
)

const (
	StockInstrStatusNotExecuted     = ""  //未执行
	StockInstrStatusReady           = "r" //就绪
	StockInstrStatusTimeout         = "t" //超时
	StockInstrStatusNew             = "0" //0-新建，挂单成功，接受
	StockInstrStatusPartiallyFilled = "1" //1-部分成交（Partially filled）
	StockInstrStatusRejected        = "8" //8-已拒绝（Rejected）
	StockInstrStatusFilled          = "2" //2-已成交（Filled）
	StockInstrStatusExpired         = "C" //C-已过期（Expired）
	StockInstrStatusCanceled        = "4" //4-已撤销（Canceled）
	StockInstrStatusReplaced        = "5" //5-已修改（Replaced）
	StockInstrStatusPendingCancel   = "6" //6-待撤销（Pending Cancel）
	StockInstrStatusSuspend         = "9" //9-已暂停（结合tag8001判断暂停类型）
	StockInstrStatusRestated        = "D" //D-重申（Restated )
)

const (
	// 成交完成状态
	DealExecuteStatusNotExecuted       = "1" // 未执行
	DealExecuteStatusPartiallyExecuted = "2" // 部分执行
	DealExecuteStatusCompleted         = "3" // 完成
	DealExecuteStatusFailed            = "4" // 失败
)*/

const (
	TradeInstrParentKey = "%v-%v-%v-%v"
	SecondaryClOrdIDPrefix  = "rhino-"
	SecondaryClOrdIDPattern = SecondaryClOrdIDPrefix + "%v-%v-%v-%v-%v"
	TransactTimeLayout = "20060102-15:04:05"
)
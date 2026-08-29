package enum

// 订单的执行方式
type OrderHandlInst string

const (
	OrderHandlInst_DMA    OrderHandlInst = "1"  //直通交易，订单通过柜台直接报送交易所
	OrderHandlInst_DSA    OrderHandlInst = "2"  //柜台侧的算法交易，即柜台原生支持算法拆单，在柜台端直接分笔成交，仅返回部分成交的结果给交易台
	OrderHandlInst_ALG    OrderHandlInst = "21" //第三方的算法拆单交易。对接第三方算法拆单服务时，算法服务仅负责拆单，不负责报送。拆单服务把拆好的子单返还交易台，交易台再报送柜台，交易台需要跟踪这些算法子单的状态
	OrderHandlInst_MANUAL OrderHandlInst = "3"  //手动订单，柜台侧人工执行订单交易动作。（生产用得较少，一般在测试阶段会用到）
)

// 用于标识SecurityID字段中提供的证券标识符的类别
type IDSource string

const (
	IDSource_CUSIP           IDSource = "1"
	IDSource_SEDOL           IDSource = "2"
	IDSource_QUIK            IDSource = "3"
	IDSource_ISIN_NUMBER     IDSource = "4"
	IDSource_RIC_CODE        IDSource = "5"
	IDSource_EXCHANGE_SYMBOL IDSource = "8"
)

// 交易方向.
type Side string

const (
	Side_Buy             Side = "1"
	Side_Sell            Side = "2"
	Side_BuyMinus        Side = "3"
	Side_SellPlus        Side = "4"
	Side_SellShort       Side = "5" // 借券卖出
	Side_SellShortExempt Side = "6"
	Side_Undisclosed     Side = "7"
	Side_Cross           Side = "8"
	Side_CrossShort      Side = "9"
)

// 订单类型
type OrdType string

const (
	OrdType_MARKET OrdType = "1" //市价单
	OrdType_LIMIT  OrdType = "2" //限价单
)

// 货币类型(全球主流金融市场和交易货币)
type Currency string

const (
	Currency_CNY = "CNY" // 人名币
	Currency_HKD = "HKD" // 港币
	Currency_USD = "USD" // 美元
	Currency_EUR = "EUR" // 欧元
	Currency_GBP = "GBP" // 英镑
	Currency_JPY = "JPY" // 日元
	Currency_CAD = "CAD" // 加元
	Currency_AUD = "AUD" // 澳元
	Currency_CHF = "CHF" // 瑞士法郎
	Currency_SGD = "SGD" // 新加坡元
	Currency_KRW = "KRW" // 韩元
	Currency_INR = "INR" // 印度卢比
	Currency_BRL = "BRL" // 巴西雷亚尔
	Currency_RUB = "RUB" // 俄罗斯卢布
	Currency_ZAR = "ZAR" // 南非兰特
	Currency_SAR = "SAR" // 沙特里亚尔
	Currency_AED = "AED" // 阿联酋迪拉姆
	Currency_QAR = "QAR" // 卡塔尔里亚尔
)

// 开平仓标记
type OpenClose string

const (
	OpenClose_CLOSE  OpenClose = "C" // 平仓
	OpenClose_OPEN   OpenClose = "O" // 开仓
	OpenClose_ROLLED OpenClose = "R" // fix 4.4
	OpenClose_FIFO   OpenClose = "F" // fix 4.4
)

// 订单状态
type OrdStatus string

const (
	OrdStatus_Draft                OrdStatus = "a"
	OrdStatus_PendingReview        OrdStatus = "b"
	OrdStatus_ReviewApproved       OrdStatus = "c"
	OrdStatus_ReviewRejected       OrdStatus = "d"
	OrdStatus_Ready                OrdStatus = "e"
	OrdStatus_Submit               OrdStatus = "f"
	OrdStatus_InternalSubmitFailed OrdStatus = "g"
	OrdStatus_Draft_Deleted        OrdStatus = "h"
	OrdStatus_ReviewTimeout        OrdStatus = "i"
	OrdStatus_ReviewCanceled       OrdStatus = "j"
	OrdStatus_New                  OrdStatus = "0"
	OrdStatus_PartiallyFilled      OrdStatus = "1"
	OrdStatus_Filled               OrdStatus = "2"
	OrdStatus_DoneForDay           OrdStatus = "3"
	OrdStatus_Canceled             OrdStatus = "4"
	OrdStatus_Replaced             OrdStatus = "5"
	OrdStatus_PendingCancel        OrdStatus = "6"
	OrdStatus_Stopped              OrdStatus = "7"
	OrdStatus_Rejected             OrdStatus = "8"
	OrdStatus_Suspended            OrdStatus = "9"
	OrdStatus_PendingNew           OrdStatus = "A"
	OrdStatus_Calculated           OrdStatus = "B"
	OrdStatus_Expired              OrdStatus = "C"
	OrdStatus_AcceptedForBidding   OrdStatus = "D"
	OrdStatus_PendingReplace       OrdStatus = "E"
)

var (
	ordStatus_DisplayName = map[string]string{
		"a": "Draft",
		"b": "PendingReview",
		"c": "ReviewApproved",
		"d": "ReviewRejected",
		"e": "Ready",
		"f": "Submit",
		"g": "InternalSubmitFailed",
		"h": "DraftDeleted",
		"i": "ReviewTimeout",
		"j": "ReviewingOrderCanceled",
		"0": "New",
		"1": "PartiallyFilled",
		"2": "Filled",
		"3": "DoneForDay",
		"4": "Canceled",
		"5": "Replaced",
		"6": "PendingCancel",
		"7": "Stopped",
		"8": "Rejected",
		"9": "Suspended",
		"A": "PendingNew",
		"B": "Calculated",
		"C": "Expired",
		"D": "AcceptedForBidding",
		"E": "PendingReplace",
	}
)

func GetOrdStatusDisplayName(status string) string {
	return ordStatus_DisplayName[status]
}

// 成交状态
type OrdFillStatus string

const (
	OrdFillStatus_Empty           OrdFillStatus = "0" // 未成交
	OrdFillStatus_PartiallyFilled OrdFillStatus = "1" // 部分成交
	OrdFillStatus_Filled          OrdFillStatus = "2" // 全部成交
)

// 订单动作
type ActionType string

const (
	ActionType_Draft             ActionType = "a"
	ActionType_SubmitForReview   ActionType = "b"
	ActionType_WithdrawForReview ActionType = "c"
	ActionType_ReviewCompleted   ActionType = "d"
	ActionType_Delete            ActionType = "e"
	ActionType_CancelForReview   ActionType = "f"
	ActionType_New               ActionType = "0"
	ActionType_Update            ActionType = "1"
	ActionType_Withdraw          ActionType = "2"
)

// 审批状态
type ApproveStatus int

const (
	ApproveStatus_NotSubmit     ApproveStatus = 0
	ApproveStatus_PendingReview ApproveStatus = 1
	ApproveStatus_Approved      ApproveStatus = 2
	ApproveStatus_Rejected      ApproveStatus = 3
)

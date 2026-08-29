package orderutils

import "strings"

const (
	DRAFT string = "待审核"
	VERIFYING string = "待审核"
	FICC_REJECTED string = "审核不通过"
	TOTRADE string = "待交易"
	ABANDONED string = "废单(意向)"
	TIMEOUT_FAIL string = "废单(未成交)"
	TIMEOUT_REJECTED string = "废单(超时未审核)"
	REJECTED string = "系统拒单"
	NEW string = "已报"
	PARTIALLY_FILLED string = "部成"
	PARTIALLY_CANCELED string = "部成部撤"
	FILLED string = "全部成交"
	PENDING_CANCEL string = "待撤单"
	CANCELED string = "已撤单"
	COMPLETED string = "部分成交"
)

func ConvertTrsOrderStatus(ordStatus string, approveStatus int, cumQty int64) string {
	switch ordStatus {
	case "a": return DRAFT
	case "b": return VERIFYING
	case "c": return "审核通过"
	case "d": return FICC_REJECTED
	case "f": return TOTRADE
	case "g": return ABANDONED
	case "0","A": return NEW
	case "1": return PARTIALLY_FILLED
	case "2": return FILLED
	case "3": //Done for day
		if approveStatus == 1 {return TIMEOUT_REJECTED}
		if cumQty == 0 {return TIMEOUT_FAIL}
		return COMPLETED
	case "4","j":
		if cumQty > 0 {return PARTIALLY_CANCELED}
		return CANCELED
	case "6": return PENDING_CANCEL
	case "7","C": return TIMEOUT_FAIL
	case "8": return REJECTED
	default: return ordStatus
	}
}

func ConvertOrdSource(ordSource string) string {
	switch strings.TrimSpace(ordSource) {
	case "goats": return "客户端"
	case "FIX": return "FIX"
	case "titans": return "线下"
	default: return ordSource
	}
}
func ConvertSide(side string) string {
	switch side {
	case "1": return "买入"
	case "2": return "卖出"
	default: return side
	}
}

func ConvertOpenClose(openClose string) string {
	switch openClose {
	case "O": return "开仓"
	case "C": return "平仓"
	default: return openClose
	}
}

func ConvertHandlInst(handlInst string) string {
	switch handlInst {
	case "1": return "快速交易"
	case "2": return "策略交易"
	case "3": return "普通交易"
	case "4": return "普通交易"
	default: return handlInst
	}
}
package domain_error

func init() {
	registerErrMsg(map[string]string{
		ORDER_EXECUTOR_ERR_CODE:                           `交易订单执行器异常`,
		ORDER_EXECUTOR_TRADE_CHANNEL_NOT_FOUND_ERR_CODE:   `无法找到匹配订单的交易通道：%v`,
		CANNOT_FIND_ORDER_BY_APP_ORD_ID_ERR_CODE:          `无法基于订单ID(AppOrdID): %s 找到交易订单`,
		ORDER_WITH_STATUS_WHICH_CANNOT_BE_CANCEL_ERR_CODE: `订单: %s 正处于状态 %s，不允许执行撤单操作`,
	})
}

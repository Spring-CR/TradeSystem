package domain_error

func init() {
	registerErrMsg(map[string]string{
		TRADE_RESP_ERR_CODE:                      `在处理交易回报的过程中发生异常`,
		TRADE_RESP_INSERT_TO_DB_ERR_CODE:         `交易回报插入数据库异常，回报报文：%s`,
		TRADE_ACTION_LATEST_RESP_UPDATE_ERR_CODE: `更新最新交易回报数据时发生异常，回报报文：%s`,
		TRADE_ORDER_LATEST_RESP_UPDATE_ERR_CODE:  `使用最新交易回报更新订单状态时发生异常，回报报文：%s`,
	})
}

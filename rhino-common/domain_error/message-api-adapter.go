package domain_error

func init() {
	registerErrMsg(map[string]string{
		API_ADAPTER_ERR_CODE:                  `API接口适配器异常`,
		DISTINGUISH_INGRESS_MSG_ERR_CODE:      `分辨交易消息类型错误，原始消息：%s`,
		CONVERT_NEW_ORDER_SINGLE_MSG_ERR_CODE: `解析新建交易订单的消息错误，原始消息：%s`,
		PROCESS_NEW_ORDER_SINGLE_ERR_CODE:     `返回新建交易订单被拒绝的消息出现错误`,
		CONVERT_TRADE_RESP_ERR_CODE:           `返回交易回报时出现错误`,
		TX_TIME_TOO_EARLY_ERR_CODE:            `订单的时间设置（TranctTime=%s）太早，可能是过期订单`,
		CONVERT_TRADE_ORDER_ERR_CODE:          `交易订单格式转换异常`,
	})
}

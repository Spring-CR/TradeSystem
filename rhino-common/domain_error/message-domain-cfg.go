package domain_error

func init() {
	registerErrMsg(map[string]string{
		DOMAIN_CONFIG_ERR_CODE:           `领域模型配置错误`,
		REQUIRED_CONFIG_MISSING_ERR_CODE: `缺少必需的配置项，配置项名称：%s，配置项描述：%s`,

		UNKNOW_CHANNEL_PROTOCOL_TYPE_ERR_CODE: `不支持的交易通道的协议类型:%v`,
		UNKNOW_CHANNEL_STATUS_ERR_CODE:        `不支持的交易通道的健康状态:%v`,

		CANNOT_INIT_MSG_SEQ_GEN_ERR_CODE: `初始化消息序号生成器失败`,
	})
}

package domain_error

const (
	// 3xxxxxx - 领域模型配置错误
	DOMAIN_CONFIG_ERR_CODE           = `300000`
	REQUIRED_CONFIG_MISSING_ERR_CODE = `300001`

	// 31xxxxx - 不符合的枚举类型
	UNKNOW_CHANNEL_PROTOCOL_TYPE_ERR_CODE = `310000`
	UNKNOW_CHANNEL_STATUS_ERR_CODE        = `310001`

	// 32xxxxx -消息序列生成器相关的错误
	CANNOT_INIT_MSG_SEQ_GEN_ERR_CODE = `320000`
)

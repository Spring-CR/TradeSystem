package domain_error

func init() {
	registerErrMsg(map[string]string{
		TRADE_CHANNEL_ERR_CODE: `交易通道错误`,
		TRADE_CHANNEL_ORDER_VALIDATE_ERR_CODE: `订单校验失败 - %s`,

		FIX_PARSE_CLI_CFG_ERR_CODE:     `解释FIX交易通道的客户端配置文件错误`,
		FIX_INIT_LOG_ERR_CODE:          `FIX交易通道的客户端创建日志文件错误`,
		FIX_INIT_CLI_ERR_CODE:          `FIX交易通道的客户端初始化错误`,
		FIX_START_CLI_ERR_CODE:         `FIX交易通道的客户端启动错误`,
		FIX_SEND_MSG_ERR_CODE:          `FIX交易通道的客户端报文发送错误`,
		FIX_SESSION_NOT_READY_ERR_CODE: `FIX交易通道的FIX会话未就绪`,

		GFGFIX_PARSE_CLI_CFG_ERR_CODE:     `解释GFGFIX交易通道的客户端配置文件错误`,
		GFGFIX_INIT_LOG_ERR_CODE:          `GFGFIX交易通道的客户端创建日志文件错误`,
		GFGFIX_INIT_CLI_ERR_CODE:          `GFGFIX交易通道的客户端初始化错误`,
		GFGFIX_START_CLI_ERR_CODE:         `GFGFIX交易通道的客户端启动错误`,
		GFGFIX_SEND_MSG_ERR_CODE:          `GFGFIX交易通道的客户端报文发送错误`,
		GFGFIX_SESSION_NOT_READY_ERR_CODE: `GFGFIX交易通道的FIX会话未就绪`,

		DATA_CONVERT_CLORDID_EMPTY_ERR_CODE:                   `由客户端指定的订单ID不能为空`,
		DATA_CONVERT_ILLEGAL_HANDINST_ERR_CODE:                `对于交易通道，订单的处理方式（HandlInst）仅能设置为1（无算法的直通交易）、2（柜台原生算法的直通交易）或3（手动交易），当前值，HandlInst=%v`,
		DATA_CONVERT_SYMBOL_EMPTY_ERR_CODE:                    `证券代码（Symbol）不能为空`,
		DATA_CONVERT_ILLEGAL_TRADE_SIDE_ERR_CODE:              `订单的交易方向（Side）仅支持设置为1（买）、2（卖），当前值，Side=%v`,
		DATA_CONVERT_ILLEGAL_ORDER_TYPE_ERR_CODE:              `订单的类型（OrdType）仅支持设置为1（市价单）、2（限价单），当前值，OrdType=%v`,
		DATA_CONVERT_ILLEGAL_CURRENCY_ERR_CODE:                `订单的交易货币类型（Currency）仅支持设置为空，或者%v，当前值，Currency=%v`,
		DATA_CONVERT_ILLEGAL_ID_SOURCE_ERR_CODE:               `订单的证券标识符类型（IDSource）仅支持设置为空，或者%v，当前值，IDSource=%v`,
		DATA_CONVERT_ILLEGAL_QTY_ERR_CODE:                     `订单的标的物数量（OrderQty）或购买金额数量（CashOrderQty）至少应该有一项大于零，当前值，OrderQty=%v，CashOrderQty=%v`,
		DATA_CONVERT_SECURITY_ID_EMPTY_ERR_CODE:               `证券ID（SecurityID）不能为空`,
		DATA_CONVERT_ILLEGAL_OPEN_CLOSE_ERR_CODE:              `订单的开平仓标记（OpenClose）仅支持设置为空，或者%v，当前值，OpenClose=%v`,
		DATA_CONVERT_ORDER_EXPIRED_ERR_CODE:                   `订单已经过期失效，当前时间=%v，过期时间=%v`,
		DATA_CONVERT_TRADE_CHANNEL_UNMATCH_ERR_CODE:           `订单指定的交易通道于当前交易通道不匹配，订单指定的交易通道：%s, 当前交易通道：%s`,
		DATA_CONVERT_SAVE_NEW_ORDER_ERR_CODE:                  `保存新订单到数据库失败`,
		DATA_CONVERT_SET_ORDER_CLIENT_ID_ERR_CODE:             `无法设置订单的客户端ID（ClOrdID），该订单的数据库ID=%v`,
		DATA_CONVERT_INSERT_TRADE_ACTION_LATEST_RESP_ERR_CODE: `无法保存订单（ClOrdID=%v）用于记录最新操作所反馈状态的数据库记录（TradeActionLatestResp）`,
		DATA_CONVERT_APPORDID_EMPTY_ERR_CODE:                  `由应用层指定的订单ID (AppOrdID) 不能为空`,
		DATA_CONVERT_ALG_PARAMS_NOT_PROVIDED_ERR_CODE:         `算法 %v 的计算参数未提供`,
		DATA_CONVERT_ALG_PARAMS_EXTRACT_ERR_CODE:              `算法 %v 的计算参数解析异常，原始输入参数:%v`,
		DATA_CONVERT_ALG_PARAMS_EXTRACT_TIME_FIELD_ERR_CODE:   `无法解析算法中的时间参数，原始输入参数:%v`,
		DATA_CONVERT_UNSUPPORT_EXCHANGE_ERR_CODE:              `不支持的证券交易所：%v`,
		DATA_CONVERT_GET_PROECTED_PRICE_ERR_CODE:              `无法计算保护限价，标的：%s，买卖方向：%s`,
	})
}

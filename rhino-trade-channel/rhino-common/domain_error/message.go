package domain_error

func init() {
	registerErrMsg(map[string]string{
		GENERIC_ERR_CODE:                             `错误发生`,
		GENERIC_WARNING_CODE:                         `%v`,
		API_CREATE_REQUEST_ERR_CODE:                  `无法创建 HTTP 请求对象, api_name:%s, method:%s, path:%s`,
		GENERIC_API_REQUEST_ERR_CODE:                 `无法调用远程API: %v`,
		API_RESPONSE_PARSE_ERR_CODE:                  `无法处理 API:%v 的返回结果, 期望的数据类型:%T`,
		API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE:           `参数 %s 不允许为空`,
		API_PARAM_PARSING_ERR_CODE:                   `解析HTTP参数=%s时发生异常，原始参数值=%s`,
		API_UNAUTHORIZED_ERR_CODE:                    `API鉴权失败`,
		DATABASE_OPERATION_ERR_CODE:                  `发生数据库异常`,
		DATABASE_OPEN_TRANS_ERR_CODE:                 `打开数据库事务时发生异常`,
		DATABASE_COMMIT_TRANS_ERR_CODE:               `提交数据库事务时发生异常`,
		DATABASE_ROLLBACK_TRANS_ERR_CODE:             `回滚数据库事务时发生异常`,
		DATABASE_ROLLBACK_TRANS_COMMIT_FAIL_ERR_CODE: `事务提交失败后的回滚操作异常`,
		DATABASE_QUERY_ERROR:                         `发生数据库查询异常，SQL=%s，参数列表=%v`,
		DATABASE_RECORD_EMPTY_ERR_CODE:               `数据库查询记录为空`,
		ILLEGAL_DATA_FORMAT_ERR_CODE:                 `错误的数据格式，数据: %s`,
		ILLEGAL_DATE_FORMAT_ERR_CODE:                 `错误的日期格式，原始值: %v`,
		KAFKA_ERR_CODE:                               `kafka异常`,
		CANNOT_CREATE_PRODUCER_ERR_CODE:              `无法创建kafka生产者`,
		CANNOT_PUBLISH_TRADE_MSG_ERR_CODE:            `无法发布交易消息`,
		CANNOT_CREATE_CONSUMER_ERR_CODE:              `无法创建kafka消费者`,
	})
}

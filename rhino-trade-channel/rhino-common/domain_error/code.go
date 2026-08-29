package domain_error

const (
	GENERIC_ERR_CODE                             = `000000` // 00xxxx - generic error code
	GENERIC_WARNING_CODE                         = `000001`
	API_CREATE_REQUEST_ERR_CODE                  = `010001` // 01xxxx - http req/resp error code
	GENERIC_API_REQUEST_ERR_CODE                 = `010002`
	API_RESPONSE_PARSE_ERR_CODE                  = `010003`
	API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE           = `010004`
	API_PARAM_PARSING_ERR_CODE                   = `010005`
	API_UNAUTHORIZED_ERR_CODE                    = `010006`
	DATABASE_OPERATION_ERR_CODE                  = `020000` // 02xxxx - database error code
	DATABASE_OPEN_TRANS_ERR_CODE                 = `020001`
	DATABASE_COMMIT_TRANS_ERR_CODE               = `020002`
	DATABASE_ROLLBACK_TRANS_ERR_CODE             = `020003`
	DATABASE_ROLLBACK_TRANS_COMMIT_FAIL_ERR_CODE = `020004`
	DATABASE_QUERY_ERROR                         = `020005`
	DATABASE_RECORD_EMPTY_ERR_CODE               = `020006`
	ILLEGAL_DATA_FORMAT_ERR_CODE                 = `030000` // 03xxxx - data format error code
	ILLEGAL_DATE_FORMAT_ERR_CODE                 = `030001`
	KAFKA_ERR_CODE                               = `040000` // 04xxxx - kafka error code
	CANNOT_CREATE_PRODUCER_ERR_CODE              = `040001`
	CANNOT_PUBLISH_TRADE_MSG_ERR_CODE            = `040002`
	CANNOT_CREATE_CONSUMER_ERR_CODE              = `040003`
)

package constant

type LangType int

const (
	Lang_CN LangType = 0
)

const (
	TradeActionRespGuidField = "system_guid" // 常规成交回报ExecutionReport或者OrderCancelReject里的记录主键
	ReqMsgSeq                = "req_msg_seq" // 对于回报ExecutionReport或者OrderCancelReject，如果其关联的报文是通过stream api发送的，可以用于标记该报文的消息序号
)

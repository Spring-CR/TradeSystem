package api_adapter

import (
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-core/schema"
	"rhino-core/types"

	//导入插件
	_ "rhino-engine/api/api_adapter/plugin"
)

// WorkerSharding: 工作者分片
// rawMsg: 原始消息
// workerCount: 工作者数目
// cumTaskCount: 程序启动至今累计已处理的消息数
// workerIndex: 返回sharding的工作者索引，从0开始编号
type WorkerSharding func(rawMsg []byte) (workerIndex int)
type GetOrderByAppOrdID func(appOrdID string) (order *schema.TradeOrder, ok bool)

type APIAdapter interface {
	// workerCount: 工作者数目
	// rawMsg: 原始消息
	// cumTaskCount: 程序启动至今累计已处理的消息数
	// workerIndex: 返回sharding的工作者索引，从0开始编号
	GetWorkerSharding(workerCount int, getOrderByAppOrdID func(appOrdID string) (order *schema.TradeOrder, ok bool)) (workerSharding func(rawMsg []byte, decodeOrder*schema.TradeOrder, cumTaskCount int) (workerIndex int))
	DistinguishIngressMessage(rawMsg []byte) (enum.ApiMessageType, map[string]interface{}, interface{}, *domain_error.Error)
	ConvertNewOrderSingleMessage(rawMsg []byte, msgProps map[string]interface{}) (*schema.TradeOrder, *domain_error.Error)
	RefineAndValidate(rawOrder *schema.TradeOrder, trade bool) (*domain_error.Error)
	ConvertTradeOrderParams(tradeOrder *schema.TradeOrder) (appOrd interface{}, de *domain_error.Error)
	ProcessNewOrderSingleError(tradeOrder *schema.TradeOrder, de *domain_error.Error) (rejectMsg []byte)
	ConvertOrderCancelRequestMessage(rawMsg []byte, msgProps map[string]interface{}) (*types.ApplicationOrderCancelRequest, *domain_error.Error)
	ProcesssOrderCancelRequestError(ordCxlReq *types.ApplicationOrderCancelRequest, order *schema.TradeOrder, de *domain_error.Error) (rejectMsg []byte)
	// 一定要实现一个system_guid的字段
	ConvertTradeResponseMessage(tradeResp *types.TradeActionRespReturn) ([]byte, *domain_error.Error)
	ErrorCouldBeIgnoreAfterReview(de *domain_error.Error) (ignoreAfterReview bool)
	// ExtractApplicationOrderID(rawMsg []byte, msgProps map[string]interface{})(appOrdID string, de *domain_error.Error)
	}

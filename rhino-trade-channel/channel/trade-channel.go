package channel

// 作者：林春泉

import (
	"rhino-common/domain_error"
	"rhino-core/order_domain/order_cache"
	"rhino-core/schema"
)

type TradeChannelInterface interface {
	// 新建订单时进行客制化校验
	ValidateOrder(order *schema.TradeOrder) (channelOrder interface{}, de *domain_error.Error)
	// 新建交易订单
	NewOrderSingle(order *schema.TradeOrder, afterTradeOrderUpsert func(tradeOrder *schema.TradeOrder), syncTradeOrderAfterError func(tradeOrder *schema.TradeOrder)) (duplicatedOrder bool, de *domain_error.Error)
	// 撤销交易订单
	OrderCancelRequest(actionUser string, actionTime int64, actionKey string, targetClOrdID string, streamInputMsgSeq int64, order *schema.TradeOrder, afterTradeActionLatestRespInsert func(tradeActionLatestResp *schema.TradeActionLatestResp), syncTradeActionLatestRespAfterError func(tradeActionLatestResp *schema.TradeActionLatestResp)) *domain_error.Error
	// 清洗成交回报
	CleanTradeActionResp(tradeActionResp *schema.TradeActionResp) (ignore bool)
	// 加工撤单相关的ID参数
	RefineOrderCancelID(rawTargetClOrdID string, rawClOrdID string, order *schema.TradeOrder) (targetClOrdID, clOrdID string)
	// 监听交易回报
	//AddTradeActionRespListener(onTradeActionResp func(tradeActionResp *schema.TradeActionResp))
	// 重置Channel，一般在归档时间重置
	Reset(force bool) *domain_error.Error
	// 收市时执行
	OnMarketClosed(orderCache *order_cache.OrderCache)
}

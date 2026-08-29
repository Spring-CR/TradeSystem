package order_manager

import (
	"rhino-core/order_domain/order_cache"
	"rhino-core/schema"
)

type DirectTradeOrderManager struct {
}

func NewDirectTradeOrderManager() *DirectTradeOrderManager {
	inst := &DirectTradeOrderManager{}
	return inst
}

//var ordStatusDraft = string(enum.OrdStatus_Draft)
func (m *DirectTradeOrderManager) OrchestrateOrder(order *schema.TradeOrder, orderCache *order_cache.OrderCache) {
	orderCache.AddRootTradeOrder(order)
}

func (m *DirectTradeOrderManager) CancelOrder(tradeActionLatestResp *schema.TradeActionLatestResp, orderCache *order_cache.OrderCache) {
	orderCache.AddTradeActionForDirectOrder(tradeActionLatestResp)
}
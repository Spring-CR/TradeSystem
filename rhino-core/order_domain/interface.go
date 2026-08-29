package order_domain

import (
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-core/order_domain/order_cache"
	"rhino-core/schema"
	"rhino-core/types"

	// 导入插件
	_ "rhino-plugins"
)

// 作者：林春泉

type OrderAcceptorInterface interface {
	AcceptOrderDraft(order *schema.TradeOrder, actionType enum.ActionType) *domain_error.Error
	AcceptOrderDraftDeletion(orderDraftDeletion *types.ApplicationOrderDraftDeleteRequest) *domain_error.Error
	AcceptOrderAttributeUpdateRequest(attrUpdateReq *types.ApplicationOrderAttributeUpdateRequest) *domain_error.Error
	AcceptNewOrderSingleRequest(order *schema.TradeOrder) (duplicatedOrder bool, de *domain_error.Error)
	AcceptOrderCancelRequest(orderCxlReq *types.ApplicationOrderCancelRequest) (*schema.TradeOrder, *domain_error.Error)
}

type OrderExecutorInterface interface {
	NewOrderSingle(order *schema.TradeOrder) (duplicatedOrder bool, de *domain_error.Error)
	OrderCancelRequest(orderCancelRequest *types.ApplicationOrderCancelRequest) (*schema.TradeOrder, *domain_error.Error)
}

type OrderManagerInterface interface {
	OrchestrateOrder(order *schema.TradeOrder, orderCache *order_cache.OrderCache)
	CancelOrder(tradeActionLatestResp *schema.TradeActionLatestResp, orderCache *order_cache.OrderCache)
}

type OrderExecutorAdapter interface {
	AfterNewOrderSingle(order *schema.TradeOrder, duplicatedOrder bool, de *domain_error.Error)
	AfterOrderCancelRequest(orderCancelRequest *types.ApplicationOrderCancelRequest, order *schema.TradeOrder, de *domain_error.Error)
}

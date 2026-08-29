package order_status_adapter

import (
	"rhino-core/schema"
	"rhino-core/types"
)

type OrderStatusAdapter interface {
	LookUpTraceableTradeActionResp(directOrderMap map[string]*types.TraceableTradeOrder, tradeActionResp *schema.TradeActionResp, tradeActionRespMap map[string]*types.TraceableTradeActionResp) (traceableTradeActionResp *types.TraceableTradeActionResp, ok bool)
	UpdateTradeActionLatestResp(tradeActionResp *schema.TradeActionResp, traceableTradeActionResp *types.TraceableTradeActionResp)
	UpdateOrderStatusForOrderCancelReject(tradeActionResp *schema.TradeActionResp, tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder)
	UpdateOrderStatusForExecutionReport(tradeActionResp *schema.TradeActionResp, tradeActionLatestResp *schema.TradeActionLatestResp, tradeActionRespList []*schema.TradeActionResp, order *schema.TradeOrder, traceableTradeOrder *types.TraceableTradeOrder) (orderUpdateAttributes map[string]interface{})
}

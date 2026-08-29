package order_cache

import (
	"rhino-core/schema"
	"rhino-core/types"
)

func defaultLookUpTraceableTradeActionResp(directOrderMap map[string]*types.TraceableTradeOrder, tradeActionResp *schema.TradeActionResp, tradeActionRespMap map[string]*types.TraceableTradeActionResp) (traceableTradeActionResp *types.TraceableTradeActionResp, ok bool) {
	traceableTradeActionResp, ok = tradeActionRespMap[tradeActionResp.ClOrdID]
	return
}

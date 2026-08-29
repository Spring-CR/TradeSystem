package types

import (
	"rhino-core/schema"
)

var defaultUpdateTradeActionLatestResp = func(tradeActionResp *schema.TradeActionResp, traceableTradeActionResp *TraceableTradeActionResp) {
	// 更新latestresp
	traceableTradeActionResp.tradeActionLatestResp.OrderID = tradeActionResp.OrderID
	traceableTradeActionResp.tradeActionLatestResp.OrigClOrdID = tradeActionResp.OrigClOrdID
	traceableTradeActionResp.tradeActionLatestResp.ExecID = tradeActionResp.ExecID
	traceableTradeActionResp.tradeActionLatestResp.ExecType = tradeActionResp.ExecType
	traceableTradeActionResp.tradeActionLatestResp.OrdStatus = tradeActionResp.OrdStatus
	traceableTradeActionResp.tradeActionLatestResp.OrdRejReason = tradeActionResp.OrdRejReason
	traceableTradeActionResp.tradeActionLatestResp.CxlRejResponseTo = tradeActionResp.CxlRejResponseTo
	traceableTradeActionResp.tradeActionLatestResp.ExecRestatementReason = tradeActionResp.ExecRestatementReason

	// 对于master，这部分其实只需要更新一次，或者并不需要更新
	if traceableTradeActionResp.slave && traceableTradeActionResp.tradeActionLatestResp.Symbol == "" {
		traceableTradeActionResp.tradeActionLatestResp.Account = tradeActionResp.Account
		traceableTradeActionResp.tradeActionLatestResp.Symbol = tradeActionResp.Symbol
		traceableTradeActionResp.tradeActionLatestResp.SymbolSfx = tradeActionResp.SymbolSfx
		traceableTradeActionResp.tradeActionLatestResp.SecurityID = tradeActionResp.SecurityID
		traceableTradeActionResp.tradeActionLatestResp.IDSource = tradeActionResp.IDSource
		traceableTradeActionResp.tradeActionLatestResp.SecurityType = tradeActionResp.SecurityType
		traceableTradeActionResp.tradeActionLatestResp.Side = tradeActionResp.Side
		traceableTradeActionResp.tradeActionLatestResp.OpenClose = tradeActionResp.OpenClose
		traceableTradeActionResp.tradeActionLatestResp.OrderQty = tradeActionResp.OrderQty
		traceableTradeActionResp.tradeActionLatestResp.CashOrderQty = tradeActionResp.CashOrderQty
		traceableTradeActionResp.tradeActionLatestResp.OrdType = tradeActionResp.OrdType
		traceableTradeActionResp.tradeActionLatestResp.Price = tradeActionResp.Price
		traceableTradeActionResp.tradeActionLatestResp.Currency = tradeActionResp.Currency
		traceableTradeActionResp.tradeActionLatestResp.EffectiveTime = tradeActionResp.EffectiveTime
		traceableTradeActionResp.tradeActionLatestResp.ExpireTime = tradeActionResp.ExpireTime
	}

	traceableTradeActionResp.tradeActionLatestResp.LastShares = tradeActionResp.LastShares
	traceableTradeActionResp.tradeActionLatestResp.LastPx = tradeActionResp.LastPx
	traceableTradeActionResp.tradeActionLatestResp.LeavesQty = tradeActionResp.LeavesQty
	traceableTradeActionResp.tradeActionLatestResp.CumQty = tradeActionResp.CumQty
	traceableTradeActionResp.tradeActionLatestResp.AvgPx = tradeActionResp.AvgPx
	traceableTradeActionResp.tradeActionLatestResp.TransactTime = tradeActionResp.TransactTime
	traceableTradeActionResp.tradeActionLatestResp.MsgTime = tradeActionResp.MsgTime
	traceableTradeActionResp.tradeActionLatestResp.MsgSeq = tradeActionResp.MsgSeq
}

var defaultUpdateOrderStatusForOrderCancelReject = func(tradeActionResp *schema.TradeActionResp, tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder) {
	//直接从tradeActionResp复制
	order.OrdStatus = tradeActionResp.OrdStatus
}

var defaultUpdateOrderStatusForExecutionReport = func(tradeActionResp *schema.TradeActionResp, tradeActionLatestResp *schema.TradeActionLatestResp, tradeActionRespList []*schema.TradeActionResp, order *schema.TradeOrder, traceableTradeOrder *TraceableTradeOrder) (orderUpdateAttributes map[string]interface{}) {
	//直接从tradeActionResp复制
	order.OrdStatus = tradeActionResp.OrdStatus
	order.LastShares = tradeActionResp.LastShares
	order.LastPx = tradeActionResp.LastPx
	order.LeavesQty = tradeActionResp.LeavesQty
	order.CumQty = tradeActionResp.CumQty
	order.AvgPx = tradeActionResp.AvgPx
	order.OrdRejReason = tradeActionResp.OrdRejReason
	return
}

// func SetFuncUpdateTradeActionLatestResp(_updateTradeActionLatestResp func(tradeActionResp *schema.TradeActionResp, traceableTradeActionResp *TraceableTradeActionResp)) {
// 	updateTradeActionLatestResp = _updateTradeActionLatestResp
// }

// func SetFuncUpdateOrderStatusForOrderCancelReject(_updateOrderStatusForOrderCancelReject func(tradeActionResp *schema.TradeActionResp, tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder)) {
// 	updateOrderStatusForOrderCancelReject = _updateOrderStatusForOrderCancelReject
// }

// func SetUpdateOrderStatusForExecutionReport(_updateOrderStatusForExecutionReport func(tradeActionResp *schema.TradeActionResp, tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder)) {
// 	updateOrderStatusForExecutionReport = _updateOrderStatusForExecutionReport
// }

package types

import "rhino-core/schema"

func ExtractTraceableTradeOrder(traceableTradeOrder *TraceableTradeOrder) (tradeOrders []*schema.TradeOrder, tradeActionLatestResps []*schema.TradeActionLatestResp, tradeActionResps []*schema.TradeActionResp) {

	tradeOrder := traceableTradeOrder.GetBasicInfo()
	tradeOrders = append(tradeOrders, tradeOrder)

	tradeActions := traceableTradeOrder.GetTradeActions()
	for _, tradeAction := range tradeActions {
		_tradeActionResps := tradeAction.GetTradeActionRespList()
		_tradeActionLatestResp := tradeAction.GetTradeActionLatestResp()
		_tradeActionLatestResp.Account = tradeOrder.Account
		_tradeActionLatestResp.Symbol = tradeOrder.Symbol
		_tradeActionLatestResp.SymbolSfx = tradeOrder.SymbolSfx
		_tradeActionLatestResp.SecurityID = tradeOrder.SecurityID
		_tradeActionLatestResp.IDSource = tradeOrder.IDSource
		_tradeActionLatestResp.SecurityType = tradeOrder.SecurityType
		_tradeActionLatestResp.Side = tradeOrder.Side
		_tradeActionLatestResp.OpenClose = tradeOrder.OpenClose
		_tradeActionLatestResp.OrderQty = tradeOrder.OrderQty
		_tradeActionLatestResp.CashOrderQty = tradeOrder.CashOrderQty
		_tradeActionLatestResp.OrdType = tradeOrder.OrdType
		_tradeActionLatestResp.Price = tradeOrder.Price
		_tradeActionLatestResp.Currency = tradeOrder.Currency
		// _tradeActionLatestResp.EffectiveTime = tradeOrder.EffectiveTime
		// _tradeActionLatestResp.ExpireTime = tradeOrder.ExpireTime

		tradeActionLatestResps = append(tradeActionLatestResps, _tradeActionLatestResp)
		tradeActionResps = append(tradeActionResps, _tradeActionResps...)
	}

	subOrders := traceableTradeOrder.GetSubOrders()
	if len(subOrders) == 0 {
		return
	} else {
		for _, subsubOrder := range subOrders {
			tradeOrders2, tradeActionLatestResps2, tradeActionResps2 := ExtractTraceableTradeOrder(subsubOrder)
			tradeOrders = append(tradeOrders, tradeOrders2...)
			tradeActionLatestResps = append(tradeActionLatestResps, tradeActionLatestResps2...)
			tradeActionResps = append(tradeActionResps, tradeActionResps2...)
		}
	}

	return
}
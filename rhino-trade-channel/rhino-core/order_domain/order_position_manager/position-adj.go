package order_position_manager

import (
	"log"
	"rhino-common/domain_error"
	"rhino-core/schema"
	"rhino-core/types"
)

func (pm *PositionManager) PreparePositionAdjustmentParams(tradeOrder *schema.TradeOrder) (mockTradeOrder *schema.TradeOrder, mockTradeActionResp *schema.TradeActionResp, de *domain_error.Error) {
	return pm.positionAdapter.PreparePositionAdjustmentParams(tradeOrder)
}

func (pm *PositionManager) AdjustPosition(mockTradeOrder *schema.TradeOrder, mockTradeActionResp *schema.TradeActionResp) {

	key := pm.positionAdapter.GetPositionCalculatorKey(mockTradeOrder)
	log.Printf("AdjustPosition, tradeOrder.AppOrdID=%s, key=%v\n", mockTradeOrder.AppOrdID, key)
	pm.orderLog.Printf(mockTradeOrder, nil, "[AdjustPosition] key=%v, symbol2=%v, qty=%v, price=%v, side=%v", key, mockTradeOrder.ExtendAttrMap["symbol2"], mockTradeOrder.OrderQty, mockTradeOrder.Price, mockTradeOrder.Side)

	pm.lock.RLock()
	pc, ok := pm.calculatorMap[key]
	pm.lock.RUnlock()

	if !ok {
		pm.lock.Lock()
		pc, ok = pm.calculatorMap[key]
		if !ok {
			pc = NewPositionCalculator(key, mockTradeOrder, pm.positionAdapter, pm.orderLog, pm.persisFunc)
			pm.calculatorMap[key] = pc
		}
		pm.lock.Unlock()
	}
	// 以上步骤确保了PositionCalculator存在（没有则创建）
	tradeActionRespReturn := types.NewTradeActionRespReturn(types.NewTraceableTradeActionResp(types.NewTraceableTradeOrder(mockTradeOrder), mockTradeOrder, nil, false), mockTradeActionResp)
	pm.updatePositionByTradeResp(tradeActionRespReturn)

	pm.lock.Lock()
	items := pm.mockTradeOrders[key]
	items = append(items, mockTradeOrder)
	pm.mockTradeOrders[key] = items
	pm.lock.Unlock()
}

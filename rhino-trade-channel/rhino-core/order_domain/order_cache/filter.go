package order_cache

import (
	"log"
	"rhino-core/types"
)

// 根据传入的条件函数，对OrderCache的母单进行过滤。分别得到匹配条件的订单列表和不匹配条件的订单列表。
func (c *OrderCache) FilterOrderByFunction(isOrderMatch func(order *types.TraceableTradeOrder) bool) (matchOrders, unmatchOrders []*types.TraceableTradeOrder) {
	c.rootOrdersLock.RLock()
	defer c.rootOrdersLock.RUnlock()

	if isOrderMatch == nil {
		matchOrders = append(matchOrders, c.rootOrders...)

		return
	}


	log.Printf("start to FilterOrderByFunction, rootOrders.Len=%d\n", len(c.rootOrders))

	for i, traceableOrder := range c.rootOrders {
		log.Printf("======> FilterOrderByFunction #%d: %s\n", i, traceableOrder.GetBasicInfo().AppOrdID)
		if isOrderMatch(traceableOrder) {
			matchOrders = append(matchOrders, traceableOrder)
		} else {
			unmatchOrders = append(unmatchOrders, traceableOrder)
		}
	}

	return
}

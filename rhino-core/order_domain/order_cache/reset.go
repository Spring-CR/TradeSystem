package order_cache

import (
	"log"
	"rhino-core/schema"
	"rhino-core/types"
	"runtime/debug"
)

func (c *OrderCache) ResetForMaster(tradeOrdersToArchive []*schema.TradeOrder, tradeActionLatestRespsToArchive []*schema.TradeActionLatestResp) {
	log.Printf("ResetForMaster...")
	appOrdIDMap := make(map[string]bool)
	clOrdIDMap := make(map[string]bool)
	composeOrdIDMap := make(map[string]bool)
	actionKeyMap := make(map[string]bool)

	for _, tradeOrderToArchive := range tradeOrdersToArchive {
		appOrdID := tradeOrderToArchive.AppOrdID
		if len(appOrdID) > 0 {
			appOrdIDMap[appOrdID] = true
		}
		clOrdID := tradeOrderToArchive.ClOrdID
		if len(clOrdID) > 0 {
			clOrdIDMap[clOrdID] = true
		}
		composeOrdID := appOrdID + clOrdID
		if len(composeOrdID) > 0 {
			composeOrdIDMap[composeOrdID] = true
		}
	}

	for _, tradeActionLatestRespToArchive := range tradeActionLatestRespsToArchive {
		clOrdID := tradeActionLatestRespToArchive.ClOrdID
		if len(clOrdID) > 0 {
			clOrdIDMap[clOrdID] = true
		}
		actionKey := tradeActionLatestRespToArchive.ActionKey
		if len(actionKey) > 0 {
			actionKeyMap[actionKey] = true
		}
	}

	var newRootOrders []*types.TraceableTradeOrder
	for i, rootOrder := range c.rootOrders {
		if rootOrder != nil && !appOrdIDMap[rootOrder.GetBasicInfo().AppOrdID] && !clOrdIDMap[rootOrder.GetBasicInfo().ClOrdID] {
			newRootOrders = append(newRootOrders, rootOrder)
		} else {
			if rootOrder != nil {
				rootOrder.Dispose()
			}
			c.rootOrders[i] = nil
		}
	}
	c.rootOrders = newRootOrders

	for k, v := range c.directOrderMap {
		if appOrdIDMap[k] {
			if v != nil {
				v.Dispose()
			}
			delete(c.directOrderMap, k)
		}
	}

	for k, v := range c.tradeActionRespMap {
		if clOrdIDMap[k] {
			if v != nil {
				v.Dispose()
			}
			delete(c.tradeActionRespMap, k)
		}
	}

	if c.tradeActionLatestRespKeyMap != nil {
		for k := range c.tradeActionLatestRespKeyMap {
			if actionKeyMap[k] {
				delete(c.tradeActionLatestRespKeyMap, k)
			}
		}
	}

	if c.rootOrderKeyMap != nil {
		for k := range c.rootOrderKeyMap {
			if composeOrdIDMap[k] {
				delete(c.rootOrderKeyMap, k)
			}
		}
	}

	// 手动触发 GC 并输出内存状态
	debug.FreeOSMemory()

	if c.IsMaster() {
		// 发送reset同步消息
		jsData, _ := newResetMessage()
		c.producer.SendMessage(jsData)

		if c._afterReset != nil {
			log.Println("invoke function after OrderCache reset")
			c._afterReset()
		}
	}
}

func (c *OrderCache) reset() {

	log.Println("start to reset OrderCache state")

	for i, rootOrder := range c.rootOrders {
		if rootOrder != nil {
			rootOrder.Dispose()
		}
		c.rootOrders[i] = nil
	}
	c.rootOrders = nil

	for k, v := range c.directOrderMap {
		if v != nil {
			v.Dispose()
		}
		delete(c.directOrderMap, k)
	}

	for k, v := range c.tradeActionRespMap {
		if v != nil {
			v.Dispose()
		}
		delete(c.tradeActionRespMap, k)
	}

	if c.tradeActionLatestRespKeyMap != nil {
		for k := range c.tradeActionLatestRespKeyMap {
			delete(c.tradeActionLatestRespKeyMap, k)
		}
	}

	if c.rootOrderKeyMap != nil {
		for k := range c.rootOrderKeyMap {
			delete(c.rootOrderKeyMap, k)
		}
	}

	// 手动触发 GC 并输出内存状态
	debug.FreeOSMemory()

	if c.IsSlave() {
		c.recoverForSlave()
	} else {
		c.strictRecover()
	}

	if c._afterReset != nil {
		log.Println("invoke function after OrderCache reset")
		c._afterReset()
	}
}


func (c *OrderCache) SetAfterResetFunc(_afterReset func()) {
	c._afterReset = _afterReset
}
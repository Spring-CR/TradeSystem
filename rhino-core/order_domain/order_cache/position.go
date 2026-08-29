package order_cache

import (
	"log"
	"rhino-common/enum"
	"rhino-core/schema"
)

var (
	submitStatus = string(enum.OrdStatus_Submit)
)

func (c *OrderCache) initPositionForNewOrder(order *schema.TradeOrder) {
	if c.positionManager == nil || order.OrdStatus != submitStatus {
		return
	}
	log.Printf("initPositionForNewOrder, appOrdID:%v, clOrdID:%v, status:%v\n", order.AppOrdID, order.ClOrdID, order.OrdStatus)
	//c.positionCalculator.AcquireOrderQuota(true, order)
	c.positionManager.FreezeQuota(true, order)
}

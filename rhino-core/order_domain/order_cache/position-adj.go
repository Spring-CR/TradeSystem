package order_cache

import "rhino-core/schema"

func (c *OrderCache) AfterAdjustPosition(mockTradeOrder *schema.TradeOrder, mockTradeActionResp *schema.TradeActionResp) {
	if c._afterAdjustPosition != nil {
		c._afterAdjustPosition(mockTradeOrder, mockTradeActionResp)
		return
	}
	if !c.IsMaster() {
		return
	}
	// 只有master才发送同步消息
	jsData, _ := newOrderCacheSyncMessageForAdjustPosition(mockTradeOrder, mockTradeActionResp)
	c.producer.SendMessage(jsData)
}

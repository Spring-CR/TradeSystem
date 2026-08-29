package order_cache

import "log"

// 由外部确保syncMsgs的顺序
func (c *OrderCache) recoverFromSyncMessages(syncMsgs []*OrderCacheSyncMessage) {

	// 对syncMsgs进行排序
	// kafka本身就是有序的，不要进行sorting处理！
	// sortOrderCacheSyncMessages(syncMsgs)

	log.Printf("start to process recover syncMsgses, msg Len=%d\n", len(syncMsgs))
	for i, syncMsg := range syncMsgs {
		//jsData, _ := json.MarshalIndent(syncMsg, "", "  ")
		//log.Printf("===> syncMsg:%s\n", jsData)
		log.Printf("===> syncMsg#%d: msgType:%v\n", i, syncMsg.MessageType)
		if syncMsg.MessageType == OrderCacheSyncMessageType_UpdateTradeOrderAttributes {
			jsData, _ := json.Marshal(syncMsg)
			log.Printf("skip UpdateTradeOrderAttributes syncMsg:%s\n", jsData)
			continue
		}
		c.updateByOrderCacheSyncMessage(syncMsg)
	}
}

func (c *OrderCache) updateByOrderCacheSyncMessage(syncMsg *OrderCacheSyncMessage) {

	// 标记为已发送
	switch syncMsg.MessageType {
	case OrderCacheSyncMessageType_AddRootTradeOrder:
		c.AddRootTradeOrder(syncMsg.TradeOrder)
	case OrderCacheSyncMessageType_UpdateByTradeActionResp:
		c.UpdateByTradeActionResp(syncMsg.TradeActionResp)
	case OrderCacheSyncMessageType_AddTradeActionForDirectOrder:
		c.AddTradeActionForDirectOrder(syncMsg.TradeActionLatestResp)
	case OrderCacheSyncMessageType_UpdateTradeOrderDraft:
		c.UpdateTradeOrderDraft(nil, syncMsg.TradeActionLatestResp, syncMsg.TradeOrder)
	case OrderCacheSyncMessageType_DeleteTradeOrderDraft:
		c.DeleteTradeOrderDraft(nil, syncMsg.TradeActionLatestResp, syncMsg.TradeOrder)
	case OrderCacheSyncMessageType_UpdateTradeOrderAttributes:
		c.UpdateTradeOrderAttributes(syncMsg.AppOrdID, syncMsg.UpdateAttributes)
	case OrderCacheSyncMessageType_Reset:
		c.reset()
	case OrderCacheSyncMessageType_SyncOrder:
		c.SyncOrder(syncMsg.TradeOrder)
	case OrderCacheSyncMessageType_SyncTradeActionLatestResp:
		c.SyncTradeActionLatestResp(syncMsg.TradeActionLatestResp)
	}
}

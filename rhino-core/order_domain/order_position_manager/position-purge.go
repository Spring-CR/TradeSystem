package order_position_manager

import (
	"log"
	"rhino-core/schema"
)

// 人工虚拟调仓 以 constant.MockPositionOrdIDPrefix (pos_adj_) 开头的订单也是要删除的，该函数又purging task调用
func (pm *PositionManager) PurgePositionByPurgingTask(tradeOrdersToKeep []*schema.TradeOrder, tradeActionLatestRespsToKeep []*schema.TradeActionLatestResp, tradeActionRespsToKeep []*schema.TradeActionResp, tradeOrdersToArchive []*schema.TradeOrder, tradeActionLatestRespsToArchive []*schema.TradeActionLatestResp, tradeActionRespsToArchive []*schema.TradeActionResp, purgingLog *schema.DataPurgingLog) {

	if purgingLog == nil {
		log.Println("purgingLog is nil for PurgePositionByPurgingTask")
		return
	}

	pm.lock.Lock()
	defer pm.lock.Unlock()

	log.Printf("start to PurgePositionByPurgingTask, purgingLog.TaskName:%v\n", purgingLog.TaskName)

	positionCalculatorKeysToPurge := pm.positionAdapter.GetPositionCalculatorKeysForPurgingTask(tradeOrdersToKeep, tradeActionLatestRespsToKeep, tradeActionRespsToKeep, tradeOrdersToArchive, tradeActionLatestRespsToArchive, tradeActionRespsToArchive, pm.mockTradeOrders, purgingLog)
	for _, key := range positionCalculatorKeysToPurge {

		c, ok := pm.calculatorMap[key]
		if ok {

			log.Printf("======> purge postion record for key:%v\n", key)

			if pm.persisFunc != nil {
				positionData := pm.positionAdapter.GeneralizePositionRecord(c.metadata, false)
				pm.persisFunc(&PositionChangeEvent{InsertOrUpdate: 2, PositionData: positionData})
			}
		}

		delete(pm.calculatorMap, key)
		delete(pm.mockTradeOrders, key)
	}
}

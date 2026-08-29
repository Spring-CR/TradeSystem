package ficc_fut

import (
	"rhino-core/domain_cfg"
	"rhino-core/order_domain/order_position_manager"
	"rhino-core/schema"
	"rhino-data/datamap"
)

func (a *FiccFutOrderPositionAdapter) GetPositionCalculatorKeysForPurgingTask(tradeOrdersToKeep []*schema.TradeOrder, tradeActionLatestRespsToKeep []*schema.TradeActionLatestResp, tradeActionRespsToKeep []*schema.TradeActionResp, tradeOrdersToArchive []*schema.TradeOrder, tradeActionLatestRespsToArchive []*schema.TradeActionLatestResp, tradeActionRespsToArchive []*schema.TradeActionResp, mockTradeOrders map[string][]*schema.TradeOrder, purgingLog *schema.DataPurgingLog) (positionCalculatorKeysToPurge []string) {

	keyMap := make(map[string]bool)

	for _, tradeOrder := range tradeOrdersToArchive {
		key := a.GetPositionCalculatorKey(tradeOrder)
		if keyMap[key] {
			continue
		}
		positionCalculatorKeysToPurge = append(positionCalculatorKeysToPurge, key)
		keyMap[key] = true
	}

	for k, v := range mockTradeOrders {

		if keyMap[k] {
			continue
		}

		if len(v) == 0 {
			continue
		}

		key := ""

		switch purgingLog.TaskName {

		case "dms":

			if v[0].ChannelCode == "olts-fut" {
				key = k
			}

		case "ovs":

			if v[0].ChannelCode == "stars-fut" {
				key = k
			}
		}

		if key == "" {
			continue
		}

		positionCalculatorKeysToPurge = append(positionCalculatorKeysToPurge, key)
		keyMap[key] = true
	}

	return
}

func (a *FiccFutOrderPositionAdapter) ReloadPositionRecordsForPurgingTask(applicationConfig *domain_cfg.ApplicationCfg, purgingLog *schema.DataPurgingLog) (params []*order_position_manager.PositionCalculatorConstructParam) {

	switch purgingLog.TaskName {

	case "dms":
		params = a.loadInitPositionRecords(applicationConfig, "PositionBaseDms")
	case "ovs":
		params = a.loadInitPositionRecords(applicationConfig, "PositionBaseOvs")
	}

	return

	//改成触发一次自动上场吧
}

func (a *FiccFutOrderPositionAdapter) ProcessPositionBaseDataSyncEvent(event *datamap.DataChangeEvent) {
	// keyList := event.AddKeys
	// keyList = append(keyList, event.ChgKeys...)
	// for _, key := range keyList {
	// 	valList, ok, _ := a.autoSyncRepo.Get(event.TableAlias, key)
	// 	if !ok || len(valList) == 0 {
	// 		continue
	// 	}
	// 	positionData := valList[0]
	// 	positionRecordBase := a.parsePositionRecordFromReposRecord("", positionData)
	// 	positionRecordCurr := a.
	// }
}

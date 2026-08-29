package order_purge

import (
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-core/schema"
	"rhino-core/store/admin_store"
)

// 第5步：重置内存模型
func (a*OrderPurger)resetMem(tradeOrdersToArchive []*schema.TradeOrder, tradeActionLatestRespsToArchive []*schema.TradeActionLatestResp, purgingLog *schema.DataPurgingLog) (err error){
	
	a.orderCache.ResetForMaster(tradeOrdersToArchive, tradeActionLatestRespsToArchive)

	purgingLog.CompletePhase = int(enum.DataPurgingLogPhase_ResetMem)
	err = admin_store.UpdateDataPurgingLogById(a.applicationCfg.GetCentralDB(), purgingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::resetKafka")
	}

	return
}
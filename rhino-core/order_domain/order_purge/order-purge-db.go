package order_purge

import (
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-core/schema"
	"rhino-core/store/admin_store"
)

func (a *OrderPurger) copyTableStructureTemporary(originalTableName string, purgingLog *schema.DataPurgingLog) (string, error) {
	targetTableName := fmt.Sprintf(originalTableName+"_%s_%s", "tmp", purgingLog.PurgingDate)
	//return targetTableName, dbutil.CopyTableStructure(a.applicationCfg.GetAppDB(), a.applicationCfg.GetCentralDB(), originalTableName, targetTableName, true)
	return targetTableName, dbutil.CopyTableStructure(a.applicationCfg.GetAppDB(), a.applicationCfg.GetAppDB(), originalTableName, targetTableName, true)
}

func (a *OrderPurger) createTemporaryTables(purgingLog *schema.DataPurgingLog) (tmpGroupTradeOrderTableName, tmpTradeActionLatestRespTableName, tmpTradeActionRespTableName, tmpTradeOrderTableName string, err error) {

	// 进行表的复制（在这个过程中，如果目标表存在，则会被清理并重新创建）
	tmpGroupTradeOrderTableName, err = a.copyTableStructureTemporary("group_trade_orders", purgingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::createTemporaryTables")
		return
	}
	tmpTradeActionLatestRespTableName, err = a.copyTableStructureTemporary("trade_action_latest_resps", purgingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::createTemporaryTables")
		return
	}
	tmpTradeActionRespTableName, err = a.copyTableStructureTemporary("trade_action_resps", purgingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::createTemporaryTables")
		return
	}
	tmpTradeOrderTableName, err = a.copyTableStructureTemporary("trade_orders", purgingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::createTemporaryTables")
		return
	}

	return
}

func (a *OrderPurger) dumpDataToTemporaryTables(tmpGroupTradeOrderTableName, tmpTradeActionLatestRespTableName, tmpTradeActionRespTableName, tmpTradeOrderTableName string, groupTradeOrders []*schema.GroupTradeOrder, tradeActionLatestResps []*schema.TradeActionLatestResp, tradeActionResps []*schema.TradeActionResp, tradeOrders []*schema.TradeOrder) (de *domain_error.Error) {
	de = a.orderArchiver.DumpTradeActionLatestResps(a.applicationCfg.GetAppDB(), tmpTradeActionLatestRespTableName, tradeActionLatestResps)
	if de != nil {
		return
	}
	de = a.orderArchiver.DumpTradeActionResps(a.applicationCfg.GetAppDB(), tmpTradeActionLatestRespTableName, tradeActionResps)
	if de != nil {
		return
	}
	de = a.orderArchiver.DumpTradeOrders(a.applicationCfg.GetAppDB(), tmpTradeOrderTableName, tradeOrders)
	if de != nil {
		return
	}
	return
}

// 第3步：分析内存订单数据，拆分待清除和继续保留的数据
func (a *OrderPurger) resetDB(tradeOrders []*schema.TradeOrder, tradeActionLatestResps []*schema.TradeActionLatestResp, tradeActionResps []*schema.TradeActionResp, purgingLog *schema.DataPurgingLog) (err error) {

	if purgingLog.CompletePhase >= int(enum.DataPurgingLogPhase_ResetDB) {
		return
	}

	tmpGroupTradeOrderTableName, tmpTradeActionLatestRespTableName, tmpTradeActionRespTableName, tmpTradeOrderTableName, err1 := a.createTemporaryTables(purgingLog)
	if err1 != nil {
		return err1
	}

	log.Printf("created tmp tables:%s, %s, %s, %s\n", tmpGroupTradeOrderTableName, tmpTradeActionLatestRespTableName, tmpTradeActionRespTableName, tmpTradeOrderTableName)

	de := a.dumpDataToTemporaryTables(tmpGroupTradeOrderTableName, tmpTradeActionLatestRespTableName, tmpTradeActionRespTableName, tmpTradeOrderTableName, nil, tradeActionLatestResps, tradeActionResps, tradeOrders)
	if de != nil {
		err = de.ToSimpleError()
		return
	}

	// 替换表
	err = dbutil.MoveTable(a.applicationCfg.GetAppDB(), tmpGroupTradeOrderTableName, "group_trade_orders", purgingLog.PurgingDate)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::resetDB")
		return
	}
	err = dbutil.MoveTable(a.applicationCfg.GetAppDB(), tmpTradeActionLatestRespTableName, "trade_action_latest_resps", purgingLog.PurgingDate)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::resetDB")
		return
	}
	err = dbutil.MoveTable(a.applicationCfg.GetAppDB(), tmpTradeActionRespTableName, "trade_action_resps", purgingLog.PurgingDate)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::resetDB")
		return
	}
	err = dbutil.MoveTable(a.applicationCfg.GetAppDB(), tmpTradeOrderTableName, "trade_orders", purgingLog.PurgingDate)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::resetDB")
		return
	}

	// Truncate FIX消息表
	err = dbutil.TruncateTableIfExist(a.applicationCfg.GetAppDB(), "util_fix_messages")
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::resetDB")
		return
	}


	purgingLog.CompletePhase = int(enum.DataPurgingLogPhase_ResetDB)
	err = admin_store.UpdateDataPurgingLogById(a.applicationCfg.GetCentralDB(), purgingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::resetDB")
		return
	}

	return
}

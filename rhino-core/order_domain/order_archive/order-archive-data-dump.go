package order_archive

import (
	"database/sql"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/utils/bean"
	"rhino-common/utils/dbutil"
	"rhino-core/schema"
	"rhino-core/store/app_store"
	"rhino-core/types"
	"sort"
)

func (a *OrderArchiver) dumpDataToDailyTables(
	dailyGroupTradeOrderTableName,
	dailyTradeActionLatestRespTableName,
	dailyTradeActionRespTableName,
	dailyTradeOrderTableName, extendDailyTradeOrderTableName, extendDailyTradeActionRespTableName string,
	archivingLog *schema.DataArchivingLog, ordersToArchive []*types.TraceableTradeOrder) (orderMap map[string]*types.TraceableTradeOrder, tradeOrders []*schema.TradeOrder, tradeActionLatestResps []*schema.TradeActionLatestResp, tradeActionResps []*schema.TradeActionResp, de *domain_error.Error) {

	orderMap, tradeOrders, tradeActionLatestResps, tradeActionResps = a.ExtractSchemaData(ordersToArchive)

	log.Printf("extractSchemaData for archiving, tradeOrders.Len=%d, tradeActionLatestResps.Len=%d, tradeActionResps.Len=%d\n", len(tradeOrders), len(tradeActionLatestResps), len(tradeActionResps))

	log.Printf("check tradeActionLatestResps of the extractSchemaData result:")
	for _, tradeActionLatestResp := range tradeActionLatestResps {
		log.Printf("===> tradeActionLatestResp, appOrdID:%s, clOrdID:%s, actionType:%v\n", tradeActionLatestResp.AppOrdID, tradeActionLatestResp.ClOrdID, tradeActionLatestResp.ActionType)
	}

	de = a.dumpTradeOrders(a.applicationCfg.GetCentralDB(), dailyTradeOrderTableName, tradeOrders)
	if de != nil {
		return
	}

	de = a.dumpTradeActionLatestResps(a.applicationCfg.GetCentralDB(), dailyTradeActionLatestRespTableName, tradeActionLatestResps)
	if de != nil {
		return
	}

	de = a.dumpTradeActionResps(a.applicationCfg.GetCentralDB(), dailyTradeActionRespTableName, tradeActionResps)
	if de != nil {
		return
	}

	de = a.dumpExtendTradeOrders(extendDailyTradeOrderTableName, tradeOrders)
	if de != nil {
		return
	}

	de = a.dumpExtendTradeActionResps(orderMap, extendDailyTradeActionRespTableName, tradeActionResps)

	return
}

func (a *OrderArchiver) dumpDataToHistoricalTables(
	historicalGroupTradeOrderTableName,
	historicalTradeActionLatestRespTableName,
	historicalTradeActionRespTableName,
	historicalTradeOrderTableName,
	extendHistoricalTradeOrderTableName,
	extendHistoricalTradeActionRespTableName string,
	archivingLog *schema.DataArchivingLog, orderMap map[string]*types.TraceableTradeOrder, tradeOrders []*schema.TradeOrder, tradeActionLatestResps []*schema.TradeActionLatestResp, tradeActionResps []*schema.TradeActionResp) (de *domain_error.Error) {

	de = a.dumpTradeOrders(a.applicationCfg.GetCentralDB(), historicalTradeOrderTableName, tradeOrders, archivingLog.ArchivingDate)
	if de != nil {
		return
	}

	de = a.dumpTradeActionLatestResps(a.applicationCfg.GetCentralDB(), historicalTradeActionLatestRespTableName, tradeActionLatestResps, archivingLog.ArchivingDate)
	if de != nil {
		return
	}

	de = a.dumpTradeActionResps(a.applicationCfg.GetCentralDB(), historicalTradeActionRespTableName, tradeActionResps, archivingLog.ArchivingDate)
	if de != nil {
		return
	}

	de = a.dumpExtendTradeOrders(extendHistoricalTradeOrderTableName, tradeOrders, archivingLog.ArchivingDate)
	if de != nil {
		return
	}

	de = a.dumpExtendTradeActionResps(orderMap, extendHistoricalTradeActionRespTableName, tradeActionResps, archivingLog.ArchivingDate)

	return
}

func (a *OrderArchiver) ExtractSchemaData(ordersToArchive []*types.TraceableTradeOrder) (orderMap map[string]*types.TraceableTradeOrder, tradeOrders []*schema.TradeOrder, tradeActionLatestResps []*schema.TradeActionLatestResp, tradeActionResps []*schema.TradeActionResp) {

	orderMap = make(map[string]*types.TraceableTradeOrder)

	for _, orderToArchive := range ordersToArchive {

		order := orderToArchive.GetBasicInfo()
		orderMap[order.AppOrdID] = orderToArchive

		tradeOrders2, tradeActionLatestResps2, tradeActionResps2 := types.ExtractTraceableTradeOrder(orderToArchive)

		//log.Printf("extractTraceableTradeOrder, AppOrdID=%s, tradeOrders2.Len=%d, tradeActionLatestResps2.Len=%d, tradeActionResps2.Len=%d\n", orderToArchive.GetBasicInfo().AppOrdID, len(tradeOrders2), len(tradeActionLatestResps2), len(tradeActionResps2))

		tradeOrders = append(tradeOrders, tradeOrders2...)
		tradeActionLatestResps = append(tradeActionLatestResps, tradeActionLatestResps2...)
		tradeActionResps = append(tradeActionResps, tradeActionResps2...)

		//log.Printf("===> accumulate counter , tradeOrders.Len=%d, tradeActionLatestResps.Len=%d, tradeActionResps.Len=%d\n", len(tradeOrders), len(tradeActionLatestResps), len(tradeActionResps))
	}

	// 排序
	sort.Slice(tradeOrders, func(i, j int) bool {
		return tradeOrders[i].DBInsertTime < tradeOrders[j].DBInsertTime
	})
	sort.Slice(tradeActionLatestResps, func(i, j int) bool {
		return tradeActionLatestResps[i].ActionTime < tradeActionLatestResps[j].ActionTime
	})
	sort.Slice(tradeActionResps, func(i, j int) bool {
		return tradeActionResps[i].MsgSeq < tradeActionResps[j].MsgSeq
	})

	return
}

var batchSize = 100000

func (a *OrderArchiver) dumpTradeOrders(db *sql.DB, tableName string, records []*schema.TradeOrder, archivingDate ...string) (de *domain_error.Error) {
	n := len(records)
	if n == 0 {
		return
	}
	log.Printf("===> dump %d trade order records to table %s\n", n, tableName)
	idxPairs := splitArray(n, batchSize)
	insertRecordStmt := app_store.GetArchiveInsertTradeOrderStmt(tableName)
	for _, idxPair := range idxPairs {
		//tx, de := dbutil.BeginTx(a.applicationCfg.GetCentralDB())
		tx, de := dbutil.BeginTx(db)
		if de != nil {
			return de
		}
		subRecords := records[idxPair[0]:idxPair[1]]
		for _, record := range subRecords {
			err := app_store.ArchiveInsertTradeOrder(tx, insertRecordStmt, record, archivingDate...)
			if dbutil.IsMysqlDuplicateEntryError(err) {
				err = nil
			}
			if err != nil {
				domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in dumpTradeOrders")
				de = domain_error.Build(domain_error.DATABASE_OPERATION_ERR_CODE, err)
				dbutil.RollbackTx(tx)
				return de
			}
		}
		de = dbutil.CommitTx(tx)
		if de != nil {
			return de
		}
	}
	return
}

func (a *OrderArchiver) dumpTradeActionLatestResps(db *sql.DB, tableName string, records []*schema.TradeActionLatestResp, archivingDate ...string) (de *domain_error.Error) {
	n := len(records)
	if n == 0 {
		return
	}
	log.Printf("===> dump %d trade action latest resp records to table %s\n", n, tableName)
	idxPairs := splitArray(n, batchSize)
	insertRecordStmt := app_store.GetArchiveInsertTradeActionLatestRespStmt(tableName)
	for _, idxPair := range idxPairs {
		//tx, de := dbutil.BeginTx(a.applicationCfg.GetCentralDB())
		tx, de := dbutil.BeginTx(db)
		if de != nil {
			return de
		}
		subRecords := records[idxPair[0]:idxPair[1]]
		for _, record := range subRecords {
			err := app_store.ArchiveInsertTradeActionLatestResp(tx, insertRecordStmt, record, archivingDate...)
			if dbutil.IsMysqlDuplicateEntryError(err) {
				err = nil
			}
			if err != nil {
				domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in dumpTradeActionLatestResps")
				de = domain_error.Build(domain_error.DATABASE_OPERATION_ERR_CODE, err)
				dbutil.RollbackTx(tx)
				return de
			}
		}
		de = dbutil.CommitTx(tx)
		if de != nil {
			return de
		}
	}
	return
}

func (a *OrderArchiver) dumpTradeActionResps(db *sql.DB, tableName string, records []*schema.TradeActionResp, archivingDate ...string) (de *domain_error.Error) {
	n := len(records)
	if n == 0 {
		return
	}
	log.Printf("===> dump %d trade action resp records to table %s\n", n, tableName)
	idxPairs := splitArray(n, batchSize)
	insertRecordStmt := app_store.GetArchiveInsertTradeActionRespStmt(tableName)
	for _, idxPair := range idxPairs {
		//tx, de := dbutil.BeginTx(a.applicationCfg.GetCentralDB())
		tx, de := dbutil.BeginTx(db)
		if de != nil {
			return de
		}
		subRecords := records[idxPair[0]:idxPair[1]]
		for _, record := range subRecords {
			err := app_store.ArchiveInsertTradeActionResp(tx, insertRecordStmt, record, archivingDate...)
			if dbutil.IsMysqlDuplicateEntryError(err) {
				err = nil
			}
			if err != nil {
				domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in dumpTradeActionResps")
				de = domain_error.Build(domain_error.DATABASE_OPERATION_ERR_CODE, err)
				dbutil.RollbackTx(tx)
				return de
			}
		}
		de = dbutil.CommitTx(tx)
		if de != nil {
			return de
		}
	}
	return
}

func (a *OrderArchiver) dumpExtendTradeOrders(tableName string, records []*schema.TradeOrder, archivingDate ...string) (de *domain_error.Error) {
	n := len(records)
	if n == 0 {
		return
	}
	log.Printf("===> dump %d extend trade order records to table %s\n", n, tableName)
	idxPairs := splitArray(n, batchSize)
	insertRecordStmt := fmt.Sprintf(a._insertDataToExtendTradeOrderTableWithoutIDSql, tableName)
	insertRecordStmt = app_store.GetArchiveInsertRecordStmt(insertRecordStmt, archivingDate...)
	for _, idxPair := range idxPairs {
		tx, de := dbutil.BeginTx(a.applicationCfg.GetCentralDB())
		if de != nil {
			return de
		}
		subRecords := records[idxPair[0]:idxPair[1]]
		for _, record := range subRecords {
			args, err := a.getInsertExtendTradeOrderArgs(record, false, "", archivingDate...)
			if err != nil {
				dbutil.RollbackTx(tx)
				return domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
			}
			_, err = tx.Exec(insertRecordStmt, args...)
			if dbutil.IsMysqlDuplicateEntryError(err) {
				err = nil
			}
			if err != nil {
				domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in dumpExtendTradeOrders")
				de = domain_error.Build(domain_error.DATABASE_OPERATION_ERR_CODE, err)
				dbutil.RollbackTx(tx)
				return de
			}
		}
		de = dbutil.CommitTx(tx)
		if de != nil {
			return de
		}
	}
	return
}

func (a *OrderArchiver) dumpExtendTradeActionResps(orderMap map[string]*types.TraceableTradeOrder, tableName string, records []*schema.TradeActionResp, archivingDate ...string) (de *domain_error.Error) {
	n := len(records)
	if n == 0 {
		return
	}
	log.Printf("===> dump %d extend trade action resp records to table %s\n", n, tableName)
	idxPairs := splitArray(n, batchSize)
	insertRecordStmt := fmt.Sprintf(a._insertDataToExtendTradeActionRespTableWithoutIDSql, tableName)
	insertRecordStmt = app_store.GetArchiveInsertRecordStmt(insertRecordStmt, archivingDate...)
	for _, idxPair := range idxPairs {
		tx, de := dbutil.BeginTx(a.applicationCfg.GetCentralDB())
		if de != nil {
			return de
		}
		subRecords := records[idxPair[0]:idxPair[1]]
		for _, record := range subRecords {

			order, ok := orderMap[record.AppOrdID]
			if !ok {
				continue
			}


			// 复制一个订单
			tradeOrder := &schema.TradeOrder{}
			err := bean.Copy(order.GetBasicInfo()).To(tradeOrder)
			if err != nil {
				dbutil.RollbackTx(tx)
				return domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
			}
			// 跟traceable-order.go的逻辑保持一致
			tradeOrder.OrdStatus = record.OrdStatus
			tradeOrder.OrdStatusUpdateTime = record.TransactTime
			tradeOrder.LastShares = record.LastShares
			tradeOrder.LastPx = record.LastPx
			tradeOrder.LeavesQty = record.LeavesQty
			tradeOrder.CumQty = record.CumQty
			tradeOrder.AvgPx = record.AvgPx
			tradeOrder.OrdRejReason = record.OrdRejReason

			args, err := a.getInsertExtendTradeActionRespArgs(tradeOrder, record, false, archivingDate...)
			if err != nil {
				dbutil.RollbackTx(tx)
				return domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
			}
			_, err = tx.Exec(insertRecordStmt, args...)
			if dbutil.IsMysqlDuplicateEntryError(err) {
				err = nil
			}
			if err != nil {
				domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in dumpExtendTradeActionResps")
				de = domain_error.Build(domain_error.DATABASE_OPERATION_ERR_CODE, err)
				dbutil.RollbackTx(tx)
				return de
			}
		}
		de = dbutil.CommitTx(tx)
		if de != nil {
			return de
		}
	}
	return
}

/*
func (a *OrderArchiver) extractTraceableTradeOrder(traceableTradeOrder *types.TraceableTradeOrder) (tradeOrders []*schema.TradeOrder, tradeActionLatestResps []*schema.TradeActionLatestResp, tradeActionResps []*schema.TradeActionResp) {

	tradeOrder := traceableTradeOrder.GetBasicInfo()
	tradeOrders = append(tradeOrders, tradeOrder)

	tradeActions := traceableTradeOrder.GetTradeActions()
	for _, tradeAction := range tradeActions {
		_tradeActionResps := tradeAction.GetTradeActionRespList()
		_tradeActionLatestResp := tradeAction.GetTradeActionLatestResp()
		_tradeActionLatestResp.Account = tradeOrder.Account
		_tradeActionLatestResp.Symbol = tradeOrder.Symbol
		_tradeActionLatestResp.SymbolSfx = tradeOrder.SymbolSfx
		_tradeActionLatestResp.SecurityID = tradeOrder.SecurityID
		_tradeActionLatestResp.IDSource = tradeOrder.IDSource
		_tradeActionLatestResp.SecurityType = tradeOrder.SecurityType
		_tradeActionLatestResp.Side = tradeOrder.Side
		_tradeActionLatestResp.OpenClose = tradeOrder.OpenClose
		_tradeActionLatestResp.OrderQty = tradeOrder.OrderQty
		_tradeActionLatestResp.CashOrderQty = tradeOrder.CashOrderQty
		_tradeActionLatestResp.OrdType = tradeOrder.OrdType
		_tradeActionLatestResp.Price = tradeOrder.Price
		_tradeActionLatestResp.Currency = tradeOrder.Currency
		// _tradeActionLatestResp.EffectiveTime = tradeOrder.EffectiveTime
		// _tradeActionLatestResp.ExpireTime = tradeOrder.ExpireTime

		tradeActionLatestResps = append(tradeActionLatestResps, _tradeActionLatestResp)
		tradeActionResps = append(tradeActionResps, _tradeActionResps...)
	}

	subOrders := traceableTradeOrder.GetSubOrders()
	if len(subOrders) == 0 {
		return
	} else {
		for _, subsubOrder := range subOrders {
			tradeOrders2, tradeActionLatestResps2, tradeActionResps2 := a.extractTraceableTradeOrder(subsubOrder)
			tradeOrders = append(tradeOrders, tradeOrders2...)
			tradeActionLatestResps = append(tradeActionLatestResps, tradeActionLatestResps2...)
			tradeActionResps = append(tradeActionResps, tradeActionResps2...)
		}
	}

	return
}*/

func splitArray(n, m int) [][]int {
	var result [][]int
	for start := 0; start < n; start += m {
		end := start + m
		if end > n {
			end = n
		}
		result = append(result, []int{start, end})
	}
	return result
}

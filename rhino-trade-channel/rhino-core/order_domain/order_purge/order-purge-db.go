package order_purge

import (
	"fmt"
	"log"
	"reflect"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	"rhino-core/store/admin_store"
	"time"
)

func (a *OrderPurger) copyTableStructureTemporary(originalTableName string, purgingLog *schema.DataPurgingLog) (string, error) {
	targetTableName := fmt.Sprintf(originalTableName+"_%s_%s", "tmp", purgingLog.PurgingDate)
	//return targetTableName, dbutil.CopyTableStructure(a.applicationCfg.GetAppDB(), a.applicationCfg.GetCentralDB(), originalTableName, targetTableName, true)
	return targetTableName, dbutil.CopyTableStructure(a.applicationCfg.GetAppDB(), a.applicationCfg.GetAppDB(), originalTableName, targetTableName, false)
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

// Todo: 如果要支持7*24，resetDB要改成删除模式，即每个purgingLog，单独删除其相关的数据，如果表空之后，就truncate
// 第3步：分析内存订单数据，拆分待清除和继续保留的数据
func (a *OrderPurger) resetDB_Old(tradeOrders []*schema.TradeOrder, tradeActionLatestResps []*schema.TradeActionLatestResp, tradeActionResps []*schema.TradeActionResp, purgingLog *schema.DataPurgingLog) (err error) {

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

	allDbTaskPrepared, err := a.checkIfAllDbTaskPrepared(purgingLog)
	if err != nil {
		return err
	}

	if allDbTaskPrepared {
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
	}

	purgingLog.CompletePhase = int(enum.DataPurgingLogPhase_ResetDB)
	err = admin_store.UpdateDataPurgingLogById(a.applicationCfg.GetCentralDB(), purgingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::resetDB")
		return
	}

	return
}

func (a *OrderPurger) checkIfAllDbTaskPrepared(purgingLog *schema.DataPurgingLog) (bool, error) {
	if a.orderArchiver.IsLast() != nil && !*a.orderArchiver.IsLast() {
		return false, nil
	} else if a.orderArchiver.IsLast() != nil {
		// 检查是否所有其他的purging任务都结束了
		archivingConfigItems := a.applicationCfg.GetApplicationArchivingCfgItems()
		for _, archivingConfigItem := range archivingConfigItems {
			if archivingConfigItem.IsLast {
				continue
			}
			purgingDate := purgingLog.PurgingDate
			otherPurgingLog, err := admin_store.GetDataPurgingLogBySystemCodeAndBusinessCodeAndPurgingDateAndTaskName(a.applicationCfg.GetCentralDB(), archivingConfigItem.SystemCode, archivingConfigItem.BusinessCode, purgingDate, purgingLog.TaskName)
			if dbutil.IsDbRecordEmptyError(err) {
				err = nil
			}
			if err != nil {
				return false, err
			}
			if otherPurgingLog == nil {

				t, err1 := time.ParseInLocation("20060102", purgingLog.PurgingDate, timeutil.CnTimeLocation)
				if err1 != nil {
					return false, err1
				}
				purgingDate = t.Add(-24 * time.Hour).In(timeutil.CnTimeLocation).Format("20060102")

				otherPurgingLog, err = admin_store.GetDataPurgingLogBySystemCodeAndBusinessCodeAndPurgingDateAndTaskName(a.applicationCfg.GetCentralDB(), archivingConfigItem.SystemCode, archivingConfigItem.BusinessCode, purgingDate, purgingLog.TaskName)
				if dbutil.IsDbRecordEmptyError(err) {
					err = nil
				}
				if err != nil {
					return false, err
				}
			}

			if otherPurgingLog == nil {
				return false, fmt.Errorf("fail to get DataPurgingLog for SystemCode:%s, BusinessCode:%s, PurgingDate:%s, TaskName:%s\n", archivingConfigItem.SystemCode, archivingConfigItem.BusinessCode, purgingDate, purgingLog.TaskName)
			}

			if !otherPurgingLog.Complete {
				return false, fmt.Errorf("DataPurgingLog for SystemCode:%s, BusinessCode:%s, PurgingDate:%s, TaskName:%s is not completed!\n", archivingConfigItem.SystemCode, archivingConfigItem.BusinessCode, purgingDate, purgingLog.TaskName)
			}
		}
	}
	return true, nil
}

func (a *OrderPurger) resetDB(tradeOrdersToKeep []*schema.TradeOrder, tradeActionLatestRespsToKeep []*schema.TradeActionLatestResp, tradeActionRespsToKeep []*schema.TradeActionResp, tradeOrdersToArchive []*schema.TradeOrder, tradeActionLatestRespsToArchive []*schema.TradeActionLatestResp, tradeActionRespsToArchive []*schema.TradeActionResp, purgingLog *schema.DataPurgingLog) (err error) {

	if purgingLog.CompletePhase >= int(enum.DataPurgingLogPhase_ResetDB) {
		return
	}

	autoTruncate := a.orderArchiver.IsLast() == nil || *a.orderArchiver.IsLast()

	// idSet := extractIDs[*schema.TradeOrder](tradeOrdersToArchive)
	// err = dbutil.DeleteByIDsInBatches(a.applicationCfg.GetAppDB(), "trade_orders", "f_id", 1000, idSet, autoTruncate)
	// if err != nil {
	//log.Printf("DeleteByIDsInBatches for trade_orders with error:%v\n", err)
	err = dbutil.DeleteByKeysInBatches(a.applicationCfg.GetAppDB(), "trade_orders", "f_app_ord_id", 1000, extractKeys[*schema.TradeOrder](tradeOrdersToArchive, "AppOrdID"), autoTruncate)
	// }
	if err != nil {
		return
	}

	// idSet = extractIDs[*schema.TradeActionLatestResp](tradeActionLatestRespsToArchive)
	// err = dbutil.DeleteByIDsInBatches(a.applicationCfg.GetAppDB(), "trade_action_latest_resps", "f_id", 1000, idSet, autoTruncate)
	// if err != nil {
	//log.Printf("DeleteByIDsInBatches for trade_action_latest_resps with error:%v\n", err)
	err = dbutil.DeleteByKeysInBatches(a.applicationCfg.GetAppDB(), "trade_action_latest_resps", "f_app_ord_id", 1000, extractKeys[*schema.TradeActionLatestResp](tradeActionLatestRespsToArchive, "AppOrdID"), autoTruncate)
	//}
	if err != nil {
		return
	}

	// idSet = extractIDs[*schema.TradeActionResp](tradeActionRespsToArchive)
	// err = dbutil.DeleteByIDsInBatches(a.applicationCfg.GetAppDB(), "trade_action_resps", "f_id", 1000, idSet, autoTruncate)
	// if err != nil {
	//log.Printf("DeleteByIDsInBatches for trade_action_resps with error:%v\n", err)
	err = dbutil.DeleteByKeysInBatches(a.applicationCfg.GetAppDB(), "trade_action_resps", "f_cl_ord_id", 1000, extractKeys[*schema.TradeActionResp](tradeActionRespsToArchive, "ClOrdID"), autoTruncate)
	if err != nil {
		return
	}

	err = dbutil.DeleteByKeysInBatches(a.applicationCfg.GetAppDB(), "trade_action_resps", "f_orig_cl_ord_id", 1000, extractKeys[*schema.TradeActionResp](tradeActionRespsToArchive, "OrigClOrdID"), autoTruncate)
	if err != nil {
		return
	}
	//}

	purgingLog.CompletePhase = int(enum.DataPurgingLogPhase_ResetDB)
	err = admin_store.UpdateDataPurgingLogById(a.applicationCfg.GetCentralDB(), purgingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::resetDB")
		return
	}

	return
}

// ExtractIDs 接受任意结构体切片（元素可为结构体或结构体指针），返回 ID 切片
func extractIDs[T any](slice []T) []int64 {
	result := make([]int64, 0, len(slice))
	for _, item := range slice {
		v := reflect.ValueOf(item)

		// 如果是指针，获取其指向的值
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		// 只处理结构体
		if v.Kind() != reflect.Struct {
			continue
		}

		field := v.FieldByName("ID")
		if !field.IsValid() {
			continue // 或 panic，根据需求决定
		}

		// 确保字段类型为 int64（可扩展支持其他整型）
		if field.Kind() == reflect.Int64 {
			result = append(result, field.Int())
		}
		// 可添加其他整型转换逻辑
	}
	return result
}

// extractIDs 接受任意结构体切片，根据 keyField 指定的字段名提取字段值，
// 将所有值转换为字符串返回（如果字段不存在或非结构体则跳过）。
func extractKeys[T any](slice []T, keyField string) []string {
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		v := reflect.ValueOf(item)

		// 如果是指针，获取其指向的值
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		// 只处理结构体
		if v.Kind() != reflect.Struct {
			continue
		}

		field := v.FieldByName(keyField)
		if !field.IsValid() {
			continue // 字段不存在则跳过
		}

		// 将字段值转换为字符串（支持所有基本类型，如 int、string、float 等）
		result = append(result, fmt.Sprint(field.Interface()))
	}
	return result
}

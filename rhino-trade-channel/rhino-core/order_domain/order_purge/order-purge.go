package order_purge

import (
	"darkpool-common/timeutil"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/jsonutil"
	"rhino-core/domain_cfg"
	"rhino-core/order_domain/order_archive"
	"rhino-core/order_domain/order_cache"
	"rhino-core/order_domain/order_position_manager"
	"rhino-core/order_domain/schedule"
	"rhino-core/schema"
	"rhino-core/store/admin_store"
	"rhino-core/store/app_store"
	"strings"
	"sync"
	"time"
)

type OrderPurger struct {
	applicationCfg  *domain_cfg.ApplicationCfg
	orderCache      *order_cache.OrderCache
	orderArchiver   *order_archive.OrderArchiver // 需要调用orderArchiver的工具方法
	scheduleAdapter schedule.ScheduleAdapter
	positionManager *order_position_manager.PositionManager
	working         bool
	lock            *sync.RWMutex
}

func NewOrderPurger(applicationCfg *domain_cfg.ApplicationCfg, orderCache *order_cache.OrderCache, orderArchiver *order_archive.OrderArchiver, scheduleAdapter schedule.ScheduleAdapter, positionManager *order_position_manager.PositionManager) *OrderPurger {
	inst := &OrderPurger{applicationCfg: applicationCfg, orderCache: orderCache, orderArchiver: orderArchiver, scheduleAdapter: scheduleAdapter, lock: &sync.RWMutex{}, positionManager: positionManager}
	return inst
}

// 启动数据清理器
func (a *OrderPurger) Start() {
	go func() {
		for {
			purgingLog, shouldPurge := a.shouldPurgeData(false)
			if shouldPurge {
				a.orderPurging(purgingLog)
			}

			time.Sleep(20 * time.Second)
		}
	}()
}

// 强制归档
func (a *OrderPurger) ForcePurging() {
	purgingLog, shouldArchive := a.shouldPurgeData(true)
	if shouldArchive {
		a.orderPurging(purgingLog)
	}
}

func (a *OrderPurger) shouldPurgeData(force bool) (purgingLog *schema.DataPurgingLog, shouldPurge bool) {

	archivingLog, beginTimeSecond, endTimeSencond, currentTimeSecond, systemCode, businessCode, dateStr, err := a.orderArchiver.GetCurrentArchiveLog()

	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::shouldPurgeData")
		return
	}

	if !force {
		log.Printf("currentTimeSecond:%d, beginTimeSecond:%d, endTimeSencond:%d\n", currentTimeSecond, beginTimeSecond, endTimeSencond)
		if archivingLog != nil {
			log.Printf("archivingLog, archivingLog.CompletePhase:%d, archivingLog.Complete:%v\n", archivingLog.CompletePhase, archivingLog.Complete)
		}
		// 需要同时满足时间约束和归档完成的条件
		shouldPurge = currentTimeSecond >= beginTimeSecond && currentTimeSecond <= endTimeSencond && archivingLog != nil && archivingLog.Complete
	} else {
		// 即使是强制模式，也需要先完成归档
		shouldPurge = archivingLog != nil && archivingLog.Complete
	}

	if shouldPurge {
		purgingLog, err = admin_store.GetDataPurgingLogBySystemCodeAndBusinessCodeAndPurgingDateAndTaskName(a.applicationCfg.GetCentralDB(), systemCode, businessCode, dateStr, archivingLog.TaskName)
		if dbutil.IsDbRecordEmptyError(err) {
			err = nil
		}
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::shouldPurgeData")
			return
		}
		// 已经完成清理时，就不必重复清理
		if purgingLog != nil && purgingLog.Complete {
			shouldPurge = false
			log.Printf("purging task complete for data %s, skip!\n", dateStr)
			return
		}
	}

	if shouldPurge && purgingLog == nil {
		currTime := timeutil.ConvertTimeToMilliseconds(time.Now())
		purgingLog = &schema.DataPurgingLog{
			SystemCode:         systemCode,
			BusinessCode:       businessCode,
			PurgingDate:        dateStr,
			TaskName:           archivingLog.TaskName,
			FirstPurgingTime:   currTime,
			CurrentPurgingTime: currTime,
			Complete:           false,
			CompletePhase:      int(enum.DataPurgingLogPhase_New),
		}
	}

	return
}

func (a *OrderPurger) orderPurging(purgingLog *schema.DataPurgingLog) {

	a.lock.Lock()
	defer func() {
		a.working = false
		a.lock.Unlock()
	}()

	// 设置工作状态标记
	a.working = true

	jsonutil.Print("======> start to purge:\n", purgingLog)

	// 新增或者更新数据清理日志
	err := a.insertOrUpdatePurgingLog(purgingLog)
	if err != nil {
		return
	}

	// 处理数据清理前的准备工作
	var tradeOrdersToKeep []*schema.TradeOrder
	var tradeActionLatestRespsToKeep []*schema.TradeActionLatestResp
	var tradeActionRespsToKeep []*schema.TradeActionResp
	var tradeOrdersToArchive []*schema.TradeOrder
	var tradeActionLatestRespsToArchive []*schema.TradeActionLatestResp
	var tradeActionRespsToArchive []*schema.TradeActionResp
	tradeOrdersToKeep, tradeActionLatestRespsToKeep, tradeActionRespsToKeep, tradeOrdersToArchive, tradeActionLatestRespsToArchive, tradeActionRespsToArchive, err = a.preparePurgingInfo(purgingLog)
	if err != nil {
		return
	}

	log.Printf("tradeOrdersToKeep.Len=%d, tradeActionLatestRespsToKeep.Len=%d, tradeActionRespsToKeep.Len=%d, tradeOrdersToArchive.Len=%d, tradeActionLatestRespsToArchive.Len=%d\n", len(tradeOrdersToKeep), len(tradeActionLatestRespsToKeep), len(tradeActionRespsToKeep), len(tradeOrdersToArchive), len(tradeActionLatestRespsToArchive))

	// 重置数据库数据
	err = a.resetDB(tradeOrdersToKeep, tradeActionLatestRespsToKeep, tradeActionRespsToKeep, tradeOrdersToArchive, tradeActionLatestRespsToArchive, tradeActionRespsToArchive, purgingLog)
	if err != nil {
		return
	}

	log.Printf("finis purging db!")

	// 清理kafka消息
	err = a.resetKafka(purgingLog)
	if err != nil {
		return
	}

	log.Printf("finis purging kafka!")

	// 重置内存模型
	err = a.resetMem(tradeOrdersToKeep, tradeActionLatestRespsToKeep, tradeActionRespsToKeep, tradeOrdersToArchive, tradeActionLatestRespsToArchive, tradeActionRespsToArchive, purgingLog)
	if err != nil {
		return
	}

	log.Printf("finis purging memory!")

	if a.scheduleAdapter != nil {
		// 当scheduleAdapter不为空时，需要执行项目自定义的清理任务
		err = a.scheduleAdapter.ExecCustomizedPurgingTask(tradeOrdersToKeep, tradeActionLatestRespsToKeep, tradeActionRespsToKeep, tradeOrdersToArchive, tradeActionLatestRespsToArchive, tradeActionRespsToArchive, purgingLog, a.positionManager)
		if err != nil {
			return
		}

		purgingLog.CompletePhase = int(enum.DataPurgingLogPhase_CustomizedTask)
		err = admin_store.UpdateDataPurgingLogById(a.applicationCfg.GetCentralDB(), purgingLog)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::customizedTask")
			return
		}
	}

	purgingLog.Complete = true
	err = admin_store.UpdateDataPurgingLogById(a.applicationCfg.GetCentralDB(), purgingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::set complete flag")
	}

	log.Printf("finis purging all!")

	systemCode, businessCode := a.applicationCfg.GetSystemAndBusinessCodes()
	cleanDailyTable(a.applicationCfg.GetCentralDB(), 31, systemCode, businessCode)
}

// 第1步：新增或更新清理日志
func (a *OrderPurger) insertOrUpdatePurgingLog(purgingLog *schema.DataPurgingLog) (err error) {
	if purgingLog.ID == 0 {
		err = admin_store.InsertDataPurgingLog(a.applicationCfg.GetCentralDB(), purgingLog)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::insertPurgingLog")
			return
		}
		return
	}

	purgingLog.CurrentPurgingTime = timeutil.ConvertTimeToMilliseconds(time.Now())
	purgingLog.ExecCount += 1
	err = admin_store.UpdateDataPurgingLogById(a.applicationCfg.GetCentralDB(), purgingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::updatePurgingLog")
		return
	}

	return
}

// 第2步：分析内存订单数据，拆分待清除和继续保留的数据
func (a *OrderPurger) preparePurgingInfo(purgingLog *schema.DataPurgingLog) (
	tradeOrdersToKeep []*schema.TradeOrder, tradeActionLatestRespsToKeep []*schema.TradeActionLatestResp, tradeActionRespsToKeep []*schema.TradeActionResp,
	tradeOrdersToArchive []*schema.TradeOrder, tradeActionLatestRespsToArchive []*schema.TradeActionLatestResp, tradeActionRespsToArchive []*schema.TradeActionResp,
	err error) {

	// if purgingLog.CompletePhase >= int(enum.DataPurgingLogPhase_Ready) {
	// 	return
	// }

	ordersToArchive, ordersToKeep, err1 := a.orderArchiver.SplitTradeOrders(nil)
	if err1 != nil {
		err = err1
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::preparePurgingInfo")
		return
	}

	// 先更新需保留订单的状态到数据库
	_, tradeOrdersToKeep, tradeActionLatestRespsToKeep, tradeActionRespsToKeep = a.orderArchiver.ExtractSchemaData(ordersToKeep)

	// 保守起见，先flush autotx，因为tradeActionLatestResp的
	a.applicationCfg.GetAutoTx().Flush()

	tx, de := dbutil.BeginTx(a.applicationCfg.GetAppDB())
	if de != nil {
		domain_error.ProcessSevereError(false, 0, de, err, "error occurs in OrderPurger::preparePurgingInfo::BeginTx")
		err = de.ToSimpleError()
		return
	}

	for _, v := range tradeOrdersToKeep {
		err = app_store.UpdateTradeOrderByAppOrdId(tx, v)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::preparePurgingInfo::UpdateTradeOrderByAppOrdId")
			return
		}
	}

	for _, v := range tradeActionLatestRespsToKeep {
		// Todo: 检查插入db的actionKey和在内存中的acitonKey的算法是否一致【针对下单、草稿保存、草稿删除、撤单的tradeActionLatestResp的actionKey已确认OK】
		err = app_store.UpdateTradeActionLatestRespByActionKey(tx, v)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::preparePurgingInfo::UpdateTradeActionLatestRespByActionKey")
			return
		}
	}

	for _, v := range tradeActionRespsToKeep {
		err = app_store.UpdateTradeActionRespByClOrdIdAndExecIdAndChannelCode(tx, v)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::preparePurgingInfo::UpdateTradeActionRespByClOrdIdAndExecIdAndChannelCode")
			return
		}
	}

	de = dbutil.CommitTx(tx)
	if de != nil {
		domain_error.ProcessSevereError(false, 0, de, err, "error occurs in OrderPurger::preparePurgingInfo::CommitTx")
		err = de.ToSimpleError()
		return
	}

	// 记录待清除记录的主键
	_, tradeOrdersToArchive, tradeActionLatestRespsToArchive, tradeActionRespsToArchive = a.orderArchiver.ExtractSchemaData(ordersToArchive)
	var tradeOrderKeys []string
	var tradeActionLatestRespKeys []string
	var tradeActionRespKeys []string
	for i, v := range tradeOrdersToArchive {
		log.Printf("#%d order to archive, id=%v, appOrdID=%v\n", i+1, v.ID, v.AppOrdID)
		tradeOrderKeys = append(tradeOrderKeys, v.AppOrdID)
	}
	for _, v := range tradeActionLatestRespsToArchive {
		tradeActionLatestRespKeys = append(tradeActionLatestRespKeys, v.ActionKey)
	}
	for _, v := range tradeActionRespsToArchive {
		tradeActionRespKeys = append(tradeActionRespKeys, v.ClOrdID+","+v.ExecID+","+v.ChannelCode)
	}

	// 记录待删除记录的主键
	purgingLog.TradeOrderPurging = strings.Join(tradeOrderKeys, ";")
	purgingLog.TradeActionLatestRespPurging = strings.Join(tradeActionLatestRespKeys, ";")
	purgingLog.TradeActionRespPurge = strings.Join(tradeActionRespKeys, ";")

	purgingLog.CompletePhase = int(enum.DataPurgingLogPhase_Ready)

	err = admin_store.UpdateDataPurgingLogById(a.applicationCfg.GetCentralDB(), purgingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::preparePurgingInfo")
	}

	return
}

package order_archive

import (
	"encoding/json"
	"errors"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/timeutil"
	"rhino-core/domain_cfg"
	"rhino-core/order_domain/order_cache"
	"rhino-core/order_domain/schedule"
	"rhino-core/schema"
	"rhino-core/store/admin_store"
	"rhino-core/store/app_store"
	"rhino-core/types"
	"strings"
	"sync"
	"time"
)

type OrderArchiver struct {
	applicationCfg  *domain_cfg.ApplicationCfg
	orderCache      *order_cache.OrderCache
	scheduleAdapter schedule.ScheduleAdapter
	working         bool
	lock            *sync.RWMutex

	_insertDataToExtendTradeOrderTableWithIDSql         string
	_insertDataToExtendTradeOrderTableWithoutIDSql      string
	_insertDataToExtendTradeActionRespTableWithIDSql    string
	_insertDataToExtendTradeActionRespTableWithoutIDSql string

	taskName                string
	matchChannels           map[string]bool
	dataArchiveCnBeginTime  string
	dataArchiveCnLatestTime string
	isDSTSensitive          *bool
	isLast                  *bool
}

func NewOrderArchiver(applicationCfg *domain_cfg.ApplicationCfg, archivingCfg *schema.ApplicationArchivingCfgItem, orderCache *order_cache.OrderCache, scheduleAdapter schedule.ScheduleAdapter) *OrderArchiver {
	inst := &OrderArchiver{applicationCfg: applicationCfg, orderCache: orderCache, lock: &sync.RWMutex{}, scheduleAdapter: scheduleAdapter, matchChannels: make(map[string]bool)}
	if archivingCfg != nil {
		inst.taskName = archivingCfg.TaskName
		inst.dataArchiveCnBeginTime = archivingCfg.DataArchiveCnBeginTime
		inst.dataArchiveCnLatestTime = archivingCfg.DataArchiveCnLatestTime
		inst.isDSTSensitive = &archivingCfg.IsDSTSensitive
		inst.isLast = &archivingCfg.IsLast
		strs := strings.Split(archivingCfg.MatchChannels, ",")
		for _, str := range strs {
			str = strings.TrimSpace(str)
			if len(str) == 0 {
				continue
			}
			inst.matchChannels[str] = true
		}
		log.Printf("NewOrderArchiver for archivingCfg, taskName:%v, dataArchiveCnBeginTime:%v, dataArchiveCnLatestTime:%v, isDSTSensitive:%v, matchChannels:%v\n", inst.taskName, inst.dataArchiveCnBeginTime, inst.dataArchiveCnLatestTime, inst.isDSTSensitive)
	}
	return inst
}

// 启动归档器
func (a *OrderArchiver) Start() {
	go func() {
		for {
			archivingLog, shouldArchive := a.shouldArchiveData(false)
			if shouldArchive {
				a.orderArchiving(archivingLog)
			}
			time.Sleep(20 * time.Second)
		}
	}()
	go func() {
		for {
			log.Printf("archiving working %v===> %v\n", a.working, a.taskName)
			time.Sleep(30 * time.Minute)
		}
	}()
}

func (a *OrderArchiver) GetCurrentArchiveLog() (archivingLog *schema.DataArchivingLog, beginTimeSecond, endTimeSencond, currentTimeSecond int, systemCode, businessCode, dateStr string, err error) {
	beginTimeSecond, endTimeSencond, err = a.applicationCfg.GetTimeRangeForDataArchiving(a.dataArchiveCnBeginTime, a.dataArchiveCnLatestTime, a.isDSTSensitive)
	log.Printf("======>GetTimeRangeForDataArchiving, beginTimeSecond: %d, endTimeSencond: %d\n", beginTimeSecond, endTimeSencond)
	if err != nil {
		return
	}

	// 以北京时间为基准
	currentTimeSecond = timeutil.GetCurrentCumulativeSeconds(timeutil.CnTimeZoneName)
	// 当原始的endTimeSencond<beginTimeSecond时，表明是跨日时段，需要定义一个向前时间平移量
	diff := 0
	// 要对跨日的情况进行调整
	if endTimeSencond < beginTimeSecond {
		diff = endTimeSencond
		beginTimeSecond -= diff
		endTimeSencond = timeutil.DateSeconds
		currentTimeSecond -= diff
		if currentTimeSecond < 0 {
			currentTimeSecond += timeutil.DateSeconds
		}
	}

	systemCode, businessCode = a.applicationCfg.GetSystemAndBusinessCodes()

	// 以北京时间为基准获取日期字符串
	// 如果diff>0，要进行时间前移
	dateStr = time.Now().Add(time.Duration(-diff) * time.Duration(time.Second)).In(timeutil.CnTimeLocation).Format("20060102")

	//archivingLog, err = admin_store.GetDataArchivingLogBySystemCodeAndBusinessCodeAndArchivingDate(a.applicationCfg.GetCentralDB(), systemCode, businessCode, dateStr)
	archivingLog, err = admin_store.GetDataArchivingLogBySystemCodeAndBusinessCodeAndArchivingDateAndTaskName(a.applicationCfg.GetCentralDB(), systemCode, businessCode, dateStr, a.taskName)
	if dbutil.IsDbRecordEmptyError(err) {
		err = nil
	}

	return
}

func (a *OrderArchiver) shouldArchiveData(force bool) (archivingLog *schema.DataArchivingLog, shouldArchive bool) {

	archivingLog, beginTimeSecond, endTimeSencond, currentTimeSecond, systemCode, businessCode, dateStr, err := a.GetCurrentArchiveLog()

	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::shouldArchiveData")
		return
	}

	if !force {
		log.Printf("currentTimeSecond:%d, beginTimeSecond:%d, endTimeSencond:%d, taskName:%v\n", currentTimeSecond, beginTimeSecond, endTimeSencond, a.taskName)
		if archivingLog != nil {
			log.Printf("archivingLog, archivingLog.CompletePhase:%d, archivingLog.Complete:%v\n", archivingLog.CompletePhase, archivingLog.Complete)
		}
		shouldArchive = currentTimeSecond >= beginTimeSecond && currentTimeSecond <= endTimeSencond && (archivingLog == nil || !archivingLog.Complete)
	} else {
		log.Printf("force archiving, need to delete existed archivingLog first! systemCode:%s, businessCode:%s, dateStr:%s\n", systemCode, businessCode, dateStr)
		//err = admin_store.DeleteDataArchivingLogBySystemCodeAndBusinessCodeAndArchivingDate(a.applicationCfg.GetCentralDB(), systemCode, businessCode, dateStr)
		err = admin_store.DeleteDataArchivingLogBySystemCodeAndBusinessCodeAndArchivingDateAndTaskName(a.applicationCfg.GetCentralDB(), systemCode, businessCode, dateStr, a.taskName)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::DeleteDataArchivingLog")
			return
		}

		shouldArchive = true
		archivingLog = nil
	}

	if archivingLog == nil {
		currTime := timeutil.ConvertTimeToMilliseconds(time.Now())
		archivingLog = &schema.DataArchivingLog{
			SystemCode:           systemCode,
			BusinessCode:         businessCode,
			ArchivingDate:        dateStr,
			TaskName:             a.taskName,
			FirstArchivingTime:   currTime,
			CurrentArchivingTime: currTime,
			Complete:             false,
			CompletePhase:        int(enum.DataArchivingLogPhase_New),
		}
	}

	return
}

// 强制归档
func (a *OrderArchiver) ForceArchiving() {
	archivingLog, shouldArchive := a.shouldArchiveData(true)
	if shouldArchive {
		a.orderArchiving(archivingLog)
	}
}

// 执行一次归档任务
func (a *OrderArchiver) orderArchiving(archivingLog *schema.DataArchivingLog) {
	a.lock.Lock()
	defer func() {
		a.working = false
		a.lock.Unlock()
	}()

	// 设置工作状态标记
	a.working = true

	jsData, _ := json.MarshalIndent(archivingLog, "", "  ")
	log.Printf("======> start to archive:\n%s\n", jsData)

	// 处理归档日志
	err := a.insertOrUpdateArchivingLog(archivingLog)
	if err != nil {
		return
	}

	// 订单拆分
	ordersToArchive, ordersToKeep, err1 := a.SplitTradeOrders(archivingLog)
	if err1 != nil {
		return
	}
	log.Printf("after SplitTradeOrders, ordersToArchive.Len=%d, ordersToKeep.Len=%d\n", len(ordersToArchive), len(ordersToKeep))

	// 处理日表
	orderMap, tradeOrders, tradeActionLatestResps, tradeActionResps, err := a.createAndDumpToDailyTables(archivingLog, ordersToArchive, ordersToKeep)
	if err != nil {
		return
	}

	// 处理历史表
	err = a.createAndDumpToHistoricalTables(archivingLog, orderMap, tradeOrders, tradeActionLatestResps, tradeActionResps)
	if err != nil {
		return
	}

	// 处理完历史表之后，应标记为已完成
	archivingLog.Complete = true
	err = admin_store.UpdateDataArchivingLogById(a.applicationCfg.GetCentralDB(), archivingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::set complete flag")
		return
	}
}

// 第1步：新增或更新归档日志
func (a *OrderArchiver) insertOrUpdateArchivingLog(archivingLog *schema.DataArchivingLog) (err error) {
	if archivingLog.ID == 0 {
		err = admin_store.InsertDataArchivingLog(a.applicationCfg.GetCentralDB(), archivingLog)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::insertArchivingLog")
			return
		}
		return
	}

	archivingLog.CurrentArchivingTime = timeutil.ConvertTimeToMilliseconds(time.Now())
	archivingLog.ExecCount += 1
	err = admin_store.UpdateDataArchivingLogById(a.applicationCfg.GetCentralDB(), archivingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::updateArchivingLog")
		return
	}

	return
}

// 第2步：拆分订单为两部分，需要归档的订单以及需要保留的订单
// Todo：当前仅考虑日内直通订单的情况，即草稿态的不归档，其他已经执行的都归档。后续需要考虑跨日订单和复杂订单拓扑的情景。
func (a *OrderArchiver) SplitTradeOrders(archivingLog *schema.DataArchivingLog) (ordersToArchive, ordersToKeep []*types.TraceableTradeOrder, err error) {
	ordersToArchive, ordersToKeep = a.orderCache.FilterOrderByFunction(func(order *types.TraceableTradeOrder) bool {
		// Todo: 当前仅考虑日内直通订单的情况，即草稿态的不归档，其他已经执行的都归档。后续需要考虑跨日订单和复杂订单拓扑的情景。
		return order.GetBasicInfo().OrdStatus != string(enum.OrdStatus_Draft)
	})

	ordersToArchive = a.filterByChannelMap(ordersToArchive)
	ordersToKeep = a.filterByChannelMap(ordersToKeep)

	if archivingLog != nil && archivingLog.CompletePhase == int(enum.DataArchivingLogPhase_New) {
		archivingLog.CompletePhase = int(enum.DataArchivingLogPhase_Ready)
		err = admin_store.UpdateDataArchivingLogById(a.applicationCfg.GetCentralDB(), archivingLog)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::SplitTradeOrders")
			return
		}
	}

	return
}

func (a *OrderArchiver) filterByChannelMap(orders []*types.TraceableTradeOrder) (ordersFilter []*types.TraceableTradeOrder) {
	if len(a.matchChannels) == 0 {
		return orders
	}
	for _, order := range orders {
		if a.matchChannels[order.GetBasicInfo().ChannelCode] {
			ordersFilter = append(ordersFilter, order)
		}
	}
	return
}

// 第3步：创建日表并完成日表归档
func (a *OrderArchiver) createAndDumpToDailyTables(archivingLog *schema.DataArchivingLog, ordersToArchive, ordersToKeep []*types.TraceableTradeOrder) (orderMap map[string]*types.TraceableTradeOrder, tradeOrders []*schema.TradeOrder, tradeActionLatestResps []*schema.TradeActionLatestResp, tradeActionResps []*schema.TradeActionResp, err error) {
	if archivingLog.CompletePhase == int(enum.DataArchivingLogPhase_CompleteDailyTable) {
		// 之前已经执行完日表处理，要处理好返回参数之后才能退出
		orderMap, tradeOrders, tradeActionLatestResps, tradeActionResps = a.ExtractSchemaData(ordersToArchive)
		return
	}

	// 进行表的复制（在这个过程中，如果目标表存在，则会被清理并重新创建）
	dailyGroupTradeOrderTableName, dailyTradeActionLatestRespTableName, dailyTradeActionRespTableName, dailyTradeOrderTableName, extendDailyTradeOrderTableName, extendDailyTradeActionRespTableName, err := a.createDailyTables(archivingLog)
	if err != nil {
		return
	}

	var de *domain_error.Error
	// 将OrderCache的数据存储到日表
	orderMap, tradeOrders, tradeActionLatestResps, tradeActionResps, de = a.dumpDataToDailyTables(dailyGroupTradeOrderTableName, dailyTradeActionLatestRespTableName, dailyTradeActionRespTableName, dailyTradeOrderTableName, extendDailyTradeOrderTableName, extendDailyTradeActionRespTableName, archivingLog, ordersToArchive)
	if de != nil {
		if de.Err != nil {
			err = de.Err
			return
		} else {
			err = errors.New(de.ErrorString())
			return
		}
	}

	// 更新归档日志
	archivingLog.CompletePhase = int(enum.DataArchivingLogPhase_CompleteDailyTable)
	err = admin_store.UpdateDataArchivingLogById(a.applicationCfg.GetCentralDB(), archivingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::updateArchivingLog")
		return
	}

	return
}

// 第4步：创建历史表并完成记录追加
func (a *OrderArchiver) createAndDumpToHistoricalTables(archivingLog *schema.DataArchivingLog, orderMap map[string]*types.TraceableTradeOrder, tradeOrders []*schema.TradeOrder, tradeActionLatestResps []*schema.TradeActionLatestResp, tradeActionResps []*schema.TradeActionResp) (err error) {
	if archivingLog.CompletePhase == int(enum.DataArchivingLogPhase_MergeToHisTable) {
		// 之前已经执行完日表处理，直接退出
		return
	}

	// 进行表的复制（在这个过程中，如果目标表存在，则会被清理并重新创建）
	historicalGroupTradeOrderTableName, historicalTradeActionLatestRespTableName, historicalTradeActionRespTableName, historicalTradeOrderTableName, extendHistoricalTradeOrderTableName, extendHistoricalTradeActionRespTableName, err := a.createHistoricalTables(archivingLog)
	if err != nil {
		return
	}

	// 删除当日已经归档的记录
	tables := []string{historicalGroupTradeOrderTableName, historicalTradeActionLatestRespTableName, historicalTradeActionRespTableName, historicalTradeOrderTableName, extendHistoricalTradeOrderTableName, extendHistoricalTradeActionRespTableName}
	for _, table := range tables {
		log.Printf("DeleteHistoricalArchiveDate for table:%s, taskName:%s\n", table, archivingLog.TaskName)
		err = app_store.DeleteHistoricalArchiveDate(a.applicationCfg.GetCentralDB(), table, archivingLog.ArchivingDate, archivingLog.TaskName)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::DeleteHistoricalArchiveDate")
			return
		}
	}

	// 将OrderCache的数据存储到历史表
	de := a.dumpDataToHistoricalTables(historicalGroupTradeOrderTableName, historicalTradeActionLatestRespTableName, historicalTradeActionRespTableName, historicalTradeOrderTableName, extendHistoricalTradeOrderTableName, extendHistoricalTradeActionRespTableName, archivingLog, orderMap, tradeOrders, tradeActionLatestResps, tradeActionResps)
	log.Printf("dumpDataToHistoricalTables, historicalGroupTradeOrderTableName:%s, historicalTradeActionLatestRespTableName:%s, historicalTradeActionRespTableName:%s, historicalTradeOrderTableName:%s, extendHistoricalTradeOrderTableName:%s, extendHistoricalTradeActionRespTableName:%s, archivingLog:{%s,%s,%v,%v}, orderMap.Len:%d, tradeOrders.Len:%d, tradeActionLatestResps.Len:%d, tradeActionResps:%d\n",
		historicalGroupTradeOrderTableName, historicalTradeActionLatestRespTableName, historicalTradeActionRespTableName, historicalTradeOrderTableName, extendHistoricalTradeOrderTableName, extendHistoricalTradeActionRespTableName, archivingLog.ArchivingDate, archivingLog.TaskName, archivingLog.Complete, archivingLog.CompletePhase, len(orderMap), len(tradeOrders), len(tradeActionLatestResps), len(tradeActionResps))
	if de != nil {
		if de.Err != nil {
			err = de.Err
			return
		} else {
			err = errors.New(de.ErrorString())
			return
		}
	}

	// 更新归档日志
	archivingLog.CompletePhase = int(enum.DataArchivingLogPhase_MergeToHisTable)
	err = admin_store.UpdateDataArchivingLogById(a.applicationCfg.GetCentralDB(), archivingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::updateArchivingLog")
		return
	}

	return
}

func (a *OrderArchiver) GetTaskName() string {
	return a.taskName
}

func (a *OrderArchiver) IsLast() *bool {
	return a.isLast
}

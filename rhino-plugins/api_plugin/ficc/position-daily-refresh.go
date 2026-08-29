package ficc

import (
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/timeutil"
	"rhino-common/utils/tradedate"
	"rhino-core/schema"
	"rhino-core/store/admin_store"
	"strings"
	"time"
)

func (a *TitansFiccAPIAdapter) positionDailyRefresh(configMap map[string]*schema.ApplicationCfgItem) {

	err := a.doPositionDailyRefresh(configMap)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to doPositionDailyRefresh")
	}

	go func() {
		for {
			nowTime := time.Now()

			nextRefreshTime, err := time.ParseInLocation(time.DateTime, nowTime.In(timeutil.CnTimeLocation).Format(time.DateOnly)+" "+configMap["ContractPositionResetTime"].ConfigItemValue, timeutil.CnTimeLocation)
			if err != nil {
				domain_error.ProcessSevereError(true, 5, nil, err, "ContractPositionResetTime is not config correct, value=%s"+configMap["ContractPositionResetTime"].ConfigItemValue)
			}

			if nextRefreshTime.Before(nowTime) {
				nextRefreshTime, err = time.ParseInLocation(time.DateTime, nowTime.Add(24*time.Hour).In(timeutil.CnTimeLocation).Format(time.DateOnly)+" "+configMap["ContractPositionResetTime"].ConfigItemValue, timeutil.CnTimeLocation)
				if err != nil {
					domain_error.ProcessSevereError(true, 5, nil, err, "ContractPositionResetTime is not config correct, value=%s"+configMap["ContractPositionResetTime"].ConfigItemValue)
				}
			}

			log.Printf("======> next position refreshTime:%v\n", nextRefreshTime)

			time.Sleep(nextRefreshTime.Sub(nowTime))

			err = a.doPositionDailyRefresh(configMap)
			i := 0
			for err != nil {
				if i >= 360 {
					break
				}
				domain_error.ProcessSevereError(false, 0, nil, err, "fail to doPositionDailyRefresh")
				time.Sleep(10 * time.Second)

				err = a.doPositionDailyRefresh(configMap)

				i++
			}
		}
	}()
}

// 参考：http://wiki.gf.com.cn/pages/viewpage.action?pageId=291973406
//
//	http://wiki.gf.com.cn/pages/viewpage.action?pageId=291238811
func (a *TitansFiccAPIAdapter) doPositionDailyRefresh(configMap map[string]*schema.ApplicationCfgItem) error {

	// 重置
	// if a.autoSyncRepo != nil {
	// 	a.autoSyncRepo.Reset()
	// }
	refreshT1DayCache(configMap)

	endTimeStr := strings.Split(configMap["TradingTime"].ConfigItemValue, "-")[1]
	syncDate, lastTradeDate, err := computeSyncDate(configMap, endTimeStr)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to computeSyncDate")
	}

	systemCode, businessCode := a.applicationCfg.GetSystemAndBusinessCodes()

	// 1、插入更新合约持仓的SyncLog
	syncParams := map[string]interface{}{
		"source": "titans",
		"url":    configMap["ContracePositionServiceUrl"].ConfigItemValue,
		"headers": map[string]interface{}{
			"AppId":     configMap["ServiceAppID"].ConfigItemValue,
			"AppSecret": configMap["ServiceAppSecret"].ConfigItemValue,
		},
		"bodyJson": map[string]interface{}{
			"tradeDate": lastTradeDate,
		},
		"method":          "POST",
		"resultType":      "json",
		"usePage":         true,
		"pageArgType":     0,
		"pageNumField":    "pageNum",
		"pageSizeField":   "pageSize",
		"pagingFrom":      1,
		"pagingSize":      20,
		"totalFieldPath":  "data.total",
		"recordFieldPath": "data.queryResults",
		"csvHeader":       "keyCtptyId,ctptyShortName,contractId,contractCode,keyPlanId,planCode,ultraContractId,ultraContractCode,keyInstrumentId,windCode,insShtDesc,currency,exchange,startDate,endDate,tradeDate,structureSettlementSpeed,structureSettlementDate,interestSettlementSpeed,interestSettlementDate,initPrice,initQuantity,quantity,notional,dynamicNotional,openTradeAmount,longShort,openBondNetPrice,openBondGrossPrice",
	}
	_syncParams, _ := json.Marshal(syncParams)
	syncLog := &schema.DataSyncLog{
		SystemCode:       systemCode,
		BusinessCode:     businessCode,
		SyncDate:         syncDate,
		TableName:        "ContractPosition",
		SyncType:         int(enum.SyncType_PAGING_HTTP_HOOK),
		SyncParams:       string(_syncParams),
		ReportTime:       timeutil.ConvertTimeToMilliseconds(time.Now()),
		FirstSyncTime:    0,
		CurrentSyncTime:  0,
		CompleteSyncTime: 0,
		ExecCount:        0,
		SyncPhase:        int(enum.DataSyncLogPhase_New),
		FailLog:          "",
	}

	err = admin_store.InsertDataSyncLog(a.applicationCfg.GetCentralDB(), syncLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to InsertDataSyncLog for ContractPosition")
		return err
	}

	// 2、插入更新底仓的SyncLog
	syncLog = &schema.DataSyncLog{
		SystemCode:       systemCode,
		BusinessCode:     businessCode,
		SyncDate:         syncDate,
		TableName:        "PositionBase",
		SyncType:         int(enum.SyncType_PAGING_HTTP_HOOK),
		SyncParams:       string(_syncParams),
		ReportTime:       timeutil.ConvertTimeToMilliseconds(time.Now()),
		FirstSyncTime:    0,
		CurrentSyncTime:  0,
		CompleteSyncTime: 0,
		ExecCount:        0,
		SyncPhase:        int(enum.DataSyncLogPhase_New),
		FailLog:          "",
	}

	err = admin_store.InsertDataSyncLog(a.applicationCfg.GetCentralDB(), syncLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to InsertDataSyncLog for ContractPosition")
		return err
	}

	// 3、插入更新资金列表
	syncParams = map[string]interface{}{
		"source": "titans",
		"url":    strings.TrimRight(configMap["CapitalServiceUrl"].ConfigItemValue, "/") + "/accounts/%v?keyCtptyId=%v",
		"headers": map[string]interface{}{
			"AppId":     configMap["ServiceAppID"].ConfigItemValue,
			"AppSecret": configMap["ServiceAppSecret"].ConfigItemValue,
		},
		"method":          "GET",
		"resultType":      "json",
		"useIteration":    true,
		"argsList":        []string{"${db.int.counterparties.KEY_CTPTY_ID}", "${db.int.counterparties.KEY_CTPTY_ID}"},
		"recordFieldPath": "data",
		"csvHeader":       "keyCtptyId,ctptyShortName,keyCapAcctId,capAcctCode,capitalAccountName,purpose,currency",
	}
	_syncParams, _ = json.Marshal(syncParams)
	syncLog = &schema.DataSyncLog{
		SystemCode:       systemCode,
		BusinessCode:     businessCode,
		SyncDate:         syncDate,
		TableName:        "CapitalAccount",
		SyncType:         int(enum.SyncType_ITERATIVE_HTTP_HOOK),
		SyncParams:       string(_syncParams),
		ReportTime:       timeutil.ConvertTimeToMilliseconds(time.Now()),
		FirstSyncTime:    0,
		CurrentSyncTime:  0,
		CompleteSyncTime: 0,
		ExecCount:        0,
		SyncPhase:        int(enum.DataSyncLogPhase_New),
		FailLog:          "",
	}
	err = admin_store.InsertDataSyncLog(a.applicationCfg.GetCentralDB(), syncLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to InsertDataSyncLog for ContractPosition")
		return err
	}

	return nil
}

const (
	dateOnly = "20060102"
	dateTime = "20060102-15:04"
)

// endTimeStr格式：HH:mm
func computeSyncDate(configMap map[string]*schema.ApplicationCfgItem, endTimeStr string) (syncDate string, lastTradeDate string, err error) {

	var date time.Time

	// 当前时间
	nowTime := time.Now().In(timeutil.CnTimeLocation)
	// 当前日期
	nowDateStr := nowTime.Format(dateOnly)
	// 当日的结束时间
	nowDateEndTime, err := time.ParseInLocation(dateTime, nowDateStr+"-"+endTimeStr, timeutil.CnTimeLocation)
	if err != nil {
		return "", "", err
	}
	// 如果当前时间早于当日的结束时间，日期字符串取当日的日期字符串
	if nowTime.Before(nowDateEndTime) {
		//return  nowDateStr, nil
		date = nowTime
	} else { // 否则，取次日的日期
		//return nowTime.Add(24 * time.Hour).Format(dateOnly), nil
		date = nowTime.Add(24 * time.Hour)
	}

	syncDate = date.Format(dateOnly)

	currDate := date.Format(time.DateOnly)

	// 计算交易日参数
	lastTradeDate = getTradeDate(configMap, currDate, 0)
	// 一定要比currDate小，currDate可以为自然日，比如周六日的日期
	if lastTradeDate >= currDate {
		lastTradeDate = getTradeDate(configMap, lastTradeDate, -1)
	}

	log.Printf("======>computeSyncDate, syncDate:%s, lastTradeDate:%s\n", syncDate, lastTradeDate)

	return
}

func getTradeDate(configMap map[string]*schema.ApplicationCfgItem, date string, day int) (tradeDate string) {

	for k, v := range configMap {
		log.Printf("configMap key=%s, v=%s\n", k, v.ConfigItemValue)
	}

	var err error
	tradeDate, err = tradedate.GetTradeDay(configMap["TradeDateServiceUrl"].ConfigItemValue, date, day, "NIB", configMap["ServiceAppID"].ConfigItemValue, configMap["ServiceAppSecret"].ConfigItemValue)
	for err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to get trade data by params, url=%s, tradeDate=%s, day=%d, serviceAppID=%s, serviceAppSecret=%s", configMap["TradeDateServiceUrl"].ConfigItemValue, date, day, configMap["ServiceAppID"].ConfigItemValue, configMap["ServiceAppSecret"].ConfigItemValue))
		tradeDate, err = tradedate.GetTradeDay(configMap["TradeDateServiceUrl"].ConfigItemValue, date, day, "NIB", configMap["ServiceAppID"].ConfigItemValue, configMap["ServiceAppSecret"].ConfigItemValue)
	}
	return
}

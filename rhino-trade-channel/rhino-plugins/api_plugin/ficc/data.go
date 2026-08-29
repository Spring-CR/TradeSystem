package ficc

import (
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	"strings"
	"sync/atomic"
	"time"
)

func (a *TitansFiccAPIAdapter) initData(configMap map[string]*schema.ApplicationCfgItem) {

	// data, err := os.ReadFile(configMap["PathConfigData"].ConfigItemValue)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Printf("config data: %s\n", data)

	// config := &datamap.DataMapConfig{}
	// err = json.Unmarshal(data, config)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, tableConfig := range config.TableConfigs {
	// 	data, _ := json.Marshal(tableConfig)
	// 	log.Printf("======>table %s config:%s\n", tableConfig.TableAlias, data)
	// }

	// autoSyncRepo := datamap.NewAutoSyncRepo(config, a.applicationCfg.GetCentralDB(), a.applicationCfg.GetAppDB(), 0, nil)
	// autoSyncRepo.Start()

	// a.autoSyncRepo = autoSyncRepo
	a.autoSyncRepo = a.applicationCfg.GetAutoSyncRepo()

	err := a.initTradeTimeRange(configMap)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to initTradeTimeRange")
	}

	go func() {
		nowTime := time.Now()
		nextRefreshTime, err := time.ParseInLocation(time.DateTime, nowTime.Add(24*time.Hour).In(timeutil.CnTimeLocation).Format(time.DateOnly)+" 00:00:01", timeutil.CnTimeLocation)
		if err != nil {
			domain_error.ProcessSevereError(true, 5, nil, err, "fail to set trade time span")
		}

		time.Sleep(nextRefreshTime.Sub(nowTime))

		err = a.initTradeTimeRange(configMap)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "fail to set trade time span")
		}
	}()
}

func (a *TitansFiccAPIAdapter) initTradeTimeRange(configMap map[string]*schema.ApplicationCfgItem) error {

	strs := strings.Split(configMap["TradingTime"].ConfigItemValue, "-")
	if len(strs) != 2 {
		return fmt.Errorf("TradingTime is not config correctly, value=%s", configMap["TradingTime"].ConfigItemValue)
	}

	a.t0SellEndTimeStr = configMap["T0SellEndTime"].ConfigItemValue
	
	dateStr := time.Now().In(timeutil.CnTimeLocation).Format(time.DateOnly)
	beginTime, err := time.ParseInLocation(time.DateTime, dateStr + " " + strings.TrimSpace(strs[0]) + ":00", timeutil.CnTimeLocation)
	if err != nil {
		return err
	}
	endTime, err := time.ParseInLocation(time.DateTime, dateStr + " " + strings.TrimSpace(strs[1]) + ":00", timeutil.CnTimeLocation)
	if err != nil {
		return err
	}
	t0SellEndTimeendTime, err := time.ParseInLocation(time.DateTime, dateStr + " " + strings.TrimSpace(a.t0SellEndTimeStr) + ":00", timeutil.CnTimeLocation)
	if err != nil {
		return err
	}

	atomic.StoreInt64(&a.tradeTimeBegin, timeutil.ConvertTimeToMilliseconds(beginTime))
	atomic.StoreInt64(&a.tradeTimeEnd, timeutil.ConvertTimeToMilliseconds(endTime))
	atomic.StoreInt64(&a.t0SellTimeEnd, timeutil.ConvertTimeToMilliseconds(t0SellEndTimeendTime))

	log.Printf("config tradeTimeBegin=%d, tradeTimeEnd=%d, t0SellTimeEnd=%d\n", a.tradeTimeBegin, a.tradeTimeEnd, a.t0SellTimeEnd)

	return nil
}

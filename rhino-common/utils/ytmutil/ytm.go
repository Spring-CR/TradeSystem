package ytmutil

import (
	"encoding/json"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/utils/timeutil"
	"rhino-common/utils/tradedate"
	"rhino-core/schema"
	"strings"
	"time"
)

type YtmUtil struct {
	t1DayCache map[string]string
	configMap  map[string]*schema.ApplicationCfgItem
}

func NewYtmUtil(configMap map[string]*schema.ApplicationCfgItem) *YtmUtil {
	inst := &YtmUtil{
		t1DayCache: make(map[string]string),
		configMap:  configMap,
	}


	inst.refreshT1DayCache(configMap)
	go func() {
		for {
			nowTime := time.Now()

			nextRefreshTime, err := time.ParseInLocation(time.DateTime, nowTime.In(timeutil.CnTimeLocation).Format(time.DateOnly)+" "+configMap["ContractPositionResetTime"].ConfigItemValue, timeutil.CnTimeLocation)
			if err != nil {
				domain_error.ProcessSevereError(true, 5, nil, err, "ContractPositionResetTime is not config correct, value=%s"+configMap["ContractPositionResetTime"].ConfigItemValue)
				continue
			}

			if nextRefreshTime.Before(nowTime) {
				nextRefreshTime, err = time.ParseInLocation(time.DateTime, nowTime.Add(24*time.Hour).In(timeutil.CnTimeLocation).Format(time.DateOnly)+" "+configMap["ContractPositionResetTime"].ConfigItemValue, timeutil.CnTimeLocation)
				if err != nil {
					domain_error.ProcessSevereError(true, 5, nil, err, "ContractPositionResetTime is not config correct, value=%s"+configMap["ContractPositionResetTime"].ConfigItemValue)
					continue
				}
			}

			time.Sleep(nextRefreshTime.Sub(nowTime))

			inst.refreshT1DayCache(configMap)
		}
	}()

	return inst
}

func (u *YtmUtil) refreshT1DayCache(configMap map[string]*schema.ApplicationCfgItem) {
	begin := time.Now()
	u.t1DayCache = make(map[string]string)
	for i := -3; i < 7; i++ {
		currDate := time.Now().Add(time.Duration(i*24) * time.Hour).In(timeutil.CnTimeLocation).Format(time.DateOnly)
		settleDate := u.getTradeDate(currDate, 1)
		settleDate = strings.ReplaceAll(settleDate, "-", "")
		u.t1DayCache[currDate] = settleDate
	}
	js, _ := json.Marshal(u.t1DayCache)
	log.Printf("After refreshT1DayCache in ytmutil, t1DayCache=%s, time cost=%v\n", js, time.Since(begin))
}

func (u *YtmUtil) getTradeDate(date string, day int) (tradeDate string) {

	for k, v := range u.configMap {
		log.Printf("configMap key=%s, v=%s\n", k, v.ConfigItemValue)
	}

	var err error
	tradeDate, err = tradedate.GetTradeDay(u.configMap["TradeDateServiceUrl"].ConfigItemValue, date, day, "NIB", u.configMap["ServiceAppID"].ConfigItemValue, u.configMap["ServiceAppSecret"].ConfigItemValue)
	for err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to get trade data by params, url=%s, tradeDate=%s, day=%d, serviceAppID=%s, serviceAppSecret=%s", u.configMap["TradeDateServiceUrl"].ConfigItemValue, date, day, u.configMap["ServiceAppID"].ConfigItemValue, u.configMap["ServiceAppSecret"].ConfigItemValue))
		tradeDate, err = tradedate.GetTradeDay(u.configMap["TradeDateServiceUrl"].ConfigItemValue, date, day, "NIB", u.configMap["ServiceAppID"].ConfigItemValue, u.configMap["ServiceAppSecret"].ConfigItemValue)
	}
	return
}

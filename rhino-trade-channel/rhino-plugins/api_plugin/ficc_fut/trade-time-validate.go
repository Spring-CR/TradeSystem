package ficc_fut

import (
	"errors"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	"strings"
	"sync/atomic"
	"time"
)

func (a *FiccFurAPIAdapter) checkTradeTime(tradeOrder *schema.TradeOrder, trade bool) (de *domain_error.Error) {
	if !trade {
		return
	}

	i := a.getTimeMinutes(time.Now())
	var tradeFlag uint32
	switch tradeOrder.ChannelCode {
	case "olts-fut":
		tradeFlag = atomic.LoadUint32(&a.dmsTradeFlagInMinutes[i])
	case "stars-fut":
		tradeFlag = atomic.LoadUint32(&a.ovsTradeFlagInMinutes[i])
	}

	if tradeFlag == 0 {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.NOT_IN_TRADING_TIME_ERR_CODE, nil)
	}

	return
}

func (a *FiccFurAPIAdapter) configDmsTradeTimeFlagFromTime(baseTime time.Time) {
	t := baseTime
	for i := 0; i < 1440; i++ {
		tradeFlag := a.getTradeFlagForTimeDms(t)
		timeMinute := a.getTimeMinutes(t)
		atomic.StoreUint32(&a.dmsTradeFlagInMinutes[timeMinute], tradeFlag)

		log.Printf("===> configDmsTradeTimeFlagFromTime, t=%v, timeMinute=%d, tradeFlag=%d\n", t, timeMinute, tradeFlag)

		// 计算下一分钟
		t = t.Add(time.Minute)
	}

	log.Printf("configDmsTradeTimeFlagFromTime from time:%v, flags:%v, dmsDayTradeTimeRange:%v, dmsDayTradeTimeRange:%v\n", baseTime, a.dmsTradeFlagInMinutes, a.dmsDayTradeTimeRange, a.dmsDayTradeTimeRange)
}

func (a *FiccFurAPIAdapter) configOvsTradeTimeFlagFromTime(baseTime time.Time) {
	t := baseTime
	for i := 0; i < 1440; i++ {
		tradeFlag := a.getTradeFlagForTimeOvs(t)
		timeMinute := a.getTimeMinutes(t)
		atomic.StoreUint32(&a.ovsTradeFlagInMinutes[timeMinute], tradeFlag)

		log.Printf("===> configOvsTradeTimeFlagFromTime, t=%v, timeMinute=%d, tradeFlag=%d\n", t, timeMinute, tradeFlag)

		// 计算下一分钟
		t = t.Add(time.Minute)
	}

	log.Printf("configOvsTradeTimeFlagFromTime from time:%v, flags:%v, ovsTradeTimeRange=%v\n", baseTime, a.ovsTradeFlagInMinutes, a.ovsTradeTimeRange)
}

func (a *FiccFurAPIAdapter) getTradeFlagForTimeDms(t time.Time) uint32 {

	todayTradable := a.isDmsDayTradable(t)
	timeMinutes := a.getTimeMinutes(t)

	if todayTradable {
		// 交易日，且位于白盘区间
		if timeMinutes >= a.dmsDayTradeTimeRange[0] && timeMinutes < a.dmsDayTradeTimeRange[1] {
			return 1
		}
		// 交易日，且在夜盘开盘后
		if timeMinutes >= a.dmsNightTradeTimeRange[0] {
			return 1
		}
	}

	// 如果昨天是交易日，且在夜盘收盘前
	yesterdayTradable := a.isDmsDayTradable(t.Add(-24 * time.Hour))
	if yesterdayTradable && timeMinutes < a.dmsNightTradeTimeRange[1] {
		return 1
	}

	return 0
}

func (a *FiccFurAPIAdapter) isDmsDayTradable(t time.Time) bool {
	// 检查是否交易日
	dateStr := t.In(timeutil.CnTimeLocation).Format("20060102")

	valList, ok, de := a.autoSyncRepo.Get("Calendar", dateStr+"-A_SHARE")
	if de != nil {
		return false
	}
	//log.Printf("dateStr=%s, Calendar value=%v\n", dateStr, valList)
	if !ok || len(valList) == 0 {
		de = domain_error.Build(domain_error.GENERIC_ERR_CODE, nil, "no data in Calendar for "+dateStr+"-A_SHARE")
		return false
	}
	dateData := valList[len(valList)-1]
	tradable, _, _ := attrutil.GetAttrValue(dateData, "Tradable", enum.AttrValueType_INT)

	return tradable.(int) == 1
}

func (a *FiccFurAPIAdapter) getTradeFlagForTimeOvs(t time.Time) uint32 {

	todayTradable := a.isOvsDayTradable(t)
	timeMinutes := a.getTimeMinutes(t)

	if todayTradable {
		// 交易日，且位于白盘区间
		if timeMinutes >= a.ovsTradeTimeRange[0] {
			return 1
		}
	}

	// 如果昨天是交易日，且在夜盘收盘前
	yesterdayTradable := a.isOvsDayTradable(t.Add(-24 * time.Hour))
	if yesterdayTradable && timeMinutes < a.ovsTradeTimeRange[1] {
		return 1
	}

	return 0
}

func (a *FiccFurAPIAdapter) isOvsDayTradable(t time.Time) bool {
	weekday := t.In(timeutil.CnTimeLocation).Weekday()
	return weekday != time.Sunday && weekday != time.Saturday
}

func (a *FiccFurAPIAdapter) configDmsTradeTimeFlag() {

	strs := strings.Split(a.configMap["TradingTimeDms"].ConfigItemValue, ",")
	if len(strs) != 2 {
		errMsg := fmt.Sprint("error TradingTimeDms:%s\n", a.configMap["TradingTimeDms"].ConfigItemValue)
		domain_error.ProcessSevereError(true, 5, nil, errors.New(errMsg), errMsg)
	}

	a.parseTimeRange(strs[0], &a.dmsDayTradeTimeRange)
	a.parseTimeRange(strs[1], &a.dmsNightTradeTimeRange)

	// 一开始就计算一次
	t := time.Now()
	a.configDmsTradeTimeFlagFromTime(t)

	// 白天收盘时间
	endTime := strings.Split(strs[0], "-")[1]

	go func() {
		for {
			//开始不断计算
			datetimeStr := time.Now().In(timeutil.CnTimeLocation).Format(time.DateOnly) + " " + endTime + ":00"
			nextTime, _ := time.ParseInLocation(time.DateTime, datetimeStr, timeutil.CnTimeLocation)

			t := time.Now()
			if nextTime.Before(t) {
				nextTime = nextTime.Add(24 * time.Hour)
			}

			sleepDuration := time.Until(nextTime)
			log.Printf("sleep %v for next configDmsTradeTimeFlagFromTime\n", sleepDuration)
			time.Sleep(sleepDuration)

			a.configDmsTradeTimeFlagFromTime(time.Now())
		}
	}()

	return
}

func (a *FiccFurAPIAdapter) configOvsTradeTimeFlag() {

	str := a.configMap["TradingTimeOvs"].ConfigItemValue

	a.parseTimeRange(str, &a.ovsTradeTimeRange)

	// 一开始就计算一次
	t := time.Now()
	a.configOvsTradeTimeFlagFromTime(t)

	// 收盘时间
	endTime := strings.Split(str, "-")[1]

	go func() {
		for {
			//开始不断计算
			datetimeStr := time.Now().In(timeutil.CnTimeLocation).Format(time.DateOnly) + " " + endTime + ":00"
			nextTime, _ := time.ParseInLocation(time.DateTime, datetimeStr, timeutil.CnTimeLocation)

			t := time.Now()
			if nextTime.Before(t) {
				nextTime = nextTime.Add(24 * time.Hour)
			}

			sleepDuration := time.Until(nextTime)
			log.Printf("sleep %v for next configOvsTradeTimeFlagFromTime\n", sleepDuration)
			time.Sleep(sleepDuration)

			a.configOvsTradeTimeFlagFromTime(time.Now())
		}
	}()

	return
}

func (a *FiccFurAPIAdapter) parseTimeRange(timeRangeStr string, tradeTimeRange *[2]int) {
	strs := strings.Split(timeRangeStr, "-")
	if len(strs) != 2 {
		errMsg := fmt.Sprint("error timeRangeStr:%s\n", timeRangeStr)
		domain_error.ProcessSevereError(true, 5, nil, errors.New(errMsg), errMsg)
	}
	tradeTimeRange[0] = a.parseTimeMinutes(strs[0])
	tradeTimeRange[1] = a.parseTimeMinutes(strs[1])
}

func (a *FiccFurAPIAdapter) parseTimeMinutes(timeStr string) int {

	log.Printf("parseTimeMinutes, timeStr=%s\n", timeStr)

	t, err := time.ParseInLocation("15:04", timeStr, timeutil.CnTimeLocation)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "error timeStr:"+timeStr)
	}
	return a.getTimeMinutes(t)
}

func (a *FiccFurAPIAdapter) getTimeMinutes(t time.Time) int {
	return timeutil.GetTimeMinutes(t.In(timeutil.CnTimeLocation))
}

func (pm *FiccFurAPIAdapter) waitForCalendarDataReady() {
	log.Println("waitForCalendarDataReady...")
	autoSyncRepo := pm.applicationCfg.GetAutoSyncRepo()
	for !autoSyncRepo.IsCollectionReady("Calendar") {
		log.Printf("PositionBase is not ready, go to sleep 10 seconds...")
		time.Sleep(10 * time.Second)
	}
	m, de := autoSyncRepo.GetMapData("Calendar")
	if de != nil {
		domain_error.ProcessSevereError(false, 0, de, nil, "init Calendar fail")
	}
	log.Printf("get Calendar, data-len=%d\n", len(m))

	log.Println("finish waitForCalendarDataReady")
}



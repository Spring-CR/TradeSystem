package domain_cfg

import (
	"log"
	"rhino-common/domain_error"
	"rhino-common/utils/envutil"
	"rhino-common/utils/timeutil"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	// Todo：默认每个小时更新一次tradeDateNum，在测试环境，如果涉及日切，人工修改日期之后，这里需要调小。因此，开发通过环境变量设置的机制。
	TradeDateUpdateSeconds = 60 * 60
)

func init() {
	val := envutil.GetEnvVarValue("TRADE_DATE_UPDATE_SECONDS", "")
	if len(val) > 0 {
		intVal, err := strconv.Atoi(val)
		if err == nil && intVal > 0 {
			TradeDateUpdateSeconds = intVal
		}
	}
}

func (d *TradeChannelDetails) keepUpdate() {

	go func() {
		for {
			dateNum, sleepForNextDate, _, _ := d.getExchangeDate(time.Now())
			swapped := atomic.CompareAndSwapInt64(&d.currentDateNum, d.currentDateNum, dateNum)
			if swapped {
				log.Printf("======> set dateNum for channel %s as %d\n", d.TradeChannel.ChannelCode, dateNum)
			}

			log.Printf("sleep %v before set dataNum next time for channel %s\n", sleepForNextDate, d.TradeChannel.ChannelCode)

			time.Sleep(sleepForNextDate)
		}
	}()

	go func() {
		for {
			// 每小时反复检测和设置一次
			time.Sleep(time.Duration(TradeDateUpdateSeconds) * time.Second)
			dateNum, _, _, _ := d.getExchangeDate(time.Now())
			swapped := atomic.CompareAndSwapInt64(&d.currentDateNum, d.currentDateNum, dateNum)
			if swapped {
				log.Printf("======> set dateNum for channel %s as %d\n", d.TradeChannel.ChannelCode, dateNum)
			}
		}
	}()
}

func (d *TradeChannelDetails) getExchangeDate(t time.Time) (dataNum int64, sleepForNextDate time.Duration, beginTime, endTime time.Time) {

	t = t.In(timeutil.GetTimeZone(d.TradeChannel.TimeZone))
	dateStr := t.Format(time.DateOnly)

	beginTime = d.getTime(dateStr, d.TradeChannel.BeginTime)
	endTime = d.getTime(dateStr, d.TradeChannel.EndTime)

	if beginTime.After(endTime) {
		endTime = endTime.Add(24*time.Hour)
	}

	if t.Before(beginTime) {
		beginTime = beginTime.Add(-24 * time.Hour)
		endTime = endTime.Add(-24*time.Hour)
	}

	sleepForNextDate = beginTime.Add(24 * time.Hour).Sub(t)
	
	dateStr = beginTime.In(timeutil.GetTimeZone(d.TradeChannel.DisplayTimeZone)).Format("20060102")
	dataNum = timeutil.Parse8BitDateStrToNum(dateStr)

	return
}

func (d *TradeChannelDetails) GetExchangeDate(t time.Time) (dataNum int64, sleepForNextDate time.Duration, beginTime, endTime time.Time) {
	return d.getExchangeDate(t)
}

func (d *TradeChannelDetails) getTime(dateStr, timeStr string) time.Time {
	dateTimeStr := dateStr + " " + timeStr
	t, err := time.ParseInLocation(time.DateTime, dateTimeStr, timeutil.GetTimeZone(d.TradeChannel.TimeZone))
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to get time for dateTimeStr:"+dateTimeStr)
	}
	return t
}

func (d *TradeChannelDetails) GetCurrentExchangeDate() int64 {
	return atomic.LoadInt64(&d.currentDateNum)
}

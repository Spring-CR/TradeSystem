package domain_cfg_test

import (
	"fmt"
	"rhino-common/utils/timeutil"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"testing"
	"time"
)

var tradeChannelDetails = domain_cfg.NewTradeChannelDetails(&schema.TradeChannel{
	TimeZone:        timeutil.UsTimeZoneName,
	DisplayTimeZone: timeutil.CnTimeZoneName,
	BeginTime:       "18:00:00",
	EndTime:         "17:30:00",
}, nil)


func getTime(timeStr, timeZone string) time.Time {
	t, err := time.ParseInLocation(time.DateTime, timeStr, timeutil.GetTimeZone(timeZone))
	if err != nil {
		panic(err)
	}
	return t
}

func Test1(t *testing.T) {
	t1 := getTime("2025-06-23 17:45:00", timeutil.UsTimeZoneName)
	dateNum, sleepForNextDate, beginTime, endTime := tradeChannelDetails.GetExchangeDate(t1)
	fmt.Printf("beginTime=%v, endTime=%v\n", tradeChannelDetails.TradeChannel.BeginTime,  tradeChannelDetails.TradeChannel.EndTime)
	fmt.Printf("Test1, dateNum=%v, sleepForNextDate=%v, beginTime=%v, endTime=%v\n", dateNum, sleepForNextDate, beginTime, endTime)
}

func Test2(t *testing.T) {
	t1 := getTime("2025-06-23 20:00:00", timeutil.UsTimeZoneName)
	dateNum, sleepForNextDate, beginTime, endTime := tradeChannelDetails.GetExchangeDate(t1)
	fmt.Printf("Test2, dateNum=%v, sleepForNextDate=%v, beginTime=%v, endTime=%v\n", dateNum, sleepForNextDate, beginTime, endTime)
}


func Test3(t *testing.T) {
	t1 := getTime("2025-06-23 02:00:00", timeutil.UsTimeZoneName)
	dateNum, sleepForNextDate, beginTime, endTime := tradeChannelDetails.GetExchangeDate(t1)
	fmt.Printf("Test3, dateNum=%v, sleepForNextDate=%v, beginTime=%v, endTime=%v\n", dateNum, sleepForNextDate, beginTime, endTime)
}
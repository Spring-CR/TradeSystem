package order_domain

import (
	"errors"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/utils/timeutil"
	"rhino-core/domain_cfg"
	"rhino-trade-channel/channel"
	"time"

	"github.com/manucorporat/try"
)

func (e *OrderExecutor) trackingMarketCloseEvent(tradeChannelCfgMap map[string]*domain_cfg.TradeChannelCfg) {

	for k, channel := range e.channelMap {
		tradeChannelCfg, ok := tradeChannelCfgMap[k]
		if !ok {
			continue
		}
		e.trackingMarketCloseEventForChannel(channel, tradeChannelCfg)
	}
}

func (e *OrderExecutor) trackingMarketCloseEventForChannel(channel channel.TradeChannelInterface, tradeChannelCfg *domain_cfg.TradeChannelCfg) {
	go func() {
		e.fireMarketCloseEvent(channel, tradeChannelCfg)
		for {
			e.fireMarketCloseEvent(channel, tradeChannelCfg)
		}
	}()
}

func (e *OrderExecutor) fireMarketCloseEvent(channel channel.TradeChannelInterface, tradeChannelCfg *domain_cfg.TradeChannelCfg) {
	duringTradingTime, err := tradeChannelCfg.IsDuringTradingTime(0)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to check is IsDuringTradingTime for channel "+tradeChannelCfg.GetTradeChannel().ChannelCode)
	}
	if !duringTradingTime {
		// 调用channel的OnMarketClosed接口
		try.This(func() {
			channel.OnMarketClosed(e.orderOrchestrator.orderCache)
		}).Catch(func(err try.E) {
			errStr:= fmt.Sprintf("error occur OnMarketClosed for channel %s, error:%v\n", tradeChannelCfg.GetTradeChannel().ChannelCode, err)
			domain_error.ProcessSevereError(false, 0, nil, errors.New(errStr), errStr)
		})
	}

	// 计算休眠时间并开始休眠
	loc := timeutil.GetTimeZone(tradeChannelCfg.GetTimeZone())
	if loc == nil {
		domain_error.ProcessSevereError(true, 5, nil, fmt.Errorf("fail to get timeZone for %s", tradeChannelCfg.GetTimeZone()), "fail get timeZone for channel "+tradeChannelCfg.GetTradeChannel().ChannelCode)
	}
	dateStr := time.Now().In(loc).Format(time.DateOnly)
	dateTimeStr := dateStr + " " + tradeChannelCfg.GetTradeChannel().EndTime
	endtime, err := time.ParseInLocation(time.DateTime, dateTimeStr, loc)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, fmt.Errorf("time ParseInLocation error"), "fail ParseInLocation for channel "+tradeChannelCfg.GetTradeChannel().ChannelCode)
	}

	nowTime := time.Now()
	if endtime.Sub(nowTime) < time.Minute {
		endtime = endtime.Add(24 * time.Hour)
	}

	sleepDuration := time.Until(endtime)
	sleepDuration += 5*time.Minute
	log.Printf("start to sleep %v for next fireMarketCloseEvent for channel %s", sleepDuration, tradeChannelCfg.GetTradeChannel().ChannelCode)
	time.Sleep(sleepDuration)
}

package domain_cfg

import (
	"database/sql"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
)

type TradeChannelDetails struct {
	currentDateNum       int64
	TradeChannel         *schema.TradeChannel
	TradeChannelCfgItems []*schema.TradeChannelCfgItem
}

func NewTradeChannelDetails(tradeChannel *schema.TradeChannel, tradeChannelCfgItems []*schema.TradeChannelCfgItem) *TradeChannelDetails {
	if tradeChannel.DisplayTimeZone == "" {
		tradeChannel.DisplayTimeZone = tradeChannel.TimeZone
	}
	inst := &TradeChannelDetails{TradeChannel: tradeChannel, TradeChannelCfgItems: tradeChannelCfgItems}
	//预加载时区
	timeutil.WarmUpTimeLocation(tradeChannel.TimeZone)
	timeutil.WarmUpTimeLocation(tradeChannel.DisplayTimeZone)
	inst.keepUpdate()
	return inst
}

func (d *TradeChannelDetails) GetTradeChannelCfgItemMap() map[string]*schema.TradeChannelCfgItem {
	m := map[string]*schema.TradeChannelCfgItem{}
	for _, tradeChannelCfgItem := range d.TradeChannelCfgItems {
		m[tradeChannelCfgItem.ConfigItemName] = tradeChannelCfgItem
	}
	return m
}

type CfgItemChecker struct {
	ConfigItemName         string
	ConfigItemDefaultValue string
	Description            string
	Required               int
	NoOverWrite            int
	SetByProgram           int
}

type TradeChannelCfgAdapter interface {
	// 返回客户端配置文件
	ToAppConfig() (configFileContent []byte, de *domain_error.Error)
}

type TradeChannelCfg struct {
	applicationCfg         *ApplicationCfg
	tradeChannelDetails    *TradeChannelDetails
	tradeChannel           *schema.TradeChannel
	tradeChannelCfgItems   []*schema.TradeChannelCfgItem
	tradeChannelCfgAdapter TradeChannelCfgAdapter
	msgSeqGen              *MsgSeqGen
}

func NewTradeChannelCfg(applicationCfg *ApplicationCfg, tradeChannelDetails *TradeChannelDetails) (cfg *TradeChannelCfg, de *domain_error.Error) {
	cfg = &TradeChannelCfg{applicationCfg: applicationCfg, tradeChannelDetails: tradeChannelDetails, tradeChannel: tradeChannelDetails.TradeChannel, tradeChannelCfgItems: tradeChannelDetails.TradeChannelCfgItems}
	switch enum.ChannelProtocolType(cfg.tradeChannel.ChannelProtocolType) {
	case enum.ChannelProtocolType_FIX42:
		cfg.tradeChannelCfgAdapter = newTradeChannelCfgAdapterForFixInitiator(tradeChannelDetails.TradeChannel, tradeChannelDetails.TradeChannelCfgItems)
	case enum.ChannelProtocolType_FIX44:
		cfg.tradeChannelCfgAdapter = newTradeChannelCfgAdapterForFixInitiator(tradeChannelDetails.TradeChannel, tradeChannelDetails.TradeChannelCfgItems)
	case enum.ChannelProtocolType_STARS:
	case enum.ChannelProtocolType_KAFKA:
	default:
		de = domain_error.Build(domain_error.UNKNOW_CHANNEL_PROTOCOL_TYPE_ERR_CODE, nil, cfg.tradeChannel.ChannelProtocolType)
		return
	}

	cfg.msgSeqGen, de = NewMsgSeqGen(cfg)

	return
}

func (c *TradeChannelCfg) GetTradeChannelCfgAdapter() TradeChannelCfgAdapter {
	return c.tradeChannelCfgAdapter
}

func (c *TradeChannelCfg) GetTradeChannel() *schema.TradeChannel {
	return c.tradeChannel
}

func (c *TradeChannelCfg) GetChannelCode() string {
	return c.tradeChannel.ChannelCode
}

func (c *TradeChannelCfg) GetTradeChannelCfgItems() []*schema.TradeChannelCfgItem {
	return c.tradeChannelCfgItems
}

func (c *TradeChannelCfg) GetAppDB() *sql.DB {
	return c.applicationCfg.GetAppDB()
}

func (c *TradeChannelCfg) GetAutoTx() *dbutil.ConcurrentAutoTx {
	return c.applicationCfg.GetAutoTx()
}

func (c *TradeChannelCfg) GetMsgSeqGen() *MsgSeqGen {
	return c.msgSeqGen
}

func (c *TradeChannelCfg) GetApplicationCfg() *ApplicationCfg {
	return c.applicationCfg
}

func (c *TradeChannelCfg) GetTimeZone() string {
	return c.tradeChannel.TimeZone
}

// 增加顺延参数，有利于处理在最后时刻重启时，能有效处理可能缺失的回报
func (c *TradeChannelCfg) IsDuringTradingTime(allowDelaySeconds int) (bool, error) {
	currSeconds := timeutil.GetCurrentCumulativeSeconds(c.tradeChannel.TimeZone)
	minSeconds, err := timeutil.GetCumulativeSecondsFromSimpleTimeString(c.tradeChannel.BeginTime)
	if err != nil {
		return false, err
	}
	maxSeconds, err := timeutil.GetCumulativeSecondsFromSimpleTimeString(c.tradeChannel.EndTime)
	if err != nil {
		return false, err
	}

	log.Printf("#1 currSeconds:%d, minSeconds:%d, maxSeconds:%d\n", currSeconds, minSeconds, maxSeconds)

	if maxSeconds < minSeconds {
		if maxSeconds+allowDelaySeconds >= minSeconds {
			maxSeconds = minSeconds - 1
		} else {
			maxSeconds += allowDelaySeconds
		}
	} else {
		maxSeconds += allowDelaySeconds
	}

	log.Printf("#2 currSeconds:%d, minSeconds:%d, maxSeconds:%d\n", currSeconds, minSeconds, maxSeconds)

	if maxSeconds > timeutil.DateSeconds {
		maxSeconds -= timeutil.DateSeconds
	}

	log.Printf("#3 currSeconds:%d, minSeconds:%d, maxSeconds:%d\n", currSeconds, minSeconds, maxSeconds)

	if maxSeconds < minSeconds { // 表明是跨日的时间段
		return currSeconds >= minSeconds && currSeconds <= timeutil.DateSeconds || currSeconds <= maxSeconds, nil
	}

	return currSeconds >= minSeconds && currSeconds <= maxSeconds, nil
}

func (c *TradeChannelCfg) GetTradeChannelDetails() *TradeChannelDetails {
	return c.tradeChannelDetails
}

/*
通过db.Ping()来维护连接的实时有效性
*/

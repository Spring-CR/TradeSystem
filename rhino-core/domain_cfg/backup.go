package domain_cfg

import (
	"database/sql"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-core/schema"
	"time"
)


type TradeChannelCfg2 struct {
	application            *schema.Application
	tradeChannel           *schema.TradeChannel
	tradeChannelCfgItems   []*schema.TradeChannelCfgItem
	tradeChannelCfgAdapter TradeChannelCfgAdapter
	appDB                  *sql.DB
	//autoTx                 *dbutil.AutoTx
	autoTx                 *dbutil.ConcurrentAutoTx
	autoTxOuputErrChan     chan *domain_error.Error
}

func NewTradeChannelCfg2(application *schema.Application, tradeChannel *schema.TradeChannel, tradeChannelCfgItems []*schema.TradeChannelCfgItem) (cfg *TradeChannelCfg2, de *domain_error.Error) {
	cfg = &TradeChannelCfg2{application: application, tradeChannel: tradeChannel, tradeChannelCfgItems: tradeChannelCfgItems}
	switch enum.ChannelProtocolType(cfg.tradeChannel.ChannelProtocolType) {
	case enum.ChannelProtocolType_FIX42:
		cfg.tradeChannelCfgAdapter = newTradeChannelCfgAdapterForFixInitiator(tradeChannel, tradeChannelCfgItems)
	case enum.ChannelProtocolType_FIX44:
		cfg.tradeChannelCfgAdapter = newTradeChannelCfgAdapterForFixInitiator(tradeChannel, tradeChannelCfgItems)
	default:
		de = domain_error.Build(domain_error.UNKNOW_CHANNEL_PROTOCOL_TYPE_ERR_CODE, nil)
	}

	db, err := sql.Open("mysql", application.DatabaseUrl)
	if err != nil {
		de = domain_error.Build(domain_error.DATABASE_OPERATION_ERR_CODE, err)
		return
	}

	// 参考：https://juejin.cn/post/6844904087427776519
	// chatgpt：https://chatgpt.com/c/671072a7-d9e0-800c-b356-d1e0ba2363e0
	db.SetMaxOpenConns(128) // 这两个值设置为一样会好一些，否则，有时会出现链接有问题
	db.SetMaxIdleConns(128) // 这两个值设置为一样会好一些，否则，有时会出现链接有问题
	db.SetConnMaxLifetime(time.Second * 600)

	cfg.appDB = db
	//cfg.autoTx, cfg.autoTxOuputErrChan = dbutil.NewAutoTx(1*time.Second, 500)
	cfg.autoTx, cfg.autoTxOuputErrChan = dbutil.NewConcurrentAutoTx(16, 1*time.Second, 256)

	cfg.startAutoTx()

	return
}

func (c *TradeChannelCfg2) startAutoTx() {
	go func() {
		for {
			de := <-c.autoTxOuputErrChan
			if de != nil {
				log.Printf("Receive error from autoTx:%s\n", de.ErrorString())
			}
		}
	}()
	// go func() {
	// 	c.autoTx.StartTicker()
	// }()
	c.autoTx.Start()
}

func (c *TradeChannelCfg2) GetTradeChannelCfgAdapter() TradeChannelCfgAdapter {
	return c.tradeChannelCfgAdapter
}

func (c *TradeChannelCfg2) GetTradeChannel() *schema.TradeChannel {
	return c.tradeChannel
}

func (c *TradeChannelCfg2) GetTradeChannelCfgItems() []*schema.TradeChannelCfgItem {
	return c.tradeChannelCfgItems
}

func (c *TradeChannelCfg2) GetAppDB() *sql.DB {
	return c.appDB
}

// func (c *TradeChannelCfg) GetAutoTx() *dbutil.AutoTx {
// 	return c.autoTx
// }
func (c *TradeChannelCfg2) GetAutoTx() *dbutil.ConcurrentAutoTx {
	return c.autoTx
}
package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"rhino-common/domain_error"
	"rhino-core/domain_cfg"
	"rhino-engine/api"
	"rhino-engine/test/engine"
	"time"

	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/goldendb"
)

var (
	f = flag.String("f", "/config/config.json", "config file")
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
}

func main() {

	flag.Parse()

	data, err := os.ReadFile(*f)
	if err != nil {
		panic(err)
	}
	log.Printf("config data: %s\n", data)

	config := &engine.Application{}
	json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}

	var tradeChannelDetailsList []*domain_cfg.TradeChannelDetails
	for _, ch := range config.TradeChannels {
		tradeChannelDetails := domain_cfg.NewTradeChannelDetails(ch.TradeChannel, ch.TradeChannelCfgItems)
		tradeChannelDetailsList = append(tradeChannelDetailsList, tradeChannelDetails)
	}

	applicationCfg, de := domain_cfg.NewApplicationCfg(config.Application, config.ApplicationArchivingCfgItems, config.ApplicationCfgItems, config.ExtendAttrItems, config.PositionAttrItems, config.TradeActionRespAttrItems, tradeChannelDetailsList)
	if de != nil {
		panic(de)
	}

	log.Println("start to create apiServer!")
	apiServer, de := api.NewTradeEngineApiServer(applicationCfg, 128, 2)
	if de != nil {
		panic(de)
	}

	scheduleDailyPanic()

	log.Println("finish create apiServer!")
	log.Println("apiServer started")
	err = apiServer.Start()
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to create api-server!")
	}
}

func scheduleDailyPanic() {
    go func() {
        for {
            // 获取当前时间
            now := time.Now()
            
            // 设置目标时间为今天的 23:59:00
            target := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 0, 0, now.Location())
            
            // 如果今天的目标时间已过，设置为明天的 23:59:00
            if now.After(target) {
                target = target.Add(24 * time.Hour)
            }
            
            // 计算需要等待的时间
            duration := target.Sub(now)
            
            // 等待到目标时间
            timer := time.NewTimer(duration)
            <-timer.C
            
            // 触发 panic
            panic("Scheduled daily panic at 23:59:00!")
        }
    }()
}
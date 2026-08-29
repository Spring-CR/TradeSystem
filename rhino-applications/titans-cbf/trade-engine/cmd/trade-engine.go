package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"rhino-core/domain_cfg"
	"rhino-engine/api"
	"rhino-engine/test/engine"

	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/goldendb"
)

var (
	f = flag.String("f", "/config.json", "config file")
)

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

	applicationCfg, de := domain_cfg.NewApplicationCfg(config.Application, nil, config.ExtendAttrItems, nil, nil, tradeChannelDetailsList)
	if de != nil {
		panic(de)
	}

	log.Println("start to create apiServer!")
	apiServer, de := api.NewTradeEngineApiServer(applicationCfg, 128, 2)
	if de != nil {
		panic(de)
	}

	log.Println("finish create apiServer!")
	log.Println("apiServer started")
	apiServer.Start()
}
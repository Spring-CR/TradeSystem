package main

import (
	"encoding/json"
	"log"
	"os"
	"rhino-core/domain_cfg"
	"rhino-engine/test/engine"
	"rhino-order-report/api/server"
	"time"

	"github.com/urfave/cli"

	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/goldendb"
	_ "github.com/mattn/go-sqlite3"

	//_ "modernc.org/sqlite" // SQLite3 dependency (pure Go)
	_ "rhino-plugins/order_capital_plugin/ficc"     // 导入插件
	_ "rhino-plugins/order_position_plugin/ficc_v2" // 导入插件
	_ "rhino-plugins/order_status_plugin/ficc"      // 导入插件
)

const (
	ENV_CONFIG_PATH = "ConfigPath"
	ENV_SERVER_PORT = "ServerPort"

	OPT_CONFIG_PATH = "config-path"
	OPT_SERVER_PORT = "server-port"
)

var flags = []cli.Flag{
	cli.StringFlag{
		EnvVar: ENV_CONFIG_PATH,
		Name:   OPT_CONFIG_PATH,
		Value:  "/config.json",
		Usage:  "config file path",
	},
	cli.IntFlag{
		EnvVar: ENV_SERVER_PORT,
		Name:   OPT_SERVER_PORT,
		Value:  9091,
		Usage:  "server http port",
	},
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
}

func main() {
	app := cli.NewApp()
	app.Name = "order-report"
	app.Version = "1.0.0"
	app.Usage = "order-report"
	app.Action = action
	app.Flags = flags
	app.Run(os.Args)
}

func action(c *cli.Context) error {

	data, err := os.ReadFile(c.String(OPT_CONFIG_PATH))
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

	scheduleDailyPanic()

	log.Println("start to create orderReport!")
	orderReport := server.NewHttpApiServer(c.Int(OPT_SERVER_PORT), applicationCfg)
	err = orderReport.Start()
	if err != nil {
		panic(err)
	}

	log.Printf("finish create orderReport server on port %v\n", c.String(OPT_SERVER_PORT))

	return nil
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

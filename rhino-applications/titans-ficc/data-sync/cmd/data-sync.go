package main

import (
	"encoding/json"
	"log"
	"os"
	"rhino-common/utils/dbutil"
	"rhino-core/domain_cfg"
	"rhino-data-sync/api/server"
	"rhino-data-sync/data_sync"

	_ "github.com/go-sql-driver/goldendb"
	"github.com/urfave/cli"
)

type Config struct {
	SystemCode          string
	BusinessCode        string
	CentralDatabaseUrl  string
	DatabaseUrl         string
	ServerPort          int
	TableConfigs        []*dbutil.TableConfig
	MrkCloseTime        string
	MrkCloseTimeZone    string
	DataSyncAdapterPath string
	TradeDateServiceUrl string
	TradeDateServiceAID string
	TradeDateServiceSec string
	ApiToken            string
}

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
		Value:  "/config/config.json",
		Usage:  "config file path",
	},
	cli.IntFlag{
		EnvVar: ENV_SERVER_PORT,
		Name:   OPT_SERVER_PORT,
		Value:  0,
		Usage:  "server http port",
	},
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	app := cli.NewApp()
	app.Name = "data-sync"
	app.Version = "1.0.0"
	app.Usage = "data-sync"
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

	config := &Config{}
	json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}

	dataSyncConfig := domain_cfg.NewDataSyncConfig(config.SystemCode, config.BusinessCode, config.CentralDatabaseUrl, config.DatabaseUrl, config.TableConfigs, config.MrkCloseTime, config.MrkCloseTimeZone, config.DataSyncAdapterPath, config.TradeDateServiceUrl, config.TradeDateServiceAID, config.TradeDateServiceSec)

	port := c.Int(OPT_SERVER_PORT)
	if port <= 0 {
		port = config.ServerPort
	}

	// 开启数据库同步
	dataSync := data_sync.NewDataSync(dataSyncConfig)
	dataSync.Start()

	dataSyncServer := server.NewHttpApiServer(port, config.ApiToken, dataSyncConfig)
	err = dataSyncServer.Start()
	if err != nil {
		panic(err)
	}

	log.Printf("finish create dataSync server on port %v\n", config.ServerPort)

	return nil
}

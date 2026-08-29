package main

import (
	"encoding/json"
	"log"
	"os"
	"rhino-data/api/server"
	"rhino-data/datamap"

	_ "github.com/go-sql-driver/goldendb"
	"github.com/urfave/cli"
)

type Config struct {
	*datamap.DataMapConfig
	ServerPort int
	ApiToken   string
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
	app.Name = "data"
	app.Version = "1.0.0"
	app.Usage = "data"
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
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}

	for _, tableConfig := range config.TableConfigs {
		data, _ := json.Marshal(tableConfig)
		log.Printf("======>table %s config:%s\n", tableConfig.TableAlias, data)
	}

	autoSyncRepo := datamap.NewAutoSyncRepo(config.DataMapConfig, nil, nil, 0, nil)
	autoSyncRepo.Start()

	port := c.Int(OPT_SERVER_PORT)
	if port <= 0 {
		port = config.ServerPort
	}

	dataServer := server.NewHttpApiServer(port, config.ApiToken, autoSyncRepo)
	err = dataServer.Start()
	if err != nil {
		panic(err)
	}

	log.Printf("finish create data server on port %v\n", config.ServerPort)

	return nil
}

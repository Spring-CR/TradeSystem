package main

import (
	"database/sql"
	"log"
	"os"
	"rhino-api/api/server"
	"rhino-instr/domain/status"
	trade_channel "rhino-instr/domain/trade-channel"
	"strings"

	"rhino-common/context"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	app := cli.NewApp()
	app.Name = "rhino-api-server"
	app.Version = "1.0.0"
	app.Usage = "rhino-api-server"
	app.Action = action
	app.Flags = flags
	app.Run(os.Args)
}

func initDB(c *cli.Context) {
	db, err := sql.Open("mysql", c.String(OPT_MAIN_DB_URL))
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Second * 500)
	context.DB = db

	db_gfqg, err := sql.Open("mysql", c.String(OPT_GFQG_DB_URL))
	if err != nil {
		log.Fatal(err)
	}
	db_gfqg.SetMaxOpenConns(10)
	db_gfqg.SetMaxIdleConns(5)
	db_gfqg.SetConnMaxLifetime(time.Second * 500)
	context.DB_GFQG = db_gfqg
}

func initTradeChannel(c*cli.Context) {
	brokers := strings.Split(c.String(OPT_KAFKA_BROKERS), ";")
	tradeChannel, de := trade_channel.NewKafkaTradeChannel(c.String(OPT_KAFKA_REQ_TOPIC), c.String(OPT_KAFKA_RSP_TOPIC), brokers)
	if de != nil {
		panic(de.ErrorString())
	}
	trade_channel.DefaultTradeChannel = tradeChannel
	status.StatusObserve(tradeChannel)
}

func action(c *cli.Context) error {
	initDB(c)
	initTradeChannel(c)
	context.DefaultATPUser = c.String(OPT_ATP_USER)
	return server.StartApiServer(c.String(OPT_SERVER_ADDR), c.String(OPT_POS_SVC_URL))
}
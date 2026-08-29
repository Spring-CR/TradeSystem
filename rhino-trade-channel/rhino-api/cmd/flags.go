package main

import (
	"rhino-common/context"

	"github.com/urfave/cli"
)

const (
	ENV_SERVER_ADDR     = "SERVER_ADDR"
	ENV_MAIN_DB_URL     = "MAIN_DB_URL"
	ENV_GFQG_DB_URL     = "GFQG_DB_URL"
	ENV_POS_SVC_URL     = "POS_SVC_URL"
	ENV_KAFKA_BROKERS   = "KAFKA_BROKERS"
	ENV_KAFKA_REQ_TOPIC = "KAFKA_REQ_TOPIC"
	ENV_KAFKA_RSP_TOPIC = "KAFKA_RSP_TOPIC"
	ENV_ATP_USER        = "ATP_USER"
	OPT_SERVER_ADDR     = "addr"
	OPT_MAIN_DB_URL     = "main-db-url"
	OPT_GFQG_DB_URL     = "gfqg-db-url"
	OPT_POS_SVC_URL     = "pos-svc-url"
	OPT_KAFKA_BROKERS   = "kafka-brokers"
	OPT_KAFKA_REQ_TOPIC = "kafka-req-topic"
	OPT_KAFKA_RSP_TOPIC = "kafka-rsp-topic"
	OPT_ATP_USER        = "atp-user"
)

var flags = []cli.Flag{
	cli.StringFlag{
		EnvVar: ENV_SERVER_ADDR,
		Name:   OPT_SERVER_ADDR,
		Value:  ":8080",
		Usage:  "server address",
	},
	cli.StringFlag{
		EnvVar: ENV_MAIN_DB_URL,
		Name:   OPT_MAIN_DB_URL,
		Value:  "root:guangfa4cool@tcp(10.51.136.72:56322)/olts_tradedesk?charset=utf8",
		Usage:  "main database url",
	},
	cli.StringFlag{
		EnvVar: ENV_GFQG_DB_URL,
		Name:   OPT_GFQG_DB_URL,
		Value:  "gfqg:Gfqg1234#@tcp(10.128.13.88:3306)/gfqgdb?charset=utf8",
		Usage:  "GFQG database url",
	},
	cli.StringFlag{
		EnvVar: ENV_POS_SVC_URL,
		Name:   OPT_POS_SVC_URL,
		Value:  "http://10.128.13.89:8888/query/posFundRedis",
		Usage:  "position service url",
	},
	cli.StringFlag{
		EnvVar: ENV_KAFKA_BROKERS,
		Name:   OPT_KAFKA_BROKERS,
		Value:  "10.128.13.230:30490;10.128.13.233:30490;10.128.13.232:30490",
		Usage:  "kafka-brokers",
	},
	cli.StringFlag{
		EnvVar: ENV_KAFKA_REQ_TOPIC,
		Name:   OPT_KAFKA_REQ_TOPIC,
		Value:  "rhino-instr-req-dev",
		Usage:  "kafka request topic",
	},
	cli.StringFlag{
		EnvVar: ENV_KAFKA_RSP_TOPIC,
		Name:   OPT_KAFKA_RSP_TOPIC,
		Value:  "rhino-instr-resp-dev",
		Usage:  "kafka response topic",
	},
	cli.StringFlag{
		EnvVar: ENV_ATP_USER,
		Name:   OPT_ATP_USER,
		Value:  context.DefaultATPUser,
		Usage: "atp-user",
	},
}

package main

import (
	"encoding/json"
	"ficc-utils/api/order_check"
	"ficc-utils/api/server"
	"ficc-utils/common/domain_error"
	"ficc-utils/common/utils/config"
	"ficc-utils/common/utils/synclogs"
	"ficc-utils/common/utils/timeutil"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var(
	c = flag.String("c", "/config.json", "path of the config file(json)")
)

func main() {

	flag.Parse()
	
	data, err := os.ReadFile(*c)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "error occurs while read config file")
	}
	
	log.Printf("===> Config:\n%s\n", data)

	config := &config.Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, fmt.Sprintf("error occurs while unmarshal json data: %s", string(data)))
	}

	if config.QueryCycleSeconds == 0 {
		config.QueryCycleSeconds = 300
	}

	synclogs.InitSyncLogs(config.SyncLogsQryUrl, config.SyncLogsInterval)
	order_check.RunOrderChecks(config)

	local := timeutil.GetTimeZone(config.TimeZoneName)
	if local == nil {
		domain_error.ProcessSevereError(true, 5, nil, fmt.Errorf("fail to get timezone: %s", config.TimeZoneName), fmt.Sprintf("fail to get timezone: %s", config.TimeZoneName))
	}

	apiServer := server.NewHttpApiServer(config.ServerPort, config.DataQryServiceUrl, config.LoginServiceUrl, config.LoginClientID, config.LoginClientKey, config.CapitalServiceUrl, config.PositionServiceUrl, config.CurrentTradeOrdersServiceUrl, config.HisTradeOrdersServiceUrl, config.CurrentExecReportServiceUrl, config.PricingHubYtmConvUrl, config.PricingHubAuth, config.DataQryUrl, config.AppId, config.AppSecret, time.Duration(config.QueryCycleSeconds)*time.Second, local)
	err = apiServer.Start()

	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to start api server")
	}
}
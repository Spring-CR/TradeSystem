package config

type Config struct {
	ServerPort                   int
	DataQryServiceUrl            string
	LoginServiceUrl              string
	LoginClientID                string
	LoginClientKey               string
	CapitalServiceUrl            string
	PositionServiceUrl           string
	CurrentTradeOrdersServiceUrl string
	HisTradeOrdersServiceUrl     string
	CurrentExecReportServiceUrl  string
	PricingHubYtmConvUrl         string
	PricingHubAuth               string
	DataQryUrl                   string
	SyncLogsQryUrl               string
	SyncLogsInterval             string
	QueryCycleSeconds            int
	TimeZoneName                 string
	AppId                        string
	AppSecret                    string
	WebhookUrl                   string
	OrderCheckTasks              []OrderCheckTaskConfig
}

type OrderCheckTaskConfig struct {
	Task      string
	StartTime string
	EndTime   string
	Interval  string
}

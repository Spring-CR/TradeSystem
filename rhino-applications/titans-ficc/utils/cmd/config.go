package main

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
	PricingHubYtmConvUrl         string
	PricingHubAuth               string
	DataQryUrl                   string
	QueryCycleSeconds            int
	TimeZoneName                 string
	AppId                        string
	AppSecret                    string
}

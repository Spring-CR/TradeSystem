package schema

type TradeChannel struct {
	ID                  int64
	ChannelCode         string `sql:"unique: pk_tc, size: 32"`
	ChannelZhName       string `sql:"size: 128"`
	ChannelEnName       string `sql:"size: 32"`
	ChannelProtocolType string `sql:"size: 32"`
	ChannelAdapterName  string `sql:"size: 64"`
	Description         string
	Addresses           string
	Exchange            string `sql:"size: 8"`
	TimeZone            string
	DisplayTimeZone     string `sql:"size: 8"`
	BeginTime           string `sql:"size: 8"`
	EndTime             string `sql:"size: 8"`
	DataNumAdj          int
	ActiveRealAddress   string `sql:"size: 64"`
	RealAddress         string `sql:"size: 64"`
	ExportAddress       string `sql:"size: 64"`
	ExportHttpPort      int
	ExportWSPort        int
	ApiToken            string `sql:"size: 256"`
	Status              int
	OfflineReason       string
	ConfigDir           string
	AdapterPath         string `sql:"size: 128"`
}

type TradeChannelCfgItem struct {
	ID                     int64
	ChannelCode            string `sql:"unique: pk_tcci, size: 32"`
	ConfigItemName         string `sql:"index: pk_tcci, size: 32"`
	ConfigItemValue        string
	ConfigItemDefaultValue string
	Description            string
	Required               int
}

package schema

type Application struct {
	ID                        int64
	SystemCode                string `sql:"unique: pk_app, size: 32"`
	SystemZhName              string `sql:"size: 64"`
	SystemEnName              string `sql:"size: 64"`
	BusinessCode              string `sql:"index: pk_app, size: 32"`
	BusinessZhName            string `sql:"size: 64"`
	BusinessEnName            string `sql:"size: 64"`
	AreaCode                  string `sql:"size: 64"`
	ActivatedTradeChannelCode string `sql:"size: 32"`
	Activated                 bool
	ApiToken                  string `sql:"size: 256"`
	DataArchiveCnBeginTime    string `sql:"size: 64"`
	DataArchiveCnLatestTime   string `sql:"size: 64"`
	IsDSTSensitive            bool
	CentralDatabaseUrl        string
	DatabaseUrl               string
	ApiAdapterPath            string `sql:"size: 128"`
	OrdStatusAdapterPath      string `sql:"size: 128"`
	OrdPositionAdapterPath    string `sql:"size: 128"`
	OrdCapitalAdapterPath     string `sql:"size: 128"`
	OrdExecutorAdapterPath    string `sql:"size: 128"`
	ScheduleAdapterPath       string `sql:"size: 128"`
	FixServerAdapterPath      string `sql:"size: 128"`
	KafkaBrokers              string `sql:"size: 128"`
	HttpAPIPort               int
	WorkingDir                string
	DataRepoConfigPath        string
	PublishDataSyncEvent      bool
}

type ApplicationCfgItem struct {
	ID                     int64
	SystemCode             string `sql:"unique: pk_aci, size: 32"`
	BusinessCode           string `sql:"index: pk_aci, size: 32"`
	ConfigItemName         string `sql:"index: pk_aci, size: 32"`
	ConfigItemValue        string
	ConfigItemDefaultValue string
	Description            string
	Required               int
}

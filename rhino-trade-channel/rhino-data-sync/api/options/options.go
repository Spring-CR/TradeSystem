package options

type DataSyncNotifyOption struct {
	SyncDate   string                 `json:"syncDate"`
	TableName  string                 `json:"tableName"`
	SyncParams map[string]interface{} `json:"syncParams"`
	SyncType   int                    `json:"syncType"`
	//ReportTime    int64                  `json:"reportTime"`
}

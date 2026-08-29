package schema

type DataSyncLog struct {
	ID               int64
	SystemCode       string `sql:"unique: pk_dsl, index: dsl_sbd;dls_sbts, size: 32"`
	BusinessCode     string `sql:"index: pk_dsl;dsl_sbd;dls_sbts, size: 32"`
	SyncDate         string `sql:"index: pk_dsl;dsl_sbd, size: 8"`
	TableName        string `sql:"index: pk_dsl;dsl_sbd;dls_sbts, size: 32"`
	SyncType         int    // 0：大数据DSP接口同步
	SyncParams       string `sql:"type: MEDIUMTEXT"`
	ReportTime       int64  `sql:"index: pk_dsl"`
	FirstSyncTime    int64
	CurrentSyncTime  int64
	CompleteSyncTime int64
	ExecCount        int
	SyncPhase        int    `sql:"index: dls_sbts"`
	FailLog          string `sql:"type: MEDIUMTEXT"`
}

package schema

type DataPurgingLog struct {
	ID                           int64
	SystemCode                   string `sql:"unique: pk_dcl, size: 32"`
	BusinessCode                 string `sql:"index: pk_dcl, size: 32"`
	PurgingDate                  string `sql:"index: pk_dcl, size: 8"`
	TaskName                     string `sql:"index: pk_dcl, size: 32"`
	GroupTradeOrderPurging       string `sql:"type: LONGTEXT"`
	TradeOrderPurging            string `sql:"type: LONGTEXT"`
	TradeActionLatestRespPurging string `sql:"type: LONGTEXT"`
	TradeActionRespPurge         string `sql:"type: LONGTEXT"`
	FirstPurgingTime             int64
	CurrentPurgingTime           int64
	ExecCount                    int
	Complete                     bool
	CompletePhase                int
}

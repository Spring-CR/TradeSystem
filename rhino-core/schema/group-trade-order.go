package schema

type GroupTradeOrder struct {
	ID                  int64
	SystemCode          string `sql:"size: 32"`
	BusinessCode        string `sql:"size: 32"`
	ClGroupOrdID        string `sql:"unique: pk_gto_clid, size: 128"`
	SecurityType        string `sql:"size: 2"`
	SubOrderDeriveType  string `sql:"size: 2"`
	TransactTime        int64
	OrdStatus           string `sql:"size: 2"`
	OrdFillStatus       string `sql:"size: 2"`
	OrdCreator          string `sql:"size: 64"`
	OrdCreateTime       int64
	OrdDraftUpdateUser  string `sql:"size: 64"`
	OrdDraftUpdateTime  string
	OrdDraftDelUser     string `sql:"size: 64"`
	OrdDraftDelTime     int64
	OrdExecUserScope    string
	OrdExecUser         string `sql:"size: 64"`
	OrdStatusUpdateTime int64
	ReviewFlag          int
	ReviewerScope       string
	Reviewer            string
	ApproveStatus       string
	ReviewTime          int64
}

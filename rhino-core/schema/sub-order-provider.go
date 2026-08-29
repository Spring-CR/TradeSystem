package schema

type SubOrderProvider struct {
	ID             int64
	SystemCode     string `sql:"unique: uq_sop, size: 32"`
	BusinessCode   string `sql:"index: uq_sop, size: 32"`
	ProviderCode   string `sql:"index: uq_sop, size: 32"`
	ProviderZhName string `sql:"size: 128"`
	ProviderEnName string `sql:"size: 32"`
	Description    string
	InvokeUrl      string
	ApiToken       string `sql:"size: 256"`
}

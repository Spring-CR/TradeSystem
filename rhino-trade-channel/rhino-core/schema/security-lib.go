package schema

type SecurityLib struct {
	ID                        int64
	SecurityLibCode           string `sql:"unique: pk_sl, size: 32"`
	SecurityType              string `sql:"size: 16"`
	SecurityLibZhName         string `sql:"size: 128"`
	SecurityLibEnName         string `sql:"size: 32"`
	PreferredSecurityIDSource string `sql:"size: 2"`
	DataSource                string `sql:"size: 128"`
	Description               string
	LastSyncDatetime          string `sql:"size: 32"`
}

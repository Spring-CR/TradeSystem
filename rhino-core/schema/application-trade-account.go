package schema

type ApplicationTradeAccount struct {
	ID           int64
	SystemCode   string `sql:"unique: pk_ata, size: 32"`
	BusinessCode string `sql:"index: pk_ata, size: 32"`
	ChannelCode  string `sql:"index: pk_ata, size: 32"`
	AccountCode  string `sql:"index: pk_ata, size: 128"`
	UserId       string `sql:"size: 256"`
}

package schema

type ApplicationTradeChannel struct {
	ID                  int64
	SystemCode          string `sql:"unique: pk_atc, size: 32"`
	BusinessCode        string `sql:"index: pk_atc, size: 32"`
	ChannelCode         string `sql:"index: pk_atc, size: 32"`
	DefaultTradeAccount string `sql:"size: 128"`
	Activated           bool
}

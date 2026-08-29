package schema

type TradeAccount struct {
	ID            int64
	AccountCode   string `sql:"unique: pk_tacct, size: 128"`
	AccountZhName string
	AccountEnName string
	ChannelCode   string `sql:"index: pk_tacct, size: 32"`
	GroupCode     string `sql:"index: pk_tacct, size: 128"`
	Description   string
}

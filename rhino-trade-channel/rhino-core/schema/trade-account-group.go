package schema

type TradeAccountGroup struct {
	ID              int64
	ParentGroupCode string `sql:"index: idx_tag_pgc, size: 128"`
	GroupCode       string `sql:"unique: pk_tag, size: 128"`
	ChannelCode     string `sql:"index: pk_tag, size: 32"`
	GroupZhName     string `sql:"size: 128"`
	GroupEnName     string `sql:"size: 32"`
	Description     string
}

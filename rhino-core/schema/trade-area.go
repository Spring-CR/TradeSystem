package schema

type TradeArea struct {
	ID             int64
	AreaCode       string `sql:"unique: pk_area, size: 64"`
	AreaZhName     string `sql:"size: 128"`
	AreaEnName     string `sql:"size: 128"`
	TradeBeginTime string `sql:"size: 64"`
	TradeEndTime   string `sql:"size: 64"`
}

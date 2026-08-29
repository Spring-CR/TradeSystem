package schema

type TradeActionRespAttrItem struct {
	ID                 int64
	SystemCode         string `sql:"unique: uq_eat, size: 32"`
	BusinessCode       string `sql:"index: uq_eat, size: 32"`
	Required           bool
	AttrName           string `sql:"index: pk_taai, size: 32"`
	AttrZhName         string `sql:"size: 128"`
	AttrValueType      int
	AttrValueLen       int
	AttrMinValue       float64
	AttrMaxValue       float64
	AttrValueRangeType int
	AttrValueRegex     string
	EnumRange          string `sql:"type: MEDIUMTEXT"`
	Index              bool
	Unique             bool
}

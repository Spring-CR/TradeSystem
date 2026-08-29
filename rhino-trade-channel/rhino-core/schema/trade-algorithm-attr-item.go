package schema

type TradeAlgorithmAttrItem struct {
	ID                 int64
	ChannelCode        string `sql:"unique: pk_taai, size: 32"`
	AlgorithmCode      string `sql:"index: pk_taai, size: 32"`
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
}

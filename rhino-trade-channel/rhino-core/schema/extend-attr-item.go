package schema

type ExtendAttrItem struct {
	ID                 int64
	SystemCode         string `sql:"unique: uq_eat, size: 32"`
	BusinessCode       string `sql:"index: uq_eat, size: 32"`
	Required           bool
	AttrName           string `sql:"index: uq_eat, size: 32"`
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
	// MapToOrderAttr     string `sql:"size: 32"`         // 映射到订单属性
	// MapToAlgAttr       string `sql:"size: 32"`         // 映射到算法参数属性
	// DatetimeStrLayout  string `sql:"size: 32"`         // 对于映射到订单属性中的时间字段，可以采用字符串值，需要定义layout，便于进行时间戳的转换。其中，时区取channel的所在时区
	// MappingValuePaires string `sql:"type: MEDIUMTEXT"` // 值对映射，只支持字符串类型，json格式: [{"AttrName": {属性名}, "ValueMap": {"值1": "映射值1", "值2": "映射值2"} }]
}

package enum

// 属性值的类型枚举
type AttrValueType int

const (
	AttrValueType_STRING AttrValueType = 0
	AttrValueType_INT    AttrValueType = 1
	AttrValueType_FLOAT  AttrValueType = 2
	AttrValueType_BOOL   AttrValueType = 3
)

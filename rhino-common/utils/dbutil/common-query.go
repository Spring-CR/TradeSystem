package dbutil

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/acsellers/inflections"
)

const (
	FieldTypeUnknown = 0
	FieldTypeFloat   = 1
	FieldTypeInt     = 2
	FieldTypeString  = 3
	FieldTypeBool    = 4
)

const (
	FieldConditionValueTypeItem  = 0
	FieldConditionValueTypeRange = 1
	FieldConditionValueTypeSet   = 2
)

type FieldType struct {
	FieldName string
	FieldType int // 1-浮点，2-整形，3-字符串，4-布尔
}

type FieldCondition struct {
	Field     string      `json:"field"`
	FieldType int         `json:"field_type"` // 1-浮点，2-整形，3-字符串，4-布尔
	ValueType int         `json:"value_type"` // 0 - 单值；1 - 区间；2 - 集合
	Value     interface{} `json:"value"`      // 对于0，直接取期值；对于 1-转换成 []interface{}，并约定长度是>=1，第一个元素代表min，第二个代表max；对于2-转换成[]interface{}，并约定长度是>=1
	IsNot     bool        `json:"is_not"`     // 是否取not，对于单值和集合的值类型有效

	Operator           string            `json:"operator"` // AND、OR
	SubFieldConditions []*FieldCondition `json:"sub_field_conditions"`
}

func (fc *FieldCondition) getValues(fieldTypeMap map[string]*FieldType) (values []interface{}) {
	vs, ok := fc.Value.([]interface{})
	if ok {
		if fieldTypeMap != nil {
			for _, v := range vs {
				switch fieldTypeMap[fc.Field].FieldType {
				case FieldTypeFloat:
					values = append(values, v.(float64))
				case FieldTypeInt:
					values = append(values, toInt(v))
				case FieldTypeString:
					values = append(values, v.(string))
				case FieldTypeBool:
					values = append(values, v.(bool))
				}
			}
		} else {
			for _, v := range vs {
				switch fc.FieldType {
				case FieldTypeFloat:
					values = append(values, v.(float64))
				case FieldTypeInt:
					values = append(values, toInt(v))
				case FieldTypeString:
					values = append(values, v.(string))
				case FieldTypeBool:
					values = append(values, v.(bool))
				default:
					values = append(values, v)
				}
			}
		}
	}
	return values
}

// toInt 将 interface{} 转换为 int 类型
func toInt(v interface{}) int {
	switch value := v.(type) {
	case float64:
		return int(value)
	case int:
		return value
	default:
		return 0 // 或者返回一个错误
	}
}

var (
	fieldTypeMapCache sync.Map
	fieldPrefix       = "f_"
)

func ConfigDBFieldPrefix(prefix string) {
	fieldPrefix = prefix
}

func PrepareQueryStatment(t reflect.Type, fieldConditions []*FieldCondition) (whereClause string, args []interface{}, err error) {

	var fieldTypeMap map[string]*FieldType
	if t != nil {
		fieldTypeMap = GetFieldTypeMap(t)
		err = CheckValueType(fieldTypeMap, fieldConditions)
		if err != nil {
			return
		}
	}

	for i, fieldCondition := range fieldConditions {
		if i > 0 {
			whereClause += " AND "
		}
		switch fieldCondition.ValueType {
		case FieldConditionValueTypeItem:
			if fieldCondition.FieldType == FieldTypeString {
				whereClause += fieldPrefix + fieldCondition.Field + " like ?"
			} else {
				whereClause += fieldPrefix + fieldCondition.Field + "=?"
			}
			args = append(args, fieldCondition.Value)
		case FieldConditionValueTypeRange:
			values := fieldCondition.getValues(fieldTypeMap)
			whereClause += fieldPrefix + fieldCondition.Field + ">=?"
			if len(values) > 1 {
				whereClause += " AND " + fieldPrefix + fieldCondition.Field + "<?"
			}
			args = append(args, values...)
		case FieldConditionValueTypeSet:
			values := fieldCondition.getValues(fieldTypeMap)
			whereClause += fieldPrefix + fieldCondition.Field + " IN ("
			for i := range values {
				if i > 0 {
					whereClause += ","
				}
				whereClause += "?"
			}
			whereClause += ")"
			args = append(args, values...)
		}
	}

	return whereClause, args, err
}

func CheckValueType(fieldTypeMap map[string]*FieldType, fieldConditions []*FieldCondition) error {
	for _, fieldCondition := range fieldConditions {
		fieldType, ok := fieldTypeMap[fieldCondition.Field]
		if !ok {
			return fmt.Errorf("query file %s is not supported", fieldCondition.Field)
		}

		var values []interface{}
		switch fieldCondition.ValueType {
		case FieldConditionValueTypeItem:
			values = append(values, fieldCondition.Value)
		case FieldConditionValueTypeRange:
			values, ok = fieldCondition.Value.([]interface{})
			if !ok {
				values = append(values, fieldCondition.Value)
			}
			if len(values) == 0 {
				return fmt.Errorf("value range cannot be empty for field %s", fieldType.FieldName)
			}
			if len(values) > 2 {
				return fmt.Errorf("value range can only have the min and max value for field %s", fieldType.FieldName)
			}
			fieldCondition.Value = values
		case FieldConditionValueTypeSet:
			values, ok = fieldCondition.Value.([]interface{})
			if !ok {
				values = append(values, fieldCondition.Value)
			}
			if len(values) == 0 {
				return fmt.Errorf("value set cannot be empty for field %s", fieldType.FieldName)
			}
			fieldCondition.Value = values
		}

		switch fieldType.FieldType {
		case FieldTypeFloat:
			for _, v := range values {
				_, ok := v.(float64)
				if !ok {
					return fmt.Errorf("value of field %s should be float, but %v is not of the expected type", fieldType.FieldName, v)
				}
			}
		case FieldTypeInt:
			for _, v := range values {
				if !isIntValue(v) {
					return fmt.Errorf("value of field %s should be int, but %v is not of the expected type", fieldType.FieldName, v)
				}
			}
		case FieldTypeString:
			for _, v := range values {
				_, ok := v.(string)
				if !ok {
					return fmt.Errorf("value of field %s should be string, but %v is not of the expected type", fieldType.FieldName, v)
				}
			}
		case FieldTypeBool:
			for _, v := range values {
				_, ok := v.(bool)
				if !ok {
					return fmt.Errorf("value of field %s should be bool, but %v is not of the expected type", fieldType.FieldName, v)
				}
			}
		}
	}

	return nil
}

// isIntValue 检查一个 interface{} 是否是一个 int 或合法的 float64（整数）
func isIntValue(v interface{}) bool {
	switch value := v.(type) {
	case float64:
		return value == float64(int(value))
	case int:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case uint:
		return true
	case uint8:
		return true
	case uint16:
		return true
	case uint32:
		return true
	case uint64:
		return true
	default:
		return false
	}
}

func GetFieldTypeMap(t reflect.Type) map[string]*FieldType {

	v, ok := fieldTypeMapCache.Load(t)
	if ok {
		return v.(map[string]*FieldType)
	}

	m := make(map[string]*FieldType)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		//fmt.Printf("字段名称: %s, 字段类型: %s\n", field.Name, field.Type)
		fieldName := inflections.Underscore(field.Name)
		var fieldType int
		typeName := field.Type.String()
		if strings.Contains(typeName, "float") {
			fieldType = FieldTypeFloat
		} else if strings.Contains(typeName, "int") {
			fieldType = FieldTypeInt
		} else if strings.Contains(typeName, "str") {
			fieldType = FieldTypeString
		} else if strings.Contains(typeName, "bool") {
			fieldType = FieldTypeBool
		}
		if fieldType == FieldTypeUnknown {
			continue
		}
		m[fieldName] = &FieldType{FieldName: fieldName, FieldType: fieldType}
	}
	fieldTypeMapCache.Store(t, m)
	return m
}

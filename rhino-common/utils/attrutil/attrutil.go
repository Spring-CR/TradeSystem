package attrutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"rhino-common/enum"
	"strconv"
	"strings"
)

func ParseAttrMapString(mapString string) (valMap map[string]interface{}, err error) {
	if mapString == "" {
		mapString = "{}"
	}
	err = json.Unmarshal([]byte(mapString), &valMap)
	if err != nil {
		return
	}
	return
}

func GetAttrValue(valMap map[string]interface{}, attrName string, valType enum.AttrValueType) (val interface{}, ok bool, err error) {
	if valMap == nil {
		err = errors.New("valMap is null")
		return
	}
	val, ok = valMap[attrName]
	if !ok {
		switch valType {
		case enum.AttrValueType_STRING:
			val = ""
		case enum.AttrValueType_INT:
			val = int(0)
		case enum.AttrValueType_FLOAT:
			val = float64(0.0)
		case enum.AttrValueType_BOOL:
			val = false
		}
		return
	}

	// 字段存在但值为 nil 时，按需处理
	if val == nil {
		switch valType {
		case enum.AttrValueType_STRING:
			val = "" // 空字符串
		case enum.AttrValueType_INT:
			val = int(0) // 默认整数
		case enum.AttrValueType_FLOAT:
			val = 0.0 // 默认浮点数
		case enum.AttrValueType_BOOL:
			val = false // 默认布尔值
		}
		ok = true // 标记为存在但值为 nil（已替换为默认值）
		return
	}

	switch valType {
	case enum.AttrValueType_STRING:
		val, ok = val.(string)
		if !ok {
			err = fmt.Errorf("attribute %s value %v is not type of string as expected", attrName, val)
			return
		}
	case enum.AttrValueType_INT:

		switch v := val.(type) {
		case int:
			return v, true, nil
		case int8:
			return int(v), true, nil
		case int16:
			return int(v), true, nil
		case int32:
			return int(v), true, nil
		case int64:
			return int(v), true, nil
		case uint:
			return int(v), true, nil
		case uint8:
			return int(v), true, nil
		case uint16:
			return int(v), true, nil
		case uint32:
			return int(v), true, nil
		case uint64:
			return int(v), true, nil
		}

		var valFloat float64
		valFloat, ok = val.(float64)
		if !ok {
			valFloat, ok = isNumInStr(val)
		}
		if ok {
			val = valFloat
		}
		if !ok {
			err = fmt.Errorf("attribute %s value %v is not type of int as expected", attrName, val)
			return
		}
		val = int(val.(float64))
	case enum.AttrValueType_FLOAT:

		switch v := val.(type) {
		case float64:
			return v, true, nil
		case float32:
			return float64(v), true, nil
		case int:
			return float64(v), true, nil
		case int8:
			return float64(v), true, nil
		case int16:
			return float64(v), true, nil
		case int32:
			return float64(v), true, nil
		case int64:
			return float64(v), true, nil
		case uint:
			return float64(v), true, nil
		case uint8:
			return float64(v), true, nil
		case uint16:
			return float64(v), true, nil
		case uint32:
			return float64(v), true, nil
		case uint64:
			return float64(v), true, nil
		}

		var valFloat float64
		valFloat, ok = val.(float64)
		if !ok {
			valFloat, ok = isNumInStr(val)
		}
		if ok {
			val = valFloat
		}
		if !ok {
			err = fmt.Errorf("attribute %s value %v is not type of float64 as expected", attrName, val)
			return
		}
	case enum.AttrValueType_BOOL:
		var boolVal bool
		boolVal, ok = val.(bool)
		if !ok {
			boolVal, ok = isBoolInStr(val)
		}
		if ok {
			val = boolVal
		}
		if !ok {
			err = fmt.Errorf("attribute %s value %v is not type of bool as expected", attrName, val)
			return
		}
	}
	return
}

func isNumInStr(val interface{}) (float64, bool) {

	//log.Printf("try to check if val is num of str, val=%v, type=%T\n", val, val)

	strVal, ok := val.(string)
	if !ok {
		return 0, false
	}

	//log.Printf("val=%v is of string type\n", val)

	strVal = strings.TrimSpace(strVal)
	if len(strVal) == 0 {
		return 0, true
	}

	floatVal, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		return 0, false
	}

	return floatVal, true
}

func isBoolInStr(val interface{}) (bool, bool) {
	strVal, ok := val.(string)
	if !ok {
		return false, false
	}
	strVal = strings.ToLower(strings.TrimSpace(strVal))
	if len(strVal) == 0 {
		return false, true
	}
	boolVal := (strVal == "true")
	ok = boolVal || strVal == "false"
	return boolVal, ok
}

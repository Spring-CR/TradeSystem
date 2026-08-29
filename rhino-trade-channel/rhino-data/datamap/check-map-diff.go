package datamap

import (
	"fmt"
)

// CheckDiff 比较两个map，返回新增、修改、删除的键
func CheckDiff(mNew, mOld map[string][]map[string]interface{}) (addKeys, chgKeys, delKeys []string, err error) {

	// 检查新增和修改的记录
	for key, newInnerMaps := range mNew {
		oldInnerMaps, exists := mOld[key]
		if !exists {
			// 新增的记录
			addKeys = append(addKeys, key)
			continue
		}

		// 如果长度不一样，则认为是发生了变化
		if len(newInnerMaps) != len(oldInnerMaps) {
			chgKeys = append(chgKeys, key)
			continue
		}

		for i := 0; i < len(newInnerMaps); i++ {
			newInnerMap := newInnerMaps[i]
			oldInnerMap := oldInnerMaps[i]
			// 检查内部map是否有变化
			isChanged, err := compareInnerMaps(newInnerMap, oldInnerMap)
			if err != nil {
				return addKeys, chgKeys, delKeys, fmt.Errorf("比较键 %s 时出错: %w", key, err)
			}
			if isChanged {
				chgKeys = append(chgKeys, key)
				break
			}
		}
	}

	// 检查删除的记录
	for key := range mOld {
		if _, exists := mNew[key]; !exists {
			delKeys = append(delKeys, key)
		}
	}

	return addKeys, chgKeys, delKeys, nil
}

// 比较两个内部map是否相同
func compareInnerMaps(map1, map2 map[string]interface{}) (bool, error) {
	// 先检查长度
	if len(map1) != len(map2) {
		return true, nil
	}

	// 检查所有键值对
	for k, v1 := range map1 {
		v2, exists := map2[k]
		if !exists {
			return true, nil
		}

		// 比较值是否相同
		switch val1 := v1.(type) {
		// 字符串类型
		case string:
			val2, ok := v2.(string)
			if !ok {
				return false, fmt.Errorf("类型不匹配: 键 %s 期望 string, 得到 %T", k, v2)
			}
			if val1 != val2 {
				return true, nil
			}

		// 有符号整数类型
		case int:
			val2, ok := v2.(int)
			if !ok {
				return false, fmt.Errorf("类型不匹配: 键 %s 期望 int, 得到 %T", k, v2)
			}
			if val1 != val2 {
				return true, nil
			}

		case int8:
			val2, ok := v2.(int8)
			if !ok {
				return false, fmt.Errorf("类型不匹配: 键 %s 期望 int8, 得到 %T", k, v2)
			}
			if val1 != val2 {
				return true, nil
			}

		case int16:
			val2, ok := v2.(int16)
			if !ok {
				return false, fmt.Errorf("类型不匹配: 键 %s 期望 int16, 得到 %T", k, v2)
			}
			if val1 != val2 {
				return true, nil
			}

		case int32:
			val2, ok := v2.(int32)
			if !ok {
				return false, fmt.Errorf("类型不匹配: 键 %s 期望 int32, 得到 %T", k, v2)
			}
			if val1 != val2 {
				return true, nil
			}

		case int64:
			val2, ok := v2.(int64)
			if !ok {
				return false, fmt.Errorf("类型不匹配: 键 %s 期望 int64, 得到 %T", k, v2)
			}
			if val1 != val2 {
				return true, nil
			}

		// 无符号整数类型
		case uint:
			val2, ok := v2.(uint)
			if !ok {
				return false, fmt.Errorf("类型不匹配: 键 %s 期望 uint, 得到 %T", k, v2)
			}
			if val1 != val2 {
				return true, nil
			}

		case uint8:
			val2, ok := v2.(uint8)
			if !ok {
				return false, fmt.Errorf("类型不匹配: 键 %s 期望 uint8, 得到 %T", k, v2)
			}
			if val1 != val2 {
				return true, nil
			}

		case uint16:
			val2, ok := v2.(uint16)
			if !ok {
				return false, fmt.Errorf("类型不匹配: 键 %s 期望 uint16, 得到 %T", k, v2)
			}
			if val1 != val2 {
				return true, nil
			}

		case uint32:
			val2, ok := v2.(uint32)
			if !ok {
				return false, fmt.Errorf("类型不匹配: 键 %s 期望 uint32, 得到 %T", k, v2)
			}
			if val1 != val2 {
				return true, nil
			}

		case uint64:
			val2, ok := v2.(uint64)
			if !ok {
				return false, fmt.Errorf("类型不匹配: 键 %s 期望 uint64, 得到 %T", k, v2)
			}
			if val1 != val2 {
				return true, nil
			}

		// 浮点数类型
		case float32:
			val2, ok := v2.(float32)
			if !ok {
				return false, fmt.Errorf("类型不匹配: 键 %s 期望 float32, 得到 %T", k, v2)
			}
			if val1 != val2 {
				return true, nil
			}

		case float64:
			val2, ok := v2.(float64)
			if !ok {
				return false, fmt.Errorf("类型不匹配: 键 %s 期望 float64, 得到 %T", k, v2)
			}
			if val1 != val2 {
				return true, nil
			}

		// 布尔类型
		case bool:
			val2, ok := v2.(bool)
			if !ok {
				return false, fmt.Errorf("类型不匹配: 键 %s 期望 bool, 得到 %T", k, v2)
			}
			if val1 != val2 {
				return true, nil
			}

		// 空值
		case nil:
			if v2 != nil {
				return true, nil
			}

		default:
			return false, fmt.Errorf("不支持的简单类型: 键 %s 类型 %T", k, v1)
		}
	}

	return false, nil
}

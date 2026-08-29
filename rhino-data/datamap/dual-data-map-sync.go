package datamap

import (
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/utils/fuzzy"
	"rhino-core/schema"
	"sort"
)

const (
	pageSize         = 1000
	maxSyncDataRetry = 100
)

func (m *DualDataMap) shouldSync(dataSyncLog *schema.DataSyncLog) bool {
	if m.dataSyncLog == nil {
		return true
	}

	if dataSyncLog.ID > m.dataSyncLog.ID || dataSyncLog.ReportTime > m.dataSyncLog.ReportTime {
		return true
	}

	return false
}

func (m *DualDataMap) SyncData(dataSyncLog *schema.DataSyncLog) (synced bool, addKeys []string, chgKeys []string, delKeys []string) {

	log.Printf("check if shouldSync for %s\n", dataSyncLog.TableName)

	if !m.shouldSync(dataSyncLog) {
		return
	}

	log.Printf("===>start1 SyncData: %s\n", dataSyncLog.TableName)

	m.lock.Lock()
	defer m.lock.Unlock()

	log.Printf("===>start2 SyncData: %s\n", dataSyncLog.TableName)

	value, _, err := CreateMapFromDB(m.db, pageSize, m.tableConfig)

	log.Println("finish CreateMapFromDB!")

	i := 1
	for err != nil && i < maxSyncDataRetry {
		domain_error.ProcessSevereError(false, 5, nil, err, fmt.Sprintf("fail to syncData for table: %s", m.tableConfig.TableAlias))
		value, _, err = CreateMapFromDB(m.db, pageSize, m.tableConfig)

		i++
	}

	log.Println("star to compare two map")

	// 两个map，才能安全
	newWriteMap := make(map[string][]map[string]interface{})
	newWriteMap2 := make(map[string][]map[string]interface{})

	for k, v := range value {
		if len(v) > 0 {
			newWriteMap[k] = v
			newWriteMap2[k] = v
		}
	}

	log.Printf("newWriteMap.Len=%d\n", len(newWriteMap))
	log.Printf("newWriteMap2.Len=%d\n", len(newWriteMap2))
	log.Printf("m.writeMap.Len=%d\n", len(m.writeMap))

	log.Println("star to check diff of two map")

	// 比较两个map，那些key的值发生了变化
	addKeys, chgKeys, delKeys, err = CheckDiff(newWriteMap, m.writeMap)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to syncData for table %s while compair two table diff", m.tableConfig.TableAlias))
	}

	log.Printf("addKeys:%d, chgKeys:%d, delKeys:%d\n", len(addKeys), len(chgKeys), len(delKeys))

	// 原子更新
	m.readMap.Store(&newWriteMap)
	m.writeMap = newWriteMap2

	// 判断是否启用fuzzyMap
	if m.enableFuzzyMap && len(m.writeMap)>0 {

		log.Println("star to enableFuzzyMap")

		keyList, valueList, err := getKeyAndValueBySorting(m.writeMap, m.tableConfig.SortRule.ColumnsForSorting, m.tableConfig.SortRule.SortDir)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to syncData for table %s while sorting the inner map data", m.tableConfig.TableAlias))
			return
		}
		log.Printf("keyList=%d, valueList=%d\n", len(keyList), len(valueList))
		
		fuzzyMap := fuzzy.NewFuzzyMap[map[string]interface{}](keyList, valueList, func(a map[string]interface{}, b map[string]interface{}) bool {
			for _, col := range m.tableConfig.ColumnsForUnique {
				if a[col] != b[col] {
					return false
				}
			}
			return true
		})

		m.fuzzyMap.Store(fuzzyMap)
	}

	log.Println("finish enableFuzzyMap")

	synced = true
	// 记录dataSyncLog
	m.dataSyncLog = dataSyncLog

	return
}

type kv struct {
	key string
	val map[string]interface{}
}

func getKeyAndValueBySorting(m map[string][]map[string]interface{}, keysForSorting []string, sortDir string) (keyList []string, valueList []map[string]interface{}, err error) {
	// 检查参数有效性
	if m == nil {
		return nil, nil, fmt.Errorf("input map is nil")
	}

	if len(keysForSorting) == 0 {
		return nil, nil, fmt.Errorf("keysForSorting is empty")
	}

	if sortDir != "ASC" && sortDir != "DESC" {
		return nil, nil, fmt.Errorf("sortDir must be 'asc' or 'desc', got: %s", sortDir)
	}

	// 检查排序键是否存在（至少在一个内层map中存在）
    keyExists := make(map[string]bool)
    for _, innerMaps := range m {
        for _, innerMap := range innerMaps {
            for _, sortKey := range keysForSorting {
                if _, exists := innerMap[sortKey]; exists {
                    keyExists[sortKey] = true
                }
            }
        }
    }
    for _, sortKey := range keysForSorting {
        if !keyExists[sortKey] {
            return nil, nil, fmt.Errorf("sort key '%s' not found in any inner map", sortKey)
        }
    }

	// 构建 kv 列表
	var kvList[]*kv
	for key, vals := range m {
		for _, val := range vals {
			kvList = append(kvList, &kv{key, val})
		}
	}

	// 排序逻辑
	sort.Slice(kvList, func(i, j int) bool {
		// 按指定的多个键进行排序
		for _, sortKey := range keysForSorting {
			valI, existsI := kvList[i].val[sortKey]
			valJ, existsJ := kvList[j].val[sortKey]

			// 处理键不存在的情况
			if !existsI && !existsJ {
				continue // 两个都不存在，比较下一个键
			}
			if !existsI {
				return sortDir == "ASC" // 如果升序，不存在的排在前面
			}
			if !existsJ {
				return sortDir != "ASC" // 如果升序，存在的排在前面
			}

			// 比较值
			compareResult := compareValues(valI, valJ)
			if compareResult != 0 {
				if sortDir == "ASC" {
					return compareResult < 0
				} else {
					return compareResult > 0
				}
			}
		}
		// 所有排序键都相等，按外层键排序
        if sortDir == "ASC" {
            return kvList[i].key < kvList[j].key 
        } else {
            return kvList[i].key > kvList[j].key
        }
	})

	// 构建结果
	keyList = make([]string, len(kvList))
	valueList = make([]map[string]interface{}, len(kvList))

	for i, kv := range kvList {
		keyList[i] = kv.key
		valueList[i] = kv.val
	}

	return keyList, valueList, nil
}

// compareValues 比较两个 interface{} 值
func compareValues(a, b interface{}) int {
	switch aVal := a.(type) {
	case int:
		if bVal, ok := b.(int); ok {
			return aVal - bVal
		}
	case int8:
		if bVal, ok := b.(int8); ok {
			return int(aVal - bVal)
		}
	case int16:
		if bVal, ok := b.(int16); ok {
			return int(aVal - bVal)
		}
	case int32:
		if bVal, ok := b.(int32); ok {
			return int(aVal - bVal)
		}
	case int64:
		if bVal, ok := b.(int64); ok {
			if aVal < bVal {
				return -1
			} else if aVal > bVal {
				return 1
			}
			return 0
		}
	case uint:
		if bVal, ok := b.(uint); ok {
			if aVal < bVal {
				return -1
			} else if aVal > bVal {
				return 1
			}
			return 0
		}
	case uint8:
		if bVal, ok := b.(uint8); ok {
			return int(aVal) - int(bVal)
		}
	case uint16:
		if bVal, ok := b.(uint16); ok {
			return int(aVal) - int(bVal)
		}
	case uint32:
		if bVal, ok := b.(uint32); ok {
			if aVal < bVal {
				return -1
			} else if aVal > bVal {
				return 1
			}
			return 0
		}
	case uint64:
		if bVal, ok := b.(uint64); ok {
			if aVal < bVal {
				return -1
			} else if aVal > bVal {
				return 1
			}
			return 0
		}
	case float32:
		if bVal, ok := b.(float32); ok {
			if aVal < bVal {
				return -1
			} else if aVal > bVal {
				return 1
			}
			return 0
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			if aVal < bVal {
				return -1
			} else if aVal > bVal {
				return 1
			}
			return 0
		}
	case string:
		if bVal, ok := b.(string); ok {
			if aVal < bVal {
				return -1
			} else if aVal > bVal {
				return 1
			}
			return 0
		}
	case bool:
		if bVal, ok := b.(bool); ok {
			if aVal == bVal {
				return 0
			} else if !aVal && bVal {
				return -1
			} else {
				return 1
			}
		}
	}

	// 如果类型不匹配或是不支持的类型，使用字符串比较
	strA := fmt.Sprintf("%v", a)
	strB := fmt.Sprintf("%v", b)
	if strA < strB {
		return -1
	} else if strA > strB {
		return 1
	}
	return 0
}

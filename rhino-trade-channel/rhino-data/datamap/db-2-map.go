package datamap

import (
	"database/sql"
	"fmt"
	"rhino-common/utils/dbutil"
	"sort"
	"strconv"
	"strings"
)

type TableConfigForMapping struct {
	*dbutil.TableConfig
	ColumnsForUnique  []string // 用列的别名
	ColumnsForMapKeys [][]string // 用列的别名                // 指定由哪些db列来构成map的key，列之间用英文中划线 - 来间隔。支持多个key。
	SortRule          *SortRuleForValueInSameKey // 同一key内多条记录的排序规则
}

type SortRuleForValueInSameKey struct {
	ColumnsForSorting []string // 需要排序的列
	SortDir           string   // 升序(ASC)或者降序(DESC)
}

type Row struct {
	Value map[string][]interface{} // map的键取ColumnAlias，value根据ColumnType转换，支持同一key多条记录
}

// CreateMapFromDB 从数据库中分页查出所有数据，然后构造一个map对象返回
func CreateMapFromDB(db *sql.DB, pageSize int, tableConfig *TableConfigForMapping) (map[string][]map[string]interface{}, map[string]*Row, error) {
	if db == nil {
		return nil, nil, fmt.Errorf("db connection is nil")
	}
	if tableConfig == nil || tableConfig.TableConfig == nil {
		return nil, nil, fmt.Errorf("tableConfig is nil")
	}
	if len(tableConfig.ColumnsForMapKeys) == 0 {
		return nil, nil, fmt.Errorf("ColumnsForMapKeys is empty")
	}
	if pageSize <= 0 {
		return nil, nil, fmt.Errorf("pageSize must be greater than 0")
	}

	// 验证排序规则
	if tableConfig.SortRule != nil {
		if tableConfig.SortRule.SortDir != "ASC" && tableConfig.SortRule.SortDir != "DESC" {
			return nil, nil, fmt.Errorf("SortDir must be ASC or DESC")
		}
		if len(tableConfig.SortRule.ColumnsForSorting) == 0 {
			return nil, nil, fmt.Errorf("ColumnsForSorting cannot be empty when SortRule is provided")
		}
	}

	// 临时存储所有数据，按多个key分组
	tempData := make(map[string][]map[string]interface{})
	offset := 0

	for {
		// 构建查询语句（不包含排序）
		query, err := buildQuery(tableConfig, pageSize, offset)
		if err != nil {
			return nil, nil, err
		}

		// 执行查询
		rows, err := db.Query(query)
		if err != nil {
			return nil, nil, fmt.Errorf("query failed: %v", err)
		}

		// 处理查询结果
		columnNames, err := rows.Columns()
		if err != nil {
			rows.Close()
			return nil, nil, fmt.Errorf("get columns failed: %v", err)
		}

		columnCount := len(columnNames)
		currentPageRowCount := 0

		for rows.Next() {
			currentPageRowCount++

			// 创建扫描值的切片
			scanValues := make([]interface{}, columnCount)
			for i := range scanValues {
				scanValues[i] = new(interface{})
			}

			// 扫描行数据
			if err := rows.Scan(scanValues...); err != nil {
				rows.Close()
				return nil, nil, fmt.Errorf("scan row failed: %v", err)
			}

			// 构建行数据map
			rowData, err := buildRowData(scanValues, columnNames, tableConfig.TableConfig.Columns)
			if err != nil {
				rows.Close()
				return nil, nil, err
			}

			// 为每个key配置构建map key并存储
			for _, keyColumns := range tableConfig.ColumnsForMapKeys {
				// 构建map key
				mapKey, err := buildMapKey(rowData, keyColumns)
				if err != nil {
					rows.Close()
					return nil, nil, err
				}

				// 存储到临时数据结构中
				tempData[mapKey] = append(tempData[mapKey], rowData)
			}
		}

		rows.Close()

		// 检查是否有错误
		if err := rows.Err(); err != nil {
			return nil, nil, fmt.Errorf("rows iteration error: %v", err)
		}

		// 如果当前页行数小于pageSize，说明已经是最后一页
		if currentPageRowCount < pageSize {
			break
		}

		offset += pageSize
	}

	// 对同一key内的数据进行排序
	if tableConfig.SortRule != nil {
		for _, rows := range tempData {
			if len(rows) > 1 {
				sort.Slice(rows, func(i, j int) bool {
					return compareRows(rows[i], rows[j], tableConfig.SortRule)
				})
			}
		}
	}

	// 构建最终结果
	row, err := buildFinalResult(tempData, tableConfig.TableConfig.Columns)
	return tempData, row, err
}

// buildQuery 构建分页查询语句（不包含排序）
func buildQuery(tableConfig *TableConfigForMapping, pageSize, offset int) (string, error) {
	if len(tableConfig.TableConfig.Columns) == 0 {
		return "", fmt.Errorf("no columns configured")
	}

	// 构建SELECT字段列表
	var selectFields []string
	for _, col := range tableConfig.TableConfig.Columns {
		field := fmt.Sprintf("`%s`", col.ColumnName)
		if col.ColumnAlias != "" && col.ColumnAlias != col.ColumnName {
			field = fmt.Sprintf("`%s` AS `%s`", col.ColumnName, col.ColumnAlias)
		}
		selectFields = append(selectFields, field)
	}

	selectClause := strings.Join(selectFields, ", ")
	fromClause := fmt.Sprintf("FROM `%s`", tableConfig.TableConfig.TableName)
	limitClause := fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, offset)

	return fmt.Sprintf("SELECT %s %s %s", selectClause, fromClause, limitClause), nil
}

// buildRowData 构建行数据map，根据ColumnType进行类型转换
func buildRowData(scanValues []interface{}, columnNames []string, columnConfigs []*dbutil.ColumnConfig) (map[string]interface{}, error) {
	rowData := make(map[string]interface{})
	columnConfigMap := make(map[string]*dbutil.ColumnConfig)

	// 创建列配置的映射，以列名为键
	for _, col := range columnConfigs {
		columnConfigMap[col.ColumnName] = col
		if col.ColumnAlias != "" {
			columnConfigMap[col.ColumnAlias] = col
		}
	}

	for i, colName := range columnNames {
		rawValue := *(scanValues[i].(*interface{}))

		// 查找列配置
		colConfig, exists := columnConfigMap[colName]
		if !exists {
			return nil, fmt.Errorf("column config not found for column: %s", colName)
		}

		// 根据列类型转换值
		convertedValue, err := convertValue(rawValue, colConfig)
		if err != nil {
			return nil, fmt.Errorf("convert value failed for column %s: %v", colName, err)
		}

		// 使用列别名作为map的key
		key := colConfig.ColumnAlias
		if key == "" {
			key = colConfig.ColumnName
		}

		rowData[key] = convertedValue
	}

	return rowData, nil
}

// convertValue 根据ColumnType转换值
func convertValue(rawValue interface{}, colConfig *dbutil.ColumnConfig) (interface{}, error) {
	if rawValue == nil {
		// 处理NULL值
		switch colConfig.ColumnType {
		case 0, 4: // 字符串和长文本
			return "", nil
		case 1, 5: // 整数和长整型
			return int64(0), nil
		case 2: // 浮点数
			return float64(0), nil
		case 3: // 布尔值
			return false, nil
		default:
			return nil, fmt.Errorf("unknown column type: %d", colConfig.ColumnType)
		}
	}

	switch colConfig.ColumnType {
	case 0, 4: // 字符串和长文本
		switch v := rawValue.(type) {
		case string:
			return v, nil
		case []byte:
			return string(v), nil
		default:
			return fmt.Sprintf("%v", v), nil
		}

	case 1, 5: // 整数和长整型
		switch v := rawValue.(type) {
		case int64:
			return v, nil
		case int32:
			return int64(v), nil
		case int:
			return int64(v), nil
		case float64:
			return int64(v), nil
		case float32:
			return int64(v), nil
		case []byte:
			strVal := string(v)
			intVal, err := strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("convert to int64 failed: %v", err)
			}
			return intVal, nil
		default:
			return nil, fmt.Errorf("unsupported type for integer: %T", v)
		}

	case 2: // 浮点数
		switch v := rawValue.(type) {
		case float64:
			return v, nil
		case float32:
			return float64(v), nil
		case int64:
			return float64(v), nil
		case int:
			return float64(v), nil
		case []byte:
			strVal := string(v)
			floatVal, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				return nil, fmt.Errorf("convert to float64 failed: %v", err)
			}
			return floatVal, nil
		default:
			return nil, fmt.Errorf("unsupported type for float: %T", v)
		}

	case 3: // 布尔值
		switch v := rawValue.(type) {
		case bool:
			return v, nil
		case int64:
			return v != 0, nil
		case int:
			return v != 0, nil
		case []byte:
			strVal := strings.ToLower(string(v))
			return strVal == "true" || strVal == "1" || strVal == "yes", nil
		default:
			return nil, fmt.Errorf("unsupported type for boolean: %T", v)
		}

	default:
		return nil, fmt.Errorf("unknown column type: %d", colConfig.ColumnType)
	}
}

// buildMapKey 构建map的key
func buildMapKey(rowData map[string]interface{}, keyColumns []string) (string, error) {
	var keyParts []string

	for _, colName := range keyColumns {
		value, exists := rowData[colName]
		if !exists {
			return "", fmt.Errorf("key column '%s' not found in row data", colName)
		}

		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case int64:
			strValue = strconv.FormatInt(v, 10)
		case float64:
			strValue = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			strValue = strconv.FormatBool(v)
		default:
			strValue = fmt.Sprintf("%v", v)
		}

		keyParts = append(keyParts, strValue)
	}

	return strings.Join(keyParts, "-"), nil
}

// buildFinalResult 构建最终结果
func buildFinalResult(tempData map[string][]map[string]interface{}, columnConfigs []*dbutil.ColumnConfig) (map[string]*Row, error) {
	resultMap := make(map[string]*Row)

	for key, rows := range tempData {
		// 转换为列切片形式
		columnSlices := make(map[string][]interface{})

		// 初始化所有列的切片
		for _, col := range columnConfigs {
			keyName := col.ColumnAlias
			if keyName == "" {
				keyName = col.ColumnName
			}
			columnSlices[keyName] = make([]interface{}, 0, len(rows))
		}

		// 填充数据
		for _, row := range rows {
			for colName, value := range row {
				columnSlices[colName] = append(columnSlices[colName], value)
			}
		}

		resultMap[key] = &Row{Value: columnSlices}
	}

	return resultMap, nil
}

// compareRows 比较两行数据，用于排序
func compareRows(row1, row2 map[string]interface{}, sortRule *SortRuleForValueInSameKey) bool {
	for _, colName := range sortRule.ColumnsForSorting {
		val1, exists1 := row1[colName]
		val2, exists2 := row2[colName]

		// 如果某一行的该列不存在，则将其视为较小值
		if !exists1 && !exists2 {
			continue
		}
		if !exists1 {
			return sortRule.SortDir == "ASC"
		}
		if !exists2 {
			return sortRule.SortDir == "DESC"
		}

		// 比较值
		comparison := compareValues(val1, val2)
		if comparison != 0 {
			if sortRule.SortDir == "ASC" {
				return comparison < 0
			} else {
				return comparison > 0
			}
		}
	}

	// 所有排序列都相等
	return false
}
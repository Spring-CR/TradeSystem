package data_sync

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// convertValue 尝试将字符串转换为最合适的类型
func convertValue(s, valType string) interface{} {
	// 去除前后空格
	trimmed := strings.TrimSpace(s)

	// 空字符串直接返回
	if trimmed == "" {
		switch valType {
		case "int":
			return 0
		case "float":
			return 0.0
		case "bool":
			return false
		}
		return trimmed
	}

	switch valType {
	case "int":
		// 首先尝试转换为整数
		if intValue, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
			return intValue
		}
	case "float":
		// 尝试转换为浮点数
		if floatValue, err := strconv.ParseFloat(trimmed, 64); err == nil {
			return floatValue
		}
	case "bool":
		// 布尔值检测
		if strings.ToLower(trimmed) == "true" || strings.ToLower(trimmed) == "false" {
			if boolValue, err := strconv.ParseBool(trimmed); err == nil {
				return boolValue
			}
		}
	}

	// 都不成功则返回字符串
	return trimmed
}

// readCSVColumnByName 通过列名读取CSV列并进行类型转换
func readCSVColumnByName(filename, columnName, valType string) ([]interface{}, error) {
	// 打开CSV文件
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 创建CSV reader
	reader := csv.NewReader(file)

	// 读取表头行
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("读取表头失败: %v", err)
	}

	// 查找目标列名的索引
	columnIndex := -1
	for i, header := range headers {
		if strings.TrimSpace(strings.ToLower(header)) == strings.ToLower(strings.TrimSpace(columnName)) {
			columnIndex = i
			break
		}
	}

	if columnIndex == -1 {
		availableColumns := strings.Join(headers, ", ")
		return nil, fmt.Errorf("未找到列名: '%s'，可用列名: %s", columnName, availableColumns)
	}

	var results []interface{}

	for {
		// 逐行读取数据
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("读取CSV数据错误: %v", err)
		}

		// 检查列索引是否有效
		if columnIndex >= len(record) {
			return nil, fmt.Errorf("列索引 %d 超出记录范围(0-%d)", columnIndex, len(record)-1)
		}

		// 转换并添加值
		value := convertValue(record[columnIndex], valType)
		results = append(results, value)
	}

	return results, nil
}

func extractVariables(template string) []string {
	// 定义正则表达式，匹配${xxx}模式
	re := regexp.MustCompile(`\$\{([^}]+)\}`)

	// 查找所有匹配项
	matches := re.FindAllStringSubmatch(template, -1)

	// 提取捕获组的内容
	var results []string
	for _, match := range matches {
		if len(match) > 1 { // match[0]是整个匹配，match[1]是第一个捕获组
			results = append(results, match[1])
		}
	}

	return results
}

func refineParam(db *sql.DB, rawParam string) (string, bool, error) {

	variables := extractVariables(rawParam)
	if len(variables) == 0 {
		return rawParam, false, nil
	}
	for _, variable := range variables {
		if strings.HasPrefix(variable, "csv.") {
			val, err := parseCsvColumnValues(variable[4:])
			if err != nil {
				return rawParam, false, err
			}
			rawParam = strings.ReplaceAll(rawParam, "\"${"+variable+"}\"", val)
		} else if strings.HasPrefix(variable, "db.") {
			val, err := parseDbColumnValues(db, variable[4:])
			if err != nil {
				return rawParam, false, err
			}
			rawParam = strings.ReplaceAll(rawParam, "\"${"+variable+"}\"", val)
		}
	}

	return rawParam, true, nil
}

func parseCsvColumnValues(variable string) (string, error) {
	valType := "string"
	var csvName, columnName string
	strs := strings.Split(variable, ".")
	if len(strs) == 4 {
		valType = strs[1]
		csvName = strs[2]
		columnName = strs[3]
	} else if len(strs) == 3 {
		csvName = strs[1]
		columnName = strs[2]
	} else {
		return "", fmt.Errorf("illegal variable:%s", variable)
	}
	valList, err := readCSVColumnByName("/tmp/"+csvName+".csv", columnName, valType)
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(valList)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func parseDbColumnValues(db *sql.DB, variable string) (string, error) {
	valType := "string"
	var tbName, columnName string
	strs := strings.Split(variable, ".")
	if len(strs) == 4 {
		valType = strs[1]
		tbName = strs[2]
		columnName = strs[3]
	} else if len(strs) == 3 {
		tbName = strs[1]
		columnName = strs[2]
	} else {
		return "", fmt.Errorf("illegal variable:%s", variable)
	}
	valList, err := readDbColumnByName(db, tbName, columnName, valType)
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(valList)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func readDbColumnByName(db *sql.DB, tableName, columnName, valType string) ([]interface{}, error) {
	// 验证输入参数
	if db == nil {
		return nil, fmt.Errorf("db connection is nil")
	}
	
	tableName = strings.TrimSpace(tableName)
	columnName = strings.TrimSpace(columnName)
	valType = strings.TrimSpace(strings.ToLower(valType))
	
	if tableName == "" || columnName == "" || valType == "" {
		return nil, fmt.Errorf("tableName, columnName and valType cannot be empty")
	}
	
	// 验证valType
	validTypes := map[string]bool{
		"int":    true,
		"float":  true,
		"bool":   true,
		"string": true,
	}
	if !validTypes[valType] {
		return nil, fmt.Errorf("valType must be one of: int, float, bool, string")
	}

	// 构建查询语句
	query := fmt.Sprintf("SELECT %s FROM %s", columnName, tableName)
	
	// 执行查询
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var results []interface{}

	// 根据不同类型处理数据
	switch valType {
	case "int":
		for rows.Next() {
			var value sql.NullInt64
			if err := rows.Scan(&value); err != nil {
				return nil, fmt.Errorf("failed to scan row: %v", err)
			}
			
			if value.Valid {
				results = append(results, int(value.Int64))
			} else {
				results = append(results, 0) // int类型的零值
			}
		}
		
	case "float":
		for rows.Next() {
			var value sql.NullFloat64
			if err := rows.Scan(&value); err != nil {
				return nil, fmt.Errorf("failed to scan row: %v", err)
			}
			
			if value.Valid {
				results = append(results, value.Float64)
			} else {
				results = append(results, 0.0) // float类型的零值
			}
		}
		
	case "bool":
		for rows.Next() {
			var value sql.NullBool
			if err := rows.Scan(&value); err != nil {
				return nil, fmt.Errorf("failed to scan row: %v", err)
			}
			
			if value.Valid {
				results = append(results, value.Bool)
			} else {
				results = append(results, false) // bool类型的零值
			}
		}
		
	case "string":
		for rows.Next() {
			var value sql.NullString
			if err := rows.Scan(&value); err != nil {
				return nil, fmt.Errorf("failed to scan row: %v", err)
			}
			
			if value.Valid {
				results = append(results, value.String)
			} else {
				results = append(results, "") // string类型的零值
			}
		}
	}

	// 检查遍历过程中是否有错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return results, nil
}
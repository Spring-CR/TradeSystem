package utils

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"rhino-common/utils/dbutil"
	"strconv"
	"strings"
)

func ConvertCsvToMap(csvBytes []byte, tableConfig *dbutil.TableConfig) (lines []map[string]interface{}, err error) {
	log.Printf("ConvertCsvToMap, TableAlias=%s\n", tableConfig.TableAlias)
	if tableConfig == nil || len(tableConfig.Columns) == 0 {
		return nil, errors.New("table config is empty")
	}

	// 创建CSV reader
	reader := csv.NewReader(strings.NewReader(string(csvBytes)))
	
	// 读取所有CSV记录
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return []map[string]interface{}{}, nil
	}

	// 第一行是表头
	headers := records[0]
	
	// 创建表头索引映射
	headerIndexMap := make(map[string]int)
	for i, header := range headers {
		headerIndexMap[header] = i
	}

	// 验证所有配置的列在CSV中是否存在
	for _, col := range tableConfig.Columns {
		if _, exists := headerIndexMap[col.ColumnName]; !exists {
			log.Printf("column defined for %s is as following:\n", tableConfig.TableAlias)
			for _, col := range tableConfig.Columns {
				log.Printf("===>%s\n", col.ColumnName)
			} 
			log.Println("headers in CSV is as following:")
			for _, header := range headers {
				log.Printf("===>%s\n", header)
			}
			return nil, errors.New("column not found in CSV: " + col.ColumnName)
		}
	}

	lines = make([]map[string]interface{}, 0, len(records)-1)

	// 从第二行开始处理数据（跳过表头）
	for rowIdx := 1; rowIdx < len(records); rowIdx++ {
		record := records[rowIdx]
		lineMap := make(map[string]interface{})
		
		// 处理每个配置的列
		for _, columnConfig := range tableConfig.Columns {

			// 根据列名找到在CSV中的索引位置
			columnIndex := headerIndexMap[columnConfig.ColumnName]
			
			if columnIndex >= len(record) {
				return nil, errors.New("record column count doesn't match header")
			}

			cellValue := record[columnIndex]
			var convertedValue interface{}
			var conversionErr error

			// 根据列类型进行转换
			switch columnConfig.ColumnType {
			case 0, 4: // 字符串类型和长文本类型
				convertedValue = cellValue
				
			case 1, 5: // 整数类型和长整型
				if cellValue == "" {
					convertedValue = int64(0)
				} else {
					convertedValue, conversionErr = strconv.ParseInt(cellValue, 10, 64)
					if conversionErr != nil {
						convertedValue, conversionErr = strconv.ParseFloat(cellValue, 64)
						if conversionErr == nil {
							convertedValue = int64(convertedValue.(float64))
						}
					}
				}
				
			case 2: // 浮点数类型
				if cellValue == "" {
					convertedValue = float64(0)
				} else {
					convertedValue, conversionErr = strconv.ParseFloat(cellValue, 64)
				}
				
			case 3: // 布尔值类型
				if cellValue == "" {
					convertedValue = false
				} else {
					convertedValue, conversionErr = parseBool(cellValue)
				}
				
			default:
				return nil, errors.New("unsupported column type: " + strconv.Itoa(columnConfig.ColumnType))
			}

			if conversionErr != nil {
				log.Printf("fail to get column value for %s\n", columnConfig.ColumnName)
				return nil, conversionErr
			}

			lineMap[columnConfig.ColumnName] = convertedValue
		}
		
		lines = append(lines, lineMap)
	}

	return lines, nil
}

// 辅助函数：解析布尔值
func parseBool(s string) (bool, error) {
	lower := strings.ToLower(strings.TrimSpace(s))
	switch lower {
	case "true", "t", "1", "yes", "y":
		return true, nil
	case "false", "f", "0", "no", "n":
		return false, nil
	default:
		return false, errors.New("invalid boolean value: " + s)
	}
}

func ConvertMapToCsv(lines []map[string]interface{}, tableConfig *dbutil.TableConfig) (csvBytes []byte, err error) {
	if tableConfig == nil || len(tableConfig.Columns) == 0 {
		return nil, errors.New("table config is empty")
	}

	// 创建 CSV writer
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// 写入表头
	headers := make([]string, len(tableConfig.Columns))
	for i, col := range tableConfig.Columns {
		headers[i] = col.ColumnName
	}
	
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// 创建列配置索引映射，用于快速查找列类型
	columnTypeMap := make(map[string]int)
	for _, col := range tableConfig.Columns {
		columnTypeMap[col.ColumnName] = col.ColumnType
	}

	// 处理每一行数据
	for _, line := range lines {
		record := make([]string, len(tableConfig.Columns))
		
		for i, col := range tableConfig.Columns {
			columnName := col.ColumnName
			value, exists := line[columnName]
			
			if !exists {
				// 如果该列不存在，使用空字符串
				record[i] = ""
				continue
			}

			// 根据列类型转换为字符串
			strValue, err := convertValueToString(value, col.ColumnType)
			if err != nil {
				return nil, fmt.Errorf("error converting column %s: %v", columnName, err)
			}
			
			record[i] = strValue
		}
		
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	// 刷新缓冲区并返回字节
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 辅助函数：将值转换为字符串
func convertValueToString(value interface{}, columnType int) (string, error) {
	if value == nil {
		return "", nil
	}

	switch columnType {
	case 0, 4: // 字符串类型和长文本类型
		switch v := value.(type) {
		case string:
			return v, nil
		default:
			return fmt.Sprintf("%v", value), nil
		}
		
	case 1, 5: // 整数类型和长整型
		switch v := value.(type) {
		case int64:
			return strconv.FormatInt(v, 10), nil
		case int:
			return strconv.Itoa(v), nil
		case int32:
			return strconv.FormatInt(int64(v), 10), nil
		default:
			return "", fmt.Errorf("expected integer type for column, got %T", value)
		}
		
	case 2: // 浮点数类型
		switch v := value.(type) {
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64), nil
		case float32:
			return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
		default:
			return "", fmt.Errorf("expected float type for column, got %T", value)
		}
		
	case 3: // 布尔值类型
		switch v := value.(type) {
		case bool:
			if v {
				return "true", nil
			}
			return "false", nil
		default:
			return "", fmt.Errorf("expected boolean type for column, got %T", value)
		}
		
	default:
		return "", fmt.Errorf("unsupported column type: %d", columnType)
	}
}
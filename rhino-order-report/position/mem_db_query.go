package position

import (
	"log"
	"math"
	"rhino-common/domain_error"
	"rhino-common/utils/dbutil"
)

func (s *MemPosition) QueryPosition(query *dbutil.StructuralQuery) ([]map[string]interface{}, int, *domain_error.Error) {

	query.Table = s.positionTableName
	
	if len(query.SortFields) == 0 {
		query.SortFields = []string{"symbol"}
		query.SortType = 1
	}

	return s.doQuery(query)
}


func (s *MemPosition) doQuery(query *dbutil.StructuralQuery) ([]map[string]interface{}, int, *domain_error.Error) {

	db := s.dbRead

	sql, args, err := query.GenerateQueryStatement()
	if err != nil {
		return nil, 0, domain_error.Build(domain_error.DATABASE_QUERY_ERROR, err, sql, args)
	}

	log.Printf("query sql:%s\n", sql)
	log.Printf("args:%v\n", args)

	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, 0, domain_error.Build(domain_error.DATABASE_QUERY_ERROR, err, sql, args)
	}
	defer rows.Close()

	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		return nil, 0, domain_error.Build(domain_error.DATABASE_QUERY_ERROR, err, sql, args)
	}

	// 准备一个切片，用来存储每一行的数据
	result := []map[string]interface{}{}

	// 生成一个切片，存储每行数据的指针
	values := make([]interface{}, len(columns))

	// 整理结果
	for i := range values {
		values[i] = new(interface{})
	}
	// 迭代每一行
	for rows.Next() {
		// 扫描每行数据
		err := rows.Scan(values...)
		if err != nil {
			return nil, 0, domain_error.Build(domain_error.DATABASE_QUERY_ERROR, err, sql, args)
		}

		// 创建一个切片来存储每一行的值
		row := make(map[string]interface{}, len(columns))

		// 处理每列数据，如果为空则设置为零值
		for i, colName := range columns {
			val := *(values[i].(*interface{}))
			if val == nil {
				// 如果值为nil，给定字段类型的零值
				row[colName] = zeroValueForColumnType(colName)
			} else {
				row[colName] = tryToAdjFloatValue(val)
			}
		}

		log.Printf("===> add row:%v\n", row)

		// 将处理后的行数据加入到结果中
		result = append(result, row)
	}
	

	// 检查迭代是否发生错误
	if err := rows.Err(); err != nil {
		return nil, 0, domain_error.Build(domain_error.DATABASE_QUERY_ERROR, err, sql, args)
	}

	sql, args = query.GenerateQueryCountStatement(sql, args)

	var count int
	row := db.QueryRow(sql, args...)
	err = row.Scan(&count)

	if err != nil && dbutil.IsDbRecordEmptyError(err){
		err = nil
	}

	if err != nil {
		return nil, 0, domain_error.Build(domain_error.DATABASE_QUERY_ERROR, err, sql, args)
	}

	return result, count, nil
}

// 新增函数：格式化浮点数值
func tryToAdjFloatValue(val interface{}) interface{} {
    switch v := val.(type) {
    case float64:
        // 四舍五入到4位小数
        return math.Round(v*10000) / 10000
    case float32:
        // 将float32转换为float64进行处理，然后再转回float32
        return float32(math.Round(float64(v)*10000) / 10000)
    }
	return val
}

// 根据字段名称返回字段类型的零值
func zeroValueForColumnType(colName string) interface{} {
	// 根据列名判断类型，给出零值
	// 这里假设你知道字段的类型，也可以扩展代码来查询表结构
	switch colName {
	case "int_field": // 如果字段类型是整数
		return 0
	case "text_field": // 如果字段类型是文本
		return ""
	case "bool_field": // 如果字段类型是布尔
		return false
	default:
		return nil
	}
}

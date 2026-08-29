package order_status

import (
	sqltype "database/sql"
	"log"
	"rhino-common/domain_error"
	"rhino-common/utils/dbutil"
	"strings"
)

func (s *OrderStatusReplica) QueryTradeOrder(query *dbutil.StructuralQuery, historical bool) ([]map[string]interface{}, int, *domain_error.Error) {

	if historical {
		query.Table = s.hisOrderTableName
	} else {
		query.Table = s.orderTableName
	}
	
	if len(query.SortFields) == 0 {
		query.SortFields = []string{"f_ord_status_update_time"}
		query.SortType = 1
	}

	return s.doQuery(query, historical)
}

func (s *OrderStatusReplica) QueryTradeActionResp(query *dbutil.StructuralQuery, historical bool) ([]map[string]interface{}, int, *domain_error.Error) {
	
	if historical {
		query.Table = s.hisOrderRespTableName
	} else {
		query.Table = s.orderRespTableName
	}

	if len(query.SortFields) == 0 {
		query.SortFields = []string{"f_msg_time"}
		query.SortType = 1
	}

	return s.doQuery(query, historical)
}

func (s *OrderStatusReplica) doQuery(query *dbutil.StructuralQuery, historical bool) ([]map[string]interface{}, int, *domain_error.Error) {

	db := s.dbRead

	if historical {
		db = s.applicationCfg.GetCentralDB()
	}

	sql, args, err := query.GenerateQueryStatement()
	if err != nil {
		return nil, 0, domain_error.Build(domain_error.DATABASE_QUERY_ERROR, err, sql, args)
	}

	rows, err := db.Query(sql, args...)
	if err != nil {

		if strings.Contains(err.Error(), "no such table") {
			s.printDBTables()
		}

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


	if !historical {
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
					row[colName] = val
				}
			}
	
			// 将处理后的行数据加入到结果中
			result = append(result, row)
		}
	} else {
		// 获取列类型信息
		colTypes, err := rows.ColumnTypes()
		if err != nil {
			return nil, 0, domain_error.Build(domain_error.DATABASE_QUERY_ERROR, err, sql, args)
		}

		for i, ct := range colTypes {
			switch ct.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "CHAR":
				var s sqltype.NullString
				values[i] = &s
			case "INT", "BIGINT":
				var n sqltype.NullInt64
				values[i] = &n
			case "DOUBLE", "FLOAT", "DECIMAL":
				var f sqltype.NullFloat64
				values[i] = &f
			case "DATETIME", "TIMESTAMP":
				var t  sqltype.NullTime
				values[i] = &t
			case "TINYINT", "BOOLEAN":
				var b sqltype.NullBool
				values[i] = &b
			default:
				log.Printf("unknow type name:%s\n", ct.DatabaseTypeName())
				var b []byte  // 未知类型用 []byte 接收
				values[i] = &b
			}
		}
	
		for rows.Next() {
			err := rows.Scan(values...)
			if err != nil {
				return nil, 0, domain_error.Build(domain_error.DATABASE_QUERY_ERROR, err, sql, args)
			}
	
			row := make(map[string]interface{})
			for i, colName := range columns {
				ct := colTypes[i]
				switch ct.DatabaseTypeName() {
				case "VARCHAR", "TEXT", "CHAR":
					ns := values[i].(*sqltype.NullString)
                    if ns.Valid {
                        row[colName] = ns.String
                    } else {
                        row[colName] = ""
                    }
				case "INT", "BIGINT":
					ni := values[i].(*sqltype.NullInt64)
                    if ni.Valid {
                        row[colName] = ni.Int64
                    } else {
                        row[colName] = 0
                    }
				case "DOUBLE", "FLOAT", "DECIMAL":
					nf := values[i].(*sqltype.NullFloat64)
                    if nf.Valid {
                        row[colName] = nf.Float64
                    } else {
                        row[colName] = 0.0
                    }
				case "DATETIME", "TIMESTAMP":
					nt := values[i].(*sqltype.NullTime)
                    if nt.Valid {
                        row[colName] = nt.Time
                    } else {
                        row[colName] = nil
                    }
				case "TINYINT", "BOOLEAN":
					nb := values[i].(*sqltype.NullBool)
                    if nb.Valid {
                        row[colName] = nb.Bool
                    } else {
                        row[colName] = false
                    }
				default:
					row[colName] = values[i]
				}
			}
			result = append(result, row)
		}
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

	// jsData, _ := json.MarshalIndent(result, "", "  ")
	// log.Printf("query result:%s\n", jsData)


	return result, count, nil
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

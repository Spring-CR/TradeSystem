package dbutil

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type TableConfig struct {
	TableName    string          // 表名
	TableAlias   string          // 表别名，创建数据库时用不到
	CharacterSet string          // 字符集
	Collate      string          // 排序规则
	Columns      []*ColumnConfig // 列的配置
	Indexes      []*IndexConfig  // 索引配置
}

type ColumnConfig struct {
	ColumnName  string // 列名
	ColumnAlias string // 列别名，创建数据库时用不到
	ColumnType  int    // 列类型，0：字符串（VARCHAR），1：整数（INTEGER），2：浮点数（DOUBLE），3：布尔值（BOOLEAN），4：长文本类型（MEDIUMTEXT），5：长整型（BIGINT）
	ColumnLen   int    // 列长度，字符串类型时（即当ColumnType=0）有效
}

type IndexConfig struct {
	IndexType int      // 0：普通索引，1：唯一索引
	Columns   []string // 索引的列名
	IndexName string   // 索引名称
}

/*
		功能：从csv内容中创建数据库表
		返回：当创建成功时，返回nil；当创建失败时，返回错误信息
		逻辑：
		1. 创建临时表，表名用TableConfig定义的表名+_tmp_+毫秒级时间戳
		2. 根据表配置创建索引
	    3. 建表的方式是先构造建表语句，然后执行建表语句
		4. 建表完成后，根据csv内容插入临时表
		5. 检查TableConfig定义的表是否存在
		6. 如果存在，则删除表，然后重命名临时表为TableConfig定义的表名
		7. 如果不存在，则重命名临时表为TableConfig定义的表名
		8. 以上过程，如果任意步骤出现异常，则立即返回error同时要尝试删除临时表
		9. 如果没有异常，则返回nil
		参数：
		- db: 数据库连接
		- csvContent: csv内容
		- tableConfig: 表配置
*/
func CreateTableFromCsv(db *sql.DB, csvContent []byte, tableConfig *TableConfig, checkReplaceTable ...func() bool) error {
	// 1. 生成临时表名
	tempTableName := generateTempTableName(tableConfig.TableName)

	// 确保在函数退出时清理临时表（如果发生错误）
	defer func() {
		if err := recover(); err != nil {
			// 发生panic时清理临时表
			cleanupTempTable(db, tempTableName)
		}
	}()

	// 2. 创建临时表
	err := createTempTableWithoutIndex(db, tempTableName, tableConfig)
	if err != nil {
		cleanupTempTable(db, tempTableName)
		return fmt.Errorf("创建临时表失败: %v", err)
	}

	// 3. 解析CSV内容并插入数据
	err = insertCsvData(db, tempTableName, csvContent, tableConfig)
	if err != nil {
		cleanupTempTable(db, tempTableName)
		return fmt.Errorf("插入CSV数据失败: %v", err)
	}

	// 4. 创建索引
	err = createIndexes(db, tempTableName, tableConfig.Indexes)
	if err != nil {
		return fmt.Errorf("创建索引失败: %v", err)
	}

	// 5. 检查原表是否存在并处理
	targetTableExists, err := checkTableExists(db, tableConfig.TableName)
	if err != nil {
		cleanupTempTable(db, tempTableName)
		return fmt.Errorf("检查表是否存在失败: %v", err)
	}

	// 6. 执行表替换操作
	if len(checkReplaceTable) == 0 || checkReplaceTable[0]() {
		log.Printf("for replace table: %s\n", strings.Split(tempTableName, "_tmp_")[0])
		err = replaceTable(db, tableConfig.TableName, tempTableName, targetTableExists)
		if err != nil {
			cleanupTempTable(db, tempTableName)
			return fmt.Errorf("表替换操作失败: %v", err)
		}
	} else {
		cleanupTempTable(db, tempTableName)
		log.Printf("for skip replace table: %s\n", strings.Split(tempTableName, "_tmp_")[0])
	}

	return nil
}

// 生成临时表名
func generateTempTableName(tableName string) string {
	timestamp := time.Now().UnixNano() / int64(time.Microsecond)
	return fmt.Sprintf("%s_tmp_%d", tableName, timestamp)
}

// 创建临时表
func createTempTableWithoutIndex(db *sql.DB, tempTableName string, tableConfig *TableConfig) error {
	// 构建建表SQL
	createTableSQL := buildCreateTableSQL(tempTableName, tableConfig)

	log.Printf("for table: %s, createTableSQL: \n%s\n", strings.Split(tempTableName, "_tmp_")[0], createTableSQL)

	// 执行建表语句
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("执行建表语句失败: %v, SQL: %s", err, createTableSQL)
	}

	return nil
}

// 构建建表SQL
func buildCreateTableSQL(tableName string, tableConfig *TableConfig) string {
	var sqlBuilder strings.Builder

	// 开始构建CREATE TABLE语句
	sqlBuilder.WriteString("CREATE TABLE `")
	sqlBuilder.WriteString(tableName)
	sqlBuilder.WriteString("` (")
	sqlBuilder.WriteString("\n")

	// 添加列定义
	for i, column := range tableConfig.Columns {
		if i > 0 {
			sqlBuilder.WriteString(", ")
		}
		sqlBuilder.WriteString("`")
		sqlBuilder.WriteString(column.ColumnName)
		sqlBuilder.WriteString("` ")
		sqlBuilder.WriteString(getColumnTypeSQL(column))
		sqlBuilder.WriteString("\n")
		// 添加NOT NULL约束（根据实际需求，这里假设所有列都不允许NULL）
		// sqlBuilder.WriteString(" NOT NULL")
	}

	// 添加字符集和排序规则
	if tableConfig.CharacterSet != "" {
		sqlBuilder.WriteString(") CHARACTER SET ")
		sqlBuilder.WriteString(tableConfig.CharacterSet)
	} else {
		sqlBuilder.WriteString(")")
	}

	if tableConfig.Collate != "" {
		sqlBuilder.WriteString(" COLLATE ")
		sqlBuilder.WriteString(tableConfig.Collate)
	}

	sqlBuilder.WriteString("\n")

	return sqlBuilder.String()
}

// 获取列类型的SQL表示
func getColumnTypeSQL(column *ColumnConfig) string {
	switch column.ColumnType {
	case 0: // 字符串（VARCHAR）
		if column.ColumnLen <= 0 {
			return "VARCHAR(255)"
		}
		return fmt.Sprintf("VARCHAR(%d)", column.ColumnLen)
	case 1: // 整数（INTEGER）
		return "INT"
	case 2: // 浮点数（DOUBLE）
		return "DOUBLE"
	case 3: // 布尔值（BOOLEAN）
		return "BOOLEAN"
	case 4: // 长文本类型（MEDIUMTEXT）
		return "MEDIUMTEXT"
	case 5: // 长整型（BIGINT）
		return "BIGINT"
	default:
		return "VARCHAR(255)" // 默认类型
	}
}

// 创建索引
func createIndexes(db *sql.DB, tableName string, indexes []*IndexConfig) error {
	for _, index := range indexes {
		indexSQL := buildIndexSQL(tableName, index)
		log.Printf("for table: %s, indexSQL: \n%s\n", strings.Split(tableName, "_tmp_")[0], indexSQL)
		_, err := db.Exec(indexSQL)
		if err != nil {
			return err
		}
	}
	return nil
}

// 构建索引SQL
func buildIndexSQL(tableName string, index *IndexConfig) string {
	var sqlBuilder strings.Builder

	// 确定索引类型
	indexType := "INDEX"
	if index.IndexType == 1 {
		indexType = "UNIQUE INDEX"
	}

	// 构建索引名称（使用列名组合）
	indexName := "idx_" + strings.Join(index.Columns, "_")
	if index.IndexName != "" {
		indexName = index.IndexName
	}

	sqlBuilder.WriteString("CREATE ")
	sqlBuilder.WriteString(indexType)
	sqlBuilder.WriteString(" `")
	sqlBuilder.WriteString(indexName)
	sqlBuilder.WriteString("` ON `")
	sqlBuilder.WriteString(tableName)
	sqlBuilder.WriteString("` (")

	for i, column := range index.Columns {
		if i > 0 {
			sqlBuilder.WriteString(", ")
		}
		sqlBuilder.WriteString("`")
		sqlBuilder.WriteString(column)
		sqlBuilder.WriteString("`")
	}

	sqlBuilder.WriteString(")")

	sqlBuilder.WriteString("\n")

	return sqlBuilder.String()
}

// 插入CSV数据
func insertCsvData2(db *sql.DB, tableName string, csvContent []byte, tableConfig *TableConfig) error {
	// 解析CSV内容
	reader := csv.NewReader(bytes.NewReader(csvContent))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("解析CSV失败: %v", err)
	}

	if len(records) < 2 { // 至少需要标题行和一行数据
		return fmt.Errorf("CSV内容为空或格式不正确")
	}

	// 获取列名（第一行）
	header := records[0]

	// 构建INSERT语句
	insertSQL := buildInsertSQL(tableName, header)
	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("准备INSERT语句失败: %v", err)
	}
	defer stmt.Close()

	// 插入数据行（从第二行开始）
	for i, record := range records[1:] {
		// 转换值为合适的类型
		values := make([]interface{}, len(record))
		for j, value := range record {
			// 根据列配置转换类型
			if j < len(tableConfig.Columns) {
				values[j] = convertValue(value, tableConfig.Columns[j])
			} else {
				values[j] = value // 默认作为字符串
			}
		}

		_, err = stmt.Exec(values...)
		if err != nil {
			return fmt.Errorf("插入第%d行数据失败: %v", i+2, err) // i+2因为从第二行开始，且行号从1开始
		}
	}

	return nil
}

// / 优化插入CSV数据的函数, 使用事务批量提交, 每1000条提交一次
func insertCsvData3(db *sql.DB, tableName string, csvContent []byte, tableConfig *TableConfig) error {
	// 解析CSV内容
	reader := csv.NewReader(bytes.NewReader(csvContent))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("解析CSV失败: %v", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV内容为空或格式不正确")
	}

	// 获取列名（第一行）
	header := records[0]

	// 优化1: 使用事务批量提交
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 优化2: 准备批量插入语句
	stmt, err := tx.Prepare(buildInsertSQL(tableName, header))
	if err != nil {
		return fmt.Errorf("准备INSERT语句失败: %v", err)
	}
	defer stmt.Close()

	// 优化3: 批量处理数据，每1000条提交一次
	batchSize := 1000
	for i, record := range records[1:] {
		values := make([]interface{}, len(record))
		for j, value := range record {
			if j < len(tableConfig.Columns) {
				values[j] = convertValue(value, tableConfig.Columns[j])
			} else {
				values[j] = value
			}
		}

		_, err = stmt.Exec(values...)
		if err != nil {
			return fmt.Errorf("插入第%d行数据失败: %v", i+2, err)
		}

		// 每batchSize条提交一次事务，避免大事务
		if (i+1)%batchSize == 0 {
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("提交事务失败: %v", err)
			}

			// 开始新事务
			tx, err = db.Begin()
			if err != nil {
				return fmt.Errorf("开始新事务失败: %v", err)
			}

			stmt, err = tx.Prepare(buildInsertSQL(tableName, header))
			if err != nil {
				return fmt.Errorf("重新准备INSERT语句失败: %v", err)
			}
		}
	}

	// 提交剩余的数据
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("提交最终事务失败: %v", err)
	}

	return nil
}

// 构建INSERT SQL语句
func buildInsertSQL(tableName string, columns []string) string {
	var sqlBuilder strings.Builder

	sqlBuilder.WriteString("INSERT INTO `")
	sqlBuilder.WriteString(tableName)
	sqlBuilder.WriteString("` (")

	// 添加列名
	for i, column := range columns {
		if i > 0 {
			sqlBuilder.WriteString(", ")
		}
		sqlBuilder.WriteString("`")
		sqlBuilder.WriteString(column)
		sqlBuilder.WriteString("`")
	}

	sqlBuilder.WriteString(") VALUES (")

	// 添加占位符
	for i := range columns {
		if i > 0 {
			sqlBuilder.WriteString(", ")
		}
		sqlBuilder.WriteString("?")
	}

	sqlBuilder.WriteString(")")

	return sqlBuilder.String()
}

func describeTable(db *sql.DB, tableName string) (string, error) {
	rows, err := db.Query("DESCRIBE " + tableName)
	if err != nil {
		return "", fmt.Errorf("无法描述表 %s: %v", tableName, err)
	}
	defer rows.Close()

	var result strings.Builder
	result.WriteString(fmt.Sprintf("表结构: %s\n", tableName))
	result.WriteString("字段名\t\t类型\t\t是否为空\t键\n")
	result.WriteString("----------------------------------------\n")

	for rows.Next() {
		var field, typ, null, key string
		var defaultValue, extra interface{}
		err := rows.Scan(&field, &typ, &null, &key, &defaultValue, &extra)
		if err != nil {
			return "", err
		}
		result.WriteString(fmt.Sprintf("%s\t\t%s\t\t%s\t\t%s\n", field, typ, null, key))
	}

	return result.String(), nil
}

func insertCsvData(db *sql.DB, tableName string, csvContent []byte, tableConfig *TableConfig) error {

	tableDesc, _ := describeTable(db, tableName)
	log.Printf("准备写入数据到表，表结构:%s\n", tableDesc)

	// 解析CSV内容
	reader := csv.NewReader(bytes.NewReader(csvContent))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("解析CSV失败: %v", err)
	}

	if len(records) < 1 {
		return fmt.Errorf("CSV内容为空或格式不正确")
	}

	header := records[0]

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 批量插入的大小
	batchSize := 1000
	rows := records[1:]

	for i := 0; i < len(rows); i += batchSize {
		end := i + batchSize
		if end > len(rows) {
			end = len(rows)
		}
		batch := rows[i:end]

		// 构建批量插入的SQL
		query := buildBatchInsertSQL(tableName, header, len(batch))
		// 准备参数
		args := make([]interface{}, 0, len(batch)*len(header))
		for _, record := range batch {
			for j, value := range record {
				if j < len(tableConfig.Columns) {
					args = append(args, convertValue(value, tableConfig.Columns[j]))
				} else {
					args = append(args, value)
				}
			}
		}

		_, err = tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("批量插入失败: %v, query:%s, args:%v, tableDesc:%s, csvContent:%s", err, query, args, tableDesc, csvContent)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// 构建批量插入SQL
func buildBatchInsertSQL(tableName string, columns []string, batchSize int) string {
	var sqlBuilder strings.Builder

	sqlBuilder.WriteString("INSERT INTO `")
	sqlBuilder.WriteString(tableName)
	sqlBuilder.WriteString("` (")

	// 添加列名
	for i, column := range columns {
		if i > 0 {
			sqlBuilder.WriteString(", ")
		}
		sqlBuilder.WriteString("`")
		sqlBuilder.WriteString(column)
		sqlBuilder.WriteString("`")
	}

	sqlBuilder.WriteString(") VALUES ")

	// 添加占位符：每行一个括号，括号内是多个占位符
	for i := 0; i < batchSize; i++ {
		if i > 0 {
			sqlBuilder.WriteString(", ")
		}
		sqlBuilder.WriteString("(")
		for j := 0; j < len(columns); j++ {
			if j > 0 {
				sqlBuilder.WriteString(", ")
			}
			sqlBuilder.WriteString("?")
		}
		sqlBuilder.WriteString(")")
	}

	return sqlBuilder.String()
}

// 转换值类型
func convertValue(value string, column *ColumnConfig) interface{} {
	if value == "" {
		return nil
	}

	switch column.ColumnType {
	case 1: // 整数（INTEGER）
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	case 2: // 浮点数（DOUBLE）
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	case 3: // 布尔值（BOOLEAN）
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	case 5: // 长整型（BIGINT）
		if int64Val, err := strconv.ParseInt(value, 10, 64); err == nil {
			return int64Val
		}
	default: // 字符串和其他类型
		return value
	}

	// 转换失败时返回原字符串
	return value
}

// 检查表是否存在
func checkTableExists(db *sql.DB, tableName string) (bool, error) {
	var exists bool
	query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"

	err := db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// 替换表操作
func replaceTable(db *sql.DB, targetTableName, tempTableName string, targetTableExists bool) error {
	// 如果目标表存在，先删除
	if targetTableExists {
		dropSQL := fmt.Sprintf("DROP TABLE `%s`", targetTableName)
		_, err := db.Exec(dropSQL)
		if err != nil {
			return fmt.Errorf("删除原表失败: %v", err)
		}
	}

	// 重命名临时表为目标表
	renameSQL := fmt.Sprintf("RENAME TABLE `%s` TO `%s`", tempTableName, targetTableName)
	_, err := db.Exec(renameSQL)
	if err != nil {
		return fmt.Errorf("重命名表失败: %v", err)
	}

	return nil
}

// 清理临时表
func cleanupTempTable(db *sql.DB, tempTableName string) {
	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", tempTableName)
	db.Exec(dropSQL) // 忽略错误，因为这只是清理操作
}

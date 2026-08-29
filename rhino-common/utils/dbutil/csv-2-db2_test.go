package dbutil_test

import (
	"database/sql"
	"fmt"
	"rhino-common/utils/dbutil"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL驱动
)

// 测试用例10：多个联合索引测试
func TestCreateTableFromCsv_MultipleCompositeIndexes(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	tableConfig := &dbutil.TableConfig{
		TableName:    "test_composite_indexes",
		CharacterSet: "utf8mb4",
		Collate:      "utf8mb4_unicode_ci",
		Columns: []*dbutil.ColumnConfig{
			{ColumnName: "user_id", ColumnType: 5},                   // BIGINT
			{ColumnName: "product_id", ColumnType: 5},                // BIGINT
			{ColumnName: "category", ColumnType: 0, ColumnLen: 50},   // VARCHAR
			{ColumnName: "price", ColumnType: 2},                     // DOUBLE
			{ColumnName: "quantity", ColumnType: 1},                  // INTEGER
			{ColumnName: "is_active", ColumnType: 3},                 // BOOLEAN
			{ColumnName: "created_at", ColumnType: 0, ColumnLen: 20}, // VARCHAR存储时间戳
		},
		Indexes: []*dbutil.IndexConfig{
			// 唯一联合索引
			{IndexType: 1, Columns: []string{"user_id", "product_id"}},
			// 普通联合索引
			{IndexType: 0, Columns: []string{"category", "price"}},
			// 三列联合索引
			{IndexType: 0, Columns: []string{"user_id", "category", "is_active"}},
			// 单列索引
			{IndexType: 0, Columns: []string{"created_at"}},
		},
	}

	csvContent := []byte(`user_id,product_id,category,price,quantity,is_active,created_at
1001,2001,电子产品,2999.99,2,true,2024-01-01 10:00:00
1001,2002,服装,199.50,1,true,2024-01-01 10:05:00
1002,2001,电子产品,2999.99,1,false,2024-01-01 10:10:00
1002,2003,图书,59.99,3,true,2024-01-01 10:15:00
1003,2001,电子产品,2999.99,1,true,2024-01-01 10:20:00
1001,2004,服装,299.00,2,true,2024-01-01 10:25:00
1003,2002,服装,199.50,1,false,2024-01-01 10:30:00
1002,2004,服装,299.00,1,true,2024-01-01 10:35:00`)

	err := dbutil.CreateTableFromCsv(db, csvContent, tableConfig)
	if err != nil {
		t.Fatalf("联合索引测试失败: %v", err)
	}

	// 验证数据完整性
	verifyDataIntegrity(t, db, tableConfig.TableName, 8, []string{
		"user_id", "product_id", "category", "price", "quantity", "is_active", "created_at",
	})

	// 验证索引创建
	verifyIndexesDetailed(t, db, tableConfig.TableName, tableConfig.Indexes)

	// 验证唯一约束
	verifyUniqueConstraint(t, db, tableConfig.TableName, []string{"user_id", "product_id"})

	t.Logf("联合索引测试通过，表 %s 包含 %d 行数据", tableConfig.TableName, 8)
}

// 测试用例11：所有列类型组合测试
func TestCreateTableFromCsv_AllColumnTypeCombinations(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	tableConfig := &dbutil.TableConfig{
		TableName:    "test_all_column_combinations",
		CharacterSet: "utf8mb4",
		Collate:      "utf8mb4_unicode_ci",
		Columns: []*dbutil.ColumnConfig{
			// 字符串类型变体
			{ColumnName: "varchar_short", ColumnType: 0, ColumnLen: 10},
			{ColumnName: "varchar_medium", ColumnType: 0, ColumnLen: 100},
			{ColumnName: "varchar_long", ColumnType: 0, ColumnLen: 500},
			{ColumnName: "text_field", ColumnType: 4},

			// 数值类型变体
			{ColumnName: "tiny_int", ColumnType: 1}, // 实际会映射为INT
			{ColumnName: "small_int", ColumnType: 1},
			{ColumnName: "big_int", ColumnType: 5},
			{ColumnName: "float_val", ColumnType: 2},
			{ColumnName: "double_val", ColumnType: 2},

			// 布尔类型
			{ColumnName: "bool_val", ColumnType: 3},
			{ColumnName: "status_flag", ColumnType: 3},

			// 混合类型
			{ColumnName: "mixed_id", ColumnType: 5},
			{ColumnName: "mixed_name", ColumnType: 0, ColumnLen: 200},
			{ColumnName: "mixed_value", ColumnType: 2},
			{ColumnName: "mixed_flag", ColumnType: 3},
			{ColumnName: "mixed_desc", ColumnType: 4},
		},
		Indexes: []*dbutil.IndexConfig{
			{IndexType: 1, Columns: []string{"mixed_id"}},
			{IndexType: 1, Columns: []string{"varchar_short"}},
			{IndexType: 0, Columns: []string{"bool_val", "status_flag"}},
			{IndexType: 0, Columns: []string{"varchar_medium", "float_val"}},
		},
	}

	// 生成包含各种数据类型的CSV
	csvContent := generateMixedTypeCSV(100) // 生成100行混合类型数据

	err := dbutil.CreateTableFromCsv(db, csvContent, tableConfig)
	if err != nil {
		t.Fatalf("所有列类型组合测试失败: %v", err)
	}

	// 详细数据验证
	verifyDetailedData(t, db, tableConfig.TableName, 100, tableConfig.Columns)

	t.Logf("所有列类型组合测试通过，表 %s 包含 %d 行数据", tableConfig.TableName, 100)
}

// 测试用例12：大数据量完整性校验
func TestCreateTableFromCsv_LargeDataIntegrity(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	tableConfig := &dbutil.TableConfig{
		TableName: "test_large_data_integrity",
		Columns: []*dbutil.ColumnConfig{
			{ColumnName: "id", ColumnType: 1},
			{ColumnName: "name", ColumnType: 0, ColumnLen: 50},
			{ColumnName: "value1", ColumnType: 2},
			{ColumnName: "value2", ColumnType: 1},
			{ColumnName: "flag", ColumnType: 3},
			{ColumnName: "description", ColumnType: 4},
		},
		Indexes: []*dbutil.IndexConfig{
			{IndexType: 1, Columns: []string{"id"}},
			{IndexType: 0, Columns: []string{"name"}},
			{IndexType: 0, Columns: []string{"value1", "flag"}},
		},
	}

	// 生成5000行测试数据
	rowCount := 5000
	csvContent := generateLargeTestData(rowCount)

	err := dbutil.CreateTableFromCsv(db, csvContent, tableConfig)
	if err != nil {
		t.Fatalf("大数据量完整性测试失败: %v", err)
	}

	// 全面数据校验
	verifyCompleteDataIntegrity(t, db, tableConfig.TableName, rowCount, tableConfig.Columns)

	t.Logf("大数据量完整性测试通过，表 %s 包含 %d 行数据", tableConfig.TableName, rowCount)
}

// 测试用例13：空值和默认值处理
func TestCreateTableFromCsv_NullAndDefaultHandling(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	tableConfig := &dbutil.TableConfig{
		TableName: "test_null_default_handling",
		Columns: []*dbutil.ColumnConfig{
			{ColumnName: "id", ColumnType: 1},
			{ColumnName: "required_str", ColumnType: 0, ColumnLen: 50},
			{ColumnName: "optional_str", ColumnType: 0, ColumnLen: 100},
			{ColumnName: "number_val", ColumnType: 2},
			{ColumnName: "int_val", ColumnType: 1},
			{ColumnName: "bool_val", ColumnType: 3},
			{ColumnName: "text_val", ColumnType: 4},
		},
	}

	csvContent := []byte(`id,required_str,optional_str,number_val,int_val,bool_val,text_val
1,必填字段1,可选字段1,100.5,50,true,有文本内容
2,必填字段2,,200.75,,false,
3,必填字段3,可选字段3,,25,,""
4,必填字段4,可选字段4,300.25,75,true,另一个文本
5,必填字段5,,,0,false,`)

	err := dbutil.CreateTableFromCsv(db, csvContent, tableConfig)
	if err != nil {
		t.Fatalf("空值处理测试失败: %v", err)
	}

	// 验证空值处理
	verifyNullValueHandling(t, db, tableConfig.TableName)

	t.Logf("空值处理测试通过，表 %s 包含 5 行数据", tableConfig.TableName)
}

// 测试用例14：数据类型边界值测试
func TestCreateTableFromCsv_DataTypeBoundaries(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	tableConfig := &dbutil.TableConfig{
		TableName: "test_data_type_boundaries",
		Columns: []*dbutil.ColumnConfig{
			{ColumnName: "id", ColumnType: 1},
			{ColumnName: "max_varchar", ColumnType: 0, ColumnLen: 255},
			{ColumnName: "min_varchar", ColumnType: 0, ColumnLen: 1},
			{ColumnName: "max_int", ColumnType: 5},
			{ColumnName: "min_int", ColumnType: 5},
			{ColumnName: "large_double", ColumnType: 2},
			{ColumnName: "small_double", ColumnType: 2},
			{ColumnName: "long_text", ColumnType: 4},
		},
	}

	csvContent := []byte(`id,max_varchar,min_varchar,max_int,min_int,large_double,small_double,long_text
1,这是一个很长的字符串用来测试VARCHAR字段的最大长度限制这是255个字符的测试字符串需要确保长度足够长以满足测试需求现在应该差不多了还需要再多一些字符来填满整个空间,a,9223372036854775807,-9223372036854775808,1.7976931348623157e+308,2.2250738585072014e-308,这是一个非常长的文本内容，用来测试MEDIUMTEXT类型的处理能力。这个字段应该能够存储大量的文本数据，包括多行内容和各种特殊字符。我们需要确保长文本能够正确存储和检索，不会出现截断或编码问题。重复这个段落多次以增加长度...
2,短字符串,b,100,-100,999999.9999,0.0001,短文本
3,中等长度字符串,c,0,0,123.456,-123.456,中等长度文本内容`)

	err := dbutil.CreateTableFromCsv(db, csvContent, tableConfig)
	if err != nil {
		t.Fatalf("数据类型边界测试失败: %v", err)
	}

	// 验证边界值处理
	verifyBoundaryValues(t, db, tableConfig.TableName)

	t.Logf("数据类型边界测试通过")
}

// 辅助函数：生成混合类型CSV数据
func generateMixedTypeCSV(rowCount int) []byte {
	var builder strings.Builder
	builder.WriteString("varchar_short,varchar_medium,varchar_long,text_field,tiny_int,small_int,big_int,float_val,double_val,bool_val,status_flag,mixed_id,mixed_name,mixed_value,mixed_flag,mixed_desc\n")

	for i := 1; i <= rowCount; i++ {
		builder.WriteString(fmt.Sprintf(
			"short%d,中文字符串%d,长字符串内容%d,文本描述字段%d,%d,%d,%d,%.2f,%.4f,%t,%t,%d,混合名称%d,%.3f,%t,混合描述内容%d\n",
			i%100, i, i, i,
			i%128, i%10000, int64(i)*1000,
			float64(i)*1.1, float64(i)*1.1111,
			i%2 == 0, i%3 == 0,
			i, i, float64(i)*2.5, i%2 == 1, i,
		))
	}

	return []byte(builder.String())
}

// 辅助函数：生成大数据量测试数据
func generateLargeTestData(rowCount int) []byte {
	var builder strings.Builder
	builder.WriteString("id,name,value1,value2,flag,description\n")

	for i := 1; i <= rowCount; i++ {
		builder.WriteString(fmt.Sprintf(
			"%d,测试用户%d,%.2f,%d,%t,这是第%d条记录的描述信息，包含一些详细的说明内容\n",
			i, i, float64(i)*1.5, i*10, i%2 == 0, i,
		))
	}

	return []byte(builder.String())
}

// 辅助函数：详细数据完整性验证
func verifyDataIntegrity(t *testing.T, db *sql.DB, tableName string, expectedCount int, expectedColumns []string) {
	// 验证行数
	var actualCount int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM `%s`", tableName)).Scan(&actualCount)
	if err != nil {
		t.Errorf("验证行数失败: %v", err)
		return
	}

	if actualCount != expectedCount {
		t.Errorf("数据完整性验证失败: 期望 %d 行，实际 %d 行", expectedCount, actualCount)
	}

	// 验证列存在
	for _, col := range expectedColumns {
		var colExists int
		err := db.QueryRow(`
            SELECT COUNT(*) FROM information_schema.COLUMNS 
            WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?`,
			tableName, col).Scan(&colExists)

		if err != nil || colExists == 0 {
			t.Errorf("列验证失败: 列 %s 不存在", col)
		}
	}
}

// 辅助函数：详细索引验证
func verifyIndexesDetailed(t *testing.T, db *sql.DB, tableName string, expectedIndexes []*dbutil.IndexConfig) {
	indexRows, err := db.Query(`
        SELECT INDEX_NAME, NON_UNIQUE, GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) 
        FROM information_schema.STATISTICS 
        WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND INDEX_NAME != 'PRIMARY'
        GROUP BY INDEX_NAME, NON_UNIQUE`, tableName)

	if err != nil {
		t.Errorf("查询索引详情失败: %v", err)
		return
	}
	defer indexRows.Close()

	foundIndexes := make(map[string]bool)

	for indexRows.Next() {
		var indexName string
		var nonUnique int
		var columns string

		err := indexRows.Scan(&indexName, &nonUnique, &columns)
		if err != nil {
			t.Errorf("扫描索引详情失败: %v", err)
			continue
		}

		foundIndexes[columns] = true
		t.Logf("找到索引: %s, 唯一性: %t, 列: %s", indexName, nonUnique == 0, columns)
	}

	// 验证所有期望的索引都存在
	for _, expectedIndex := range expectedIndexes {
		expectedColumns := strings.Join(expectedIndex.Columns, ",")
		if !foundIndexes[expectedColumns] {
			t.Errorf("索引验证失败: 期望的索引列组合不存在: %s", expectedColumns)
		}
	}
}

// 辅助函数：唯一约束验证
func verifyUniqueConstraint(t *testing.T, db *sql.DB, tableName string, uniqueColumns []string) {
	// 尝试插入重复数据来验证唯一约束
	columns := strings.Join(uniqueColumns, ",")
	placeholders := strings.Repeat("?,", len(uniqueColumns)-1) + "?"

	query := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", tableName, columns, placeholders)

	// 创建一些测试值（使用表中可能已存在的值）
	testValues := make([]interface{}, len(uniqueColumns))
	for i := range uniqueColumns {
		testValues[i] = 1001 // 使用已知存在的值
	}
	if len(uniqueColumns) > 1 {
		testValues[1] = 2001
	}

	_, err := db.Exec(query, testValues...)
	if err == nil {
		t.Errorf("唯一约束验证失败: 应该拒绝重复数据插入")
	} else {
		t.Logf("唯一约束验证通过: 正确拒绝重复数据插入，错误: %v", err)
	}
}

// 辅助函数：详细数据验证
func verifyDetailedData(t *testing.T, db *sql.DB, tableName string, expectedCount int, columns []*dbutil.ColumnConfig) {
	verifyDataIntegrity(t, db, tableName, expectedCount, getColumnNames(columns))

	// 验证数据类型
	for _, col := range columns {
		verifyColumnDataType(t, db, tableName, col)
	}
}

// 辅助函数：获取列名列表
func getColumnNames(columns []*dbutil.ColumnConfig) []string {
	names := make([]string, len(columns))
	for i, col := range columns {
		names[i] = col.ColumnName
	}
	return names
}

// 辅助函数：验证列数据类型
func verifyColumnDataType(t *testing.T, db *sql.DB, tableName string, column *dbutil.ColumnConfig) {
	var dataType string
	err := db.QueryRow(`
        SELECT DATA_TYPE FROM information_schema.COLUMNS 
        WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?`,
		tableName, column.ColumnName).Scan(&dataType)

	if err != nil {
		t.Errorf("查询列数据类型失败: %v", err)
		return
	}

	expectedType := getExpectedDataType(column)
	if dataType != expectedType {
		t.Errorf("数据类型验证失败: 列 %s 期望类型 %s，实际类型 %s",
			column.ColumnName, expectedType, dataType)
	}
}

// 辅助函数：获取期望的数据类型
func getExpectedDataType(column *dbutil.ColumnConfig) string {
	switch column.ColumnType {
	case 0:
		return "varchar"
	case 1:
		return "int"
	case 2:
		return "double"
	case 3:
		return "tinyint" // MySQL中BOOLEAN是tinyint(1)的别名
	case 4:
		return "mediumtext"
	case 5:
		return "bigint"
	default:
		return "varchar"
	}
}

// 辅助函数：完整数据完整性验证
func verifyCompleteDataIntegrity(t *testing.T, db *sql.DB, tableName string, expectedCount int, columns []*dbutil.ColumnConfig) {
	verifyDataIntegrity(t, db, tableName, expectedCount, getColumnNames(columns))

	// 验证没有空的主键
	var nullKeyCount int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE id IS NULL", tableName)).Scan(&nullKeyCount)
	if err == nil && nullKeyCount > 0 {
		t.Errorf("数据完整性验证失败: 发现 %d 个空主键", nullKeyCount)
	}

	// 验证数据一致性（示例：检查数值范围）
	var minVal, maxVal float64
	err = db.QueryRow(fmt.Sprintf("SELECT MIN(value1), MAX(value1) FROM `%s`", tableName)).Scan(&minVal, &maxVal)
	if err == nil {
		if minVal < 0 || maxVal > 1000000 {
			t.Errorf("数据范围验证失败: value1 范围异常 (min: %f, max: %f)", minVal, maxVal)
		}
	}
}

// 辅助函数：空值处理验证
func verifyNullValueHandling(t *testing.T, db *sql.DB, tableName string) {
	// 验证必填字段没有空值
	var nullRequiredCount int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE required_str IS NULL OR required_str = ''", tableName)).Scan(&nullRequiredCount)
	if err == nil && nullRequiredCount > 0 {
		t.Errorf("空值处理验证失败: 发现 %d 个空必填字段", nullRequiredCount)
	}

	// 验证可选字段可以有空值（但应该正确处理）
	var optionalNullCount int
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE optional_str IS NULL OR optional_str = ''", tableName)).Scan(&optionalNullCount)
	if err != nil {
		t.Errorf("空值处理验证失败: 无法统计可选字段空值")
	}
	t.Logf("可选字段空值统计: %d", optionalNullCount)
}

// 辅助函数：边界值验证
func verifyBoundaryValues(t *testing.T, db *sql.DB, tableName string) {
	// 验证最大整数值
	var maxInt int64
	err := db.QueryRow(fmt.Sprintf("SELECT max_int FROM `%s` WHERE id = 1", tableName)).Scan(&maxInt)
	if err == nil && maxInt != 9223372036854775807 {
		t.Errorf("边界值验证失败: 最大整数值不正确")
	}

	// 验证最小整数值
	var minInt int64
	err = db.QueryRow(fmt.Sprintf("SELECT min_int FROM `%s` WHERE id = 1", tableName)).Scan(&minInt)
	if err == nil && minInt != -9223372036854775808 {
		t.Errorf("边界值验证失败: 最小整数值不正确")
	}

	t.Logf("边界值验证通过")
}

// 获取测试数据库连接
func getTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:guangfa4cool@tcp(10.51.136.72:56322)/testdb")
	if err != nil {
		t.Fatalf("数据库连接失败: %v", err)
	}

	db.SetMaxOpenConns(128) // 这两个值设置为一样会好一些，否则，有时会出现链接有问题
	db.SetMaxIdleConns(128) // 这两个值设置为一样会好一些，否则，有时会出现链接有问题
	db.SetConnMaxLifetime(time.Second * 600)

	err = db.Ping()
	if err != nil {
		t.Fatalf("数据库连接测试失败: %v", err)
	}

	return db
}

func TestCreateTableFromCsv_CheckError(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	tableConfig := &dbutil.TableConfig{
		TableName: "counterpartyAuthority",
		Columns: []*dbutil.ColumnConfig{
			{ColumnName: "CTPTY_ID", ColumnType: 5},
			{ColumnName: "XBOND_TRS_FLAG", ColumnType: 5},
			{ColumnName: "UPDATED_DATETIME", ColumnType: 0, ColumnLen: 24},
		},
		Indexes: []*dbutil.IndexConfig{
			{IndexType: 1, Columns: []string{"CTPTY_ID"}},
		},
	}

	// 生成5000行测试数据
	csvContent := []byte(`CTPTY_ID,XBOND_TRS_FLAG,UPDATED_DATETIME
    15346,1,2025-10-09 11:09:21.000
    15348,1,2025-10-13 12:13:38.000
    15349,1,2025-09-29 17:02:57.000
    15353,1,2025-10-10 11:28:00.000
    15358,1,2025-10-09 15:47:30.000`)

	err := dbutil.CreateTableFromCsv(db, csvContent, tableConfig)
	if err != nil {
		t.Fatalf("测试出错: %v", err)
	}
}

func TestCreateTableFromCsv_Comprehensive2(t *testing.T) {

	// 新增的测试用例
	t.Run("多个联合索引测试", func(t *testing.T) {
		TestCreateTableFromCsv_MultipleCompositeIndexes(t)
	})
	t.Run("所有列类型组合测试", func(t *testing.T) {
		TestCreateTableFromCsv_AllColumnTypeCombinations(t)
	})
	t.Run("大数据量完整性校验", func(t *testing.T) {
		TestCreateTableFromCsv_LargeDataIntegrity(t)
	})
	t.Run("空值和默认值处理", func(t *testing.T) {
		TestCreateTableFromCsv_NullAndDefaultHandling(t)
	})
	t.Run("数据类型边界值测试", func(t *testing.T) {
		TestCreateTableFromCsv_DataTypeBoundaries(t)
	})
	t.Run("补充测试", func(t *testing.T) {
		TestCreateTableFromCsv_CheckError(t)
	})
}

func cleanupTestTables2(db *sql.DB) {
	tables := []string{
		"test_all_column_types",
		"test_boundary_values",
		"test_special_chars",
		"test_null_values",
		"test_large_data",
		"test_table_replacement",
		"test_error_cases",
		"test_performance",
	}

	for _, table := range tables {
		db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table))
	}

	// 清理字符集测试表
	charsetTables := []string{
		"test_charset_utf8mb4",
		"test_charset_utf8",
		"test_charset_latin1",
		"test_charset_invalid",
	}

	for _, table := range charsetTables {
		db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table))
	}
}

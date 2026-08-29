package dbutil_test

import (
	"database/sql"
	"fmt"
	"log"
	"rhino-common/utils/dbutil"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL驱动
)

// 测试函数
func TestCreateTableFromCsv_Comprehensive(t *testing.T) {
	// 1. 建立数据库连接
	db, err := sql.Open("mysql", "root:guangfa4cool@tcp(127.0.0.1:56322)/testdb")
	if err != nil {
		t.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(128) // 这两个值设置为一样会好一些，否则，有时会出现链接有问题
	db.SetMaxIdleConns(128) // 这两个值设置为一样会好一些，否则，有时会出现链接有问题
	db.SetConnMaxLifetime(time.Second * 600)

	err = db.Ping()
	if err != nil {
		t.Fatalf("数据库连接测试失败: %v", err)
	}

	// 2. 测试用例1：基本功能测试（所有列类型）
	t.Run("所有列类型测试", func(t *testing.T) {
		tableConfig := &dbutil.TableConfig{
			TableName:    "test_all_column_types",
			CharacterSet: "utf8mb4",
			Collate:      "utf8mb4_unicode_ci",
			Columns: []*dbutil.ColumnConfig{
				// 字符串类型
				{ColumnName: "id", ColumnType: 5}, // BIGINT
				{ColumnName: "name", ColumnType: 0, ColumnLen: 100}, // VARCHAR(100)
				{ColumnName: "short_string", ColumnType: 0, ColumnLen: 10}, // VARCHAR(10)
				// 数值类型
				{ColumnName: "age", ColumnType: 1}, // INTEGER
				{ColumnName: "salary", ColumnType: 2}, // DOUBLE
				{ColumnName: "big_number", ColumnType: 5}, // BIGINT
				// 布尔类型
				{ColumnName: "is_active", ColumnType: 3}, // BOOLEAN
				{ColumnName: "is_deleted", ColumnType: 3}, // BOOLEAN
				// 文本类型
				{ColumnName: "description", ColumnType: 4}, // MEDIUMTEXT
				{ColumnName: "notes", ColumnType: 4}, // MEDIUMTEXT
			},
			Indexes: []*dbutil.IndexConfig{
				{IndexType: 1, Columns: []string{"id"}}, // 唯一索引
				{IndexType: 0, Columns: []string{"name"}}, // 普通索引
				{IndexType: 0, Columns: []string{"age", "is_active"}}, // 复合索引
			},
		}

		csvContent := []byte(`id,name,short_string,age,salary,big_number,is_active,is_deleted,description,notes
1,张三,short1,25,5000.50,123456789,true,false,"这是一个描述","备注1"
2,李四,short2,30,6000.00,987654321,false,true,"另一个描述","备注2"
3,王五,short3,35,7500.75,111222333,true,false,"第三个描述","备注3"
4,赵六,short4,28,5500.25,444555666,true,true,"第四个描述","备注4"
5,钱七,short5,42,8000.00,777888999,false,false,"第五个描述","备注5"`)

		err := dbutil.CreateTableFromCsv(db, csvContent, tableConfig)
		if err != nil {
			t.Errorf("所有列类型测试失败: %v", err)
		}

		// 验证表结构和数据
		verifyTable(t, db, tableConfig.TableName, 5)
	})

	// 3. 测试用例2：边界值测试
	t.Run("边界值测试", func(t *testing.T) {
		tableConfig := &dbutil.TableConfig{
			TableName:    "test_boundary_values",
			CharacterSet: "utf8",
			Collate:      "utf8_general_ci",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnType: 1},
				{ColumnName: "max_string", ColumnType: 0, ColumnLen: 255}, // 最大长度字符串
				{ColumnName: "min_string", ColumnType: 0, ColumnLen: 1}, // 最小长度字符串
				{ColumnName: "large_int", ColumnType: 5},
				{ColumnName: "small_int", ColumnType: 1},
				{ColumnName: "decimal_val", ColumnType: 2},
			},
		}

		csvContent := []byte(`id,max_string,min_string,large_int,small_int,decimal_val
1,这是一个很长的字符串用来测试最大长度限制,a,9223372036854775807,2147483647,999999.9999
2,短字符串,b,-9223372036854775808,-2147483648,0.0001
3,中等长度字符串,c,0,0,123.456`)

		err := dbutil.CreateTableFromCsv(db, csvContent, tableConfig)
		if err != nil {
			t.Errorf("边界值测试失败: %v", err)
		}

		verifyTable(t, db, tableConfig.TableName, 3)
	})

	// 4. 测试用例3：特殊字符和编码测试
	t.Run("特殊字符测试", func(t *testing.T) {
		tableConfig := &dbutil.TableConfig{
			TableName:    "test_special_chars",
			CharacterSet: "utf8mb4",
			Collate:      "utf8mb4_unicode_ci",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnType: 1},
				{ColumnName: "text_content", ColumnType: 4},
				{ColumnName: "emoji", ColumnType: 0, ColumnLen: 50},
			},
		}

		csvContent := []byte(`id,text_content,emoji
1,"特殊字符: !@#$%^&*()_+-=[]{}|;:',.<>?/","😀🎉🌟"
2,"中文测试：你好世界","❤️🔥🚀"
3,"SQL注入测试: ' OR '1'='1","📱💻🖥️"
4,"换行符测试: 第一行\n第二行","🐱🐶🐯"`)

		err := dbutil.CreateTableFromCsv(db, csvContent, tableConfig)
		if err != nil {
			t.Errorf("特殊字符测试失败: %v", err)
		}

		verifyTable(t, db, tableConfig.TableName, 4)
	})

	// 5. 测试用例4：空值和默认值测试
	t.Run("空值处理测试", func(t *testing.T) {
		tableConfig := &dbutil.TableConfig{
			TableName: "test_null_values",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnType: 1},
				{ColumnName: "name", ColumnType: 0, ColumnLen: 50},
				{ColumnName: "age", ColumnType: 1},
				{ColumnName: "salary", ColumnType: 2},
				{ColumnName: "description", ColumnType: 4},
			},
		}

		csvContent := []byte(`id,name,age,salary,description
1,张三,25,5000.50,"有描述"
2,,30,,  // 空名字和空薪水
3,王五,,6000.50,  // 空年龄
4,赵六,28,,"空描述"
5,,,,"全部为空"`)

		err := dbutil.CreateTableFromCsv(db, csvContent, tableConfig)
		if err != nil {
			t.Errorf("空值处理测试失败: %v", err)
		}

		verifyTable(t, db, tableConfig.TableName, 5)
	})

	// 6. 测试用例5：大数据量测试
	t.Run("大数据量测试", func(t *testing.T) {
		tableConfig := &dbutil.TableConfig{
			TableName: "test_large_data",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnType: 1},
				{ColumnName: "data", ColumnType: 0, ColumnLen: 100},
			},
			Indexes: []*dbutil.IndexConfig{
				{IndexType: 1, Columns: []string{"id"}},
			},
		}

		// 生成1000行测试数据
		csvBuilder := "id,data\n"
		for i := 1; i <= 1000; i++ {
			csvBuilder += fmt.Sprintf("%d,测试数据行%d\n", i, i)
		}
		csvContent := []byte(csvBuilder)

		err := dbutil.CreateTableFromCsv(db, csvContent, tableConfig)
		if err != nil {
			t.Errorf("大数据量测试失败: %v", err)
		}

		verifyTable(t, db, tableConfig.TableName, 1000)
	})

	// 7. 测试用例6：表替换测试（测试表已存在的情况）
	t.Run("表替换测试", func(t *testing.T) {
		// 先创建一个表
		initialTableConfig := &dbutil.TableConfig{
			TableName: "test_table_replacement",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnType: 1},
				{ColumnName: "old_data", ColumnType: 0, ColumnLen: 50},
			},
		}

		initialCsv := []byte(`id,old_data
1,旧数据1
2,旧数据2`)

		err := dbutil.CreateTableFromCsv(db, initialCsv, initialTableConfig)
		if err != nil {
			t.Errorf("创建初始表失败: %v", err)
		}

		// 现在用新数据替换表
		newTableConfig := &dbutil.TableConfig{
			TableName: "test_table_replacement",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnType: 1},
				{ColumnName: "new_data", ColumnType: 0, ColumnLen: 50},
				{ColumnName: "additional_col", ColumnType: 1},
			},
		}

		newCsv := []byte(`id,new_data,additional_col
1,新数据1,100
2,新数据2,200
3,新数据3,300`)

		err = dbutil.CreateTableFromCsv(db, newCsv, newTableConfig)
		if err != nil {
			t.Errorf("表替换测试失败: %v", err)
		}

		verifyTable(t, db, newTableConfig.TableName, 3)
	})

	// 8. 测试用例7：错误情况测试
	t.Run("错误情况测试", func(t *testing.T) {
		// 测试空的CSV内容
		tableConfig := &dbutil.TableConfig{
			TableName: "test_error_cases",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnType: 1},
			},
		}

		// 测试空CSV
		emptyCsv := []byte("")
		err := dbutil.CreateTableFromCsv(db, emptyCsv, tableConfig)
		if err == nil {
			t.Error("空CSV应该返回错误")
		}

		// 测试格式错误的CSV
		invalidCsv := []byte("id,name\n1,张三\n2") // 缺少列
		err = dbutil.CreateTableFromCsv(db, invalidCsv, tableConfig)
		if err == nil {
			t.Error("格式错误的CSV应该返回错误")
		}

		// 测试列数不匹配
		mismatchCsv := []byte("id,name,age\n1,张三") // 列数不匹配
		err = dbutil.CreateTableFromCsv(db, mismatchCsv, tableConfig)
		if err == nil {
			t.Error("列数不匹配应该返回错误")
		}
	})

	// 9. 测试用例8：性能测试
	t.Run("性能测试", func(t *testing.T) {
		tableConfig := &dbutil.TableConfig{
			TableName: "test_performance",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnType: 1},
				{ColumnName: "data1", ColumnType: 0, ColumnLen: 50},
				{ColumnName: "data2", ColumnType: 1},
				{ColumnName: "data3", ColumnType: 2},
			},
			Indexes: []*dbutil.IndexConfig{
				{IndexType: 1, Columns: []string{"id"}},
				{IndexType: 0, Columns: []string{"data2"}},
				{IndexType: 1, Columns: []string{"data2", "data3"}},
			},
		}

		startTime := time.Now()

		// 生成性能测试数据
		csvBuilder := "id,data1,data2,data3\n"
		n := 50505
		for i := 1; i <= n; i++ {
			csvBuilder += fmt.Sprintf("%d,性能测试数据%d,%d,%.2f\n", i, i, i%100, float64(i)*1.5)
		}
		csvContent := []byte(csvBuilder)

		err := dbutil.CreateTableFromCsv(db, csvContent, tableConfig)
		if err != nil {
			t.Errorf("性能测试失败: %v", err)
		}

		elapsed := time.Since(startTime)
		t.Logf("性能测试完成，处理%d行数据耗时: %v", n, elapsed)

		if elapsed > 15*time.Second {
			t.Errorf("性能测试超时，耗时过长: %v", elapsed)
		}
	})

	// 10. 测试用例9：不同字符集测试
	t.Run("字符集测试", func(t *testing.T) {
		charsetTests := []struct {
			name        string
			charset     string
			collate     string
			shouldWork  bool
		}{
			{"utf8mb4", "utf8mb4", "utf8mb4_unicode_ci", true},
			{"utf8", "utf8", "utf8_general_ci", true},
			{"latin1", "latin1", "latin1_swedish_ci", true},
			{"invalid", "invalid_charset", "invalid_collate", false},
		}

		for _, tt := range charsetTests {
			t.Run(tt.name, func(t *testing.T) {
				tableConfig := &dbutil.TableConfig{
					TableName:    "test_charset_" + tt.name,
					CharacterSet: tt.charset,
					Collate:      tt.collate,
					Columns: []*dbutil.ColumnConfig{
						{ColumnName: "id", ColumnType: 1},
						{ColumnName: "data", ColumnType: 0, ColumnLen: 50},
					},
				}

				csvContent := []byte(`id,data
1,字符集测试数据`)

				err := dbutil.CreateTableFromCsv(db, csvContent, tableConfig)
				if tt.shouldWork && err != nil {
					t.Errorf("字符集 %s 应该工作但失败: %v", tt.name, err)
				}
				if !tt.shouldWork && err == nil {
					t.Errorf("无效字符集 %s 应该失败但成功了", tt.name)
				}

				if tt.shouldWork && err == nil {
					verifyTable(t, db, tableConfig.TableName, 1)
				}
			})
		}
	})
}

// 验证表结构和数据的辅助函数
func verifyTable(t *testing.T, db *sql.DB, tableName string, expectedRowCount int) {
	// 验证表存在
	var tableExists bool
	err := db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", tableName).Scan(&tableExists)
	if err != nil || !tableExists {
		t.Errorf("表 %s 不存在", tableName)
		return
	}

	// 验证行数
	var rowCount int
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM `%s`", tableName)).Scan(&rowCount)
	if err != nil {
		t.Errorf("查询表 %s 行数失败: %v", tableName, err)
		return
	}

	if rowCount != expectedRowCount {
		t.Errorf("表 %s 行数不正确，期望 %d，实际 %d", tableName, expectedRowCount, rowCount)
	}

	// 验证索引
	indexRows, err := db.Query(`
		SELECT INDEX_NAME, NON_UNIQUE 
		FROM information_schema.STATISTICS 
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? 
		ORDER BY INDEX_NAME`, tableName)
	if err == nil {
		defer indexRows.Close()
		
		var indexCount int
		for indexRows.Next() {
			indexCount++
		}
		t.Logf("表 %s 有 %d 个索引", tableName, indexCount)
	}

	t.Logf("表 %s 验证通过，包含 %d 行数据", tableName, rowCount)
}

// 清理测试环境的函数
func cleanupTestTables(db *sql.DB) {
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

// 主测试函数
func TestMain(m *testing.M) {
	// 建立测试数据库连接
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/testdb")
	if err != nil {
		log.Fatalf("测试数据库连接失败: %v", err)
	}
	defer db.Close()

	// 运行测试前清理环境
	cleanupTestTables(db)

	// 运行测试
	exitCode := m.Run()

	// 测试完成后清理环境
	cleanupTestTables(db)

	// 退出
	if exitCode != 0 {
		log.Fatal("测试失败")
	}
}
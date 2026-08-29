package datamap_test

import (
	"database/sql"
	"fmt"
	"rhino-data/datamap"
	"testing"
	"time"

	"rhino-common/utils/dbutil"

	_ "github.com/go-sql-driver/mysql" // MySQL驱动
)

// 测试数据库连接
func getTestDB() (*sql.DB, error) {
	return sql.Open("mysql", "root:guangfa4cool@tcp(127.0.0.1:56322)/testdb")
}

// 初始化测试表
func initTestTable(db *sql.DB) error {
	// 删除已存在的测试表
	_, _ = db.Exec("DROP TABLE IF EXISTS test_users")
	_, _ = db.Exec("DROP TABLE IF EXISTS test_orders")

	// 创建用户测试表
	createUserTableSQL := `
	CREATE TABLE test_users (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		username VARCHAR(50) NOT NULL,
		email VARCHAR(100),
		age INT,
		salary DOUBLE,
		is_active TINYINT(1),
		description MEDIUMTEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(createUserTableSQL)
	if err != nil {
		return err
	}

	// 插入测试用户数据
	insertUsersSQL := `
	INSERT INTO test_users (username, email, age, salary, is_active, description) VALUES
	('user1', 'user1@example.com', 25, 50000.50, 1, 'First user'),
	('user2', 'user2@example.com', 30, 60000.75, 1, 'Second user'),
	('user3', 'user3@example.com', 35, 70000.25, 0, 'Third user'),
	('user4', 'user4@example.com', 28, 55000.00, 1, NULL),
	('user5', 'user5@example.com', 32, 65000.50, 0, 'Fifth user')
	`
	_, err = db.Exec(insertUsersSQL)
	if err != nil {
		return err
	}

	// 创建订单测试表（用于测试同一key多条记录）
	createOrdersTableSQL := `
	CREATE TABLE test_orders (
		order_id BIGINT PRIMARY KEY AUTO_INCREMENT,
		customer_id BIGINT NOT NULL,
		product_name VARCHAR(100) NOT NULL,
		quantity INT,
		amount DOUBLE,
		order_date DATE
	)`
	_, err = db.Exec(createOrdersTableSQL)
	if err != nil {
		return err
	}

	// 插入测试订单数据
	insertOrdersSQL := `
	INSERT INTO test_orders (customer_id, product_name, quantity, amount, order_date) VALUES
	(1, 'Product A', 2, 100.50, '2023-01-15'),
	(1, 'Product B', 1, 50.25, '2023-02-20'),
	(1, 'Product C', 3, 150.75, '2023-01-05'),
	(2, 'Product A', 1, 100.50, '2023-03-10'),
	(2, 'Product D', 2, 75.30, '2023-03-15'),
	(3, 'Product B', 5, 251.25, '2023-04-01')
	`
	_, err = db.Exec(insertOrdersSQL)
	return err
}

// 清理测试表
func cleanupTestTable(db *sql.DB) {
	_, _ = db.Exec("DROP TABLE IF EXISTS test_users")
	_, _ = db.Exec("DROP TABLE IF EXISTS test_orders")
}

// TestCreateMapFromDB_MultiKeys 多个key配置测试
func TestCreateMapFromDB_MultiKeys(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = initTestTable(db)
	if err != nil {
		t.Fatalf("Failed to init test table: %v", err)
	}
	defer cleanupTestTable(db)

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "test_orders",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "order_id", ColumnAlias: "order_id", ColumnType: 5},
				{ColumnName: "customer_id", ColumnAlias: "customer_id", ColumnType: 5},
				{ColumnName: "product_name", ColumnAlias: "product", ColumnType: 0, ColumnLen: 100},
				{ColumnName: "quantity", ColumnAlias: "qty", ColumnType: 1},
				{ColumnName: "amount", ColumnAlias: "amount", ColumnType: 2},
				{ColumnName: "order_date", ColumnAlias: "date", ColumnType: 0, ColumnLen: 20},
			},
		},
		ColumnsForMapKeys: [][]string{
			{"customer_id"},                    // 按客户ID分组
			{"customer_id", "product"},         // 按客户ID和产品分组
			{"date"},                           // 按日期分组
			{"product"},                        // 按产品分组
		},
	}

	tempData, result, err := datamap.CreateMapFromDB(db, 100, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	// 验证tempData和result的大小应该相同
	if len(tempData) != len(result) {
		t.Errorf("tempData and result should have same size, tempData: %d, result: %d", len(tempData), len(result))
	}

	// 验证不同维度的key都存在
	expectedKeys := []string{
		"1",                    // customer_id=1
		"2",                    // customer_id=2  
		"3",                    // customer_id=3
		"1-Product A",          // customer_id=1, product=Product A
		"1-Product B",          // customer_id=1, product=Product B
		"1-Product C",          // customer_id=1, product=Product C
		"2-Product A",          // customer_id=2, product=Product A
		"2-Product D",          // customer_id=2, product=Product D
		"3-Product B",          // customer_id=3, product=Product B
		"2023-01-15",           // date=2023-01-15
		"2023-02-20",           // date=2023-02-20
		"2023-01-05",           // date=2023-01-05
		"2023-03-10",           // date=2023-03-10
		"2023-03-15",           // date=2023-03-15
		"2023-04-01",           // date=2023-04-01
		"Product A",            // product=Product A
		"Product B",            // product=Product B
		"Product C",            // product=Product C
		"Product D",            // product=Product D
	}

	for _, expectedKey := range expectedKeys {
		if _, exists := tempData[expectedKey]; !exists {
			t.Errorf("Expected key '%s' not found in tempData", expectedKey)
		}
		if _, exists := result[expectedKey]; !exists {
			t.Errorf("Expected key '%s' not found in result", expectedKey)
		}
	}

	// 验证客户1有3个订单
	if rows, exists := tempData["1"]; exists {
		if len(rows) != 3 {
			t.Errorf("Expected 3 orders for customer 1, got %d", len(rows))
		}
	}

	// 验证客户1和产品A的组合有1个订单
	if rows, exists := tempData["1-Product A"]; exists {
		if len(rows) != 1 {
			t.Errorf("Expected 1 order for customer 1 and Product A, got %d", len(rows))
		}
	}

	// 验证产品A有2个订单（来自客户1和客户2）
	if rows, exists := tempData["Product A"]; exists {
		if len(rows) != 2 {
			t.Errorf("Expected 2 orders for Product A, got %d", len(rows))
		}
	}

	// 验证日期2023-01-15有1个订单
	if rows, exists := tempData["2023-01-15"]; exists {
		if len(rows) != 1 {
			t.Errorf("Expected 1 order for date 2023-01-15, got %d", len(rows))
		}
	}
}

// TestCreateMapFromDB_Basic 基础功能测试
func TestCreateMapFromDB_Basic(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = initTestTable(db)
	if err != nil {
		t.Fatalf("Failed to init test table: %v", err)
	}
	defer cleanupTestTable(db)

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "test_users",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "user_id", ColumnType: 5},
				{ColumnName: "username", ColumnAlias: "name", ColumnType: 0, ColumnLen: 50},
				{ColumnName: "email", ColumnAlias: "email", ColumnType: 0, ColumnLen: 100},
				{ColumnName: "age", ColumnAlias: "age", ColumnType: 1},
				{ColumnName: "salary", ColumnAlias: "salary", ColumnType: 2},
				{ColumnName: "is_active", ColumnAlias: "active", ColumnType: 3},
				{ColumnName: "description", ColumnAlias: "desc", ColumnType: 4},
			},
		},
		ColumnsForMapKeys: [][]string{{"user_id"}},
	}

	_, result, err := datamap.CreateMapFromDB(db, 100, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	// 验证结果数量
	if len(result) != 5 {
		t.Errorf("Expected 5 records, got %d", len(result))
	}

	// 验证第一条记录的数据类型和值
	firstKey := "1"
	if row, exists := result[firstKey]; exists {
		// 验证用户ID
		if userIDs, ok := row.Value["user_id"]; ok && len(userIDs) > 0 {
			if userID, ok := userIDs[0].(int64); ok {
				if userID != 1 {
					t.Errorf("Expected user_id 1, got %d", userID)
				}
			} else {
				t.Error("user_id is not int64")
			}
		} else {
			t.Error("user_id not found or not a slice")
		}

		// 验证用户名
		if names, ok := row.Value["name"]; ok && len(names) > 0 {
			if name, ok := names[0].(string); ok {
				if name != "user1" {
					t.Errorf("Expected name 'user1', got '%s'", name)
				}
			} else {
				t.Error("name is not string")
			}
		}

		// 验证年龄
		if ages, ok := row.Value["age"]; ok && len(ages) > 0 {
			if age, ok := ages[0].(int64); ok {
				if age != 25 {
					t.Errorf("Expected age 25, got %d", age)
				}
			} else {
				t.Error("age is not int64")
			}
		}

		// 验证薪水
		if salaries, ok := row.Value["salary"]; ok && len(salaries) > 0 {
			if salary, ok := salaries[0].(float64); ok {
				if salary != 50000.50 {
					t.Errorf("Expected salary 50000.50, got %f", salary)
				}
			} else {
				t.Error("salary is not float64")
			}
		}

		// 验证激活状态
		if actives, ok := row.Value["active"]; ok && len(actives) > 0 {
			if active, ok := actives[0].(bool); ok {
				if !active {
					t.Error("Expected active to be true")
				}
			} else {
				t.Error("active is not bool")
			}
		}
	} else {
		t.Error("First record not found")
	}
}

// TestCreateMapFromDB_MultiColumnKey 多列作为key的测试
func TestCreateMapFromDB_MultiColumnKey(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = initTestTable(db)
	if err != nil {
		t.Fatalf("Failed to init test table: %v", err)
	}
	defer cleanupTestTable(db)

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "test_users",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "user_id", ColumnType: 5},
				{ColumnName: "username", ColumnAlias: "name", ColumnType: 0, ColumnLen: 50},
				{ColumnName: "age", ColumnAlias: "age", ColumnType: 1},
			},
		},
		ColumnsForMapKeys: [][]string{{"user_id", "name"}},
	}

	_, result, err := datamap.CreateMapFromDB(db, 100, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	// 验证组合key
	expectedKey := "1-user1"
	if _, exists := result[expectedKey]; !exists {
		t.Errorf("Expected key '%s' not found", expectedKey)
	}

	if len(result) != 5 {
		t.Errorf("Expected 5 records, got %d", len(result))
	}
}

// TestCreateMapFromDB_SameKeyMultipleRecords 同一key多条记录的测试
func TestCreateMapFromDB_SameKeyMultipleRecords(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = initTestTable(db)
	if err != nil {
		t.Fatalf("Failed to init test table: %v", err)
	}
	defer cleanupTestTable(db)

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "test_orders",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "order_id", ColumnAlias: "order_id", ColumnType: 5},
				{ColumnName: "customer_id", ColumnAlias: "customer_id", ColumnType: 5},
				{ColumnName: "product_name", ColumnAlias: "product", ColumnType: 0, ColumnLen: 100},
				{ColumnName: "quantity", ColumnAlias: "qty", ColumnType: 1},
				{ColumnName: "amount", ColumnAlias: "amount", ColumnType: 2},
				{ColumnName: "order_date", ColumnAlias: "date", ColumnType: 0, ColumnLen: 20},
			},
		},
		ColumnsForMapKeys: [][]string{{"customer_id"}},
	}

	_, result, err := datamap.CreateMapFromDB(db, 100, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	// 验证客户1有3个订单
	customer1Key := "1"
	if row, exists := result[customer1Key]; exists {
		if orderIDs, ok := row.Value["order_id"]; ok {
			if len(orderIDs) != 3 {
				t.Errorf("Expected 3 orders for customer 1, got %d", len(orderIDs))
			}
		} else {
			t.Error("order_id not found or not a slice")
		}
	} else {
		t.Error("Customer 1 not found")
	}

	// 验证客户2有2个订单
	customer2Key := "2"
	if row, exists := result[customer2Key]; exists {
		if orderIDs, ok := row.Value["order_id"]; ok {
			if len(orderIDs) != 2 {
				t.Errorf("Expected 2 orders for customer 2, got %d", len(orderIDs))
			}
		}
	} else {
		t.Error("Customer 2 not found")
	}
}

// TestCreateMapFromDB_WithSorting 带排序的测试
func TestCreateMapFromDB_WithSorting(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = initTestTable(db)
	if err != nil {
		t.Fatalf("Failed to init test table: %v", err)
	}
	defer cleanupTestTable(db)

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "test_orders",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "order_id", ColumnAlias: "order_id", ColumnType: 5},
				{ColumnName: "customer_id", ColumnAlias: "customer_id", ColumnType: 5},
				{ColumnName: "product_name", ColumnAlias: "product", ColumnType: 0, ColumnLen: 100},
				{ColumnName: "quantity", ColumnAlias: "qty", ColumnType: 1},
				{ColumnName: "amount", ColumnAlias: "amount", ColumnType: 2},
				{ColumnName: "order_date", ColumnAlias: "date", ColumnType: 0, ColumnLen: 20},
			},
		},
		ColumnsForMapKeys: [][]string{{"customer_id"}},
		SortRule: &datamap.SortRuleForValueInSameKey{
			ColumnsForSorting: []string{"amount"},
			SortDir:           "DESC",
		},
	}

	_, result, err := datamap.CreateMapFromDB(db, 100, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	// 验证客户1的订单按金额降序排列
	customer1Key := "1"
	if row, exists := result[customer1Key]; exists {
		if amounts, ok := row.Value["amount"]; ok && len(amounts) == 3 {
			// 验证是否是降序排列
			amount1, _ := amounts[0].(float64)
			amount2, _ := amounts[1].(float64)
			amount3, _ := amounts[2].(float64)
			
			if amount1 < amount2 || amount2 < amount3 {
				t.Errorf("Amounts are not in descending order: %f, %f, %f", amount1, amount2, amount3)
			}
		} else {
			t.Error("Amounts not found or incorrect count")
		}
	}
}

// TestCreateMapFromDB_Pagination 分页测试
func TestCreateMapFromDB_Pagination(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = initTestTable(db)
	if err != nil {
		t.Fatalf("Failed to init test table: %v", err)
	}
	defer cleanupTestTable(db)

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "test_users",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "user_id", ColumnType: 5},
				{ColumnName: "username", ColumnAlias: "name", ColumnType: 0, ColumnLen: 50},
			},
		},
		ColumnsForMapKeys: [][]string{{"user_id"}},
	}

	// 使用小页大小测试分页
	_, result, err := datamap.CreateMapFromDB(db, 2, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	// 应该仍然获取所有5条记录
	if len(result) != 5 {
		t.Errorf("Expected 5 records with pagination, got %d", len(result))
	}
}

// TestCreateMapFromDB_NullValues NULL值处理测试
func TestCreateMapFromDB_NullValues(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = initTestTable(db)
	if err != nil {
		t.Fatalf("Failed to init test table: %v", err)
	}
	defer cleanupTestTable(db)

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "test_users",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "user_id", ColumnType: 5},
				{ColumnName: "username", ColumnAlias: "name", ColumnType: 0, ColumnLen: 50},
				{ColumnName: "description", ColumnAlias: "desc", ColumnType: 4},
			},
		},
		ColumnsForMapKeys: [][]string{{"user_id"}},
	}

	_, result, err := datamap.CreateMapFromDB(db, 100, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	// 验证用户4的描述为NULL
	user4Key := "4"
	if row, exists := result[user4Key]; exists {
		if descs, ok := row.Value["desc"]; ok && len(descs) > 0 {
			if desc, ok := descs[0].(string); ok {
				if desc != "" {
					t.Errorf("Expected empty string for NULL description, got '%s'", desc)
				}
			} else {
				t.Error("desc is not string")
			}
		}
	} else {
		t.Error("User 4 not found")
	}
}

// TestCreateMapFromDB_ErrorCases 错误情况测试
func TestCreateMapFromDB_ErrorCases(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 测试1: 空配置
	_, _, err = datamap.CreateMapFromDB(db, 100, nil)
	if err == nil {
		t.Error("Expected error for nil tableConfig")
	}

	// 测试2: 空TableConfig
	tableConfig := &datamap.TableConfigForMapping{}
	_, _, err = datamap.CreateMapFromDB(db, 100, tableConfig)
	if err == nil {
		t.Error("Expected error for nil TableConfig")
	}

	// 测试3: 空ColumnsForMapKeys
	tableConfig = &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "test_users",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "user_id", ColumnType: 5},
			},
		},
		ColumnsForMapKeys: [][]string{},
	}
	_, _, err = datamap.CreateMapFromDB(db, 100, tableConfig)
	if err == nil {
		t.Error("Expected error for empty ColumnsForMapKeys")
	}

	// 测试4: 无效页大小
	tableConfig.ColumnsForMapKeys = [][]string{{"user_id"}}
	_, _, err = datamap.CreateMapFromDB(db, 0, tableConfig)
	if err == nil {
		t.Error("Expected error for zero pageSize")
	}

	// 测试5: 不存在的表
	tableConfig.TableConfig.TableName = "non_existent_table"
	_, _, err = datamap.CreateMapFromDB(db, 100, tableConfig)
	if err == nil {
		t.Error("Expected error for non-existent table")
	}
}

// TestCreateMapFromDB_ColumnAlias 列别名测试
func TestCreateMapFromDB_ColumnAlias(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = initTestTable(db)
	if err != nil {
		t.Fatalf("Failed to init test table: %v", err)
	}
	defer cleanupTestTable(db)

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "test_users",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "user_identifier", ColumnType: 5},
				{ColumnName: "username", ColumnAlias: "", ColumnType: 0, ColumnLen: 50}, // 空别名，应使用列名
				{ColumnName: "email", ColumnAlias: "email_address", ColumnType: 0, ColumnLen: 100},
			},
		},
		ColumnsForMapKeys: [][]string{{"user_identifier"}},
	}

	_, result, err := datamap.CreateMapFromDB(db, 100, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	// 验证别名是否正确应用
	firstKey := "1"
	if row, exists := result[firstKey]; exists {
		// 验证自定义别名
		if _, ok := row.Value["user_identifier"]; !ok {
			t.Error("user_identifier alias not found")
		}
		
		// 验证空别名时使用列名
		if _, ok := row.Value["username"]; !ok {
			t.Error("username (column name) not found when alias is empty")
		}
		
		// 验证自定义别名
		if _, ok := row.Value["email_address"]; !ok {
			t.Error("email_address alias not found")
		}
	} else {
		t.Error("First record not found")
	}
}

// TestCreateMapFromDB_AllDataTypes 所有数据类型测试
func TestCreateMapFromDB_AllDataTypes(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = initTestTable(db)
	if err != nil {
		t.Fatalf("Failed to init test table: %v", err)
	}
	defer cleanupTestTable(db)

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "test_users",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "user_id", ColumnType: 5},           // BIGINT -> int64
				{ColumnName: "username", ColumnAlias: "name", ColumnType: 0, ColumnLen: 50}, // VARCHAR -> string
				{ColumnName: "age", ColumnAlias: "age", ColumnType: 1},              // INTEGER -> int64
				{ColumnName: "salary", ColumnAlias: "salary", ColumnType: 2},        // DOUBLE -> float64
				{ColumnName: "is_active", ColumnAlias: "active", ColumnType: 3},     // BOOLEAN -> bool
				{ColumnName: "description", ColumnAlias: "desc", ColumnType: 4},     // MEDIUMTEXT -> string
			},
		},
		ColumnsForMapKeys: [][]string{{"user_id"}},
	}

	_, result, err := datamap.CreateMapFromDB(db, 100, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	// 验证所有数据类型转换正确
	for key, row := range result {
		t.Logf("Testing data types for key: %s", key)
		
		for colName, values := range row.Value {
			if len(values) == 0 {
				t.Errorf("No values for column %s", colName)
				continue
			}
			
			value := values[0]
			
			switch colName {
			case "user_id", "age":
				if _, ok := value.(int64); !ok {
					t.Errorf("Column %s should be int64, got %T", colName, value)
				}
			case "name", "desc":
				if _, ok := value.(string); !ok {
					t.Errorf("Column %s should be string, got %T", colName, value)
				}
			case "salary":
				if _, ok := value.(float64); !ok {
					t.Errorf("Column %s should be float64, got %T", colName, value)
				}
			case "active":
				if _, ok := value.(bool); !ok {
					t.Errorf("Column %s should be bool, got %T", colName, value)
				}
			}
		}
	}
}

// TestCreateMapFromDB_EmptyTable 空表测试
func TestCreateMapFromDB_EmptyTable(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 创建空表
	_, err = db.Exec("DROP TABLE IF EXISTS empty_table")
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
	_, err = db.Exec("CREATE TABLE empty_table (id INT PRIMARY KEY, name VARCHAR(50))")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	defer db.Exec("DROP TABLE IF EXISTS empty_table")

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "empty_table",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "id", ColumnType: 1},
				{ColumnName: "name", ColumnAlias: "name", ColumnType: 0, ColumnLen: 50},
			},
		},
		ColumnsForMapKeys: [][]string{{"id"}},
	}

	_, result, err := datamap.CreateMapFromDB(db, 100, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty result for empty table, got %d records", len(result))
	}
}

// TestCreateMapFromDB_SingleRecord 单条记录测试
func TestCreateMapFromDB_SingleRecord(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 创建单记录表
	_, err = db.Exec("DROP TABLE IF EXISTS single_table")
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
	_, err = db.Exec("CREATE TABLE single_table (id INT PRIMARY KEY, data VARCHAR(50))")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	_, err = db.Exec("INSERT INTO single_table VALUES (1, 'single record')")
	if err != nil {
		t.Fatalf("Failed to insert record: %v", err)
	}
	defer db.Exec("DROP TABLE IF EXISTS single_table")

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "single_table",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "id", ColumnType: 1},
				{ColumnName: "data", ColumnAlias: "data", ColumnType: 0, ColumnLen: 50},
			},
		},
		ColumnsForMapKeys: [][]string{{"id"}},
	}

	_, result, err := datamap.CreateMapFromDB(db, 100, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 record, got %d", len(result))
	}

	if row, exists := result["1"]; exists {
		if data, ok := row.Value["data"]; ok && len(data) == 1 {
			if data[0].(string) != "single record" {
				t.Errorf("Expected 'single record', got '%v'", data[0])
			}
		} else {
			t.Error("Data not found or incorrect format")
		}
	} else {
		t.Error("Record not found")
	}
}

// TestCreateMapFromDB_LargePageSize 大页大小测试
func TestCreateMapFromDB_LargePageSize(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = initTestTable(db)
	if err != nil {
		t.Fatalf("Failed to init test table: %v", err)
	}
	defer cleanupTestTable(db)

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "test_users",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "user_id", ColumnType: 5},
				{ColumnName: "username", ColumnAlias: "name", ColumnType: 0, ColumnLen: 50},
			},
		},
		ColumnsForMapKeys: [][]string{{"user_id"}},
	}

	// 使用非常大的页大小
	_, result, err := datamap.CreateMapFromDB(db, 10000, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	if len(result) != 5 {
		t.Errorf("Expected 5 records, got %d", len(result))
	}
}

// TestCreateMapFromDB_ComplexMultiColumnSort 复杂多列排序测试
func TestCreateMapFromDB_ComplexMultiColumnSort(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = initTestTable(db)
	if err != nil {
		t.Fatalf("Failed to init test table: %v", err)
	}
	defer cleanupTestTable(db)

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "test_orders",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "order_id", ColumnAlias: "order_id", ColumnType: 5},
				{ColumnName: "customer_id", ColumnAlias: "customer_id", ColumnType: 5},
				{ColumnName: "product_name", ColumnAlias: "product", ColumnType: 0, ColumnLen: 100},
				{ColumnName: "quantity", ColumnAlias: "qty", ColumnType: 1},
				{ColumnName: "amount", ColumnAlias: "amount", ColumnType: 2},
				{ColumnName: "order_date", ColumnAlias: "date", ColumnType: 0, ColumnLen: 20},
			},
		},
		ColumnsForMapKeys: [][]string{{"customer_id"}},
		SortRule: &datamap.SortRuleForValueInSameKey{
			ColumnsForSorting: []string{"amount", "order_id"}, // 多列排序
			SortDir:           "ASC",
		},
	}

	_, result, err := datamap.CreateMapFromDB(db, 100, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	// 验证客户1的订单按金额和订单ID升序排列
	customer1Key := "1"
	if row, exists := result[customer1Key]; exists {
		if amounts, ok := row.Value["amount"]; ok && len(amounts) == 3 {
			// 验证是否是升序排列
			amount1, _ := amounts[0].(float64)
			amount2, _ := amounts[1].(float64)
			amount3, _ := amounts[2].(float64)
			
			if amount1 > amount2 || amount2 > amount3 {
				t.Errorf("Amounts are not in ascending order: %f, %f, %f", amount1, amount2, amount3)
			}
		}
	}
}

// TestCreateMapFromDB_DuplicateKeys 重复key测试
func TestCreateMapFromDB_DuplicateKeys(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 创建有重复key的表
	_, err = db.Exec("DROP TABLE IF EXISTS duplicate_keys")
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
	_, err = db.Exec("CREATE TABLE duplicate_keys (group_id INT, item_name VARCHAR(50), value INT)")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	_, err = db.Exec(`
		INSERT INTO duplicate_keys VALUES 
		(1, 'A', 10),
		(1, 'B', 20),
		(1, 'C', 30),
		(2, 'D', 40),
		(2, 'E', 50)
	`)
	if err != nil {
		t.Fatalf("Failed to insert records: %v", err)
	}
	defer db.Exec("DROP TABLE IF EXISTS duplicate_keys")

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "duplicate_keys",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "group_id", ColumnAlias: "group_id", ColumnType: 1},
				{ColumnName: "item_name", ColumnAlias: "item", ColumnType: 0, ColumnLen: 50},
				{ColumnName: "value", ColumnAlias: "value", ColumnType: 1},
			},
		},
		ColumnsForMapKeys: [][]string{{"group_id"}},
		SortRule: &datamap.SortRuleForValueInSameKey{
			ColumnsForSorting: []string{"value"},
			SortDir:           "DESC",
		},
	}

	_, result, err := datamap.CreateMapFromDB(db, 100, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	// 验证分组1有3条记录
	if row, exists := result["1"]; exists {
		if items, ok := row.Value["item"]; ok {
			if len(items) != 3 {
				t.Errorf("Expected 3 items for group 1, got %d", len(items))
			}
			
			// 验证排序
			if values, ok := row.Value["value"]; ok && len(values) == 3 {
				val1, _ := values[0].(int64)
				val2, _ := values[1].(int64)
				val3, _ := values[2].(int64)
				
				if val1 < val2 || val2 < val3 {
					t.Errorf("Values are not in descending order: %d, %d, %d", val1, val2, val3)
				}
			}
		}
	}

	// 验证分组2有2条记录
	if row, exists := result["2"]; exists {
		if items, ok := row.Value["item"]; ok {
			if len(items) != 2 {
				t.Errorf("Expected 2 items for group 2, got %d", len(items))
			}
		}
	}
}

// TestCreateMapFromDB_SpecialCharacters 特殊字符测试
func TestCreateMapFromDB_SpecialCharacters(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 创建包含特殊字符的表
	_, err = db.Exec("DROP TABLE IF EXISTS special_chars")
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
	_, err = db.Exec("CREATE TABLE special_chars (id INT PRIMARY KEY, text_data VARCHAR(100))")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	_, err = db.Exec(`
		INSERT INTO special_chars VALUES 
		(1, 'Normal text'),
		(2, 'Text with ''quotes'''),
		(3, 'Text with "double quotes"'),
		(4, 'Text with \\backslashes\\'),
		(5, 'Text with -dashes-'),
		(6, 'Text with %percent%'),
		(7, 'Text with _underscores_'),
		(8, 'Text with 😀 emoji'),
		(9, 'Text with 中文'),
		(10, '')
	`)
	if err != nil {
		t.Fatalf("Failed to insert records: %v", err)
	}
	defer db.Exec("DROP TABLE IF EXISTS special_chars")

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "special_chars",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "id", ColumnType: 1},
				{ColumnName: "text_data", ColumnAlias: "text", ColumnType: 0, ColumnLen: 100},
			},
		},
		ColumnsForMapKeys: [][]string{{"id"}},
	}

	_, result, err := datamap.CreateMapFromDB(db, 100, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	// 验证特殊字符处理
	testCases := []struct {
		id   string
		text string
	}{
		{"1", "Normal text"},
		{"2", "Text with 'quotes'"},
		{"3", "Text with \"double quotes\""},
		{"4", "Text with \\backslashes\\"},
		{"5", "Text with -dashes-"},
		{"6", "Text with %percent%"},
		{"7", "Text with _underscores_"},
		{"8", "Text with 😀 emoji"},
		{"9", "Text with 中文"},
		{"10", ""},
	}

	for _, tc := range testCases {
		if row, exists := result[tc.id]; exists {
			if texts, ok := row.Value["text"]; ok && len(texts) > 0 {
				if text, ok := texts[0].(string); ok {
					if text != tc.text {
						t.Errorf("For ID %s, expected '%s', got '%s'", tc.id, tc.text, text)
					}
				} else {
					t.Errorf("For ID %s, text is not string: %T", tc.id, texts[0])
				}
			} else {
				t.Errorf("For ID %s, text not found", tc.id)
			}
		} else {
			t.Errorf("ID %s not found", tc.id)
		}
	}
}

// TestCreateMapFromDB_ColumnNameVsAlias 列名与别名映射测试
func TestCreateMapFromDB_ColumnNameVsAlias(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = initTestTable(db)
	if err != nil {
		t.Fatalf("Failed to init test table: %v", err)
	}
	defer cleanupTestTable(db)

	// 测试列名和别名混合使用的情况
	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "test_users",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "user_id", ColumnType: 5},
				{ColumnName: "username", ColumnAlias: "", ColumnType: 0, ColumnLen: 50}, // 空别名
				{ColumnName: "email", ColumnAlias: "email_address", ColumnType: 0, ColumnLen: 100},
				// 注意：这里没有包含 age 列，测试缺失列的情况
			},
		},
		ColumnsForMapKeys: [][]string{{"user_id"}},
	}

	_, result, err := datamap.CreateMapFromDB(db, 100, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	// 验证映射是否正确
	if row, exists := result["1"]; exists {
		// 验证使用别名的列
		if _, ok := row.Value["user_id"]; !ok {
			t.Error("user_id (alias) not found")
		}
		
		// 验证使用列名的列（空别名）
		if _, ok := row.Value["username"]; !ok {
			t.Error("username (column name) not found when alias is empty")
		}
		
		// 验证使用别名的列
		if _, ok := row.Value["email_address"]; !ok {
			t.Error("email_address (alias) not found")
		}
		
		// 验证未包含的列不应该存在
		if _, ok := row.Value["age"]; ok {
			t.Error("age should not be present as it's not in column config")
		}
	} else {
		t.Error("First record not found")
	}
}

// TestCreateMapFromDB_InvalidSortConfig 无效排序配置测试
func TestCreateMapFromDB_InvalidSortConfig(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = initTestTable(db)
	if err != nil {
		t.Fatalf("Failed to init test table: %v", err)
	}
	defer cleanupTestTable(db)

	// 测试1: 无效的排序方向
	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "test_users",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "user_id", ColumnType: 5},
				{ColumnName: "username", ColumnAlias: "name", ColumnType: 0, ColumnLen: 50},
			},
		},
		ColumnsForMapKeys: [][]string{{"user_id"}},
		SortRule: &datamap.SortRuleForValueInSameKey{
			ColumnsForSorting: []string{"name"},
			SortDir:           "INVALID", // 无效的排序方向
		},
	}

	_, _, err = datamap.CreateMapFromDB(db, 100, tableConfig)
	if err == nil {
		t.Error("Expected error for invalid sort direction")
	}

	// 测试2: 空的排序列
	tableConfig.SortRule.SortDir = "ASC"
	tableConfig.SortRule.ColumnsForSorting = []string{}
	_, _, err = datamap.CreateMapFromDB(db, 100, tableConfig)
	if err == nil {
		t.Error("Expected error for empty sort columns")
	}

	// 测试3: 排序列不存在于配置中
	tableConfig.SortRule.ColumnsForSorting = []string{"nonexistent_column"}
	_, _, err = datamap.CreateMapFromDB(db, 100, tableConfig)
	// 注意：这个测试可能不会立即失败，但在排序时会遇到问题
	if err != nil {
		t.Logf("Got expected error for nonexistent sort column: %v", err)
	}
}

// TestCreateMapFromDB_ConcurrentAccess 并发访问测试
func TestCreateMapFromDB_ConcurrentAccess(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = initTestTable(db)
	if err != nil {
		t.Fatalf("Failed to init test table: %v", err)
	}
	defer cleanupTestTable(db)

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "test_users",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "user_id", ColumnType: 5},
				{ColumnName: "username", ColumnAlias: "name", ColumnType: 0, ColumnLen: 50},
			},
		},
		ColumnsForMapKeys: [][]string{{"user_id"}},
	}

	// 并发调用测试
	const concurrentCalls = 5
	results := make(chan map[string]*datamap.Row, concurrentCalls)
	errors := make(chan error, concurrentCalls)

	for i := 0; i < concurrentCalls; i++ {
		go func() {
			_, result, err := datamap.CreateMapFromDB(db, 100, tableConfig)
			if err != nil {
				errors <- err
				return
			}
			results <- result
		}()
	}

	// 收集结果
	successCount := 0
	errorCount := 0

	timeout := time.After(30 * time.Second)
	for i := 0; i < concurrentCalls; i++ {
		select {
		case result := <-results:
			if len(result) == 5 {
				successCount++
			} else {
				t.Errorf("Concurrent call returned %d records, expected 5", len(result))
			}
		case err := <-errors:
			t.Errorf("Concurrent call failed: %v", err)
			errorCount++
		case <-timeout:
			t.Fatal("Timeout waiting for concurrent calls")
		}
	}

	if successCount != concurrentCalls {
		t.Errorf("Expected %d successful concurrent calls, got %d", concurrentCalls, successCount)
	}
}

// TestCreateMapFromDB_Performance 性能测试（大数据量）
func TestCreateMapFromDB_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 创建大数据量表
	tableName := "performance_test"
	_, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE %s (
			id INT PRIMARY KEY AUTO_INCREMENT,
			group_id INT,
			data_value DOUBLE,
			description VARCHAR(100)
		)`, tableName))
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	defer db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))

	// 插入1000条测试数据
	const recordCount = 1000
	for i := 0; i < recordCount; i++ {
		_, err = db.Exec(fmt.Sprintf(`
			INSERT INTO %s (group_id, data_value, description) 
			VALUES (?, ?, ?)`, tableName), 
			i%10, // 10个不同的组
			float64(i)*1.5, 
			fmt.Sprintf("Record %d", i))
		if err != nil {
			t.Fatalf("Failed to insert record %d: %v", i, err)
		}
	}

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: tableName,
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "id", ColumnType: 1},
				{ColumnName: "group_id", ColumnAlias: "group", ColumnType: 1},
				{ColumnName: "data_value", ColumnAlias: "value", ColumnType: 2},
				{ColumnName: "description", ColumnAlias: "desc", ColumnType: 0, ColumnLen: 100},
			},
		},
		ColumnsForMapKeys: [][]string{{"group"}},
		SortRule: &datamap.SortRuleForValueInSameKey{
			ColumnsForSorting: []string{"value"},
			SortDir:           "ASC",
		},
	}

	// 测量执行时间
	start := time.Now()
	_, result, err := datamap.CreateMapFromDB(db, 500, tableConfig) // 使用中等页大小
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	// 验证结果
	totalRecords := 0
	for _, row := range result {
		if ids, ok := row.Value["id"]; ok {
			totalRecords += len(ids)
		}
	}

	if totalRecords != recordCount {
		t.Errorf("Expected %d total records, got %d", recordCount, totalRecords)
	}

	// 验证分组数量
	if len(result) != 10 {
		t.Errorf("Expected 10 groups, got %d", len(result))
	}

	t.Logf("Processed %d records in %d groups in %v", totalRecords, len(result), elapsed)

	// 性能断言（可根据实际情况调整）
	if elapsed > 10*time.Second {
		t.Errorf("Performance test took too long: %v", elapsed)
	}
}

// TestCreateMapFromDB_MemoryUsage 内存使用测试
func TestCreateMapFromDB_MemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory usage test in short mode")
	}

	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 创建一个包含各种数据类型的表来测试内存使用
	_, err = db.Exec("DROP TABLE IF EXISTS memory_test")
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE memory_test (
			id INT PRIMARY KEY,
			short_text VARCHAR(10),
			long_text MEDIUMTEXT,
			int_value INT,
			bigint_value BIGINT,
			float_value DOUBLE,
			bool_value TINYINT(1)
		)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	defer db.Exec("DROP TABLE IF EXISTS memory_test")

	// 插入包含各种数据类型的记录
	for i := 1; i <= 100; i++ {
		_, err = db.Exec(`
			INSERT INTO memory_test VALUES (?, ?, ?, ?, ?, ?, ?)`,
			i,
			fmt.Sprintf("Text%d", i),
			fmt.Sprintf("This is a longer text for record %d with some content to test memory usage", i),
			i*10,
			int64(i)*100,
			float64(i)*1.2345,
			i%2,
		)
		if err != nil {
			t.Fatalf("Failed to insert record %d: %v", i, err)
		}
	}

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "memory_test",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "id", ColumnType: 1},
				{ColumnName: "short_text", ColumnAlias: "short", ColumnType: 0, ColumnLen: 10},
				{ColumnName: "long_text", ColumnAlias: "long", ColumnType: 4},
				{ColumnName: "int_value", ColumnAlias: "int_val", ColumnType: 1},
				{ColumnName: "bigint_value", ColumnAlias: "bigint_val", ColumnType: 5},
				{ColumnName: "float_value", ColumnAlias: "float_val", ColumnType: 2},
				{ColumnName: "bool_value", ColumnAlias: "bool_val", ColumnType: 3},
			},
		},
		ColumnsForMapKeys: [][]string{{"id"}},
	}

	// 运行多次测试稳定性
	for i := 0; i < 3; i++ {
		_, result, err := datamap.CreateMapFromDB(db, 50, tableConfig)
		if err != nil {
			t.Fatalf("datamap.CreateMapFromDB failed on iteration %d: %v", i, err)
		}

		if len(result) != 100 {
			t.Errorf("Iteration %d: Expected 100 records, got %d", i, len(result))
		}
	}
}

// TestCreateMapFromDB_EdgeCases 边界情况测试
func TestCreateMapFromDB_EdgeCases(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 测试1: 零值和极值
	_, err = db.Exec("DROP TABLE IF EXISTS edge_cases")
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE edge_cases (
			id INT PRIMARY KEY,
			zero_int INT,
			max_int BIGINT,
			zero_float DOUBLE,
			large_float DOUBLE,
			empty_string VARCHAR(50),
			null_value VARCHAR(50)
		)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	_, err = db.Exec(`
		INSERT INTO edge_cases VALUES 
		(1, 0, 9223372036854775807, 0.0, 1.7976931348623157e+308, '', NULL)
	`)
	if err != nil {
		t.Fatalf("Failed to insert records: %v", err)
	}
	defer db.Exec("DROP TABLE IF EXISTS edge_cases")

	tableConfig := &datamap.TableConfigForMapping{
		TableConfig: &dbutil.TableConfig{
			TableName: "edge_cases",
			Columns: []*dbutil.ColumnConfig{
				{ColumnName: "id", ColumnAlias: "id", ColumnType: 1},
				{ColumnName: "zero_int", ColumnAlias: "zero_int", ColumnType: 1},
				{ColumnName: "max_int", ColumnAlias: "max_int", ColumnType: 5},
				{ColumnName: "zero_float", ColumnAlias: "zero_float", ColumnType: 2},
				{ColumnName: "large_float", ColumnAlias: "large_float", ColumnType: 2},
				{ColumnName: "empty_string", ColumnAlias: "empty_str", ColumnType: 0, ColumnLen: 50},
				{ColumnName: "null_value", ColumnAlias: "null_val", ColumnType: 0, ColumnLen: 50},
			},
		},
		ColumnsForMapKeys: [][]string{{"id"}},
	}

	_, result, err := datamap.CreateMapFromDB(db, 100, tableConfig)
	if err != nil {
		t.Fatalf("datamap.CreateMapFromDB failed: %v", err)
	}

	if row, exists := result["1"]; exists {
		// 验证零值
		if zeroInts, ok := row.Value["zero_int"]; ok && len(zeroInts) > 0 {
			if zeroInt, ok := zeroInts[0].(int64); ok {
				if zeroInt != 0 {
					t.Errorf("Expected zero_int to be 0, got %d", zeroInt)
				}
			}
		}

		// 验证最大值
		if maxInts, ok := row.Value["max_int"]; ok && len(maxInts) > 0 {
			if maxInt, ok := maxInts[0].(int64); ok {
				if maxInt != 9223372036854775807 {
					t.Errorf("Expected max_int to be 9223372036854775807, got %d", maxInt)
				}
			}
		}

		// 验证空字符串
		if emptyStrs, ok := row.Value["empty_str"]; ok && len(emptyStrs) > 0 {
			if emptyStr, ok := emptyStrs[0].(string); ok {
				if emptyStr != "" {
					t.Errorf("Expected empty_str to be empty, got '%s'", emptyStr)
				}
			}
		}

		// 验证NULL值
		if nullVals, ok := row.Value["null_val"]; ok && len(nullVals) > 0 {
			if nullVal, ok := nullVals[0].(string); ok {
				if nullVal != "" {
					t.Errorf("Expected null_val to be empty string for NULL, got '%s'", nullVal)
				}
			}
		}
	}
}
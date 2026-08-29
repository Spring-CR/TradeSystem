package dbutil

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"rhino-common/domain_error"
	"strings"
)

func CreateSimpleTableWithSysID(db *sql.DB, tableName string, recreate bool) error {

	// 检查并删除已存在的目标表
	var exists int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?",
		tableName,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check table existence: %w", err)
	}

	if exists > 0 {
		if recreate {
			if _, err := db.Exec(fmt.Sprintf("DROP TABLE `%s`", tableName)); err != nil {
				return fmt.Errorf("failed to drop existing table: %w", err)
			}
		} else {
			return nil
		}
	}

	createTableSql := `CREATE TABLE IF NOT EXISTS %s (
    f_id BIGINT PRIMARY KEY AUTO_INCREMENT
);`
	createTableSql = fmt.Sprintf(createTableSql, tableName)

	// 初始化DB
	_, err = db.Exec(createTableSql)

	return err

}

func CopyTableStructure(originalDB *sql.DB, targetDB *sql.DB, originalTableName string, targetTableName string, recreate bool) error {
	// 检查并删除已存在的目标表
	var exists int
	err := targetDB.QueryRow(
		"SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?",
		targetTableName,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check table existence: %w", err)
	}

	if exists > 0 {
		if recreate {
			if _, err := targetDB.Exec(fmt.Sprintf("DROP TABLE `%s`", targetTableName)); err != nil {
				return fmt.Errorf("failed to drop existing table: %w", err)
			}
		} else {
			// 如果不是recreate，更新表
			return forceSyncTableStructure(originalDB, targetDB, originalTableName, targetTableName)
		}
	}

	// 获取原表创建语句
	var createTableSQL string
	query := fmt.Sprintf("SHOW CREATE TABLE `%s`", originalTableName)
	if err := originalDB.QueryRow(query).Scan(&originalTableName, &createTableSQL); err != nil {
		return fmt.Errorf("failed to get create table statement: %w", err)
	}

	// 修改建表语句
	newCreateSQL := createTableSQL

	// 步骤1：替换表名
	newCreateSQL = strings.Replace(
		newCreateSQL,
		fmt.Sprintf("`%s`", originalTableName),
		fmt.Sprintf("`%s`", targetTableName),
		1,
	)

	// 步骤2：移除AUTO_INCREMENT定义
	autoIncrementRegex := regexp.MustCompile(`(?i)\s*AUTO_INCREMENT\s*=\s*\d+`)
	newCreateSQL = autoIncrementRegex.ReplaceAllString(newCreateSQL, "")

	log.Printf("Create target table SQL:\n%s\n", newCreateSQL)

	if _, err := targetDB.Exec(newCreateSQL); err != nil {
		return fmt.Errorf("failed to create new table: %w", err)
	}

	return nil
}

// 核心函数：强制同步表结构（独立执行每个操作）
func forceSyncTableStructure(originalDB *sql.DB, targetDB *sql.DB, srcTable, dstTable string) error {
    // 获取原表结构
    srcCreateSQL, err := getCreateTableSQL(originalDB, srcTable)
    if err != nil {
        return err
    }

    // 获取目标表结构
    dstCreateSQL, err := getCreateTableSQL(targetDB, dstTable)
    if err != nil {
        return err
    }

    // 提取字段定义
    srcColumns, err := extractColumnDefinitions(srcCreateSQL)
    if err != nil {
        return fmt.Errorf("failed to parse source table columns: %w", err)
    }

    dstColumns, err := extractColumnDefinitions(dstCreateSQL)
    if err != nil {
        return fmt.Errorf("failed to parse target table columns: %w", err)
    }

    // 收集所有需要执行的SQL语句
    var sqlStatements []string
    
    // 1. 添加缺失字段
    for colName, def := range srcColumns {
        if _, exists := dstColumns[colName]; !exists {
            sqlStatements = append(sqlStatements, 
                fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s", dstTable, def))
        }
    }
    
    // 2. 强制更新所有存在的字段
    for colName, srcDef := range srcColumns {
        if _, exists := dstColumns[colName]; exists {
            sqlStatements = append(sqlStatements, 
                fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN %s", dstTable, srcDef))
        }
    }
    
    // 如果没有变更，直接返回
    if len(sqlStatements) == 0 {
        log.Println("Table structure is already up-to-date")
        return nil
    }
    
    // 3. 独立执行每个SQL语句
    var errors []error
    
    for _, sql := range sqlStatements {
        log.Printf("Executing: %s\n", sql)
        if _, err := targetDB.Exec(sql); err != nil {
            // 记录错误但继续执行
            errMsg := fmt.Errorf("failed to execute '%s': %w", sql, err)
            log.Printf("Error: %v\n", errMsg)
            errors = append(errors, errMsg)
        }
    }
    
    // 返回第一个错误（如果有）
    if len(errors) > 0 {
        err = fmt.Errorf("%d errors occurred during sync, first error: %w", len(errors), errors[0])
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs while forceSyncTableStructure")
		return err
    }
    
    return nil
}

func getCreateTableSQL(db *sql.DB, tableName string) (string, error) {
    var table, createSQL string
    query := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
    if err := db.QueryRow(query).Scan(&table, &createSQL); err != nil {
        return "", fmt.Errorf("failed to get create table statement for %s: %w", tableName, err)
    }
    return createSQL, nil
}


var(
	ignoreField = map[string]bool{
		"PRIMARY":true,
		"UNIQUE":true,
		"KEY":true,
	}
)
func extractColumnDefinitions(createSQL string) (map[string]string, error) {
	start := strings.Index(createSQL, "(")
	end := strings.LastIndex(createSQL, ")")
	if start == -1 || end == -1 || start >= end {
		return nil, fmt.Errorf("invalid CREATE TABLE statement")
	}

	inner := createSQL[start+1 : end]
	definitions := splitColumnDefinitions(inner)

	result := make(map[string]string)
	for _, def := range definitions {
		if def == "" {
			continue
		}

		parts := strings.Fields(def)
		if len(parts) < 1 {
			continue
		}

		colName := strings.Trim(parts[0], "`")

		f := strings.ToUpper(colName)
		if ignoreField[f] {
			return result, nil
		}

		result[colName] = def
	}

	return result, nil
}

func splitColumnDefinitions(inner string) []string {
	var result []string
	var current strings.Builder
	depth := 0

	for _, r := range inner {
		switch r {
		case '(':
			depth++
			current.WriteRune(r)
		case ')':
			depth--
			current.WriteRune(r)
		case ',':
			if depth == 0 {
				result = append(result, strings.TrimSpace(current.String()))
				current.Reset()
				continue
			}
			current.WriteRune(r)
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		result = append(result, strings.TrimSpace(current.String()))
	}

	return result
}
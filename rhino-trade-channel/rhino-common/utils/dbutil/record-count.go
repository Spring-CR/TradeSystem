package dbutil

import (
	"database/sql"
	"fmt"
	"regexp"
)

// getRecordCount 返回指定表的记录总数。
// 注意：tableName 应当来自可信源，或经过验证，以避免 SQL 注入风险。
// 本实现通过正则表达式简单校验表名仅包含字母、数字和下划线，拒绝其他字符。
func GetRecordCount(db *sql.DB, tableName string) (int64, error) {
	// 校验表名是否合法（仅允许字母、数字、下划线）
	validTableName := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validTableName.MatchString(tableName) {
		return 0, fmt.Errorf("invalid table name: %s", tableName)
	}

	var count int64
	// 使用 fmt.Sprintf 拼接 SQL，因为表名不能作为参数化查询的占位符
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
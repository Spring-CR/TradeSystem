package dbutil

import (
	"database/sql"
	"fmt"
	"strings"
)

func TruncateTableIfExist(db *sql.DB, tableName string) error {

	// 1. 验证表是否存在
	var exists int
	err := db.QueryRow(
		"SELECT 1 FROM information_schema.tables WHERE table_name = ?",
		tableName,
	).Scan(&exists)

	if err == sql.ErrNoRows || err == nil && exists != 1 { // 不存在的表，直接返回
		return nil
	} else if err != nil {
		return fmt.Errorf("check table %s existence failed, error: %v", tableName, err)
	}

	// 2. 执行TRUNCATE
	if _, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", tableName)); err != nil {
		switch {
		case strings.Contains(err.Error(), "DROP privilege"):
			return fmt.Errorf("insufficient privileges to truncate table %s", tableName)
		case strings.Contains(err.Error(), "LOCK"):
			return fmt.Errorf("table %s is locked by another process", tableName)
		default:
			return fmt.Errorf("truncate table %s failed, error: %v", tableName, err)
		}
	}

	return nil
}

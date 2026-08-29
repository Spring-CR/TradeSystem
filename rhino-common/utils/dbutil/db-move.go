package dbutil

import (
	"database/sql"
	"fmt"
)

func MoveTable(db *sql.DB, srcTable string, destTable string, dateStr string) error {

	tmpTable := destTable + "_backup_" + dateStr

	// 0. 先检查备份表是否已存在
	var tmpExists int
	err := db.QueryRow("SELECT 1 FROM information_schema.tables WHERE table_name = ?", tmpTable).Scan(&tmpExists)

	// 1. 仅当备份表不存在时才创建备份
	if err == sql.ErrNoRows || err == nil && tmpExists != 1 {
		if _, err := db.Exec(fmt.Sprintf("RENAME TABLE %s TO %s", destTable, tmpTable)); err != nil {
			return fmt.Errorf("backup table %s failed, error: %v", destTable, err)
		}
	} else if err != nil {
		return fmt.Errorf("check backup table %s exists failed, error: %v", tmpTable, err)
	}

	// 2. 备份完成后，正式删除原表
	if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", destTable)); err != nil {
		return fmt.Errorf("drop table %s failed, error: %v", destTable, err)
	}

	// 3. 删除完成后，从命名表，最终达到替换表的目标
	if _, err := db.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", srcTable, destTable)); err != nil {
		// 恢复逻辑需要判断是否有实际创建过备份表
		if _, recoverErr := db.Exec(fmt.Sprintf("RENAME TABLE %s TO %s", tmpTable, destTable)); recoverErr != nil {
			return fmt.Errorf("CRITICAL: failed to recover table after rename table failed: %v", err)
		} else {
			return fmt.Errorf("rename table %s as %s (so recover table %s as %s) failed, error:%v", srcTable, destTable, tmpTable, destTable, err)
		}
	}

	// 4. 仅在成功创建备份时清理
	if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tmpTable)); err != nil {
		return fmt.Errorf("drop backup table %s failed, error: %v", tmpTable, err)
	}

	return nil
}

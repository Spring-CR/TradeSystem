package dbutil

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"
)

var identifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// validateIdentifier 检查标识符是否安全，防止SQL注入。
func validateIdentifier(name string) error {
	if name == "" {
		return fmt.Errorf("identifier cannot be empty")
	}
	// 只允许字母、数字、下划线，且不能以数字开头
	matched := identifierPattern.MatchString(name)
	if !matched {
		return fmt.Errorf("invalid identifier: %s", name)
	}
	return nil
}

// placeholders 生成 n 个 "?" 占位符，用逗号连接。
func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat("?,", n-1) + "?"
}

// batchSlice 将 []T 按指定大小切分成多个批次。
func batchSlice[T any](slice []T, size int) [][]T {
	if size <= 0 || len(slice) == 0 {
		return nil
	}
	var batches [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		batches = append(batches, slice[i:end])
	}
	return batches
}

// validateIDSet 检查 idSet 中是否包含 0，避免无效删除请求。
func validateIDSet(idSet []int64) error {
	for _, id := range idSet {
		if id == 0 {
			return fmt.Errorf("idSet contains invalid zero value: %d", id)
		}
	}
	return nil
}

// validateKeySet 检查 keySet 中是否包含空字符串，避免无效删除请求。
func validateKeySet(keySet []string) error {
	for _, key := range keySet {
		if key == "" {
			return fmt.Errorf("keySet contains invalid empty value")
		}
	}
	return nil
}

// DeleteByIDsInBatches 根据指定的 id 列分批删除数据。
// 不存在的 id 会被自动忽略；如果删除后表为空，则自动执行 TRUNCATE。
func DeleteByIDsInBatches(db *sql.DB, tableName string, idColumnName string, batchSize int, idSet []int64, autoTruncate bool) error {

	log.Printf("DeleteByIDsInBatches, tableName:%s, idColumnName:%s, batchSize:%v, autoTruncate:%v, idSet:%v", tableName, idColumnName, batchSize, autoTruncate, idSet)

	// 1. 参数校验
	if err := validateIdentifier(tableName); err != nil {
		return fmt.Errorf("invalid table name: %w", err)
	}
	if err := validateIdentifier(idColumnName); err != nil {
		return fmt.Errorf("invalid id column name: %w", err)
	}
	if batchSize <= 0 {
		return fmt.Errorf("batchSize must be greater than 0")
	}
	if len(idSet) == 0 {
		return nil
	}
	if err := validateIDSet(idSet); err != nil {
		return fmt.Errorf("invalid idSet: %w", err)
	}

	// 2. 分批删除
	batches := batchSlice(idSet, batchSize)
	// 预编译 DELETE 模板，留下 %s 给占位符
	deleteTpl := fmt.Sprintf("DELETE FROM %s WHERE %s IN (%%s)", tableName, idColumnName)
	countQueryTpl := fmt.Sprintf("Select count(*) FROM %s WHERE %s IN (%%s)", tableName, idColumnName)

	for _, batch := range batches {
		ph := placeholders(len(batch))
		query := fmt.Sprintf(deleteTpl, ph)

		log.Printf("Delete for sub-batch, tableName:%s, batch:%v", tableName, batch)

		args := make([]interface{}, len(batch))
		for i, id := range batch {
			args[i] = id
		}

		if result, err := db.Exec(query, args...); err != nil {
			return fmt.Errorf("delete batch failed, query: %s, args: %v, error: %w", query, args, err)
		} else {
			// 检查删除的行数是否与预期一致，如果不一致，记录日志但不返回错误
			rowsAffected, err := result.RowsAffected()
			if err == nil && rowsAffected != int64(len(batch)) {
				log.Printf("expected %d rows affected, got %d, query: %s, args: %v", len(batch), rowsAffected, query, args)
			}
		}

		// 检查是否有数据未被删除，如果有，返回错误
		var rowCount int
		countQuery := fmt.Sprintf(countQueryTpl, ph)
		if err := db.QueryRow(countQuery, args...).Scan(&rowCount); err != nil {
			return fmt.Errorf("count batch failed, query: %s, args: %v, error: %w", countQuery, args, err)
		}
		if rowCount > 0 {
			return fmt.Errorf("delete batch failed, some ids still exist, query: %s, args: %v", countQuery, args)
		}
	}

	if !autoTruncate {
		return nil
	}

	// 3. 全部删除后，如果表为空则执行 TRUNCATE
	var rowCount int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	if err := db.QueryRow(countQuery).Scan(&rowCount); err != nil {
		return fmt.Errorf("count rows failed, table name: %s, error: %w", tableName, err)
	}

	if rowCount == 0 {
		truncateQuery := fmt.Sprintf("TRUNCATE TABLE %s", tableName)
		if _, err := db.Exec(truncateQuery); err != nil {
			return fmt.Errorf("truncate table failed, table name: %s, error: %w", tableName, err)
		}
		log.Printf("Table truncated: %s", tableName)
	}

	return nil
}

// trimKeySet 去除 keySet 中的元素的空格，并返回一个新的切片。
func trimKeySet(keySet []string) []string {
	if len(keySet) == 0 {
		return keySet
	}
	var trimmed []string
	for _, key := range keySet {
		str := strings.TrimSpace(key)
		if str == "" {
			continue
		}
		trimmed = append(trimmed, str)
	}
	return trimmed
}

// DeleteByKeysInBatches 根据指定的 key 列分批删除数据。
// 不存在的 key 会被自动忽略；如果删除后表为空，则自动执行 TRUNCATE。
func DeleteByKeysInBatches(db *sql.DB, tableName string, columnName string, batchSize int, keySet []string, autoTruncate bool) error {

	log.Printf("DeleteByKeysInBatches, tableName:%s, columnName:%s, batchSize:%v, autoTruncate:%v, keySet:%v", tableName, columnName, batchSize, autoTruncate, keySet)

	// 1. 参数校验
	if err := validateIdentifier(tableName); err != nil {
		return fmt.Errorf("invalid table name: %w", err)
	}
	if err := validateIdentifier(columnName); err != nil {
		return fmt.Errorf("invalid id column name: %w", err)
	}
	if batchSize <= 0 {
		return fmt.Errorf("batchSize must be greater than 0")
	}
	if len(keySet) == 0 {
		return nil
	}
	keySet = trimKeySet(keySet)
	if err := validateKeySet(keySet); err != nil {
		return fmt.Errorf("invalid keySet: %w", err)
	}

	// 2. 分批删除
	batches := batchSlice(keySet, batchSize)
	// 预编译 DELETE 模板，留下 %s 给占位符
	deleteTpl := fmt.Sprintf("DELETE FROM %s WHERE %s IN (%%s)", tableName, columnName)
	countQueryTpl := fmt.Sprintf("Select count(*) FROM %s WHERE %s IN (%%s)", tableName, columnName)

	for _, batch := range batches {
		ph := placeholders(len(batch))
		query := fmt.Sprintf(deleteTpl, ph)

		log.Printf("Delete for sub-batch, tableName:%s, batch:%v", tableName, batch)

		args := make([]interface{}, len(batch))
		for i, id := range batch {
			args[i] = id
		}

		if result, err := db.Exec(query, args...); err != nil {
			return fmt.Errorf("delete batch failed, query: %s, args: %v, error: %w", query, args, err)
		} else {
			// 检查删除的行数是否与预期一致，如果不一致，记录日志但不返回错误
			rowsAffected, err := result.RowsAffected()
			if err == nil && rowsAffected != int64(len(batch)) {
				log.Printf("expected %d rows affected, got %d, query: %s, args: %v", len(batch), rowsAffected, query, args)
			}
		}

		// 检查是否有数据未被删除，如果有，返回错误
		var rowCount int
		countQuery := fmt.Sprintf(countQueryTpl, ph)
		if err := db.QueryRow(countQuery, args...).Scan(&rowCount); err != nil {
			return fmt.Errorf("count batch failed, query: %s, args: %v, error: %w", countQuery, args, err)
		}
		if rowCount > 0 {
			return fmt.Errorf("delete batch failed, some ids still exist, query: %s, args: %v", countQuery, args)
		}
	}

	if !autoTruncate {
		return nil
	}

	// 3. 全部删除后，如果表为空则执行 TRUNCATE
	var rowCount int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	if err := db.QueryRow(countQuery).Scan(&rowCount); err != nil {
		return fmt.Errorf("count rows failed, table name: %s, error: %w", tableName, err)
	}

	if rowCount == 0 {
		truncateQuery := fmt.Sprintf("TRUNCATE TABLE %s", tableName)
		if _, err := db.Exec(truncateQuery); err != nil {
			return fmt.Errorf("truncate table failed, table name: %s, error: %w", tableName, err)
		}
		log.Printf("Table truncated: %s", tableName)
	}

	return nil
}

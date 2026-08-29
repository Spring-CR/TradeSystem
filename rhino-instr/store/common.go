package store

import "strings"

const (
	PaggingSuffix = ` LIMIT ? OFFSET ? `
)

func isPagging(limit, offset int) bool {
	return limit > 0 && offset >= 0 
}

func getCountStmtFromOriginalStmt(originalStmt string) string {
	tmpStr := strings.ToUpper(originalStmt)
	fromIndex := strings.Index(tmpStr, "FROM")
	if fromIndex >= 0 {
		return "SELECT COUNT(1) as total_records " + originalStmt[fromIndex:]
	}
	return ""
}
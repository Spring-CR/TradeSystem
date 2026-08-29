package app_store

import (
	"rhino-core/schema"
	"strings"

	"github.com/linchunquan/sqlgen/db"
)

func GetArchiveInsertTradeOrderStmt(tableName string) string {
	return strings.Replace(InsertTradeOrderStmt, "INSERT INTO trade_orders", "INSERT INTO "+tableName, 1)
}

func ArchiveInsertTradeOrder(db db.SimpleDB, insertTradeOrderStmt string, v *schema.TradeOrder, archivingLog *schema.DataArchivingLog) error {
	args := sliceTradeOrder(v)
	if archivingLog != nil {
		args = append(args, archivingLog.ArchivingDate, archivingLog.TaskName)
	}
	res, err := db.Exec(GetArchiveInsertRecordStmt(insertTradeOrderStmt, archivingLog), args[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func GetArchiveInsertTradeActionLatestRespStmt(tableName string) string {
	return strings.Replace(InsertTradeActionLatestRespStmt, "INSERT INTO trade_action_latest_resps", "INSERT INTO "+tableName, 1)
}

func ArchiveInsertTradeActionLatestResp(db db.SimpleDB, insertTradeActionLatestRespStmt string, v *schema.TradeActionLatestResp, archivingLog *schema.DataArchivingLog) error {
	args := sliceTradeActionLatestResp(v)
	if archivingLog != nil {
		args = append(args, archivingLog.ArchivingDate, archivingLog.TaskName)
	}
	res, err := db.Exec(GetArchiveInsertRecordStmt(insertTradeActionLatestRespStmt, archivingLog), args[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func GetArchiveInsertTradeActionRespStmt(tableName string) string {
	return strings.Replace(InsertTradeActionRespStmt, "INSERT INTO trade_action_resps", "INSERT INTO "+tableName, 1)
}

func ArchiveInsertTradeActionResp(db db.SimpleDB, insertTradeActionRespStmt string, v *schema.TradeActionResp, archivingLog *schema.DataArchivingLog) error {
	args := sliceTradeActionResp(v)
	if archivingLog != nil {
		args = append(args, archivingLog.ArchivingDate, archivingLog.TaskName)
	}
	res, err := db.Exec(GetArchiveInsertRecordStmt(insertTradeActionRespStmt, archivingLog), args[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func GetArchiveInsertRecordStmt(rawStmt string, archivingLog *schema.DataArchivingLog) string {
	if archivingLog == nil {
		return rawStmt
	}

	rawStmt = strings.Replace(rawStmt, ") VALUES (", ", f_archiving_date, f_task_name) VALUES (", 1)
	rawStmt = strings.Replace(rawStmt, ",?)", ",?,?,?)", 1)

	return rawStmt
}

func DeleteHistoricalArchiveDate(db db.SimpleDB, tableName, archivingDate string, taskName string) error {
	sql := `DELETE FROM ` + tableName + ` WHERE f_archiving_date=? and f_task_name=?`
	args := []interface{}{archivingDate, taskName}
	_, err := db.Exec(sql, args...)
	return err
}

package app_store

import (
	"rhino-core/schema"
	"strings"

	"github.com/linchunquan/sqlgen/db"
)

func GetArchiveInsertTradeOrderStmt(tableName string) string {
	return strings.Replace(InsertTradeOrderStmt, "INSERT INTO trade_orders", "INSERT INTO "+tableName, 1)
}

func ArchiveInsertTradeOrder(db db.SimpleDB, insertTradeOrderStmt string, v *schema.TradeOrder, archivingDate ...string) error {
	args := sliceTradeOrder(v)
	if len(archivingDate) > 0 {
		args = append(args, archivingDate[0])
	}
	res, err := db.Exec(GetArchiveInsertRecordStmt(insertTradeOrderStmt, archivingDate...), args[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func GetArchiveInsertTradeActionLatestRespStmt(tableName string) string {
	return strings.Replace(InsertTradeActionLatestRespStmt, "INSERT INTO trade_action_latest_resps", "INSERT INTO "+tableName, 1)
}

func ArchiveInsertTradeActionLatestResp(db db.SimpleDB, insertTradeActionLatestRespStmt string, v *schema.TradeActionLatestResp, archivingDate ...string) error {
	args := sliceTradeActionLatestResp(v)
	if len(archivingDate) > 0 {
		args = append(args, archivingDate[0])
	}
	res, err := db.Exec(GetArchiveInsertRecordStmt(insertTradeActionLatestRespStmt, archivingDate...), args[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func GetArchiveInsertTradeActionRespStmt(tableName string) string {
	return strings.Replace(InsertTradeActionRespStmt, "INSERT INTO trade_action_resps", "INSERT INTO "+tableName, 1)
}

func ArchiveInsertTradeActionResp(db db.SimpleDB, insertTradeActionRespStmt string, v *schema.TradeActionResp, archivingDate ...string) error {
	args := sliceTradeActionResp(v)
	if len(archivingDate) > 0 {
		args = append(args, archivingDate[0])
	}
	res, err := db.Exec(GetArchiveInsertRecordStmt(insertTradeActionRespStmt, archivingDate...), args[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func GetArchiveInsertRecordStmt(rawStmt string, archivingDate ...string) string {
	if len(archivingDate) == 0 {
		return rawStmt
	}
	rawStmt = strings.Replace(rawStmt, ") VALUES (", ", f_archiving_date) VALUES (", 1)
	rawStmt = strings.Replace(rawStmt, ",?)", ",?,?)", 1)
	return rawStmt
}

func DeleteHistoricalArchiveDate(db db.SimpleDB, tableName, archivingDate string) error {
	sql := `DELETE FROM ` + tableName + ` WHERE f_archiving_date=?`
	args := []interface{}{archivingDate}
	_, err := db.Exec(sql, args...)
	return err
}

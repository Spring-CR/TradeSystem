package store

import (
	"database/sql"
	"rhino-common/utils/dbutil"
	"rhino-instr/schema"
	"strings"

	"github.com/linchunquan/sqlgen/db"
)

var (
	LockTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNoStmt = SelectTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNoStmt + ` for update`
)

func GetAndLockTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNo(db db.SimpleDB, date int, dailyInstrNo int64, indexDailyModify int64, stockSerialNo int64) (*schema.TaskInstrStock, error) {
	args := []interface{}{date, dailyInstrNo, indexDailyModify, stockSerialNo}
	v, err := genericSelectTaskInstrStock(db, LockTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNoStmt, args...)
	return v, err
}

func GetMaxKafkaOffset(db db.SimpleDB) (int64, error) {
	var maxOffset sql.NullInt64
	sql := "select max(f_status_kafka_offset) from trade_instr_resps"
	row := db.QueryRow(sql)
	err := row.Scan(&maxOffset)
	if dbutil.IsDbRecordEmptyError(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	if maxOffset.Valid {
		return maxOffset.Int64, nil
	}
	return 0, nil
}

// func FindTradeInstrsBeginWithParentKey(db db.SimpleDB, parentKey string) ([]*schema.TradeInstr, error) {
// 	args := []interface{}{parentKey}
// 	v, err := genericSelectTradeInstrs(db, SelectTradeInstrByParentKeyStmt, args...)
// 	return v, err
// }

func findTradeInstrsBeginWithParentKey(db db.SimpleDB, parentKey string) ([]*schema.TradeInstr, error) {
	args := []interface{}{parentKey+"%"}
	v, err := genericSelectTradeInstrs(db, strings.ReplaceAll(SelectTradeInstrByParentKeyStmt, "f_parent_key=?", "f_parent_key like ?"), args...)
	return v, err
}

func FindTradeDeskOrderIdsByParentKey(db db.SimpleDB, parentKey string) (ids []string, err error) {

	var tradeInstrs []*schema.TradeInstr

	strs := strings.Split(parentKey,"-")
	if len(strs) == 4 {
		tradeInstrs, err = FindTradeInstrsByParentKey(db, parentKey)
	} else {
		tradeInstrs, err = findTradeInstrsBeginWithParentKey(db, parentKey)
	}
	
	if dbutil.IsDbRecordEmptyError(err) {
		err = nil
	}
	if err != nil {
		return
	}
	for _, tradeInstr := range tradeInstrs {
		id := tradeInstr.ClOrdID
		// 异常的委托也需要展示出来
		//if id != "" && strings.HasPrefix(id, "EQ:") {
			ids = append(ids, id)
		//}
	}
	return
}

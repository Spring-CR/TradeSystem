package store

import (
	"rhino-instr/schema"
	"strings"

	"github.com/linchunquan/sqlgen/db"
)

var (
	CreateInsertTaskInstrStmt = ""
)

func init() {
	i := strings.Index(InsertTaskInstrStmt, "VALUES")
	if i > 0 {
		insertIntoPart := InsertTaskInstrStmt[:i]
		placeholdersPart := strings.ReplaceAll(InsertTaskInstrStmt[i+6:], "(", "")
		placeholdersPart = strings.ReplaceAll(placeholdersPart, ")", "")
		placeholdersPart = strings.ReplaceAll(placeholdersPart, " ", "")
		placeholdersPart = strings.TrimLeft(placeholdersPart, "?,")

		CreateInsertTaskInstrStmt = insertIntoPart + ` SELECT COUNT(f_id) + 1, ` + placeholdersPart + " from task_instrs WHERE f_date = ? and f_index_daily_modify=1"

	} else {
		panic("InsertTaskInstrStmt is not correct!")
	}
}

// 创建指令时插入数据，要确保，指令的ID，每天从1开始，全系统当日唯一
func CreateInsertTaskInstr(db db.SimpleDB, v *schema.TaskInstr) error {
	args := append(sliceTaskInstr(v)[2:], v.Date)
	res, err := db.Exec(CreateInsertTaskInstrStmt, args...)
	if err != nil {
		return err
	}
	v.ID, err = res.LastInsertId()
	return err
}

func LockCountForCreateInsertTaskInstr(db db.SimpleDB, date int) (int, error) {
	var count int
	//err := db.QueryRow(`SELECT COUNT(f_id) FROM task_instrs WHERE f_date = ? and f_index_daily_modify=1 FOR UPDATE`, date).Scan(&count)
	err := db.QueryRow(`SELECT COUNT(f_id) FROM task_instrs WHERE f_date = ? and f_index_daily_modify=1`, date).Scan(&count)
	return count, err
}

func LockCountOfStocksForInsertTaskInstr(db db.SimpleDB, date int, dailyInstrNo, indexDailyModify int64) (int, error) {
	var count int
	//err := db.QueryRow(`SELECT COUNT(f_id) FROM task_instr_stocks WHERE f_date = ? FOR UPDATE`, date).Scan(&count)
	err := db.QueryRow(`SELECT COUNT(f_id) FROM task_instr_stocks WHERE f_date = ? AND f_daily_instr_no = ? AND f_index_daily_modify = ?`, date, dailyInstrNo, indexDailyModify).Scan(&count)
	return count, err
}

func UpdateTaskInstrStockStatus(db db.SimpleDB, date int, dailyInstrNo, indexDailyModify, stockSerialNo int64,
	stockEntrustExecuteStatus, stockDealExecuteStatus string,
	totalDealAmount, totalDealBalance, cumAvgPrice, totalEntrustAmount, totalEntrustBalance float64,
	dealCompleteDateTime int64) error {
	args := []interface{}{
		stockEntrustExecuteStatus, stockDealExecuteStatus, 
		totalDealAmount, totalDealBalance, cumAvgPrice, totalEntrustAmount, totalEntrustBalance,
		dealCompleteDateTime,
		date, dailyInstrNo, indexDailyModify, stockSerialNo}
	_, err := db.Exec(`UPDATE task_instr_stocks SET 
		f_stock_entrust_execute_status=?, f_stock_deal_execute_status=?, 
		f_total_deal_amount=?, f_total_deal_balance=?, f_cum_avg_price=?, f_total_entrust_amount=?, f_total_entrust_balance=?, 
		f_deal_complete_date_time=? WHERE f_date=? AND f_daily_instr_no=? AND f_index_daily_modify=? AND f_stock_serial_no=?`, args...)
	return err
}

func UpdateOperatorOfTaskInstrStock(db db.SimpleDB, operator string, date int, dailyInstrNo, indexDailyModify, stockSerialNo int64) error {
	args := []interface{}{operator, date, dailyInstrNo, indexDailyModify, stockSerialNo}
	_, err := db.Exec("UPDATE task_instr_stocks SET f_stock_instr_operator=? WHERE f_date=? AND f_daily_instr_no=? AND f_index_daily_modify=? AND f_stock_serial_no=?", args...)
	return err
}

func UpdateOperatorOfTaskInstr(db db.SimpleDB, operator string, date int, dailyInstrNo, indexDailyModify int64) error {
	args := []interface{}{operator, date, dailyInstrNo, indexDailyModify}
	_, err := db.Exec("UPDATE task_instrs SET f_operator=? WHERE f_date=? AND f_daily_instr_no=? AND f_index_daily_modify=?", args...)
	return err
}

package store

import (
	"errors"
	"log"
	"reflect"
	"rhino-common/utils/dbutil"
	"rhino-instr/schema"

	"github.com/linchunquan/sqlgen/db"
)



var (
	FindTaskInstrsByDateAndDirectOperatorStmt = SelectTaskInstrStmt + ` where f_date>=? and f_date<? and f_direct_operator=? `
	FindTaskInstrsByDateAndOperatorStmt = SelectTaskInstrStmt + ` where f_date>=? and f_date<? and f_operator=? `
	FindTaskInstrsByDateStmt = SelectTaskInstrStmt + ` where f_date>=? and f_date<?`
)

func FindTaskInstrsByDateAndDirectOperator(db db.SimpleDB, beginDate, endDate int, directOperator string, limit, offset int) (result []*schema.TaskInstr, err error) {
	args := []interface{}{beginDate, endDate, directOperator}
	sqlStmt := FindTaskInstrsByDateAndDirectOperatorStmt
	return findTaskInstrs(db, sqlStmt, args, limit, offset)
}

func FindTaskInstrsByDateAndOperator(db db.SimpleDB, beginDate, endDate int, operator string, limit, offset int) (result []*schema.TaskInstr, err error) {
	args := []interface{}{beginDate, endDate, operator}
	sqlStmt := FindTaskInstrsByDateAndOperatorStmt
	return findTaskInstrs(db, sqlStmt, args, limit, offset)
}

func FindTaskInstrsByDateRange(db db.SimpleDB, beginDate, endDate int, limit, offset int) (result []*schema.TaskInstr, err error) {
	args := []interface{}{beginDate, endDate}
	sqlStmt := FindTaskInstrsByDateStmt
	return findTaskInstrs(db, sqlStmt, args, limit, offset)
}

func findTaskInstrs(db db.SimpleDB, sqlStmt string, args[]interface{}, limit, offset int) (result []*schema.TaskInstr, err error) {
	if isPagging(limit, offset) {
		sqlStmt += PaggingSuffix
		args = append(args, limit, offset)
	}
	result, err = genericSelectTaskInstrs(db, sqlStmt, args...)
	return
}

func FindTaskInstrs(db db.SimpleDB, fieldConditions []*dbutil.FieldCondition, limit, offset int) (result []*schema.TaskInstrView, err error) {
	whereClause, args, err := dbutil.PrepareQueryStatment(reflect.TypeOf((*schema.TaskInstrView)(nil)).Elem(), fieldConditions)
	if err != nil {
		return nil, err
	}
	sqlStmt := SelectTaskInstrViewStmt + " WHERE " + whereClause + " order by f_date desc, f_daily_instr_no desc, f_index_daily_modify desc, f_stock_serial_no desc"

	if isPagging(limit, offset) {
		sqlStmt += PaggingSuffix
		args = append(args, limit, offset)
	}

	log.Printf("sqlStmt:%s\n", sqlStmt)
	log.Printf("args:%+v\n", args)

	result, err = genericSelectTaskInstrViews(db, sqlStmt, args...)
	
	return
}


func FindTaskInstrsCount(db db.SimpleDB, fieldConditions []*dbutil.FieldCondition) (count int, err error) {
	whereClause, args, err := dbutil.PrepareQueryStatment(reflect.TypeOf((*schema.TaskInstrView)(nil)).Elem(), fieldConditions)
	if err != nil {
		return count, err
	}

	selectCountStmt := getCountStmtFromOriginalStmt(SelectTaskInstrViewStmt)
	if selectCountStmt == "" {
		return 0, errors.New("cannnot create SELECT COUNT statment from " + SelectTaskInstrViewStmt)
	}

	sqlStmt := selectCountStmt + " WHERE " + whereClause

	log.Printf("sqlStmt:%s\n", sqlStmt)
	log.Printf("args:%+v\n", args)

	row := db.QueryRow(sqlStmt, args...)
	err = row.Scan(&count)
	return count, err
}

var (
	LockTaskInstrByDateAndDailyInstrNoAndIndexDailyModifyStmt = SelectTaskInstrByDateAndDailyInstrNoAndIndexDailyModifyStmt + ` for update`
)

func GetAndLockTaskInstrByDateAndDailyInstrNoAndIndexDailyModify(db db.SimpleDB, date int, dailyInstrNo int64, indexDailyModify int64) (*schema.TaskInstr, error) {
	args := []interface{}{date, dailyInstrNo, indexDailyModify}
	v, err := genericSelectTaskInstr(db, LockTaskInstrByDateAndDailyInstrNoAndIndexDailyModifyStmt, args...)
	return v, err
}
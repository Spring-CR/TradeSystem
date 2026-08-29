package domain_task_instr

import (
	"rhino-common/context"
	"rhino-common/domain_error"
	"rhino-common/utils/bean"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/timeutil"
	"rhino-instr/domain/status"
	"rhino-instr/schema"
	"rhino-instr/store"
	"strings"
	"time"
)

func CreateInsertTaskInstr(
	accountNo, combiNo, instrType string,
	beginDate, endDate, beginTime, endTime int,
	directOperator string,
	businessType string,
	limitOperator string,
	stocks []*schema.TaskStock,
	datetimes ...int64,
) (dailyInstrNo, indexDailyModify, batchSerialNo int64, dateNum int, de *domain_error.Error) {

	dailyInstrNo, indexDailyModify, batchSerialNo, dateNum, de = _createInsertTaskInstr(accountNo, combiNo, instrType, beginDate, endDate, beginTime, endTime, directOperator, businessType, limitOperator, stocks, datetimes...)
	
	for de != nil {
		if strings.Contains(de.Err.Error(), "Duplicate entry") {
			time.Sleep(time.Millisecond * 100)
			dailyInstrNo, indexDailyModify, batchSerialNo, dateNum, de = _createInsertTaskInstr(accountNo, combiNo, instrType, beginDate, endDate, beginTime, endTime, directOperator, businessType, limitOperator, stocks, datetimes...)
		} else {
			return dailyInstrNo, indexDailyModify, batchSerialNo, dateNum, de
		}
	}

	return dailyInstrNo, indexDailyModify, batchSerialNo, dateNum, nil
}

func _createInsertTaskInstr(
	accountNo, combiNo, instrType string,
	beginDate, endDate, beginTime, endTime int,
	directOperator string,
	businessType string,
	limitOperator string,
	stocks []*schema.TaskStock,
	datetimes ...int64,
) (dailyInstrNo, indexDailyModify, batchSerialNo int64, dateNum int, de *domain_error.Error) {

	// 开始事务
	tx, de := dbutil.BeginTx(context.DB)
	if de != nil {
		return dailyInstrNo, indexDailyModify, batchSerialNo, dateNum, de
	}

	var createTime int64
	var _createTime time.Time
	if len(datetimes) > 0 {
		createTime = datetimes[0]
		_createTime = timeutil.ConvertMicrosecondsToTime(createTime)
	} else {
		_createTime = time.Now()
		createTime = timeutil.ConvertTimeToMicroseconds(_createTime)
	}

	dateNum = timeutil.GetDateNum(_createTime)
	timeNum := timeutil.GetTimeNum(_createTime)

	taskInstrCount, err := store.LockCountForCreateInsertTaskInstr(tx, dateNum)
	if err != nil {
		dbutil.RollbackTx(tx)
		return dailyInstrNo, indexDailyModify, batchSerialNo, dateNum, domain_error.Build(domain_error.CANNOT_LOCK_INSTR_MAIN_RECORD_ERR_CODE, err)
	}

	taskInstr := &schema.TaskInstr{

		DailyInstrNo:     int64(taskInstrCount) + 1,
		Date:             dateNum,
		IndexDailyModify: 1,
		IndexLastModify:  int64(taskInstrCount),

		// 前端录入
		AccountNo:     accountNo,
		CombiNo:       combiNo,
		InstrType:     instrType,
		BeginDate:     beginDate,
		EndDate:       endDate,
		BeginTime:     beginTime,
		EndTime:       endTime,
		BusinessType:  businessType,
		LimitOperator: limitOperator,
		DirectOperator: directOperator,

		DirectDate:     dateNum,
		DirectTime:     timeNum,

		// 前期不做分发逻辑，直接已分发
		DispenseDate:     dateNum,
		DispenseTime:     timeNum,
		DispenseOperator: "admin",
		DispenseStatus:   status.DispenseStatusProcessed,

		InstrStatus:          status.InstrStatusEffective,
		EntrustExecuteStatus: schema.TaskInstrEntrustExecuteStatusNotBegin,
		DealExecuteStatus:    schema.TaskInstrDealExecuteStatusNotBegin,

		CreateTime: createTime,
	}

	err = store.InsertTaskInstr(tx, taskInstr)
	if err != nil {
		dbutil.RollbackTx(tx)
		return dailyInstrNo, indexDailyModify, batchSerialNo, dateNum, domain_error.Build(domain_error.CANNOT_INSERT_MAIN_INSTR_RECORD_ERR_CODE, err)
	}

	taskInstrStockCount, err := store.LockCountOfStocksForInsertTaskInstr(tx, dateNum, taskInstr.DailyInstrNo, taskInstr.IndexDailyModify)
	if err != nil {
		dbutil.RollbackTx(tx)
		return dailyInstrNo, indexDailyModify, batchSerialNo, dateNum, domain_error.Build(domain_error.CANNOT_LOCK_INSTR_SECONDLY_RECORD_ERR_CODE, err)
	}

	for i, stock := range stocks {
		taskInstrStock := &schema.TaskInstrStock{
			DailyInstrNo:     taskInstr.DailyInstrNo,
			IndexDailyModify: taskInstr.IndexDailyModify,
			Date:             taskInstr.Date,
			StockSerialNo:    int64(taskInstrStockCount + i + 1),

			StockEntrustExecuteStatus: schema.TaskInstrEntrustExecuteStatusNotBegin,
			StockDealExecuteStatus:    schema.TaskInstrDealExecuteStatusNotBegin,
		}
		err = bean.Copy(stock).To(taskInstrStock)
		if err != nil {
			dbutil.RollbackTx(tx)
			return dailyInstrNo, indexDailyModify, batchSerialNo, dateNum, domain_error.Build(domain_error.INSTR_ERR_CODE, err)
		}

		err = store.InsertTaskInstrStock(tx, taskInstrStock)
		if err != nil {
			dbutil.RollbackTx(tx)
			return dailyInstrNo, indexDailyModify, batchSerialNo, dateNum, domain_error.Build(domain_error.CANNOT_INSERT_SECONDLY_INSTR_RECORD_ERR_CODE, err)
		}
	}

	// 提交事务
	de = dbutil.CommitTx(tx)
	if de != nil {
		return dailyInstrNo, indexDailyModify, batchSerialNo, dateNum, de
	}

	return taskInstr.DailyInstrNo, taskInstr.IndexDailyModify, taskInstr.BatchSerialNo, dateNum, de
}

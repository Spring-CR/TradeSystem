package domain_trade_instr

import (
	"encoding/json"
	"fmt"
	"rhino-common/context"
	"rhino-common/domain_error"
	"rhino-common/utils/dbutil"
	"rhino-instr/domain/status"
	trade_channel "rhino-instr/domain/trade-channel"
	"rhino-instr/schema"
	"rhino-instr/store"
	"time"
)

// 1、user 要作为opertator反写到指令
// 2、确保不能重复执行
// 3、有个定时任务，定期去查委托和接收委托消息，并及时更新数据库
// 其他补充：
// 1. 对instrstock锁行（锁task instr效果更好）
// 2. 分析tradeinstr获取各种统计数据并更新，可以封装在独立函数，传入tx
// 3. 执行交易
// 4. 提交tx
// 5. tradeinstr增加定位串，原来的key不再加唯一性约束，定位串加序号作为新key
// 6. 委托状态，未执行：无tradeinstr记录；部分执行: 有tradeinstr记录，但是剩余可委托数量大于0；剩余可委托数量为0
// 7. 成交状态：无成交：成交状态的数量是0；部分成交：成交数量大于0但小于指令数量；全部成交：成交数量等于指令数量。
func ExecuteSingleTradeInstr(date int, dailyInstrNo int64, indexDailyModify int64, stockSerialNo int64, ordType string, price float64, amount float64, user string, apiOperator string, c *trade_channel.KafkaTradeChannel) (*schema.TradeInstr, *domain_error.Error) {
	// Step1：开始数据库事务
	tx, de := dbutil.BeginTx(context.DB)
	if de != nil {
		// 失败直接退出，不影响交易安全
		return nil, de
	}

	// Step2：基于select for update机制，先锁住目标行的状态，防止其他异步操作，进而引起超买超卖的风险（对视图锁行是无效的）
	instr, err := store.GetAndLockTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNo(tx, date, dailyInstrNo, indexDailyModify, stockSerialNo)
	if instr == nil || err != nil {
		dbutil.RollbackTx(tx)
		// 失败直接退出，不影响交易安全
		return nil, domain_error.Build(domain_error.CANNOT_GET_INSTR_STOCK_ERR_CODE, err, date, dailyInstrNo, indexDailyModify, stockSerialNo)
	}

	// Step3：状态判断，阻止重复执行
	// if instr.StockDealExecuteStatus != status.StockInstrStatusNotExecuted {
	// 	dbutil.RollbackTx(tx)
	// 	// 基于状态判断阻止重复执行不影响交易安全
	// 	return domain_error.Build(domain_error.ONLY_NONE_EXE_STOCK_INSTR_CAN_BE_EXE, nil, instr.ReportCode)
	// }

	// Step3：统计最新数据，判断是否有可委托数量剩余
	totalTradeInstrs, _, _, _, _, _, totalEntrustAmount, _, _, err := store.StatisTaskInstrStock(tx, date, dailyInstrNo, indexDailyModify, stockSerialNo, false)
	if err != nil {
		dbutil.RollbackTx(tx)
		// 失败直接退出，不影响交易安全
		return nil, domain_error.Build(domain_error.CANNOT_STATIS_INSTR_STOCK_ERR_CODE, err, date, dailyInstrNo, indexDailyModify, stockSerialNo)
	}
	if totalEntrustAmount >= instr.Amount {
		dbutil.RollbackTx(tx)
		// 失败直接退出，不影响交易安全
		return nil, domain_error.Build(domain_error.CANNOT_OVER_ENTRUST_ERR_CODE, err)
	}

	// Step4：开始构建tradeInstr对象
	strategyParameters := GetSingleTradeStrategyParameters()
	strategyParametersText, _ := json.Marshal(strategyParameters)
	userInfo := GetTradeInstrUserInfo(user)

	// 为了获得组合编号，先取出任务指令
	taskInstr, err := store.GetTaskInstrByDateAndDailyInstrNoAndIndexDailyModify(tx, date, dailyInstrNo, indexDailyModify)
	if taskInstr == nil || err != nil {
		dbutil.RollbackTx(tx)
		// 失败直接退出，不影响交易安全
		return nil, domain_error.Build(domain_error.DATABASE_OPERATION_ERR_CODE, err)
	}

	// 开始构建trade-instr对象
	secondaryClOrdID := fmt.Sprintf(status.SecondaryClOrdIDPattern, date, dailyInstrNo, indexDailyModify, stockSerialNo, totalTradeInstrs)
	// 重置价格及数量
	//if price <= 0 || price > instr.Price { // price > instr.Price使用指令价格是在买入的时候才有效，对于卖出，应该是相反的
	if price <= 0 {
		price = instr.Price
	}
	if amount <= 0 || amount > instr.Amount-totalEntrustAmount {
		amount = instr.Amount - totalEntrustAmount
	}

	tradeInstr := &schema.TradeInstr{
		ParentKey:              fmt.Sprintf(status.TradeInstrParentKey, date, dailyInstrNo, indexDailyModify, stockSerialNo),
		MsgType:                "D",
		ClientID:               taskInstr.CombiNo,
		SecondaryClOrdID:       secondaryClOrdID,
		SecurityID:             instr.ReportCode,
		Symbol:                 instr.ReportCode,
		Side:                   instr.EntrustDirection,
		TransactTime:           time.Now().Format(status.TransactTimeLayout),
		OrderQty:               amount,
		OrdType:                ordType,
		Price:                  price,
		TimeInForce:            "0",
		TargetStrategy:         "0",
		StrategyParametersText: string(strategyParametersText),
		StrategyParameters:     strategyParameters,
		MarketCode:             instr.MarketNo,
		UserText:               user,
		User:                   userInfo,
		OpenClose:              instr.OpenClose,
		ApiOperator:            apiOperator,
	}

	// Step 5：保存tradeInstr到数据库, 唯一性约束同样可以阻止重复的执行操作
	err = store.InsertTradeInstr(tx, tradeInstr)
	if err != nil {
		dbutil.RollbackTx(tx)
		return nil, domain_error.Build(domain_error.DATABASE_OPERATION_ERR_CODE, err)
	}

	// 即时更新一次状态
	_, _, _, _, _, _, _, _, _, err = store.StatisTaskInstrStock(tx, date, dailyInstrNo, indexDailyModify, stockSerialNo, true)
	if err != nil {
		dbutil.RollbackTx(tx)
		return nil, domain_error.Build(domain_error.CANNOT_STATIS_INSTR_STOCK_ERR_CODE, err, date, dailyInstrNo, indexDailyModify, stockSerialNo)
	}

	// Step 6：先提交transaction更新状态，这样更安全。如果后续kafka失败了，只要把instr-stock设置为失败状态即可
	de = dbutil.CommitTx(tx)
	if de != nil {
		dbutil.RollbackTx(tx)
		// 失败直接退出，不影响交易安全
		return nil, de
	}

	// Step 7：提交到kafka
	_, de = c.PublishTradeInstr(tradeInstr)
	if de != nil {
		// 提交失败，尝试还原数据库
		// 错误容忍的，即时这步失败了，再重新执行时，还是会先把记录删除
		store.DeleteTradeInstrBySecondaryClOrdId(context.DB, secondaryClOrdID)
		// 先删除，再更新，才是准确的。错误容忍的，即时这步失败了，系统的状态跟踪程序还是会将其设置为超时状态。
		store.StatisTaskInstrStock(context.DB, date, dailyInstrNo, indexDailyModify, stockSerialNo, true)
		return nil, de
	}

	// Step 8：更新操作员信息
	err = store.UpdateOperatorOfTaskInstrStock(context.DB, user, date, dailyInstrNo, indexDailyModify, stockSerialNo)
	if err != nil {
		return nil, domain_error.Build(domain_error.DATABASE_OPERATION_ERR_CODE, err)
	}
	// 反写user作为operator到task-instr
	err = store.UpdateOperatorOfTaskInstr(context.DB, user, date, dailyInstrNo, indexDailyModify)
	if err != nil {
		return nil, domain_error.Build(domain_error.DATABASE_OPERATION_ERR_CODE, err)
	}

	// Todo
	// 后续，需要有独立的状态跟踪部件，不断往柜台查询状态来进行更新
	// 更新时，使用key+offset做去去重阻挡
	// TaskInstr的对象状态更新

	return tradeInstr, nil
}

func GetSingleTradeStrategyParameters() map[string]interface{} {
	parameters := make(map[string]interface{})
	dateStr := time.Now().Format("20060102")
	parameters["TimeStart"] = dateStr + "-00:00:00"
	parameters["TimeEnd"] = dateStr + "-23:59:59"
	return parameters
}

func GetTradeInstrUserInfo(user string) map[string]interface{} {
	parameters := make(map[string]interface{})
	parameters["name"] = user
	parameters["userID"] = user
	parameters["source"] = "rhino"
	return parameters
}

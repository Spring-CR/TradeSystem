package store

import (
	"encoding/json"
	"fmt"
	"log"
	"rhino-instr/schema"

	"github.com/linchunquan/sqlgen/db"
)

var (
	ignoreTaskInstrStatus = map[string]bool{
		//schema.TradeInstrOrdStatusCanceled: true,
		//schema.TradeInstrOrdStatusExpired:  true,
		schema.TradeInstrOrdStatusRejected: true,
		schema.TradeInstrOrdStatusTimeout:  true,
	}
)

func StatisTaskInstrStock(db db.SimpleDB, date int, dailyInstrNo int64, indexDailyModify int64, stockSerialNo int64, update bool) (
	totalTradeInstrs int,
	stockEntrustExecuteStatus,
	stockDealExecuteStatus string,
	totalDealAmount, totalDealBalance, cumAvgPrice, totalEntrustAmount, totalEntrustBalance float64,
	dealCompleteDateTime int64, err error) {

	// 获取inst-stock记录
	var taskInstrStock *schema.TaskInstrStock
	taskInstrStock, err = GetTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNo(db, date, dailyInstrNo, indexDailyModify, stockSerialNo)
	if err != nil {
		return
	}

	// 价格乘数（确定前端是否这样传的）
	priceRatio := taskInstrStock.ContractSize
	if priceRatio == 0 { // 取不了就算
		priceRatio = taskInstrStock.Balance / taskInstrStock.Amount / taskInstrStock.Price
	}

	expectedAmount := taskInstrStock.Amount

	parentKey := fmt.Sprintf("%v-%v-%v-%v", date, dailyInstrNo, indexDailyModify, stockSerialNo)
	log.Printf("parentKey:%s\n", parentKey)
	var tradeInstrs []*schema.TradeInstr
	tradeInstrs, err = FindTradeInstrsByParentKey(db, parentKey)
	if err != nil {
		return
	}

	totalTradeInstrs = len(tradeInstrs)

	log.Printf("totalTradeInstrs:%d\n", totalTradeInstrs)

	var entrustList []*schema.TradeInstr
	var filledList []*schema.TradeInstr

	for _, tradeInstr := range tradeInstrs {
		status := tradeInstr.OrdStatus
		if ignoreTaskInstrStatus[status] {
			continue
		}
		
		entrustList = append(entrustList, tradeInstr)

		if status == schema.TradeInstrOrdStatusFilled || status == schema.TradeInstrOrdStatusCanceled && tradeInstr.CumQty>0 || status == schema.TradeInstrOrdStatusExpired && tradeInstr.CumQty>0  {
			filledList = append(filledList, tradeInstr)
		}
	}

	jsData, _ := json.MarshalIndent(filledList, "", "  ")
	log.Printf("filledList:%s\n", jsData)


	// 开始根据委托中和已成交的交易记录来更新统计信息
	// 计算累计成交量和累计成交金额
	for _, filledRecord := range filledList {
		totalDealAmount += filledRecord.CumQty
		totalDealBalance += filledRecord.CumAmt
	}
	// 计算累计成交均价（成交量占比加权法）
	for _, filledRecord := range filledList {
		cumAvgPrice += filledRecord.AvgPx * filledRecord.CumQty / totalDealAmount
	}
	// 计算累计委托数量和委托金额
	for _, entrustRecord := range entrustList {
		if entrustRecord.OrdStatus == schema.TradeInstrOrdStatusCanceled || entrustRecord.OrdStatus == schema.TradeInstrOrdStatusExpired {
			totalEntrustAmount += entrustRecord.CumQty
			totalEntrustBalance += entrustRecord.Price * entrustRecord.CumQty * priceRatio // 这里不是真实的成交价格
		} else {
			totalEntrustAmount += entrustRecord.OrderQty
			totalEntrustBalance += entrustRecord.Price * entrustRecord.OrderQty * priceRatio // 这里不是真实的成交价格
		}
	}
	// 设置成交状态
	if totalDealAmount == 0 {
		stockDealExecuteStatus = schema.TaskInstrDealExecuteStatusNotBegin
	} else if totalDealAmount < expectedAmount {
		stockDealExecuteStatus = schema.TaskInstrDealExecuteStatusRunning
	} else {
		stockDealExecuteStatus = schema.TaskInstrDealExecuteStatusFinish
	}
	// 设置委托状态
	if totalEntrustAmount == 0 {
		stockEntrustExecuteStatus = schema.TaskInstrEntrustExecuteStatusNotBegin
	} else if totalEntrustAmount < expectedAmount {
		stockEntrustExecuteStatus = schema.TaskInstrEntrustExecuteStatusRunning
	} else {
		stockEntrustExecuteStatus = schema.TaskInstrEntrustExecuteStatusFinish
	}
	// 设置成交完成时间（只有在全部完成成交时才设置）
	if totalDealAmount >= expectedAmount {
		for _, filledRecord := range filledList {
			if filledRecord.MessageTime > dealCompleteDateTime {
				dealCompleteDateTime = filledRecord.MessageTime
			}
		}
	}

	// 更新
	if update {
		err = UpdateTaskInstrStockStatus(db, date, dailyInstrNo, indexDailyModify, stockSerialNo,
			stockEntrustExecuteStatus, stockDealExecuteStatus,
			totalDealAmount, totalDealBalance, cumAvgPrice, totalEntrustAmount, totalEntrustBalance,
			dealCompleteDateTime)

		//Todo: TaskInstr的对象状态更新
	}

	return
}

package status

import (
	"bytes"
	"encoding/json"
	"farm/util/bean"
	"log"
	"rhino-common/context"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/timeutil"
	trade_channel "rhino-instr/domain/trade-channel"
	"rhino-instr/schema"
	"rhino-instr/store"
	"runtime"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
)

/*
1、对task_instr_stocks，执行select for updtae锁行
2、trade_instr_resps，key+offset唯一性，可以确保去重
3、确保kafka 应答topic使用的是单分区
4、开始时从TradeInstrResp的最大offset来获取偏移量
*/
func StatusObserve(c *trade_channel.KafkaTradeChannel) {
	c.KeepListening(func(m *sarama.ConsumerMessage) bool {

		log.Printf("start to process of message %d\n", m.Offset)

		data := m.Value

		// Step 1: json格式化
		log.Println("Step 1")
		tradeInstrResp := &schema.TradeInstrResp{}
		err := json.Unmarshal(data, tradeInstrResp)
		if logError(data, err) {
			return false
		}

		// 如果msgType!=8, 不必考虑
		if tradeInstrResp.MsgType != "8" {
			return true
		}

		tradeInstrResp.StatusKafkaOffset = m.Offset
		tradeInstrResp.MessageTime = timeutil.ConvertTimeToMicroseconds(m.Timestamp)

		log.Println("Step 2")
		// Step 2: 从id字符串解析元数据
		date, dailyInstrNo, indexDailyModify, stockSerialNo, _, ok, _ := parseSecondaryClOrdIDPattern(data, tradeInstrResp.SecondaryClOrdID)
		if !ok {
			log.Printf("cannot parse second id:%s\n", tradeInstrResp.SecondaryClOrdID)
			return false
		}

		log.Println("Step 3")
		// Step 3: 开启trasaction
		tx, de := dbutil.BeginTx(context.DB)
		if de != nil  {
			log.Println("cannot begin tx")
			dbutil.RollbackTx(tx)
			return false
		}

		log.Println("Step 4")
		// Step 4: 对task_instrs，执行select for updtae锁行
		_, err = store.GetAndLockTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNo(tx, date, int64(dailyInstrNo), int64(indexDailyModify), int64(stockSerialNo))
		if logError(data, err) {
			dbutil.RollbackTx(tx)
			return false
		}
		
		log.Println("Step 5")
		// Step 5: 插入trade_instr_resps记录，唯一性检查保护
		err = store.InsertTradeInstrResp(tx, tradeInstrResp)
		// 如果是duplicate error，则应该要忽略的
		if err != nil && strings.Contains(err.Error(), "Duplicate entry") {
			log.Printf("err.Error():%s, tradeInstrResp.SecondaryClOrdID:%s, tradeInstrResp.StatusKafkaOffset:%d\n", err.Error(), tradeInstrResp.SecondaryClOrdID, tradeInstrResp.StatusKafkaOffset)
			dbutil.RollbackTx(tx)
			err = nil
			return true
		}

		log.Printf("======> Success insert tradeInstrResp.SecondaryClOrdID:%s, tradeInstrResp.StatusKafkaOffset:%d\n", tradeInstrResp.SecondaryClOrdID, tradeInstrResp.StatusKafkaOffset)
			

		if logError(data, err) {
			dbutil.RollbackTx(tx)
			return false
		}

		log.Println("Step 6")
		// Step 6: 更新关联的trade_instr记录（注意：不要把id也copyt了！）
		tradeInstr ,err := store.GetTradeInstrBySecondaryClOrdId(tx, tradeInstrResp.SecondaryClOrdID)
		if logError(data, err) || tradeInstr==nil {
			dbutil.RollbackTx(tx)
			return false
		}
		updateId := tradeInstr.ID
		err = bean.Copy(tradeInstrResp).To(tradeInstr)
		if logError(data, err) {
			dbutil.RollbackTx(tx)
			return false
		}
		tradeInstr.ID = updateId
	
		log.Printf("get tradeInstr to update, ID:%d\n, SecondaryClOrdID:%v, OrdStatus:%v, CumAmt:%v, CumQty:%v\n", tradeInstr.ID, tradeInstr.SecondaryClOrdID, tradeInstr.OrdStatus, tradeInstr.CumAmt, tradeInstr.CumQty)
		
		err = store.UpdateTradeInstrById(tx, tradeInstr)
		if logError(data, err) {
			dbutil.RollbackTx(tx)
			return false
		}

		// Step 7: 更新状态
		_, _, _, _, _, _, _, _, _, err = store.StatisTaskInstrStock (tx, date, int64(dailyInstrNo), int64(indexDailyModify), int64(stockSerialNo), true)
		if logError(data, err) {
			dbutil.RollbackTx(tx)
			return false
		}

		// Step 8: 提交事务
		de = dbutil.CommitTx(tx)
		if de != nil  {
			dbutil.RollbackTx(tx)
			return false
		}

		return true
	})
	// 再开一个gorouting，检查超过2分钟没有返回的消息，是否有交易状态，是否超时
}

func parseSecondaryClOrdIDPattern(data []byte, secondaryClOrdID string)(date, dailyInstrNo, indexDailyModify, stockSerialNo, seriNo int, ok bool, err error){

	if !strings.HasPrefix(secondaryClOrdID, SecondaryClOrdIDPrefix) {
		return
	}

	strs := strings.Split(secondaryClOrdID, "-")
	if len(strs) != 6 {
		return
	}

	date, err = strconv.Atoi(strs[1])
	if logError(data, err) {
		return
	}

	dailyInstrNo, err = strconv.Atoi(strs[2])
	if logError(data, err) {
		return
	}

	indexDailyModify, err = strconv.Atoi(strs[3])
	if logError(data, err) {
		return
	}

	stockSerialNo, err = strconv.Atoi(strs[4])
	if logError(data, err) {
		return
	}

	seriNo, err = strconv.Atoi(strs[5])
	if logError(data, err) {
		return
	}

	ok = true

	return
}

func logError(data []byte, err error) (errHappen bool) {
	if err == nil {
		return false
	}
	stackBuf := &bytes.Buffer{}
	for i := 1; i <= 4; i++ {
		_, fn, line, _ := runtime.Caller(i + 1)
		stackBuf.WriteString(fn)
		stackBuf.WriteByte(':')
		stackBuf.WriteString(strconv.Itoa(line))
		stackBuf.WriteByte('\n')
	}
	stack := stackBuf.Bytes()
	log.Printf("error occurs while listening kafka message, data:%s, error:%+v, stack:%s\n", data, err, stack)
	return true
}

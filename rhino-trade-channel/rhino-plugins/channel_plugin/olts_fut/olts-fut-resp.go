package olts_fut

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"rhino-common/context/constant"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/byteutils"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/jsonutil"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	"rhino-core/store/app_store"
	"rhino-core/types"
	"strconv"
	"strings"
	"time"

	json "github.com/bytedance/sonic"
	"github.com/manucorporat/try"
)

var (
	bytesMsgType              = []byte("\"msgType\"")
	bytesMsgExecutionReport   = []byte("\"8\"")
	bytesMsgOrderCancelReject = []byte("\"9\"")
)

func (c *OltsFutChannel) getMessageType(msg []byte) string {
	const key = `"msgType"` // 固定匹配字符串，长度9
	const keyLen = 9
	n := len(msg)
	if n < keyLen+4 { // 最小有效长度: `"msgType":"8"` 共13字节，但粗略过滤
		return ""
	}

	// 遍历查找 `"msgType"`
	for i := 0; i <= n-keyLen; i++ {
		// 逐字节精确匹配（编译器会优化为快速比较）
		if msg[i] == '"' && msg[i+1] == 'm' && msg[i+2] == 's' &&
			msg[i+3] == 'g' && msg[i+4] == 'T' && msg[i+5] == 'y' &&
			msg[i+6] == 'p' && msg[i+7] == 'e' && msg[i+8] == '"' {

			pos := i + keyLen // 跳过 `"msgType"`

			// 跳过冒号和可能存在的空格（处理 `": "` 或 `":"`）
			for pos < n && (msg[pos] == ':' || msg[pos] == ' ') {
				pos++
			}
			if pos >= n || msg[pos] != '"' {
				return "" // 期望值的引号
			}

			// 提取值内容（起始引号之后）
			start := pos + 1
			end := start
			for end < n && msg[end] != '"' {
				end++
			}
			if end >= n || end-start != 1 { // 只允许单个字符的值
				return ""
			}

			// 精确匹配 "8" 或 "9"
			switch msg[start] {
			case '8':
				return "8"
			case '9':
				return "9"
			default:
				return ""
			}
		}
	}
	return ""
}

func (c *OltsFutChannel) processExecutionReport(data []byte, executionReport *ExecutionReport, msgSeq int64, msgTime time.Time) {
	orderID := executionReport.OrdID
	clOrdID := executionReport.ClOrdID
	origClOrdID := executionReport.OrigClOrdID
	if executionReport.OrdStatus == "C" {
		executionReport.OrdStatus = "8"
	}
	execID := executionReport.ExecID
	if execID == "" {
		errMsg := fmt.Sprintf("execID not found for order:%v", clOrdID)
		domain_error.ProcessSevereError(false, 0, nil, errors.New(errMsg), errMsg)
		execID = clOrdID + executionReport.TransactTime + executionReport.OrdStatus
	}
	execRefID := executionReport.ExecRefID
	execTransType := executionReport.ExecTransType
	execType := executionReport.ExecType
	ordStatus := executionReport.OrdStatus
	ordRejReason := executionReport.OrdRejReason

	var text string
	if len(executionReport.Text) > 0 && strings.Contains(executionReport.Text, "\\u") {
		var err error
		log.Printf("start to Unquote:%v\n", executionReport.Text)
		text, err = strconv.Unquote("\"" + executionReport.Text + "\"")
		if err != nil {
			log.Printf("Unquote error:%v\n", err)
			text = executionReport.Text
		}
	} else {
		text = executionReport.Text
	}
	log.Printf("======> response text:%s\n", text)
	execRestatementReason := ""
	account := executionReport.Account
	symbol := executionReport.Symbol
	symbolSfx := ""
	securityID := ""
	idSource := ""
	securityType := executionReport.SecurityType
	side := oltsToFixSideMap[executionReport.Side]
	orderQty := executionReport.OrdQty
	cashOrderQty := 0.0
	ordType := executionReport.OrdType
	price := executionReport.Price
	currency := executionReport.Currency
	effectiveTime := ""
	expireTime := ""

	lastShares := executionReport.LastQty
	lastPx := executionReport.LastPx
	leavesQty := executionReport.LeavesQty
	cumQty := 0.0
	if orderQty > 0 {
		cumQty = orderQty - leavesQty
	}
	avgPx := 0.0
	// log.Printf("_lastShares:%s, _lastPx:%s, leavesQty:%s, _cumQty:%s, _avgPx:%s\n", _lastShares, _lastPx, _leavesQty, _cumQty, _avgPx)

	transactTime := msgTime
	_transactTime := executionReport.TransactTime
	if len(_transactTime) > 0 {
		t, err := timeutil.ParseTimeStrToTimeByTimeLocation(TransactTimeLayout, _transactTime, timeutil.CnTimeLocation)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("illegal time string:%s\n", _transactTime))
		} else {
			transactTime = t
		}
	}

	exchangeTradeDate := transactTime.Format(time.DateOnly)

	if ordStatus == string(enum.OrdStatus_Rejected) && len(text) > 0 {
		if len(ordRejReason) > 0 {
			ordRejReason += "-" + text
		} else {
			ordRejReason = text
		}
	}

	extendAttrMap := make(map[string]interface{})
	c.refineExtendAttrMap(extendAttrMap, executionReport)
	extendAttr, _ := json.Marshal(extendAttrMap)

	tradeActionResp := &schema.TradeActionResp{
		OrderID:               orderID,
		ClOrdID:               clOrdID,
		OrigClOrdID:           origClOrdID,
		ExecID:                execID,
		ExecRefID:             execRefID,
		ExecTransType:         execTransType,
		ExecType:              execType,
		OrdStatus:             ordStatus,
		OrdRejReason:          ordRejReason,
		ExecRestatementReason: execRestatementReason,
		Account:               account,
		Symbol:                symbol,
		SymbolSfx:             symbolSfx,
		SecurityID:            securityID,
		IDSource:              idSource,
		SecurityType:          securityType,
		Side:                  side,
		OrderQty:              orderQty,
		CashOrderQty:          cashOrderQty,
		OrdType:               ordType,
		Price:                 price,
		Currency:              currency,
		EffectiveTime:         effectiveTime,
		ExpireTime:            expireTime,
		LastShares:            int64(lastShares),
		LastPx:                lastPx,
		LeavesQty:             int64(leavesQty),
		CumQty:                int64(cumQty),
		AvgPx:                 avgPx,
		TransactTime:          timeutil.ConvertTimeToMilliseconds(transactTime),
		ExchangeTradeDate:     exchangeTradeDate,
		MsgTime:               timeutil.ConvertTimeToMilliseconds(msgTime),
		MsgSeq:                msgSeq,
		ChannelCode:           c.cfg.GetTradeChannel().ChannelCode,
		RawMsg:                byteutils.GetZeroCopyString(data),
		ExtendAttr:            byteutils.GetZeroCopyString(extendAttr),
		ExtendAttrMap:         extendAttrMap,
	}

	if ordStatus == string(enum.OrdStatus_DoneForDay) {
		tradeActionResp.AppOrdID = execRefID
		log.Printf("tradeActionResp.AppOrdID:%s\n", tradeActionResp.AppOrdID)
	}

	// 对成交回报进行清洗
	ignore := c.CleanTradeActionResp(tradeActionResp)
	if ignore {
		log.Printf("======> ignore tradeActionResp: %s\n", tradeActionResp.RawMsg)
		return
	}

	jsonutil.PrintSimple("======> received tradeActionResp:\n", tradeActionResp)

	if len(tradeActionResp.OrdRejReason) > constant.MaxRejectMsgLen {
		tradeActionResp.OrdRejReason = tradeActionResp.OrdRejReason[:constant.MaxRejectMsgLen]
	}

	c.tradeActionRespBuf <- tradeActionResp
}

func (c *OltsFutChannel) processOrderCancelReject(data []byte, orderCancelReject *OrderCancelReject, msgSeq int64, msgTime time.Time) {
	orderID := orderCancelReject.OrdID
	clOrdID := orderCancelReject.ClOrdID
	origClOrdID := orderCancelReject.OrigClOrdID
	ordStatus := orderCancelReject.OrdStatus
	account := ""
	cxlRejReason := orderCancelReject.CxlRejReason
	cxlRejResponseTo := string(enum.CxlRejResponseTo_Cancel) // 必须，代表对Order Cancel Request 的响应

	transactTime := msgTime
	_transactTime := orderCancelReject.TransactTime
	if len(_transactTime) > 0 {
		t, err := timeutil.ParseTimeStrToTimeByTimeLocation(TransactTimeLayout, _transactTime, timeutil.CnTimeLocation)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("illegal time string:%s\n", _transactTime))
		} else {
			transactTime = t
		}
	}

	tradeActionResp := &schema.TradeActionResp{
		OrderID:          orderID,
		ClOrdID:          clOrdID,
		OrigClOrdID:      origClOrdID,
		OrdStatus:        ordStatus,
		OrdRejReason:     types.OrderCancelRejectPrefix + strings.TrimSpace(cxlRejReason),
		CxlRejResponseTo: cxlRejResponseTo,
		Account:          account,
		TransactTime:     timeutil.ConvertTimeToMilliseconds(transactTime),
		MsgTime:          timeutil.ConvertTimeToMilliseconds(time.Now()),
		MsgSeq:           int64(msgSeq),
		ChannelCode:      c.cfg.GetTradeChannel().ChannelCode,
		RawMsg:           byteutils.GetZeroCopyString(data),
	}

	// 对成交回报进行清洗
	c.CleanTradeActionResp(tradeActionResp)

	if len(tradeActionResp.OrdRejReason) > constant.MaxRejectMsgLen {
		tradeActionResp.OrdRejReason = tradeActionResp.OrdRejReason[:constant.MaxRejectMsgLen]
	}

	c.tradeActionRespBuf <- tradeActionResp
}

func (c *OltsFutChannel) startProcessExecutionReportFromBuf() {
	go func() {
		for {
			tradeActionResp := <-c.tradeActionRespBuf
			try.This(func() {
				// 插入tradeActionResp
				c.cfg.GetAutoTx().Input(c.cfg.GetAppDB(), func(tx *sql.Tx) (de *domain_error.Error) {

					// 设置插入时间
					tradeActionResp.DBInsertTime = timeutil.ConvertTimeToMilliseconds(time.Now())

					if len(tradeActionResp.OrdRejReason) > 512 {
						tradeActionResp.OrdRejReason = tradeActionResp.OrdRejReason[:512]
					}

					err := app_store.InsertTradeActionResp(tx, tradeActionResp)
					if err != nil && !dbutil.IsMysqlDuplicateEntryError(err) { // 排除重复插入的错误

						de = domain_error.Build(domain_error.TRADE_RESP_INSERT_TO_DB_ERR_CODE, err, tradeActionResp.RawMsg)
						return
					}
					return
				}, "", "")

				if c.onTradeActionResp != nil {
					c.onTradeActionResp(tradeActionResp)
				} else {
					domain_error.ProcessSevereError(true, 5, nil, nil, "onTradeActionResp function is not defined!")
				}
			}).Catch(func(err try.E) {
				log.Printf("error occur while processing execution report! error:%v\n", err)
				de := domain_error.Build(domain_error.GENERIC_ERR_CODE, fmt.Errorf("error occur while processing execution report! error:%v", err))
				domain_error.ReportIfErrorHappen(de)
			})
		}
	}()
}

func (c *OltsFutChannel) refineExtendAttrMap(extendAttrMap map[string]interface{}, executionReport *ExecutionReport) {
	// 处理成交回报的扩展属性
}

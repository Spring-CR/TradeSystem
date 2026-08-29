package ficc

import (
	"bytes"
	"database/sql"
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
	"strings"
	"time"

	"github.com/manucorporat/try"
)

var (
	bytesMsgType              = []byte("\"MsgType\"")
	bytesMsgExecutionReport   = []byte("\"8\"")
	bytesMsgOrderCancelReject = []byte("\"9\"")
)

func (c *FiccChannel) getMessageType(msg []byte) string {
	i := bytes.Index(msg, bytesMsgType)
	if i < 0 {
		return ""
	}
	msg = msg[9:]
	i = bytes.IndexByte(msg, ',')
	if i < 0 {
		i = bytes.IndexByte(msg, '}')
		if i < 0 {
			return ""
		}
	}

	msg = msg[:i]

	if bytes.Contains(msg, bytesMsgExecutionReport) {
		return "8"
	}

	if bytes.Contains(msg, bytesMsgOrderCancelReject) {
		return "9"
	}

	return ""
}

func (c *FiccChannel) getMessageType2(msg []byte) string {
    const key = `"MsgType"`   // 固定匹配字符串，长度9
    const keyLen = 9
    n := len(msg)
    if n < keyLen+4 {         // 最小有效长度: `"MsgType":"8"` 共13字节，但粗略过滤
        return ""
    }

    // 遍历查找 `"MsgType"`
    for i := 0; i <= n-keyLen; i++ {
        // 逐字节精确匹配（编译器会优化为快速比较）
        if msg[i] == '"' && msg[i+1] == 'M' && msg[i+2] == 's' &&
            msg[i+3] == 'g' && msg[i+4] == 'T' && msg[i+5] == 'y' &&
            msg[i+6] == 'p' && msg[i+7] == 'e' && msg[i+8] == '"' {

            pos := i + keyLen // 跳过 `"MsgType"`

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

func (c *FiccChannel) processExecutionReport(data []byte, executionReport *ExecutionReport, msgSeq int64, msgTime time.Time) {
	orderID := executionReport.OrdID
	clOrdID := executionReport.ClOrdID
	origClOrdID := executionReport.OrigClOrdID
	execID := executionReport.ExecID
	if execID == "" {
		//execID = clOrdID + executionReport.TransactTime + executionReport.OrdStatus
		//execID = executionReport.MsgID
		if executionReport.OrdStatus == string(enum.OrdStatus_DoneForDay) {
			execID = executionReport.MsgID
		} else {
			execID = clOrdID + executionReport.TransactTime + executionReport.OrdStatus + executionReport.MsgID
		}
	}
	execRefID := executionReport.ExecRefID
	execTransType := executionReport.ExecTransType
	execType := executionReport.ExecType
	ordStatus := executionReport.OrdStatus
	ordRejReason := executionReport.OrdRejReason
	text := executionReport.Text
	execRestatementReason := ""
	account := ""
	symbol := executionReport.Symbol
	symbolSfx := ""
	securityID := ""
	idSource := ""
	securityType := ""
	side := executionReport.Side
	orderQty := executionReport.Quantity
	cashOrderQty := 0.0
	ordType := ""
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
		t, err := timeutil.ParseTimeStrToTimeByTimeLocation(timeutil.TransactTimeLayout, _transactTime, time.UTC)
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

func (c *FiccChannel) processOrderCancelReject(data []byte, orderCancelReject *OrderCancelReject, msgSeq int64, msgTime time.Time) {
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
		t, err := timeutil.ParseTimeStrToTimeByTimeLocation(timeutil.TransactTimeLayout, _transactTime, time.UTC)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("illegal time string:%s\n", _transactTime))
		} else {
			transactTime = t
		}
	}

	tradeActionResp := &schema.TradeActionResp{
		OrderID:          orderID,
		ClOrdID:          clOrdID, // 这里的clOrdID，并不是订单号，而是OrderCancelRequest的请求ID，具有唯一性
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

func (c *FiccChannel) startProcessExecutionReportFromBuf() {
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

func (c *FiccChannel) refineExtendAttrMap(extendAttrMap map[string]interface{}, executionReport *ExecutionReport) {
	execID := executionReport.ExecID
	extendAttrMap["respExecID"] = execID
	handlInst := executionReport.HandlInst
	if execID != "" {
		if handlInst == "" {
			if strings.HasPrefix(execID, "CBT") {
				handlInst = "1"
			} else if strings.HasPrefix(execID, "TPB") {
				handlInst = "3"
			} else {
				handlInst = "2"
			}
		}
	} else {
		handlInst = ""
	}
	switch handlInst {
	case "2":
		if executionReport.Broker != "" {
			handlInst = "3"
		}
	case "3":
		handlInst = "4"
	}
	extendAttrMap["respHandlInst"] = handlInst
}

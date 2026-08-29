package fix

import (
	"bytes"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-trade-channel/adapter/store/fix_store"

	"github.com/manucorporat/try"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
)

func (e *TradeClient) processRejectMessage(msg *quickfix.Message) {

	refSeqNum, err := msg.Body.GetInt(tag.RefSeqNum)
	if refSeqNum <= 0 || err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to get refSeqNum from the reject message")
		return
	}

	if fix_store.FixStore == nil {
		return
	}

	var message []byte
	try.This(func() {
		messages, err := fix_store.FixStore.GetMessages(refSeqNum, refSeqNum)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to get message from fixStore, refSeqNum:%d, error:%v\n", refSeqNum, err))
			return
		}
		if len(messages) > 0 {
			message = messages[0]
		}
	}).Catch(func(err try.E) {
		log.Printf("fail to get message from fixStore, error:%v\n", err)
	})
	if len(message) == 0 {
		return
	}

	buf := bytes.NewBuffer(nil)
	buf.Write(message)

	fixMsg := quickfix.NewMessage()
	err1 := quickfix.ParseMessage(fixMsg, buf)
	if err1 != nil {
		domain_error.ProcessSevereError(false, 0, nil, err1, fmt.Sprintf("fail to parse fix message, error:%v\n", err1))
		return
	}

	log.Printf("fetch message from store: %s\n", message)

	msgType, _ := fixMsg.MsgType()

	if msgType == "" {
		return
	}

	refMsgType, _ := msg.Body.GetString(tag.RefMsgType)
	if refMsgType != "" && refMsgType != msgType {
		return
	}

	// 目前只处理针对D报文的Reject类型消息
	if msgType != string(enum.MsgType_ORDER_SINGLE) {
		return
	}

	clOrdID, _ := fixMsg.Body.GetString(tag.ClOrdID)
	if clOrdID == "" {
		return
	}

	msgSeqNum, _ := msg.Header.GetString(tag.MsgSeqNum)
	fixMsg.Header.SetString(tag.MsgType, string(enum.MsgType_EXECUTION_REPORT))
	fixMsg.Body.SetString(tag.OrderID, clOrdID+"_reject")
	fixMsg.Body.SetString(tag.ExecID, clOrdID+"_reject_"+msgSeqNum)
	fixMsg.Body.SetString(tag.ExecTransType, string(enum.ExecTransType_NEW))
	fixMsg.Body.SetString(tag.ExecType, string(enum.ExecType_REJECTED))
	fixMsg.Body.SetString(tag.OrdStatus, string(enum.OrdStatus_REJECTED))

	rejectReason, _ := msg.Body.GetString(tag.SessionRejectReason)
	if rejectReason != "" {
		fixMsg.Body.SetString(tag.OrdRejReason, rejectReason)
	}

	fixMsg.Body.SetInt(tag.LastShares, 0)
	fixMsg.Body.SetInt(tag.LastPx, 0)
	fixMsg.Body.SetInt(tag.LeavesQty, 0)
	fixMsg.Body.SetInt(tag.CumQty, 0)
	fixMsg.Body.SetInt(tag.AvgPx, 0)

	text, _ := msg.Body.GetString(tag.Text)
	if text != "" {
		fixMsg.Body.SetString(tag.Text, text)
	}

	// 按执行回报的逻辑来执行
	e.processExecutionReport(fixMsg)
}

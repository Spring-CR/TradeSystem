package app_store

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"sort"
	"strconv"

	"github.com/linchunquan/sqlgen/db"
)

var (
	resetFlagTarget    = []byte("141=Y\x01")
	resetFlagEndTarget = []byte("141=Y")
	seqNumTarget       = []byte("\x0134=")
	lenSeqNumTarget    = len(seqNumTarget)
	soh                = byte(1) // FIX分隔符
)

// 高性能版本：直接操作字节数组，避免内存分配
func hasResetSeqNumFlagY(fixMsg []byte) bool {
	// 直接使用字节搜索，无需分割字段
	if bytes.Contains(fixMsg, resetFlagTarget) {
		return true
	}

	// 回退检查字段可能出现在消息末尾的情况（无结尾SOH）
	return bytes.HasSuffix(fixMsg, resetFlagEndTarget)
}

func getFixMsgSeqNum(fixMsg []byte) (seqNum int, err error) {

	idx := bytes.Index(fixMsg, seqNumTarget)
	if idx == -1 {
		return seqNum, fmt.Errorf("the fix message does not contain seqNum, FIX message:%s", fixMsg)
	}

	// 定位值的起始位置（跳过"34="）
	valueStart := idx + lenSeqNumTarget
	if valueStart >= len(fixMsg) {
		return 0, errors.New("invalid tag 34 format")
	}

	// 查找值结束位置（下一个SOH或消息末尾）
	valueEnd := valueStart
	for ; valueEnd < len(fixMsg); valueEnd++ {
		if fixMsg[valueEnd] == soh {
			break
		}
	}

	// 提取值部分
	valueBytes := fixMsg[valueStart:valueEnd]

	// 转换为整数
	seqNum, err = strconv.Atoi(string(valueBytes))
	if err != nil {
		err = fmt.Errorf("invalid sequence number format, FIX message:%s", fixMsg)
		return
	}

	return
}

func ParseMsgSeqFromFixMessage(db db.SimpleDB, channelCode string) (latestRestMsgSeqTime int64, maxSenderMsgSeq, maxTargetMsgSeq int, err error) {

	args := []interface{}{channelCode}
	fixMessages, err1 := genericSelectUtilFixMessages(db, SelectUtilFixMessageStmt+" where f_channel_code = ?", args...)

	if dbutil.IsDbRecordEmptyError(err1) {
		return
	}

	if err1 != nil {
		err = err1
		return
	}

	// 按MsgTime倒序排序
	sort.Slice(fixMessages, func(i, j int) bool {
		t1 := fixMessages[i].MsgTime
		t2 := fixMessages[j].MsgTime
		if t1 != t2 {
			return t2 < t1
		}
		return fixMessages[j].ID < fixMessages[i].ID
	})

	for _, msg := range fixMessages {

		data := msg.Data

		//log.Printf("FIX msg data:%s\n", msg.Data)

		switch msg.MsgType {
		case "A":
			if hasResetSeqNumFlagY(data) {
				latestRestMsgSeqTime = msg.MsgTime
				return
			}
		case "4":
			latestRestMsgSeqTime = msg.MsgTime
			return
		}

		seqNum, err1 := getFixMsgSeqNum(data)
		if err1 != nil {
			err = err1
			log.Printf("error occurs for msg, id=%d\n", msg.ID)
			return
		}

		switch enum.UtilFixMessageSide(msg.MsgSide) {
		case enum.UtilFixMessageSide_FromAdmin:
			if maxTargetMsgSeq < seqNum {
				maxTargetMsgSeq = seqNum
			}
		case enum.UtilFixMessageSide_FromApp:
			if maxTargetMsgSeq < seqNum {
				maxTargetMsgSeq = seqNum
			}
		case enum.UtilFixMessageSide_ToAdmin:
			if maxSenderMsgSeq < seqNum {
				maxSenderMsgSeq = seqNum
			}
		case enum.UtilFixMessageSide_ToApp:
			if maxSenderMsgSeq < seqNum {
				maxSenderMsgSeq = seqNum
			}
		}
	}

	return
}

package domain_cfg

import (
	"bytes"
	"errors"
	"log"
	"rhino-common/utils/kafka"
	"rhino-common/utils/timeutil"
	"strconv"
	"time"
)


func (s *MsgSeqGen) getFixMsgSeqNum(rawMsg[]byte, msgOffset int64) (seqNum int, sendTime time.Time, err error) {

	// 如果当前的时间，不在trade channel的[Begin, End]范围内，seqNum要立即返回0；当seqNum=0时，fix session需要重置
	duringTradingTime, err := s.cfg.IsDuringTradingTime(60 * 60)
	log.Printf("IsDuringTradingTime: %v\n", duringTradingTime)
	if err != nil {
		return 0, sendTime, err
	}

	// 如果不在trade channel的[Begin, End]范围内
	if !duringTradingTime {
		return 0, sendTime, nil
	}

	// 将字节消息转换为字符串，按SOH（0x01）分割
	fields := bytes.Split(rawMsg, []byte{'\001'})

	// tag34存放消息序号
	prefixTagSeqNum := []byte("34=")
	// 遍历字段查找MsgSeqNum
	for _, field := range fields {
		// 检查字段是否包含MsgSeqNum的tag
		if bytes.HasPrefix(field, prefixTagSeqNum) {
			// 提取MsgSeqNum值
			seqNumByte := bytes.TrimPrefix(field, prefixTagSeqNum)
			seqNum, err = strconv.Atoi(string(seqNumByte))
			if err != nil {
				return 0, sendTime, errors.New("FIX message '"+string(rawMsg)+"' without illegal tag 34 value '" + string(seqNumByte) + "'")
			}
			break
		}
	}
	if seqNum == 0 {
		return 0, sendTime, errors.New("illegal FIX message '"+string(rawMsg)+"' without tag 52")
	}

	// tag52存放消息的SendingTime，在Go的实现中，这个字段是被quickfix/go强制注入的
	prefixTagSendingTime := []byte("52=")
	// 确保找到一条包含了tag52的最新消息
	for {
		if bytes.Contains(rawMsg, prefixTagSendingTime) {
			break
		}
		msgOffset = msgOffset - 1
		var exist bool
		rawMsg, _, exist, err = kafka.GetMessage(s.cfg.GetApplicationCfg().GetKafkaBrokers(), s.cfg.GetApplicationCfg().GetTradeChannelReqTopic(), msgOffset)
		if err != nil {
			return 0, sendTime, err
		}
		if !exist {
			return 0, sendTime, nil
		}
	}

	var sendingTime[]byte
	// 遍历字段查找SendingTime
	for _, field := range fields {
		// 检查字段是否包含SendingTime的tag
		if bytes.HasPrefix(field, prefixTagSendingTime) {
			// 提取MsgSeqNum值
			sendingTime = bytes.TrimPrefix(field, prefixTagSendingTime)
			break
		}
	}

	// 考虑交易通道的时区问题
	// 如果通过SendingTime判断交易日期，发现是隔日的时区，需要将seqNum置为0
	// 需要首先以交易通道时区为基准，判断当前的日期
	// 格式如: 20241204-01:38:17.938, 前17位的layout: 20060102-15:04:05
	// 注意：FIX的时间戳是以0时区为基准的
	if len(sendingTime) >= 17 {
		sendingTime = sendingTime[:17]
		sendTime, err = time.ParseInLocation("20060102-15:04:05", string(sendingTime), time.UTC) // fix的时间是0时区的
		if err != nil {
			return 0, sendTime, err
		}
		// 如果不在同一天，则直接重制
		sendDate := sendTime.In(timeutil.GetTimeZone(s.cfg.GetTimeZone())).Format(time.DateOnly)
		currDate := time.Now().In(timeutil.GetTimeZone(s.cfg.GetTimeZone())).Format(time.DateOnly)
		if sendDate != currDate {
			seqNum = 0
		}
	}

	return
}

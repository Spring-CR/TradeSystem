package stream

import (
	"time"
)

type IngressMessage struct {
	Data           []byte
	MsgTime        time.Time
	MsgSeq         int64
	WorkerAffinity int
}

type StreamClient interface {
	PrepareIngressMessageChannel(buffer int) (msgChan <-chan *IngressMessage)
	// 对于kafka，改用TenaciousProducerWithBuffer工具类，msgSeq永远返回0，err永远返回nil
	SendMessage(data []byte) (msgSeq int64, err error)
	GetHistoricalSentKeysAndReqMsgSeqs() (keys []string, reqMsgSeqs []int64)
}

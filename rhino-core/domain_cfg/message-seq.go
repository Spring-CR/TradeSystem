package domain_cfg

import (
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/kafka"
	"rhino-common/utils/timeutil"
	"rhino-core/store/app_store"
	"sync/atomic"
	"time"
)

// 消息序号生成器
type MsgSeqGen struct {
	cfg                    *TradeChannelCfg
	senderMsgSeqNum        int64 // 来自TradeOrder和TradeActionLatestResp的最大值
	targetMsgSeqNum        int64 // 来自TradeActionResp
	reachMaxRetryLogonFail int   // 0-初始值，1-logon之后设置，2-连续多次登录失败时设置
}

func NewMsgSeqGen(cfg *TradeChannelCfg) (msgSeqGen *MsgSeqGen, de *domain_error.Error) {
	msgSeqGen = &MsgSeqGen{
		cfg: cfg,
	}
	de = msgSeqGen.setMsgSeqNum()
	return
}

func (s *MsgSeqGen) setMsgSeqNum() (de *domain_error.Error) {

	latestRestMsgSeqTime, dbMaxSenderMsgSeq, dbMaxTargetMsgSeq, err := app_store.ParseMsgSeqFromFixMessage(s.cfg.GetAppDB(), s.cfg.GetChannelCode())
	if err != nil {
		de = domain_error.Build(domain_error.CANNOT_INIT_MSG_SEQ_GEN_ERR_CODE, err)
		return
	}
	log.Printf("======>ParseMsgSeqFromFixMessage, latestRestMsgSeqTime:%v, dbMaxSenderMsgSeq:%v, dbMaxTargetMsgSeq:%v\n", latestRestMsgSeqTime, dbMaxSenderMsgSeq, dbMaxTargetMsgSeq)

	// get sender msg seq
	var senderMsgSeq int
	var sendTime time.Time
	brokers := s.cfg.GetApplicationCfg().GetKafkaBrokers()
	if len(brokers) > 0 {
		newestMsg, offset, exist, err := kafka.GetNewestMessage(brokers, s.cfg.GetApplicationCfg().GetTradeChannelReqTopic())
		log.Printf("======>newestMsg:%s\n", newestMsg)
		if err != nil {
			de = domain_error.Build(domain_error.CANNOT_INIT_MSG_SEQ_GEN_ERR_CODE, err)
			return
		}
		if exist {
			// 判断协议类型再进行解析
			switch enum.ChannelProtocolType(s.cfg.tradeChannel.ChannelProtocolType) {
			case enum.ChannelProtocolType_FIX42:
				senderMsgSeq, sendTime, err = s.getFixMsgSeqNum(newestMsg, offset)
				if err != nil {
					de = domain_error.Build(domain_error.CANNOT_INIT_MSG_SEQ_GEN_ERR_CODE, err)
					return
				}
			case enum.ChannelProtocolType_FIX44:
				senderMsgSeq, sendTime, err = s.getFixMsgSeqNum(newestMsg, offset)
				if err != nil {
					de = domain_error.Build(domain_error.CANNOT_INIT_MSG_SEQ_GEN_ERR_CODE, err)
					return
				}
			default:
				de = domain_error.Build(domain_error.UNKNOW_CHANNEL_PROTOCOL_TYPE_ERR_CODE, nil)
				return
			}
		}
	}

	// get target msg seq
	begin := time.Now()
	systemCode, businessCode := s.cfg.GetApplicationCfg().GetSystemAndBusinessCodes()
	targetMsgSeq, err := app_store.GetMaxMsgSeqOfTradeActionResp(s.cfg.GetAppDB(), systemCode, businessCode, s.cfg.GetChannelCode(), latestRestMsgSeqTime)
	if err != nil {
		de = domain_error.Build(domain_error.CANNOT_INIT_MSG_SEQ_GEN_ERR_CODE, err)
		return
	}

	if senderMsgSeq > 0 { // 如果senderMsgSeq==0，表示不在可交易的时间，fix是需要强制重置的
		_sendTime := timeutil.ConvertTimeToMilliseconds(sendTime)
		if _sendTime <= latestRestMsgSeqTime {
			senderMsgSeq = dbMaxSenderMsgSeq
		}
		if senderMsgSeq < dbMaxSenderMsgSeq {
			log.Println("senderMsgSeq < dbMaxSenderMsgSeq")
			senderMsgSeq = dbMaxSenderMsgSeq
		}
	}

	// 为了确保不丢失成交回报，应该仅以TradeActionResp所记录的MsgSeq为基准，FixMessage表只是提供了latestRestMsgSeqTime的约束
	// if targetMsgSeq < dbMaxTargetMsgSeq {
	// 	log.Println("targetMsgSeq < dbMaxTargetMsgSeq")
	// 	targetMsgSeq = dbMaxTargetMsgSeq
	// }

	s.senderMsgSeqNum = int64(senderMsgSeq)
	s.targetMsgSeqNum = int64(targetMsgSeq)
	//s.targetMsgSeqNum = int64(targetMsgSeq/targetMsgSeq) // 经调试，Resend Request是有效果的，但是，它只会返回回报之类的应用层报文，心跳包是不返回的！

	log.Printf("======> get senderMsgSeqNum: %d, targetMsgSeqNum: %d, timeCost:%v\n", s.senderMsgSeqNum, s.targetMsgSeqNum, time.Since(begin))

	return
}

func (s *MsgSeqGen) NextSenderMsgSeqNum() int {
	return int(atomic.LoadInt64(&s.senderMsgSeqNum) + 1)
}

func (s *MsgSeqGen) NextTargetMsgSeqNum() int {
	return int(atomic.LoadInt64(&s.targetMsgSeqNum) + 1)
}

func (s *MsgSeqGen) IncrNextSenderMsgSeqNum() error {
	atomic.AddInt64(&s.senderMsgSeqNum, 1)
	return nil
}

func (s *MsgSeqGen) IncrNextTargetMsgSeqNum() error {
	atomic.AddInt64(&s.targetMsgSeqNum, 1)
	return nil
}

func (s *MsgSeqGen) SetNextSenderMsgSeqNum(nextSeqNum int) error {
	//store.senderMsgSeqNum = nextSeqNum - 1
	atomic.StoreInt64(&s.senderMsgSeqNum, int64(nextSeqNum-1))
	return nil
}
func (s *MsgSeqGen) SetNextTargetMsgSeqNum(nextSeqNum int) error {
	//store.targetMsgSeqNum = nextSeqNum - 1
	atomic.StoreInt64(&s.targetMsgSeqNum, int64(nextSeqNum-1))
	return nil
}
func (s *MsgSeqGen) Reset() {
	atomic.StoreInt64(&s.senderMsgSeqNum, 0)
	atomic.StoreInt64(&s.targetMsgSeqNum, 0)
}
func (s *MsgSeqGen) GetMsgSeqNum() (sender int64, target int64) {
	sender = s.senderMsgSeqNum
	target = s.targetMsgSeqNum
	return
}
func (s *MsgSeqGen) SetReachMaxRetryLogonFail(reachMaxRetryLogonFail int) {
	s.reachMaxRetryLogonFail = reachMaxRetryLogonFail
}
func (s *MsgSeqGen) GetReachMaxRetryLogonFail() int {
	return s.reachMaxRetryLogonFail
}
package fix

// 作者：林春泉

import (
	"database/sql"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/fixutil"
	"rhino-common/utils/kafka"
	"rhino-common/utils/timeutil"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/store/app_store"
	"sync/atomic"
	"time"

	"github.com/quickfixgo/quickfix"
)

const (
	max_logout_count_for_reset = 20
)

// TradeClient implements the quickfix.Application interface
type TradeClient struct {
	sessionReady             int32 // 0 - not ready; 1 - ready
	processExecutionReport   func(msg *quickfix.Message)
	processOrderCancelReject func(msg *quickfix.Message)
	kafkaProducer            *kafka.BlocklessProducer
	kafkaTopic               string
	logoutCount              int
	parent                   *GenericFIXChannel
	username                 string
	password                 string
}

func NewTradeClient(cfg *domain_cfg.TradeChannelCfg, parent *GenericFIXChannel) (*TradeClient, error) {
	//kafkaProducer := kafka.NewBlocklessProducer(cfg.)
	client := &TradeClient{parent: parent}
	brokers := cfg.GetApplicationCfg().GetKafkaBrokers()
	if len(brokers) > 0 {
		var err error
		client.kafkaProducer, err = kafka.NewBlocklessProducer(brokers, 32768)
		if err != nil {
			return nil, err
		}
		client.kafkaTopic = cfg.GetApplicationCfg().GetTradeChannelReqTopic()
	}
	client.setCredentials(cfg)
	return client, nil
}

// OnCreate implemented as part of Application interface
func (e *TradeClient) OnCreate(sessionID quickfix.SessionID) {
	log.Printf("session create %s\n", sessionID.String())
}

// OnLogon implemented as part of Application interface
func (e *TradeClient) OnLogon(sessionID quickfix.SessionID) {
	atomic.StoreInt32(&e.sessionReady, 1)
	e.parent.msgSeqGen.SetReachMaxRetryLogonFail(1)
	e.logoutCount = 0
	log.Printf("======> session logon %s, e.sessionReady:%d, logoutCount:%d\n", sessionID.String(), e.sessionReady, e.logoutCount)
}

// OnLogout implemented as part of Application interface
func (e *TradeClient) OnLogout(sessionID quickfix.SessionID) {
	log.Printf("======> session logout %s, logoutCount:%d\n\n", sessionID.String(), e.logoutCount)
	e.logoutCount++
	atomic.StoreInt32(&e.sessionReady, 0)
	if e.logoutCount > max_logout_count_for_reset {

		log.Printf("======> reach max_logout_count_for_reset:%d for reset, logoutCount:%d\n", max_logout_count_for_reset, e.logoutCount)

		e.parent.msgSeqGen.SetReachMaxRetryLogonFail(2)
		de := e.parent.Reset(true)
		if de != nil {
			domain_error.ProcessSevereError(false, 0, de, nil, "fail to reset FIX client after session logout!")
		}
	}
}

// FromAdmin implemented as part of Application interface
func (e *TradeClient) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	log.Printf("FromAdmin: %s", fixutil.ConvertFIXDataToString(msg.Bytes()))
	e.insertFixMessageToDB(msg, enum.UtilFixMessageSide_FromAdmin)
	msgType, _ := msg.MsgType()
	if msgType == "3" {
		e.processRejectMessage(msg)
	}
	return nil
}

func (e *TradeClient) insertFixMessageToDB(msg *quickfix.Message, side enum.UtilFixMessageSide) {
	msgType, err := msg.MsgType()
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs while parsing FIX message type")
	}
	fixMsg := &schema.UtilFixMessage{
		MsgSide:     int(side),
		MsgType:     msgType,
		MsgTime:     timeutil.ConvertTimeToMilliseconds(time.Now()),
		Data:        msg.Bytes(),
		ChannelCode: e.parent.GetChannelConfig().GetChannelCode(),
	}
	e.parent.cfg.GetAutoTx().Input(e.parent.cfg.GetAppDB(), func(tx *sql.Tx) (de *domain_error.Error) {
		err := app_store.InsertUtilFixMessage(tx, fixMsg)
		if err != nil && !dbutil.IsMysqlDuplicateEntryError(err) { // 排除重复插入的错误
			de = domain_error.Build(domain_error.DATABASE_OPERATION_ERR_CODE, err)
			return
		}
		return
	}, "", "")

	if side == enum.UtilFixMessageSide_FromAdmin || side == enum.UtilFixMessageSide_ToAdmin {
		e.parent.cfg.GetAutoTx().Flush()
	}
}

// ToAdmin implemented as part of Application interface
func (e *TradeClient) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {

	if msg.IsMsgTypeOf("A") { // Logon消息类型为A

		if e.username != "" {
			msg.Body.SetString(553, e.username) // Username
			msg.Body.SetString(554, e.password) // Password
		}

		log.Printf("detect logon message to Admin, e.parent.msgSeqGen.GetReachMaxRetryLogonFail:%v\n", e.parent.msgSeqGen.GetReachMaxRetryLogonFail())
		if e.parent.msgSeqGen.GetReachMaxRetryLogonFail() == 1 {
			restNum, err := msg.Body.GetString(141)
			log.Printf("restNum:%s, err:%v\n", restNum, err)
			if restNum != "N" {
				msg.Body.SetField(141, quickfix.FIXString("N"))
				log.Println("===> send login mesage to server with 141=N by force!!!")
			}
		}
	}

	log.Printf("ToAdmin: %s", fixutil.ConvertFIXDataToString(msg.Bytes()))
	if e.kafkaProducer != nil {
		e.kafkaProducer.SendMessage(e.kafkaTopic, msg.Bytes())
	}
	e.insertFixMessageToDB(msg, enum.UtilFixMessageSide_ToAdmin)
}

// ToApp implemented as part of Application interface
func (e *TradeClient) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
	log.Printf("ToApp: %s", fixutil.ConvertFIXDataToString(msg.Bytes()))
	if e.kafkaProducer != nil {
		e.kafkaProducer.SendMessage(e.kafkaTopic, msg.Bytes())
	}
	return
}

// FromApp implemented as part of Application interface. This is the callback for all Application level messages from the counter party.
func (e *TradeClient) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	log.Printf("FromApp: %s", fixutil.ConvertFIXDataToString(msg.Bytes()))
	// 检查消息类型
	msgType, _ := msg.MsgType()
	if msgType == "8" {
		e.processExecutionReport(msg)
	} else if msgType == "9" { // 还需要对OrderCancelReject等消息进行响应
		e.processOrderCancelReject(msg)
	} else if msgType != ""{
		e.insertFixMessageToDB(msg, enum.UtilFixMessageSide_FromApp)
	} else {
		log.Printf("Error getting MsgType, fix message: %s\n", msg.Bytes())
	}
	return
}

func (e *TradeClient) IsSessionReady() bool {
	sessionReady := atomic.LoadInt32(&e.sessionReady)
	if sessionReady > 0 {
		return true
	} else {
		return false
	}
}

func (e *TradeClient) setCredentials(cfg *domain_cfg.TradeChannelCfg) {
	for _, v := range cfg.GetTradeChannelCfgItems(){
		if v.ConfigItemName == "UserName" {
			e.username = v.ConfigItemValue
		}
		if v.ConfigItemName == "Password" {
			e.password = v.ConfigItemValue
		}
	}
}
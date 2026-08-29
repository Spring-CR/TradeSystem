package olts_fut

import (
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/timeutil"
	"rhino-core/domain_cfg"
	"rhino-core/order_domain/order_cache"
	"rhino-core/schema"
	"rhino-core/store/app_store"
	"rhino-core/types"
	"rhino-trade-channel/adapter/data_convert/domain_order_validate"
	"strings"
	"time"

	positionplugin "rhino-plugins/order_position_plugin/ficc"

	"github.com/IBM/sarama"
	json "github.com/bytedance/sonic"
)

const (
	TransactTimeLayout = "20060102-15:04:05"
)

var (
	oltsToFixSideMap = map[string]string{
		"1": "2",
		"3": "1",
	}
	fixToOltsMap = map[string]string{
		"2": "1",
		"1": "3",
	}
)

type OltsFutChannel struct {
	cfg                *domain_cfg.TradeChannelCfg
	kafkaClient        *kafkaClient
	tradeActionRespBuf chan *schema.TradeActionResp
	onTradeActionResp  func(tradeActionResp *schema.TradeActionResp)
	stopReceivingResp  bool
	tradeAccount       map[string]string
	operator           string
}

func NewOltsFutChannel(cfg *domain_cfg.TradeChannelCfg, onTradeActionResp func(tradeActionResp *schema.TradeActionResp)) (channel *OltsFutChannel, de *domain_error.Error) {

	// 测试配置
	_, err := timeutil.Until(cfg.GetTradeChannel().BeginTime, time.TimeOnly, cfg.GetTradeChannel().TimeZone)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, fmt.Sprintf("beginTime or timeZone of channel %s is not config correctly", cfg.GetChannelCode()))
	}

	channel = &OltsFutChannel{tradeActionRespBuf: make(chan *schema.TradeActionResp, 4096*2), onTradeActionResp: onTradeActionResp}
	channel.cfg = cfg

	var brokers []string
	var reqTopic string
	var respTopic string
	tradeAccount := make(map[string]string)
	for _, item := range cfg.GetTradeChannelCfgItems() {
		if item.ConfigItemName == "KafkaBroker" {
			brokers = strings.Split(item.ConfigItemValue, ",")
			for i, v := range brokers {
				brokers[i] = strings.TrimSpace(v)
			}
		}
		if item.ConfigItemName == "RequestTopic" {
			reqTopic = strings.TrimSpace(item.ConfigItemValue)
		}
		if item.ConfigItemName == "ResponseTopic" {
			respTopic = strings.TrimSpace(item.ConfigItemValue)
		}
		if item.ConfigItemName == "TradeAccount" {
			//channel.tradeAccount = strings.TrimSpace(item.ConfigItemValue)
			tradeAccounts := strings.Split(strings.TrimSpace(item.ConfigItemValue), ",")
			for _, str := range tradeAccounts {
				if len(str) == 0 {
					continue
				}
				kv := strings.Split(str, ":")
				if len(kv) != 2 {
					continue
				}
				if len(kv[0]) == 0 || len(kv[1]) == 0 {
					continue
				}
				tradeAccount[kv[0]] = kv[1]
			}
		}
		if item.ConfigItemName == "TradeOperator" {
			channel.operator = strings.TrimSpace(item.ConfigItemValue)
		}
	}

	channel.tradeAccount = tradeAccount

	kafkaClient := newKafkaClient(brokers, reqTopic, respTopic, sarama.OffsetOldest, channel.onDataReceived)
	channel.kafkaClient = kafkaClient

	// 开启回报消费协程
	channel.startProcessExecutionReportFromBuf()

	return
}

func (c *OltsFutChannel) GetChannelConfig() *domain_cfg.TradeChannelCfg {
	return c.cfg
}

func (c *OltsFutChannel) ValidateOrder(order *schema.TradeOrder) (channelOrder interface{}, de *domain_error.Error) {

	extendAttrMap := order.ExtendAttrMap
	if extendAttrMap == nil {
		extendAttrMap = map[string]interface{}{}
	}

	transactTime := timeutil.ConvertMillisecondsToTime(order.TransactTime).In(timeutil.CnTimeLocation).Format(TransactTimeLayout)

	log.Printf("ValidateOrder, transactTime:%v\n", transactTime)

	begin := time.Now()
	end := time.Now().Add(24 * time.Hour)
	strategyParameters := make(map[string]interface{})
	strategyParameters["TimeStart"] = begin.In(timeutil.CnTimeLocation).Format("20060102-15:04:05")
	strategyParameters["TimeEnd"] = end.In(timeutil.CnTimeLocation).Format("20060102-15:04:05")

	rawOrder := &NewOrderSingle{
		MsgType:            "D",
		Account:            c.tradeAccount[order.SecurityExchange],
		HandlInst:          order.HandlInst,
		Symbol:             order.Symbol2,
		SecurityID:         order.Symbol2,
		SecurityExchange:   order.SecurityExchange,
		Side:               fixToOltsMap[order.Side],
		TransactTime:       transactTime,
		OrderQty:           int(order.OrderQty),
		OrdType:            order.OrdType,
		Price:              order.Price,
		Currency:           order.Currency,
		CashOrderQty:       0,
		TimeInForce:        "0",
		OpenClose:          "AUTO",
		Operator:           c.operator,
		StrategyParameters: strategyParameters,
		TradePurpose:       "互换其他",
	}

	if order.AlgName == "" {
		rawOrder.AlgoName = "0"
	}

	businessType, _, _ := attrutil.GetAttrValue(order.ExtendAttrMap, "businessType", enum.AttrValueType_STRING)
	switch businessType {
	case "N":
		rawOrder.BusinessType = "N_CROSS_FUTURE_SWAP"
	case "C":
		rawOrder.BusinessType = "CHN_FUTURE_SWAP"
	}

	return rawOrder, nil
}

func (c *OltsFutChannel) NewOrderSingle(order *schema.TradeOrder, afterTradeOrderUpsert func(tradeOrder *schema.TradeOrder), syncTradeOrderAfterError func(tradeOrder *schema.TradeOrder)) (duplicatedOrder bool, de *domain_error.Error) {

	var channelOrder interface{}
	duplicatedOrder, channelOrder, de = domain_order_validate.CheckAndBaseSetDomainOrder(order, c.GetChannelConfig(), c.ValidateOrder, afterTradeOrderUpsert)
	if de != nil {
		return duplicatedOrder, de
	}
	// do nothing!
	if duplicatedOrder {
		return duplicatedOrder, nil
	}
	msg := channelOrder.(*NewOrderSingle)

	// 要在CheckAndBaseSetDomainOrder方法之后设置clOrdID
	msg.ClOrdID = order.ClOrdID

	err := c.kafkaClient.send(msg)
	if err != nil {
		de = domain_error.Build(domain_error.FIX_SEND_MSG_ERR_CODE, err)
		// 更新订单状态，设置为内部提交失败的状态
		_, ordStatus, ordStatusUpdateTime, err := app_store.UpdateTradeOrderOnSubmitFailedById(c.GetChannelConfig().GetAppDB(), order.ID, de.SimpleErrorString())
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "fail to UpdateTradeOrderOnSubmitFailedById")
		}
		order.OrdStatus = ordStatus
		order.OrdStatusUpdateTime = ordStatusUpdateTime
		order.OrderSubmitFailReason = de.SimpleErrorString()
		if syncTradeOrderAfterError != nil {
			syncTradeOrderAfterError(order)
		}
		return duplicatedOrder, de
	}

	return duplicatedOrder, nil
}

func (c *OltsFutChannel) onDataReceived(data []byte, msgSeq int64, msgTime time.Time) {
	msgType := c.getMessageType(data)

	log.Printf("======> onDataReceived:%s, msgType=%v\n", data, msgType)
	if c.stopReceivingResp {
		if duringTradingTime, _ := c.cfg.IsDuringTradingTime(0); duringTradingTime {
			c.stopReceivingResp = false
		}
	}

	switch msgType {
	case "8":
		executionReport := &ExecutionReport{}
		err := json.Unmarshal(data, executionReport)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to Unmarshal 8 message:%s\n", data))
		}

		if c.stopReceivingResp && executionReport.OrdStatus != string(enum.OrdStatus_DoneForDay) {
			log.Printf("======> ignore:%s, stopReceivingResp:%v, msgType=%v\n", data, c.stopReceivingResp, msgType)
			return
		}

		c.processExecutionReport(data, executionReport, msgSeq, msgTime)

	case "9":

		if c.stopReceivingResp {
			log.Printf("======> ignore:%s, stopReceivingResp:%v, msgType=%v\n", data, c.stopReceivingResp, msgType)
			return
		}

		orderCancelReject := &OrderCancelReject{}
		err := json.Unmarshal(data, orderCancelReject)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to Unmarshal 9 message:%s\n", data))
		}
		// 不能启用，因为在star返回pending时，理论上还不一定有委托状态，此时如果发生cancel request，它可能还是会返回OrdStatus==""的情形
		// 因此，在订单tradeOrder里要记录一下上次的非pendingcancel的状态，如果出现OrdStatus==“”时，就采用上次记录的状态
		if orderCancelReject.OrdStatus == "" {
			// 如果返回的msg9不带状态，需要打印出来
			log.Printf("received orderCancelReject without setting OrdStatus! msg=%s\n", data)
		}
		c.processOrderCancelReject(data, orderCancelReject, msgSeq, msgTime)
	default:
		log.Printf("detect unsupport message type for data: %s\n", data)
	}
}

func (c *OltsFutChannel) CleanTradeActionResp(tradeActionResp *schema.TradeActionResp) (ignore bool) {
	return
}

func (c *OltsFutChannel) RefineOrderCancelID(rawTargetClOrdID string, rawClOrdID string, order *schema.TradeOrder) (targetClOrdID, clOrdID string) {
	// 不改
	return rawTargetClOrdID, rawClOrdID
}

var statusUnderProcessing = map[string]bool{
	string(enum.OrdStatus_PendingReview):   true,
	string(enum.OrdStatus_Submit):          true,
	string(enum.OrdStatus_New):             true,
	string(enum.OrdStatus_PartiallyFilled): true,
	string(enum.OrdStatus_PendingCancel):   true,
	string(enum.OrdStatus_PendingNew):      true,
	string(enum.OrdStatus_PendingReplace):  true,
}

func (c *OltsFutChannel) OnMarketClosed(orderCache *order_cache.OrderCache) {

	// 收市后休眠5分钟才触发DoneForDay
	time.Sleep(5 * time.Minute)
	duringTradingTime, err := c.cfg.IsDuringTradingTime(0)
	if err != nil || duringTradingTime {
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "fail to determine if IsDuringTradingTime")
		}
		return
	}

	// 停止接收成交回报
	c.stopReceivingResp = true
	defer func() {
		// 重新接收成交回报
		log.Printf("start to receive trade resp at %v\n", time.Now())
		c.stopReceivingResp = false
	}()

	matchOrders, _ := orderCache.FilterOrderByFunction(func(order *types.TraceableTradeOrder) bool {
		if order.GetBasicInfo().ChannelCode != c.cfg.GetTradeChannel().ChannelCode {
			return false
		}

		basicInfo := order.GetBasicInfo()
		jsData, _ := json.Marshal(basicInfo)
		log.Printf("OnMarketClosed, check order:%s\n", jsData)

		if basicInfo.OrdStatus == string(enum.OrdStatus_PendingCancel) {
			basicInfo.OrdStatus = basicInfo.OrdStatus2
		}

		return statusUnderProcessing[basicInfo.OrdStatus]
	})

	log.Printf("pickup %d matched orders\n", len(matchOrders))
	for _, matchOrder := range matchOrders {
		order := matchOrder.GetBasicInfo()
		log.Printf("===> ClOrdID:%s, Status:%s\n", order.ClOrdID, order.OrdStatus)
		// 构造DoneForDay的执行回报
		var msgID string
		if order.ClOrdID != "" {
			msgID = order.ClOrdID + "-DoneForDay"
		} else {
			msgID = order.AppOrdID + "-DoneForDay"
		}
		executionReport := &ExecutionReport{
			MsgType:          "8",
			SecurityExchange: order.SecurityExchange,
			ClOrdID:          order.ClOrdID,
			OrigClOrdID:      "",
			OrdID:            order.OrdID,
			Symbol:           order.Symbol,
			ExecType:         "",
			ExecRefID:        order.AppOrdID,
			ExecTransType:    "",
			OrdStatus:        string(enum.OrdStatus_DoneForDay),
			LastQty:          0,
			LastPx:           0,
			OrdQty:           order.OrderQty,
			Price:            order.Price,
			Currency:         order.Currency,
			LeavesQty:        0,
			OrdType:          order.OrdType,
			ExecID:           msgID,
			TransactTime:     time.Now().In(time.UTC).Format(timeutil.TransactTimeLayout),
			Side:             order.Side,
			OrdRejReason:     "",
			Text:             "",
		}
		err := c.kafkaClient.sendDFD(executionReport)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "fail to send DFD for "+order.AppOrdID)
		}
	}

	d, _ := timeutil.Until(c.cfg.GetTradeChannel().BeginTime, time.TimeOnly, c.cfg.GetTradeChannel().TimeZone)
	if d > 23*time.Hour {
		d = 0
	}
	if d > 5*time.Minute {
		d = 5 * time.Minute
	}
	time.Sleep(d)

	positionplugin.CapitalController.UnFreezeAllCapital()

	d, _ = timeutil.Until(c.cfg.GetTradeChannel().BeginTime, time.TimeOnly, c.cfg.GetTradeChannel().TimeZone)
	if d > 23*time.Hour {
		d = 0
	}

	time.Sleep(d)
}

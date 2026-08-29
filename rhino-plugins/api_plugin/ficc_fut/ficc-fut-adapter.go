package ficc_fut

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"rhino-common/context/constant"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/timeutil"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/types"
	"rhino-data/datamap"
	"rhino-engine/api/api_adapter/plugin/util"
	"strconv"
	"time"

	json "github.com/bytedance/sonic"
)

type FiccFurAPIAdapter struct {
	configMap      map[string]*schema.ApplicationCfgItem
	applicationCfg *domain_cfg.ApplicationCfg
	autoSyncRepo   *datamap.AutoSyncRepo
	// 与交易时间检查有关
	dmsDayTradeTimeRange   [2]int
	dmsNightTradeTimeRange [2]int
	ovsTradeTimeRange      [2]int
	dmsTradeFlagInMinutes  [1440]uint32
	ovsTradeFlagInMinutes  [1440]uint32
	channelMap             map[string]*domain_cfg.TradeChannelDetails
}

var (
	fieldBytesAppOrdID    = []byte(`"origClOrdID":"`)
	lenFieldBytesAppOrdID = len(fieldBytesAppOrdID)
	orderCxlReqBytes      = []byte(`"msgType":"F"`)
)

func NewFiccFurAPIAdapter(applicationCfg *domain_cfg.ApplicationCfg) (inst *FiccFurAPIAdapter, de *domain_error.Error) {
	configMap := applicationCfg.GetApplicationCfgItemMap()
	inst = &FiccFurAPIAdapter{configMap: configMap, applicationCfg: applicationCfg, autoSyncRepo: applicationCfg.GetAutoSyncRepo()}
	inst.syncAppConfigData(applicationCfg, configMap)
	inst.waitForCalendarDataReady()
	inst.configDmsTradeTimeFlag()
	inst.configOvsTradeTimeFlag()
	inst.initChannelMap()
	return
}

func (a *FiccFurAPIAdapter) GetWorkerSharding(workerCount int, getOrderByAppOrdID func(appOrdID string) (order *schema.TradeOrder, ok bool)) func([]byte, *schema.TradeOrder, int) int {
	n := workerCount - 1
	return func(rawMsg []byte, decodeOrder *schema.TradeOrder, cumTaskCount int) (workerIndex int) {
		idx := bytes.Index(rawMsg, orderCxlReqBytes)
		if idx >= 0 {
			// 撤单，根据id找到原订单
			idx = bytes.Index(rawMsg, fieldBytesAppOrdID)
			if idx >= 0 {
				subBytes := rawMsg[idx+lenFieldBytesAppOrdID:]
				idx = bytes.IndexByte(subBytes, '"')
				if idx >= 0 {
					appOrdID := subBytes[:idx]
					tradeOrder, ok := getOrderByAppOrdID(string(appOrdID))
					if ok {
						log.Printf("WorkerAffinity: %d\n", tradeOrder.WorkerAffinity)
						return tradeOrder.WorkerAffinity
					}
				}
			}
		} else if decodeOrder != nil {
			tradeOrder, ok := getOrderByAppOrdID(decodeOrder.AppOrdID)
			if ok {
				log.Printf("WorkerAffinity: %d\n", tradeOrder.WorkerAffinity)
				return tradeOrder.WorkerAffinity
			}
		}
		return cumTaskCount & n
	}
}

func (a *FiccFurAPIAdapter) DistinguishIngressMessage(rawMsg []byte) (msgType enum.ApiMessageType, msgProps map[string]interface{}, msg interface{}, de *domain_error.Error) {

	err := json.Unmarshal(rawMsg, &msgProps)
	if err != nil {
		de = domain_error.Build(domain_error.DISTINGUISH_INGRESS_MSG_ERR_CODE, err, rawMsg)
		return
	}

	v, ok := msgProps["msgType"]
	if !ok {
		de = domain_error.Build(domain_error.DISTINGUISH_INGRESS_MSG_ERR_CODE, errors.New("原始消息不存在属性是msgType的字段"), rawMsg)
		return
	}

	_msgType, ok := v.(string)
	if !ok {
		de = domain_error.Build(domain_error.DISTINGUISH_INGRESS_MSG_ERR_CODE, errors.New("msgType字段的值必须为字符串"), rawMsg)
		return
	}

	switch _msgType {
	case "D":
		msgType = enum.ApiMessageType_NewOrderSingle
	case "F":
		msgType = enum.ApiMessageType_OrderCancelRequest
	default:
		de = domain_error.Build(domain_error.DISTINGUISH_INGRESS_MSG_ERR_CODE, errors.New("不支持的消息类型，msgType="+_msgType), rawMsg)
		return
	}

	return
}

func (a *FiccFurAPIAdapter) returnOfConvertNewOrderSingleMessageOnError(extendAttr []byte, msgProps map[string]interface{}, de1 *domain_error.Error) (*schema.TradeOrder, *domain_error.Error) {
	if extendAttr == nil {
		extendAttr, _ = json.Marshal(msgProps)
	}
	tradeOrder := &schema.TradeOrder{ExtendAttr: string(extendAttr)}
	de1.Refine(domain_error.ERROR, tradeOrder)
	return tradeOrder, de1
}

// 生成默认参数
func (a *FiccFurAPIAdapter) generateDefaultParams(msgProps map[string]interface{}) {

}

func (a *FiccFurAPIAdapter) ConvertNewOrderSingleMessage(rawMsg []byte, msgProps map[string]interface{}) (tradeOrder *schema.TradeOrder, de *domain_error.Error) {

	a.generateDefaultParams(msgProps)

	var extendAttr []byte

	account, _, de1 := util.GetIntValueInField(msgProps, "account")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
	}

	clOrdID, _, de1 := util.GetStringValueInField(msgProps, "clOrdID")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
	}

	// 固定
	handlInst := "1"
	if _, ok := msgProps["handlInst"]; ok {
		handlInst, _, de1 = util.GetStringValueInField(msgProps, "handlInst")
		if de1 != nil {
			return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
		}
	}

	symbol, _, de1 := util.GetStringValueInField(msgProps, "symbol2")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
	}

	side, _, de1 := util.GetStringValueInField(msgProps, "side")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
	}

	transactTime, de1 := getTimeStamp(rawMsg, msgProps, "transactTime")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
	}

	orderQty, _, de1 := util.GetFloatValueInField(msgProps, "orderQty")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
	}

	// 固定
	ordType := "2"
	if _, ok := msgProps["ordType"]; ok {
		ordType, _, de1 = util.GetStringValueInField(msgProps, "ordType")
		if de1 != nil {
			return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
		}
	}

	price, _, de1 := util.GetFloatValueInField(msgProps, "price")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
	}

	tradeOrder = &schema.TradeOrder{
		Account:            strconv.Itoa(account),
		AppOrdID:           clOrdID,
		HandlInst:          handlInst,
		Symbol:             symbol,
		Side:               side,
		TransactTime:       transactTime,
		OrderQty:           orderQty,
		OrdType:            ordType,
		Price:              price,
		ReviewFlag:         "1",
		Reviewer:           "system",
		IsDirectOrd:        true,
		ExtendAttrMap:      msgProps,
		OrdDraftUpdateUser: "system",
		OrdExecUser:        "system",
	}

	return
}

func (a *FiccFurAPIAdapter) ProcessNewOrderSingleError(tradeOrder *schema.TradeOrder, de *domain_error.Error) (rejectMsg []byte) {
	msg := make(map[string]interface{})
	if len(tradeOrder.ExtendAttr) > 0 {
		err := json.Unmarshal([]byte(tradeOrder.ExtendAttr), &msg)
		if err != nil {
			log.Printf("ProcessNewOrderSingleError:: fail to parse tradeOrder.ExtendAttr:%s\n", tradeOrder.ExtendAttr)
			domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("ProcessNewOrderSingleError:: fail to parse tradeOrder.ExtendAttr:%s\n", tradeOrder.ExtendAttr))
			msg = make(map[string]interface{})
		}
	}
	msg["msgType"] = "8"
	msg["account"] = tradeOrder.Account
	msg["execTransType"] = "0"
	msg["execType"] = "8"
	msg["ordStatus"] = "8"
	msg["lastQty"] = 0
	msg["lastPx"] = 0
	msg["leavesQty"] = 0
	msg["cumQty"] = 0
	msg["avgPx"] = 0
	msg["ordRejReason"] = "交易委托被拒绝: " + de.ErrorString()
	msg["transactTime"] = time.Now().In(time.UTC).Format(timeutil.TransactTimeLayout)

	// Todo：在生成错误回报时，必须设置
	msg[constant.ReqMsgSeq] = tradeOrder.MsgSeq

	//msg["text"] = fmt.Sprintf("交易委托已经被拒绝，原因：%s\n", de.ErrorString())
	rejectMsg, _ = json.Marshal(msg)
	return
}

// ExecutionReport 和 OrderCancelReject其实都在这里处理了。区别可以根据 CxlRejResponseTo 字段是否有值
func (a *FiccFurAPIAdapter) ConvertTradeResponseMessage(tradeResp *types.TradeActionRespReturn) (tradeRespMsg []byte, de *domain_error.Error) {

	msg := make(map[string]interface{})
	tradeOrder := tradeResp.GetTradeOrder()
	if len(tradeOrder.ExtendAttr) > 0 {
		err := json.Unmarshal([]byte(tradeOrder.ExtendAttr), &msg)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("ConvertTradeResponseMessage:: fail to parse tradeOrder.ExtendAttr:%s\n", tradeOrder.ExtendAttr))
			log.Printf("ConvertTradeResponseMessage:: fail to parse tradeOrder.ExtendAttr:%s\n", tradeOrder.ExtendAttr)
			msg = make(map[string]interface{})
		}
	}

	resp := tradeResp.CurrentTradeActionResp
	tradeActionLatestResp := tradeResp.GetTradeActionLatestResp()

	if resp.CxlRejResponseTo == "" { // ExecutionReport
		//常规回报，一定要支持system_guid字段
		msg[constant.TradeActionRespGuidField] = resp.GetCacheKey()

		msg["msgType"] = "8"
		msg["account"] = tradeOrder.Account
		msg["handlInst"] = tradeOrder.HandlInst
		msg["ordQty"] = tradeOrder.OrderQty
		msg["price"] = tradeOrder.Price
		msg["side"] = tradeOrder.Side
		msg["symbol"] = tradeOrder.Symbol
		msg["securityType"] = tradeOrder.SecurityType
		msg["securityExchange"] = tradeOrder.SecurityExchange
		msg["ordType"] = tradeOrder.OrdType
		msg["currency"] = resp.Currency
		msg["timeInForce"] = "0"
		msg["clOrdID"] = tradeActionLatestResp.ActionKey //tradeOrder.AppOrdID     // 应用层订单编号
		msg["ordID"] = tradeOrder.ClOrdID                // 用内部的clOrdID
		msg["origClOrdID"] = tradeOrder.AppOrdID         // 应用层订单编号
		msg["execID"] = resp.ExecID
		msg["execType"] = resp.ExecType
		msg["execTransType"] = resp.ExecTransType
		msg["execRefID"] = resp.ExecRefID
		msg["ordStatus"] = resp.OrdStatus
		msg["ordRejReason"] = resp.OrdRejReason
		msg["lastQty"] = resp.LastShares
		msg["lastPx"] = resp.LastPx
		msg["avgPx"] = resp.AvgPx
		msg["leavesQty"] = resp.LeavesQty
		msg["cumQty"] = resp.CumQty
		//msg["cumAmt"] = float64(resp.CumQty) * resp.AvgPx
		msg["transactTime"] = timeutil.ConvertMillisecondsToTime(resp.TransactTime).In(time.UTC).Format(timeutil.TransactTimeLayout)
	} else { // OrderCancelReject
		//常规回报，一定要支持system_guid字段
		msg[constant.TradeActionRespGuidField] = resp.GetCacheKey()

		msg["msgType"] = "9"
		msg["ordID"] = tradeOrder.ClOrdID                // 用内部的clOrdID
		msg["clOrdID"] = tradeActionLatestResp.ActionKey //tradeOrder.AppOrdID     // 应用层订单编号
		msg["origClOrdID"] = tradeOrder.AppOrdID         // 应用层订单编号
		msg["ordStatus"] = tradeOrder.OrdStatus
		msg["cxlRejResponseTo"] = resp.CxlRejResponseTo
		msg["cxlRejReason"] = resp.OrdRejReason
		msg["account"] = tradeOrder.Account
		msg["transactTime"] = timeutil.ConvertMillisecondsToTime(resp.TransactTime).In(time.UTC).Format(timeutil.TransactTimeLayout)
	}

	//text 不设置
	tradeRespMsg, _ = json.Marshal(msg)

	return
}

func (a *FiccFurAPIAdapter) ConvertOrderCancelRequestMessage(rawMsg []byte, msgProps map[string]interface{}) (ordCxlReq *types.ApplicationOrderCancelRequest, de *domain_error.Error) {

	ordCxlReq = &types.ApplicationOrderCancelRequest{}

	appOrdID, _, de1 := util.GetStringValueInField(msgProps, "OrigClOrdID")
	if de1 != nil {
		return ordCxlReq, de1
	}
	// 设置AppOrdID
	ordCxlReq.AppOrdID = appOrdID

	transactTime, de1 := getTimeStamp(rawMsg, msgProps, "TransactTime")
	if de1 != nil {
		return ordCxlReq, de1
	}

	actionKey, ok, de1 := util.GetStringValueInField(msgProps, "ClOrdID")
	if de1 != nil {
		return ordCxlReq, de1
	}
	if !ok {
		// 构造actionKey
		actionKey = appOrdID + "_cancel_" + strconv.Itoa(int(timeutil.ConvertTimeToMilliseconds(time.Now())))
	}

	ordCxlReq = &types.ApplicationOrderCancelRequest{AppOrdID: appOrdID, ActionTime: transactTime, ActionKey: actionKey}

	return
}

func (a *FiccFurAPIAdapter) ProcesssOrderCancelRequestError(ordCxlReq *types.ApplicationOrderCancelRequest, tradeOrder *schema.TradeOrder, de *domain_error.Error) (rejectMsg []byte) {
	msg := make(map[string]interface{})
	msg["msgType"] = "9"
	//msg["clOrdID"] = ordCxlReq.AppOrdID          // 应用层订单编号
	//msg["origClOrdID"] = ordCxlReq.AppOrdID      // 应用层订单编号
	//msg["secondaryClOrdID"] = ordCxlReq.AppOrdID // 应用层订单编号
	msg["cxlRejResponseTo"] = string(enum.CxlRejResponseTo_Cancel)
	msg["cxlRejReason"] = "撤单拒绝: " + de.ErrorString()
	msg["transactTime"] = time.Now().In(time.UTC).Format(timeutil.TransactTimeLayout)
	if tradeOrder != nil {
		msg["ordID"] = tradeOrder.ClOrdID // 用内部的clOrdID
		msg["ordStatus"] = tradeOrder.OrdStatus
		msg["account"] = tradeOrder.Account

		// Todo：在生成错误回报时，必须设置
		msg[constant.ReqMsgSeq] = tradeOrder.MsgSeq
	}
	rejectMsg, _ = json.Marshal(msg)
	return
}

func getTimeStamp(rawMsg []byte, msgProps map[string]interface{}, timeField string) (transactTime int64, de *domain_error.Error) {
	transactTimeStr, _, de1 := util.GetStringValueInField(msgProps, timeField)
	//log.Printf("======>transactTimeStr:%s, msgProps:%v\n", transactTimeStr, msgProps)
	if de1 != nil {
		return transactTime, de1
	}
	if transactTimeStr != "" {

		t, err := timeutil.ParseTimeStrToTimeByTimeLocation(timeutil.TransactTimeLayout, transactTimeStr, time.UTC) //20241107-09:46:10
		if err != nil {
			de = domain_error.Build(domain_error.TX_TIME_TOO_EARLY_ERR_CODE, err, transactTimeStr)
			return
		}
		if time.Since(t) > 60*time.Second {
			err = fmt.Errorf("expire transactTime:%s", transactTimeStr)
			de = domain_error.Build(domain_error.TX_TIME_TOO_EARLY_ERR_CODE, err, transactTimeStr)
			return
		}
		transactTime = timeutil.ConvertTimeToMilliseconds(t)
		return
	}
	return timeutil.ConvertTimeToMilliseconds(time.Now()), nil
}

func (a *FiccFurAPIAdapter) ErrorCouldBeIgnoreAfterReview(de *domain_error.Error) (ignoreAfterReview bool) {
	return
}

func (a *FiccFurAPIAdapter) ConvertTradeOrderParams(tradeOrder *schema.TradeOrder) (appOrd interface{}, de *domain_error.Error) {
	return tradeOrder, nil
}

func (a *FiccFurAPIAdapter) initChannelMap() {
	a.channelMap = make(map[string]*domain_cfg.TradeChannelDetails)
	for _, channelDetails := range a.applicationCfg.GetTradeChannels() {
		a.channelMap[channelDetails.TradeChannel.ChannelCode] = channelDetails
	}
}
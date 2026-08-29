package ficc

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"rhino-common/context/constant"
	"rhino-common/domain_error"
	"rhino-common/domain_error/notify_provider"
	"rhino-common/enum"
	"rhino-common/utils/timeutil"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/types"
	"rhino-data/datamap"
	"rhino-engine/api/api_adapter/plugin/util"
	"strconv"
	"strings"
	"time"
)

type TitansFiccAPIAdapter struct {
	configMap        map[string]*schema.ApplicationCfgItem
	applicationCfg   *domain_cfg.ApplicationCfg
	autoSyncRepo     *datamap.AutoSyncRepo
	tradeTimeBegin   int64
	tradeTimeEnd     int64
	t0SellTimeEnd    int64
	t0SellEndTimeStr string
	// 从配置文件设置的固定值
	execUser        string
	traderID        string
	twoInvestorID   [2]string
	initMarginRatio float64

	errNotifier *notify_provider.WechatErrorNotifyProvider
}

func NewTitansFiccAPIAdapter(applicationCfg *domain_cfg.ApplicationCfg) (inst *TitansFiccAPIAdapter, de *domain_error.Error) {
	configMap := applicationCfg.GetApplicationCfgItemMap()
	inst = &TitansFiccAPIAdapter{configMap: configMap, applicationCfg: applicationCfg}
	inst.syncAppConfigData(applicationCfg, configMap)
	inst.positionDailyRefresh(configMap)
	inst.initOrderConst(configMap)
	inst.initData(configMap)
	inst.initErrNotifier(configMap)
	return
}

var (
	fieldBytesAppOrdID    = []byte(`"origClOrdID":"`)
	lenFieldBytesAppOrdID = len(fieldBytesAppOrdID)
	orderCxlReqBytes      = []byte(`"msgType":"F"`)
)

func (a *TitansFiccAPIAdapter) GetWorkerSharding(workerCount int, getOrderByAppOrdID func(appOrdID string) (order *schema.TradeOrder, ok bool)) func([]byte, *schema.TradeOrder, int) int {
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

func (a *TitansFiccAPIAdapter) DistinguishIngressMessage(rawMsg []byte) (msgType enum.ApiMessageType, msgProps map[string]interface{}, msg interface{}, de *domain_error.Error) {

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

func (a *TitansFiccAPIAdapter) returnOfConvertNewOrderSingleMessageOnError(extendAttr []byte, msgProps map[string]interface{}, de1 *domain_error.Error) (*schema.TradeOrder, *domain_error.Error) {
	if extendAttr == nil {
		extendAttr, _ = json.Marshal(msgProps)
	}
	tradeOrder := &schema.TradeOrder{ExtendAttr: string(extendAttr)}
	de1.Refine(domain_error.ERROR, tradeOrder)
	return tradeOrder, de1
}

// 生成默认参数
func (a *TitansFiccAPIAdapter) generateDefaultParams(msgProps map[string]interface{}) {

	// ctpyShortName, _, _ := util.GetStringValueInField(msgProps, "counterparty")
	// if ctpyShortName != "" {
	// 	// 根据交易对手短名，设置InvestorID
	// 	if ctpyShortName == "广发全球资本FICC" {
	// 		msgProps["investorID"] = a.twoInvestorID[0]
	// 	} else {
	// 		msgProps["investorID"] = a.twoInvestorID[1]
	// 	}
	// }

	// 设置execUser
	msgProps["execUser"] = a.execUser
	// 设置traderID
	msgProps["traderID"] = a.traderID
	// 设置初保比率
	msgProps["initMarginRatio"] = a.initMarginRatio
	// 业务类型
	_, ok := msgProps["businessType"]
	if !ok {
		msgProps["businessType"] = "收益互换"
	}
	// 市场
	_, ok = msgProps["securityExchange"]
	if !ok {
		msgProps["securityExchange"] = "8"
	}
	// 合约方向
	_, ok = msgProps["side"]
	if !ok {
		msgProps["side"] = "1"
	}
	// 未来需要支持空头
	// 价格类型
	_, ok = msgProps["ordType"]
	if !ok {
		msgProps["ordType"] = "2"
	}
}

func (a *TitansFiccAPIAdapter) ConvertNewOrderSingleMessage(rawMsg []byte, msgProps map[string]interface{}) (tradeOrder *schema.TradeOrder, de *domain_error.Error) {

	a.generateDefaultParams(msgProps)
	// 不能返回空，否则，ProcessNewOrderSingleError 会报错

	// 原始参数，作为整个扩展属性
	// extendAttr, err := json.Marshal(msgProps)
	// if err != nil {
	// 	de = domain_error.Build(domain_error.CONVERT_NEW_ORDER_SINGLE_MSG_ERR_CODE, err, rawMsg)
	// 	return
	// }
	var extendAttr []byte

	id, _, de1 := util.GetIntValueInField(msgProps, "id")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
	}

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

	symbol, _, de1 := util.GetStringValueInField(msgProps, "symbol")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
	}

	securityExchange, _, de1 := util.GetStringValueInField(msgProps, "securityExchange")
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

	orderQty, _, de1 := util.GetFloatValueInField(msgProps, "quantity")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
	}

	cashOrderQty, _, de1 := util.GetFloatValueInField(msgProps, "cashOrderQty")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
	}

	ordType, _, de1 := util.GetStringValueInField(msgProps, "ordType")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
	}

	price, _, de1 := util.GetFloatValueInField(msgProps, "price")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
	}

	currency, _, de1 := util.GetStringValueInField(msgProps, "currency")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
	}

	execUser, _, de1 := util.GetStringValueInField(msgProps, "execUser")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
	}
	if execUser == "" {
		execUser, _, de1 = util.GetStringValueInField(msgProps, "counterparty")
		if de1 != nil {
			return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, msgProps, de1)
		}
	}

	draftUpdateUser := execUser

	channelCode := "ficc"

	tradeOrder = &schema.TradeOrder{
		ID:               int64(id),
		Account:          strconv.Itoa(account),
		AppOrdID:         clOrdID,
		HandlInst:        handlInst,
		Symbol:           symbol,
		SecurityExchange: securityExchange,
		Side:             side,
		TransactTime:     transactTime,
		OrderQty:         orderQty,
		CashOrderQty:     cashOrderQty,
		OrdType:          ordType,
		Price:            price,
		Currency:         currency,
		ReviewFlag:       "1",
		Reviewer:         "system",
		//ApproveStatus:      int(enum.ApproveStatus_Approved),
		IsDirectOrd:        true,
		ChannelCode:        channelCode,
		ExtendAttrMap:      msgProps,
		OrdDraftUpdateUser: draftUpdateUser,
		OrdExecUser:        execUser,
	}

	return
}

func (a *TitansFiccAPIAdapter) processExecUserInfo(extendAttrMap map[string]interface{}, tradeOrder *schema.TradeOrder) {
	execUser := tradeOrder.OrdExecUser
	execUserName := ""
	idx := strings.Index(execUser, "/")
	if idx > 0 {
		execUser = execUser[:idx]
		execUserName = tradeOrder.OrdExecUser[idx+1:]
		tradeOrder.OrdExecUser = execUser
		extendAttrMap["execUserName"] = execUserName
	}

	extendAttrMap["execUser"] = execUser
}

func (a *TitansFiccAPIAdapter) ConvertTradeOrderParams(tradeOrder *schema.TradeOrder) (appOrd interface{}, de *domain_error.Error) {
	if len(tradeOrder.ExtendAttrMap) > 0 {
		tradeOrder.ExtendAttrMap["id"] = tradeOrder.ID

		a.processExecUserInfo(tradeOrder.ExtendAttrMap, tradeOrder)

		return tradeOrder.ExtendAttrMap, nil
	} else if len(tradeOrder.ExtendAttr) > 0 {
		msgProps := map[string]interface{}{}
		err := json.Unmarshal([]byte(tradeOrder.ExtendAttr), &msgProps)
		if err != nil {
			de = domain_error.Build(domain_error.CONVERT_TRADE_ORDER_ERR_CODE, err)
			return
		}
		msgProps["id"] = tradeOrder.ID

		tradeOrder.ExtendAttrMap = msgProps

		a.processExecUserInfo(tradeOrder.ExtendAttrMap, tradeOrder)

		return msgProps, nil
	}
	return tradeOrder, nil
}

func (a *TitansFiccAPIAdapter) ProcessNewOrderSingleError(tradeOrder *schema.TradeOrder, de *domain_error.Error) (rejectMsg []byte) {
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
func (a *TitansFiccAPIAdapter) ConvertTradeResponseMessage(tradeResp *types.TradeActionRespReturn) (tradeRespMsg []byte, de *domain_error.Error) {

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

	if resp.CxlRejResponseTo == "" { // ExecutionReport
		//常规回报，一定要支持system_guid字段
		msg[constant.TradeActionRespGuidField] = resp.GetCacheKey()

		msg["msgType"] = "8"
		msg["clOrdID"] = tradeOrder.AppOrdID // 应用层订单编号
		//msg["origClOrdID"] = tradeOrder.AppOrdID      // 应用层订单编号
		//msg["secondaryClOrdID"] = tradeOrder.AppOrdID // 应用层订单编号
		msg["ordID"] = tradeOrder.ClOrdID // 用内部的clOrdID
		msg["execID"] = resp.ExecID
		//msg["execTransType"] = resp.ExecTransType
		msg["execType"] = resp.OrdStatus
		msg["ordStatus"] = resp.OrdStatus
		msg["ordRejReason"] = resp.OrdRejReason
		msg["account"] = tradeOrder.Account // 有问题，可能不是fix里的account
		msg["securityID"] = tradeOrder.SecurityID
		msg["symbol"] = tradeOrder.Symbol
		msg["side"] = tradeOrder.Side
		msg["quantity"] = tradeOrder.OrderQty
		//msg["cashOrderQty"] = resp.CashOrderQty
		msg["ordType"] = tradeOrder.OrdType
		msg["price"] = tradeOrder.Price
		//timeInForce 不设置
		msg["lastQty"] = resp.LastShares
		msg["lastPx"] = resp.LastPx
		msg["leavesQty"] = resp.LeavesQty
		msg["cumQty"] = resp.CumQty
		msg["avgPx"] = resp.AvgPx
		//msg["cumAmt"] = float64(resp.CumQty) * resp.AvgPx
		msg["transactTime"] = timeutil.ConvertMillisecondsToTime(resp.TransactTime).In(time.UTC).Format(timeutil.TransactTimeLayout)
	} else { // OrderCancelReject
		//常规回报，一定要支持system_guid字段
		msg[constant.TradeActionRespGuidField] = resp.GetCacheKey()

		msg["msgType"] = "9"
		msg["ordID"] = tradeOrder.ClOrdID    // 用内部的clOrdID
		msg["clOrdID"] = tradeOrder.AppOrdID // 应用层订单编号
		//msg["origClOrdID"] = tradeOrder.AppOrdID      // 应用层订单编号
		//msg["secondaryClOrdID"] = tradeOrder.AppOrdID // 应用层订单编号
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

func (a *TitansFiccAPIAdapter) ConvertOrderCancelRequestMessage(rawMsg []byte, msgProps map[string]interface{}) (ordCxlReq *types.ApplicationOrderCancelRequest, de *domain_error.Error) {

	ordCxlReq = &types.ApplicationOrderCancelRequest{}

	appOrdID, _, de1 := util.GetStringValueInField(msgProps, "origClOrdID")
	if de1 != nil {
		return ordCxlReq, de1
	}
	// 设置AppOrdID
	ordCxlReq.AppOrdID = appOrdID

	transactTime, de1 := getTimeStamp(rawMsg, msgProps, "transactTime")
	if de1 != nil {
		return ordCxlReq, de1
	}

	actionKey, ok, de1 := util.GetStringValueInField(msgProps, "actionKey")
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

func (a *TitansFiccAPIAdapter) ProcesssOrderCancelRequestError(ordCxlReq *types.ApplicationOrderCancelRequest, tradeOrder *schema.TradeOrder, de *domain_error.Error) (rejectMsg []byte) {
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

func (a *TitansFiccAPIAdapter) ErrorCouldBeIgnoreAfterReview(de *domain_error.Error) (ignoreAfterReview bool) {
	if de != nil && (de.Code == domain_error.QUOTA_LIMIT_EXCEEDED_ERR_CODE || de.Code == domain_error.CAPITAL_AMOUNT_NOT_ENOUGH_ERR_CODE) {
		de.Msg += "，订单已转为待审批，请联系业务人员"
		return true
	}
	return
}

// func (a *TitansFiccAPIAdapter) ExtractApplicationOrderID(rawMsg []byte, msgProps map[string]interface{})(appOrdID string, de *domain_error.Error){
// 	appOrdID, _, de = util.GetStringValueInField(msgProps, "clOrdID")
// 	if de != nil {
// 		return
// 	}
// 	return
// }
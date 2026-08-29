package titansnorth

import (
	//"encoding/json"
	"bytes"
	"errors"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/byteutils"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	"rhino-core/types"
	"rhino-engine/api/api_adapter/plugin/util"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// titans北向交易业务适配器

type TitansNAPIAdapter struct {
}

func NewTitansNAPIAdapter() (inst *TitansNAPIAdapter, de *domain_error.Error) {
	inst = &TitansNAPIAdapter{}
	return
}

var (
	fieldBytesForSharding    = []byte(`"account":"`)
	lenFieldBytesForSharding = len(fieldBytesForSharding)
	fieldBytesAppOrdID       = []byte(`"origClOrdID":"`)
	lenFieldBytesAppOrdID    = len(fieldBytesAppOrdID)
)

func (a *TitansNAPIAdapter) GetWorkerSharding(workerCount int, getOrderByAppOrdID func(appOrdID string) (order *schema.TradeOrder, ok bool)) func([]byte, *schema.TradeOrder, int) int {
	//return nil
	n := workerCount - 1
	return func(rawMsg []byte, decodeOrder *schema.TradeOrder, cumTaskCount int) (workerIndex int) {
		idx := bytes.Index(rawMsg, fieldBytesForSharding)
		if idx >= 0 {
			subBytes := rawMsg[idx+lenFieldBytesForSharding:]
			idx = bytes.IndexByte(subBytes, '"')
			if idx >= 0 {
				workerIndex := byteutils.PositiveDecimalBytesToInt(subBytes[:idx]) & n
				return workerIndex
			}
		} else {
			// 可能是撤单，根据id找到原订单
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
		}
		return cumTaskCount & n
	}
}

func (a *TitansNAPIAdapter) DistinguishIngressMessage(rawMsg []byte) (msgType enum.ApiMessageType, msgProps map[string]interface{}, msg interface{}, de *domain_error.Error) {

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
func (a *TitansNAPIAdapter) ConvertNewOrderSingleMessage(rawMsg []byte, msgProps map[string]interface{}) (tradeOrder *schema.TradeOrder, de *domain_error.Error) {

	// account, _ := msgProps["account"].(string)                   // as Account
	// clOrdID, _ := msgProps["clOrdID"].(string)                   // as AppOrdID
	// secondaryClOrdID, _ := msgProps["secondaryClOrdID"].(string) // as ExtendAttr.secondaryClOrdID
	// clientID, _ := msgProps["clientID"].(string)                 // as ExtendAttr.clientID
	// handlInst, _ := msgProps["handlInst"].(string)               // as HandlInst
	// securityID, _ := msgProps["securityID"].(string)             // as SecurityID
	// symbol, _ := msgProps["symbol"].(string)                     // as Symbol
	// securityExchange, _ := msgProps["securityExchange"].(string) // as SecurityExchange
	// side, _ := msgProps["side"].(string)                         // as Side
	// openClose, _ := msgProps["openClose"].(string)               // as OpenClose
	// transactTime, _ := msgProps["transactTime"].(string)            // 8时区，20241107-09:46:10 as TransactTime
	// orderQty, _ := msgProps["orderQty"].(float64)                // as OrderQty
	// cashOrderQty, _ := msgProps["cashOrderQty"].(float64)        // as CashOrderQty
	// ordType, _ := msgProps["OrdType"].(string)                   // as OrdType
	// price, _ := msgProps["price"].(float64)                      // as Price
	// currency, _ := msgProps["currency"].(string)                 // as Currency
	// targetStrategy, _ := msgProps["targetStrategy"].(string)     // as ExtendAttr.targetStrategy；order.AlgName，需要按http://wiki.gf.com.cn/pages/viewpage.action?spaceKey=ZYTZGL&title=New+Order+-+Single 转码；中信qfii的算法参数还不清楚如何设置
	// securityType, _ := msgProps["securityType"].(string)         // as SecurityType
	// businessType, _ := msgProps["businessType"].(string)         // as ExtendAttr.businessType
	// counterParty, _ := msgProps["counterParty"].(string)         // as ExtendAttr.counterParty
	// tradePurpose, _ := msgProps["tradePurpose"].(string)         // as ExtendAttr.tradePurpose
	// tradeChannel, _ := msgProps["tradeChannel"].(string)         // as ExtendAttr.tradeChannel
	// counterName, _ := msgProps["counterName"].(string)           // as ExtendAttr.counterName
	// orderUserId, _ := msgProps["orderUserId"].(string)           // as ExtendAttr.orderUserId
	// orderUserName, _ := msgProps["orderUserName"].(string)       // as ExtendAttr.orderUserName
	// internalTradeId, _ := msgProps["internalTradeId"].(string)   // as ExtendAttr.internalTradeId
	// hedgeFlag, _ := msgProps["hedgeFlag"].(string)               // as ExtendAttr.hedgeFlag
	// isInternalExec, _ := msgProps["isInternalExec"].(bool)       // as ExtendAttr.isInternalExec
	// timeInForce, _ := msgProps["timeInForce"].(string) 	        // 不需要设置，会被强制设置为当日

	// 只需要挑选订单需要的字段，其他字段直接填到扩展字段
	// 参考：http://wiki.gf.com.cn/pages/viewpage.action?pageId=257391093
	//      http://wiki.gf.com.cn/pages/viewpage.action?spaceKey=ZYTZGL&title=New+Order+-+Single
	// clientID/account、externalEntrustCode/clOrdID 是等价字段

	account, _, de1 := util.GetStringValuePerhapsInTwoFields(msgProps, "account", "clientID")
	if de1 != nil {
		return tradeOrder, de1
	}

	clOrdID, _, de1 := util.GetStringValuePerhapsInTwoFields(msgProps, "externalEntrustCode", "clOrdID")
	if de1 != nil {
		return tradeOrder, de1
	}

	handlInst, _, de1 := util.GetStringValueInField(msgProps, "handlInst")
	if de1 != nil {
		return tradeOrder, de1
	}
	if handlInst == "" {
		handlInst = "1"
	}

	securityID, _, de1 := util.GetStringValueInField(msgProps, "securityID")
	if de1 != nil {
		return tradeOrder, de1
	}

	symbol, _, de1 := util.GetStringValueInField(msgProps, "symbol")
	if de1 != nil {
		return tradeOrder, de1
	}

	securityExchange, _, de1 := util.GetStringValueInField(msgProps, "securityExchange")
	if de1 != nil {
		return tradeOrder, de1
	}

	side, _, de1 := util.GetStringValueInField(msgProps, "side")
	if de1 != nil {
		return tradeOrder, de1
	}

	openClose, _, de1 := util.GetStringValueInField(msgProps, "openClose")
	if de1 != nil {
		return tradeOrder, de1
	}

	transactTime, de1 := getTimeStamp(msgProps, "transactTime")
	if de1 != nil {
		return tradeOrder, de1
	}

	orderQty, _, de1 := util.GetFloatValueInField(msgProps, "orderQty")
	if de1 != nil {
		return tradeOrder, de1
	}

	cashOrderQty, _, de1 := util.GetFloatValueInField(msgProps, "cashOrderQty")
	if de1 != nil {
		return tradeOrder, de1
	}

	ordType, _, de1 := util.GetStringValueInField(msgProps, "ordType")
	if de1 != nil {
		return tradeOrder, de1
	}

	price, _, de1 := util.GetFloatValueInField(msgProps, "price")
	if de1 != nil {
		return tradeOrder, de1
	}

	currency, _, de1 := util.GetStringValueInField(msgProps, "currency")
	if de1 != nil {
		return tradeOrder, de1
	}

	// securityType, 有些问题，建议不要传了

	// 算法参数，集中处理
	targetStrategy, ok, de1 := util.GetStringValueInField(msgProps, "targetStrategy")
	if de1 != nil {
		return tradeOrder, de1
	}
	var algParams string
	if ok {
		rawParams, ok := msgProps["strategyParameters"]
		if ok {
			algParamsJson, err := json.Marshal(rawParams)
			if err != nil {
				de = domain_error.Build(domain_error.CONVERT_NEW_ORDER_SINGLE_MSG_ERR_CODE, err)
				return
			}
			algParams = string(algParamsJson)
		}
	}

	// 原始参数，作为整个扩展属性
	extendAttr, err := json.Marshal(msgProps)
	if err != nil {
		de = domain_error.Build(domain_error.CONVERT_NEW_ORDER_SINGLE_MSG_ERR_CODE, err)
		return
	}

	tradeOrder = &schema.TradeOrder{
		Account:          account,
		AppOrdID:         clOrdID,
		HandlInst:        handlInst,
		SecurityID:       securityID,
		Symbol:           symbol,
		SecurityExchange: securityExchange,
		Side:             side,
		OpenClose:        openClose,
		TransactTime:     transactTime,
		OrderQty:         orderQty,
		CashOrderQty:     cashOrderQty,
		OrdType:          ordType,
		Price:            price,
		Currency:         currency,
		AlgName:          targetStrategy,
		AlgParams:        algParams,
		ExtendAttr:       string(extendAttr),
		IsDirectOrd:      true,
		ChannelCode:      "citic_qfii",
	}

	return
}

func (a *TitansNAPIAdapter) RefineAndValidate(tradeOrder *schema.TradeOrder, trade bool) *domain_error.Error {
	return nil
}

func (a *TitansNAPIAdapter) ProcessNewOrderSingleError(tradeOrder *schema.TradeOrder, de *domain_error.Error) (rejectMsg []byte) {
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
	msg["lastShares"] = 0
	msg["lastPx"] = 0
	msg["leavesQty"] = tradeOrder.OrderQty
	msg["cumQty"] = 0
	msg["avgPx"] = 0
	msg["ordRejReason"] = "交易委托被拒绝: " + de.ErrorString()
	msg["transactTime"] = time.Now().In(timeutil.CnTimeLocation).Format("20060102-15:04:05")
	//msg["text"] = fmt.Sprintf("交易委托已经被拒绝，原因：%s\n", de.ErrorString())
	rejectMsg, _ = json.Marshal(msg)
	return
}

// ExecutionReport 和 OrderCancelReject其实都在这里处理了。区别可以根据 CxlRejResponseTo 字段是否有值
func (a *TitansNAPIAdapter) ConvertTradeResponseMessage(tradeResp *types.TradeActionRespReturn) (tradeRespMsg []byte, de *domain_error.Error) {

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
		msg["msgType"] = "8"
		msg["clOrdID"] = tradeOrder.AppOrdID          // 应用层订单编号
		msg["origClOrdID"] = tradeOrder.AppOrdID      // 应用层订单编号
		msg["secondaryClOrdID"] = tradeOrder.AppOrdID // 应用层订单编号
		msg["orderID"] = tradeOrder.ClOrdID           // 用内部的clOrdID
		msg["execID"] = resp.ExecID
		msg["execTransType"] = resp.ExecTransType
		msg["execType"] = resp.ExecType
		msg["ordStatus"] = resp.OrdStatus
		msg["ordRejReason"] = resp.OrdRejReason
		msg["account"] = resp.Account // 有问题，可能不是fix里的account
		msg["securityID"] = resp.SecurityID
		msg["symbol"] = resp.Symbol
		msg["side"] = resp.Side
		msg["openClose"] = resp.OpenClose
		msg["orderQty"] = resp.OrderQty
		msg["cashOrderQty"] = resp.CashOrderQty
		msg["ordType"] = resp.OrdType
		msg["price"] = resp.Price
		//timeInForce 不设置
		msg["lastShares"] = resp.LastShares
		msg["lastPx"] = resp.LastPx
		msg["leavesQty"] = resp.LeavesQty
		msg["cumQty"] = resp.CumQty
		msg["avgPx"] = resp.AvgPx
		//msg["cumAmt"] = float64(resp.CumQty) * resp.AvgPx
		msg["transactTime"] = timeutil.ConvertMillisecondsToTime(resp.TransactTime).In(timeutil.CnTimeLocation).Format("20060102-15:04:05")
	} else { // OrderCancelReject
		msg["msgType"] = "9"
		msg["orderID"] = tradeOrder.ClOrdID           // 用内部的clOrdID
		msg["clOrdID"] = tradeOrder.AppOrdID          // 应用层订单编号
		msg["origClOrdID"] = tradeOrder.AppOrdID      // 应用层订单编号
		msg["secondaryClOrdID"] = tradeOrder.AppOrdID // 应用层订单编号
		msg["ordStatus"] = resp.OrdStatus
		msg["cxlRejResponseTo"] = resp.CxlRejResponseTo
		msg["cxlRejReason"] = resp.OrdRejReason
		msg["account"] = resp.Account
		msg["transactTime"] = timeutil.ConvertMillisecondsToTime(resp.TransactTime).In(timeutil.CnTimeLocation).Format("20060102-15:04:05")
	}

	//text 不设置
	tradeRespMsg, _ = json.Marshal(msg)

	return
}

func (a *TitansNAPIAdapter) ConvertOrderCancelRequestMessage(rawMsg []byte, msgProps map[string]interface{}) (ordCxlReq *types.ApplicationOrderCancelRequest, de *domain_error.Error) {

	appOrdID, _, de1 := util.GetStringValueInField(msgProps, "origClOrdID")
	if de1 != nil {
		return ordCxlReq, de1
	}

	transactTime, de1 := getTimeStamp(msgProps, "transactTime")
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

func (a *TitansNAPIAdapter) ProcesssOrderCancelRequestError(ordCxlReq *types.ApplicationOrderCancelRequest, tradeOrder *schema.TradeOrder, de *domain_error.Error) (rejectMsg []byte) {
	msg := make(map[string]interface{})
	msg["msgType"] = "9"
	msg["clOrdID"] = ordCxlReq.AppOrdID          // 应用层订单编号
	msg["origClOrdID"] = ordCxlReq.AppOrdID      // 应用层订单编号
	msg["secondaryClOrdID"] = ordCxlReq.AppOrdID // 应用层订单编号
	msg["cxlRejResponseTo"] = string(enum.CxlRejResponseTo_Cancel)
	msg["cxlRejReason"] = "撤单拒绝: " + de.ErrorString()
	msg["transactTime"] = time.Now().In(timeutil.CnTimeLocation).Format("20060102-15:04:05")
	if tradeOrder != nil {
		msg["orderID"] = tradeOrder.ClOrdID // 用内部的clOrdID
		msg["ordStatus"] = tradeOrder.OrdStatus
		msg["account"] = tradeOrder.Account
	}
	rejectMsg, _ = json.Marshal(msg)
	return
}

func getTimeStamp(msgProps map[string]interface{}, timeField string) (transactTime int64, de *domain_error.Error) {
	transactTimeStr, _, de1 := util.GetStringValueInField(msgProps, timeField)
	if de1 != nil {
		return transactTime, de1
	}
	if transactTimeStr != "" {

		t, err := timeutil.ParseTimeStrToTimeByTimeLocation("20060102-15:04:05", transactTimeStr, timeutil.CnTimeLocation) //20241107-09:46:10
		if err != nil {
			de = domain_error.Build(domain_error.CONVERT_NEW_ORDER_SINGLE_MSG_ERR_CODE, err)
			return
		}
		if time.Since(t) > 6*time.Hour {
			err = fmt.Errorf("expire transactTime:%s", transactTimeStr)
			de = domain_error.Build(domain_error.CONVERT_NEW_ORDER_SINGLE_MSG_ERR_CODE, err)
			return
		}
		transactTime = timeutil.ConvertTimeToMilliseconds(t)
	}
	return
}

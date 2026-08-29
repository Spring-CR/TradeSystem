package atp_cbf

import (
	//"encoding/json"
	"bytes"
	"encoding/json"
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
	"rhino-engine/api/api_adapter/plugin/util"
	"strconv"
	"strings"
	"time"
)

// titans北向交易业务适配器

type TitansCrossBorderFutureAPIAdapter struct {
}

func NewTitansCrossBorderFutureAPIAdapter(applicationCfg *domain_cfg.ApplicationCfg) (inst *TitansCrossBorderFutureAPIAdapter, de *domain_error.Error) {
	inst = &TitansCrossBorderFutureAPIAdapter{}
	return
}

var (
	fieldBytesAppOrdIDCBF    = []byte(`"keyOrderId":`)
	lenFieldBytesAppOrdIDCBF = len(fieldBytesAppOrdIDCBF)
	orderCxlReqBytesCBF      = []byte(`"WITHDRAW"`)
)

func (a *TitansCrossBorderFutureAPIAdapter) GetWorkerSharding(workerCount int, getOrderByAppOrdID func(appOrdID string) (order *schema.TradeOrder, ok bool)) func([]byte, *schema.TradeOrder, int) int {
	n := workerCount - 1
	return func(rawMsg []byte, decodeOrder *schema.TradeOrder, cumTaskCount int) (workerIndex int) {
		idx := bytes.Index(rawMsg, orderCxlReqBytesCBF)
		if idx >= 0 {
			// 撤单，根据id找到原订单
			idx = bytes.Index(rawMsg, fieldBytesAppOrdIDCBF)
			if idx >= 0 {
				subBytes := rawMsg[idx+lenFieldBytesAppOrdIDCBF:]
				idx = bytes.IndexByte(subBytes, ',')
				if idx >= 0 {
					appOrdID := bytes.TrimSpace(subBytes[:idx])
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

func (a *TitansCrossBorderFutureAPIAdapter) DistinguishIngressMessage(rawMsg []byte) (msgType enum.ApiMessageType, msgProps map[string]interface{}, msg interface{}, de *domain_error.Error) {

	err := json.Unmarshal(rawMsg, &msgProps)
	if err != nil {
		de = domain_error.Build(domain_error.DISTINGUISH_INGRESS_MSG_ERR_CODE, err, rawMsg)
		return
	}

	v, ok := msgProps["commandType"]
	if !ok {
		de = domain_error.Build(domain_error.DISTINGUISH_INGRESS_MSG_ERR_CODE, errors.New("原始消息不存在属性是commandType的字段"), rawMsg)
		return
	}

	_msgType, ok := v.(string)
	if !ok {
		de = domain_error.Build(domain_error.DISTINGUISH_INGRESS_MSG_ERR_CODE, errors.New("commandType字段的值必须为字符串"), rawMsg)
		return
	}

	switch _msgType {
	case "ENTRUST":
		msgType = enum.ApiMessageType_NewOrderSingle
	case "WITHDRAW":
		msgType = enum.ApiMessageType_OrderCancelRequest
	default:
		de = domain_error.Build(domain_error.DISTINGUISH_INGRESS_MSG_ERR_CODE, errors.New("不支持的消息类型，commandType="+_msgType), rawMsg)
		return
	}

	log.Printf("======> DistinguishIngressMessage, msgType:%v\n", msgType)

	return
}

func (a *TitansCrossBorderFutureAPIAdapter) returnOfConvertNewOrderSingleMessageOnError(extendAttr []byte, de1 *domain_error.Error) (*schema.TradeOrder, *domain_error.Error) {
	tradeOrder := &schema.TradeOrder{ExtendAttr: string(extendAttr)}
	return tradeOrder, de1
}

func (a *TitansCrossBorderFutureAPIAdapter) ConvertNewOrderSingleMessage(rawMsg []byte, msgProps map[string]interface{}) (tradeOrder *schema.TradeOrder, de *domain_error.Error) {

	log.Printf("======> received NewOrderSingleMessage:%s\n", rawMsg)
	// 不能返回空，否则，ProcessNewOrderSingleError 会报错

	// 原始参数，作为整个扩展属性
	extendAttr, err := json.Marshal(msgProps)
	if err != nil {
		de = domain_error.Build(domain_error.CONVERT_NEW_ORDER_SINGLE_MSG_ERR_CODE, err, rawMsg)
		return
	}

	id, _, de1 := util.GetIntValueInField(msgProps, "id")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	account, _, de1 := util.GetStringValueInField(msgProps, "subAccountNumber")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	// originalProps := msgProps
	// var ok bool
	// msgProps, ok = msgProps["data"].(map[string]interface{})
	// if !ok {
	// 	msgProps = map[string]interface{}{}
	// }

	_clOrdID, _, de1 := util.GetIntValueInField(msgProps, "keyOrderId")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}
	clOrdID := strconv.Itoa(_clOrdID)

	// handlInst, _, de1 := util.GetStringValueInField(msgProps, "handlInst")
	// if de1 != nil {
	// 	return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	// }
	// if handlInst == "" {
	// 	handlInst = "1"
	// }
	// 固定，以后有需要再做
	handlInst := "1"

	// securityID, _, de1 := util.GetStringValueInField(msgProps, "securityID")
	// if de1 != nil {
	// 	return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	// }

	symbol, _, de1 := util.GetStringValueInField(msgProps, "atpInstrumentId")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	securityExchange, _, de1 := util.GetStringValueInField(msgProps, "securityExchange")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	side, _, de1 := util.GetStringValueInField(msgProps, "orderDirection")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}
	switch side {
	case "SELL":
		side = string(enum.Side_Sell)
	case "BUY":
		side = string(enum.Side_Buy)
	}

	transactTime, de1 := getTimeStamp(rawMsg, msgProps, "transactTime")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}
	if transactTime == 0 {
		transactTime = timeutil.ConvertTimeToMilliseconds(time.Now())
	}

	orderQty, _, de1 := util.GetFloatValueInField(msgProps, "quantity")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	ordType, _, de1 := util.GetStringValueInField(msgProps, "priceType")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}
	switch ordType {
	case "LimitOrder":
		ordType = string(enum.OrdType_LIMIT)
	case "MarketOrder":
		ordType = string(enum.OrdType_MARKET)
	default:
		ordType = string(enum.OrdType_LIMIT)
	}

	price, _, de1 := util.GetFloatValueInField(msgProps, "limitPrice")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	draftUpdateUser, _, de1 := util.GetStringValueInField(msgProps, "draftUpdateUser")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	execUser, _, de1 := util.GetStringValueInField(msgProps, "execUser")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	channelCode, ok, de1 := util.GetStringValueInField(msgProps, "channelCode")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}
	if !ok {
		channelCode = "atp_cbf" // cross border future
	}

	// securityType, 有些问题，建议不要传了

	// 算法参数，集中处理
	var algoParams string
	var algoName string
	var rawParams interface{}
	var algParamsMap map[string]interface{}
	rawParams, ok = msgProps["algo"]
	if ok {
		algParamsJson, err := json.Marshal(rawParams)
		if err != nil {
			de = domain_error.Build(domain_error.CONVERT_NEW_ORDER_SINGLE_MSG_ERR_CODE, err, rawMsg)
			return
		}
		algoParams = string(algParamsJson)
		algParamsMap, ok = rawParams.(map[string]interface{})
		if ok {
			algoName, _ = algParamsMap["algoName"].(string)
		}
	}

	tradeOrder = &schema.TradeOrder{
		ID:                 int64(id),
		Account:            account,
		AppOrdID:           clOrdID,
		HandlInst:          handlInst,
		Symbol:             symbol,
		SecurityExchange:   securityExchange,
		Side:               side,
		TransactTime:       transactTime,
		OrderQty:           orderQty,
		OrdType:            ordType,
		Price:              price,
		AlgName:            algoName,
		AlgParams:          algoParams,
		ExtendAttr:         string(extendAttr),
		IsDirectOrd:        true,
		ChannelCode:        channelCode,
		ExtendAttrMap:      msgProps,
		AlgParamsMap:       algParamsMap,
		OrdDraftUpdateUser: draftUpdateUser,
		OrdExecUser:        execUser,
	}

	log.Printf("======>Received new order %s\n", rawMsg)
	return
}

func (a *TitansCrossBorderFutureAPIAdapter) RefineAndValidate(tradeOrder *schema.TradeOrder, trade bool) *domain_error.Error {
	return nil
}

func (a *TitansCrossBorderFutureAPIAdapter) processExecUserInfo(extendAttrMap map[string]interface{}, tradeOrder *schema.TradeOrder) {
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

// 将TradeOrder模型转化为应用层的静态交易参数数据结构
func (a *TitansCrossBorderFutureAPIAdapter) ConvertTradeOrderParams(tradeOrder *schema.TradeOrder) (appOrd interface{}, de *domain_error.Error) {
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

func (a *TitansCrossBorderFutureAPIAdapter) ProcessNewOrderSingleError(tradeOrder *schema.TradeOrder, de *domain_error.Error) (rejectMsg []byte) {
	log.Printf("start to ProcessNewOrderSingleError, tradeOrder:%v\n", tradeOrder)
	msg := make(map[string]interface{})
	log.Printf("msg1:%+v, tradeOrder.ExtendAttr:%s\n", msg == nil, tradeOrder.ExtendAttr)
	if len(tradeOrder.ExtendAttr) > 0 {
		err := json.Unmarshal([]byte(tradeOrder.ExtendAttr), &msg)
		if err != nil {
			log.Printf("ProcessNewOrderSingleError:: fail to parse tradeOrder.ExtendAttr:%s\n", tradeOrder.ExtendAttr)
			domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("ProcessNewOrderSingleError:: fail to parse tradeOrder.ExtendAttr:%s\n", tradeOrder.ExtendAttr))
			msg = make(map[string]interface{})
		}
	}
	log.Printf("msg2:%+v\n", msg == nil)
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
	msg["transactTime"] = timeutil.ConvertTimeToMilliseconds(time.Now())

	// Todo：在生成错误回报时，必须设置
	msg[constant.ReqMsgSeq] = tradeOrder.MsgSeq

	//msg["text"] = fmt.Sprintf("交易委托已经被拒绝，原因：%s\n", de.ErrorString())
	rejectMsg, _ = json.Marshal(msg)
	return
}

// ExecutionReport 和 OrderCancelReject其实都在这里处理了。区别可以根据 CxlRejResponseTo 字段是否有值
func (a *TitansCrossBorderFutureAPIAdapter) ConvertTradeResponseMessage(tradeResp *types.TradeActionRespReturn) (tradeRespMsg []byte, de *domain_error.Error) {

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
		msg["transactTime"] = resp.TransactTime

		dateStr := resp.ExchangeTradeDate
		if dateStr == "" {
			dateStr = strconv.Itoa(int(tradeOrder.TradeDate))
		}
		if len(dateStr) == 8 {
			dateStr = dateStr[0:4] + "-" + dateStr[4:6] + "-" + dateStr[6:]
		}
		msg["exchangeTradeDate"] = dateStr

		msg["ordCreateTime"] = tradeOrder.OrdCreateTime
	} else { // OrderCancelReject
		//常规回报，一定要支持system_guid字段
		msg[constant.TradeActionRespGuidField] = resp.GetCacheKey()

		msg["msgType"] = "9"
		msg["clOrdID"] = tradeOrder.AppOrdID          // 应用层订单编号
		msg["origClOrdID"] = tradeOrder.AppOrdID      // 应用层订单编号
		msg["secondaryClOrdID"] = tradeOrder.AppOrdID // 应用层订单编号
		msg["orderID"] = tradeOrder.ClOrdID           // 用内部的clOrdID
		msg["execID"] = resp.ClOrdID + "_" + resp.OrigClOrdID + "_exec"
		msg["execTransType"] = ""
		msg["execType"] = ""
		msg["ordStatus"] = tradeOrder.OrdStatus
		msg["ordRejReason"] = tradeOrder.OrdRejReason
		msg["account"] = tradeOrder.Account // 有问题，可能不是fix里的account
		msg["securityID"] = tradeOrder.SecurityID
		msg["symbol"] = tradeOrder.Symbol
		msg["side"] = tradeOrder.Side
		msg["openClose"] = tradeOrder.OpenClose
		msg["orderQty"] = tradeOrder.OrderQty
		msg["cashOrderQty"] = tradeOrder.CashOrderQty
		msg["ordType"] = tradeOrder.OrdType
		msg["price"] = tradeOrder.Price
		//timeInForce 不设置
		msg["lastShares"] = tradeOrder.LastShares
		msg["lastPx"] = tradeOrder.LastPx
		msg["leavesQty"] = tradeOrder.LeavesQty
		msg["cumQty"] = tradeOrder.CumQty
		msg["avgPx"] = tradeOrder.AvgPx
		//msg["cumAmt"] = float64(resp.CumQty) * resp.AvgPx
		msg["transactTime"] = tradeOrder.OrdStatusUpdateTime

		dateStr := resp.ExchangeTradeDate
		if dateStr == "" {
			dateStr = strconv.Itoa(int(tradeOrder.TradeDate))
		}
		if len(dateStr) == 8 {
			dateStr = dateStr[0:4] + "-" + dateStr[4:6] + "-" + dateStr[6:]
		}
		msg["exchangeTradeDate"] = dateStr

		msg["ordCreateTime"] = tradeOrder.OrdCreateTime

		msg["cxlRejResponseTo"] = resp.CxlRejResponseTo
		msg["cxlRejReason"] = resp.OrdRejReason
	}

	//text 不设置
	tradeRespMsg, _ = json.Marshal(msg)

	return
}

func (a *TitansCrossBorderFutureAPIAdapter) ConvertOrderCancelRequestMessage(rawMsg []byte, msgProps map[string]interface{}) (ordCxlReq *types.ApplicationOrderCancelRequest, de *domain_error.Error) {

	log.Printf("======> received OrderCancelRequestMessage:%s\n", rawMsg)

	ordCxlReq = &types.ApplicationOrderCancelRequest{}

	// var ok bool
	// msgProps, ok = msgProps["data"].(map[string]interface{})
	// if !ok {
	// 	msgProps = map[string]interface{}{}
	// }

	_appOrdID, _, de1 := util.GetIntValueInField(msgProps, "keyOrderId")
	if de1 != nil {
		return ordCxlReq, de1
	}
	if de1 != nil {
		return ordCxlReq, de1
	}
	appOrdID := strconv.Itoa(_appOrdID)

	// 设置AppOrdID
	ordCxlReq.AppOrdID = appOrdID

	transactTime, de1 := getTimeStamp(rawMsg, msgProps, "transactTime")
	if de1 != nil {
		return ordCxlReq, de1
	}
	if transactTime == 0 {
		transactTime = timeutil.ConvertTimeToMilliseconds(time.Now())
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

func (a *TitansCrossBorderFutureAPIAdapter) ProcesssOrderCancelRequestError(ordCxlReq *types.ApplicationOrderCancelRequest, tradeOrder *schema.TradeOrder, de *domain_error.Error) (rejectMsg []byte) {
	msg := make(map[string]interface{})
	msg["msgType"] = "9"
	msg["clOrdID"] = ordCxlReq.AppOrdID          // 应用层订单编号
	msg["origClOrdID"] = ordCxlReq.AppOrdID      // 应用层订单编号
	msg["secondaryClOrdID"] = ordCxlReq.AppOrdID // 应用层订单编号
	msg["cxlRejResponseTo"] = string(enum.CxlRejResponseTo_Cancel)
	msg["cxlRejReason"] = "撤单拒绝: " + de.ErrorString()
	msg["transactTime"] = timeutil.ConvertTimeToMilliseconds(time.Now())
	if tradeOrder != nil {
		msg["orderID"] = tradeOrder.ClOrdID // 用内部的clOrdID
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

		t, err := timeutil.ParseTimeStrToTimeByTimeLocation("20060102-15:04:05", transactTimeStr, timeutil.CnTimeLocation) //20241107-09:46:10
		if err != nil {
			de = domain_error.Build(domain_error.TX_TIME_TOO_EARLY_ERR_CODE, err, transactTimeStr)
			return
		}
		if time.Since(t) > 6*time.Hour {
			err = fmt.Errorf("expire transactTime:%s", transactTimeStr)
			de = domain_error.Build(domain_error.TX_TIME_TOO_EARLY_ERR_CODE, err, transactTimeStr)
			return
		}
		transactTime = timeutil.ConvertTimeToMilliseconds(t)
	}
	return
}

func (a *TitansCrossBorderFutureAPIAdapter) ErrorCouldBeIgnoreAfterReview(de *domain_error.Error) (ignoreAfterReview bool) {
	return
}

// func (a *TitansCrossBorderFutureAPIAdapter) ExtractApplicationOrderID(rawMsg []byte, msgProps map[string]interface{})(appOrdID string, de *domain_error.Error){
// 	return
// }

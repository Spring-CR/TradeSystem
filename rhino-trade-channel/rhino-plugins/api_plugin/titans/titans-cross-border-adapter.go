package titans

import (
	//"encoding/json"
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
	"rhino-engine/api/api_adapter/plugin/util"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// titans北向交易业务适配器

type TitansCrossBorderAPIAdapter struct {
	titansApiBase     string
	titansApiSecret   string
	softCheckPosition bool
}

func NewTitansCrossBorderAPIAdapter(applicationCfg *domain_cfg.ApplicationCfg) (inst *TitansCrossBorderAPIAdapter, de *domain_error.Error) {
	inst = &TitansCrossBorderAPIAdapter{}

	for _, item := range applicationCfg.GetApplicationCfgItems() {
		if item.ConfigItemName == "TitansApiBase" {
			inst.titansApiBase = item.ConfigItemValue
			inst.titansApiBase = strings.TrimRight(inst.titansApiBase, "/")
		}
		if item.ConfigItemName == "TitansApiSecret" {
			inst.titansApiSecret = item.ConfigItemValue
		}
		if item.ConfigItemName == "SoftCheckPosition" && item.ConfigItemValue == "1" {
			inst.softCheckPosition = true
		}
	}

	return
}

var (
	fieldBytesAppOrdID    = []byte(`"origClOrdID":"`)
	lenFieldBytesAppOrdID = len(fieldBytesAppOrdID)
	orderCxlReqBytes      = []byte(`"msgType":"F"`)

	qfiiAccounts = map[string]bool{
		"GF2CITICS": true,
		"HH2CITICS": true,
		"414":       true,
		"322":       true,
	}
)

func (a *TitansCrossBorderAPIAdapter) GetWorkerSharding(workerCount int, getOrderByAppOrdID func(appOrdID string) (order *schema.TradeOrder, ok bool)) func([]byte, *schema.TradeOrder, int) int {
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

func (a *TitansCrossBorderAPIAdapter) DistinguishIngressMessage(rawMsg []byte) (msgType enum.ApiMessageType, msgProps map[string]interface{}, msg interface{}, de *domain_error.Error) {

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

func (a *TitansCrossBorderAPIAdapter) returnOfConvertNewOrderSingleMessageOnError(extendAttr []byte, de1 *domain_error.Error) (*schema.TradeOrder, *domain_error.Error) {
	tradeOrder := &schema.TradeOrder{ExtendAttr: string(extendAttr)}
	return tradeOrder, de1
}

func (a *TitansCrossBorderAPIAdapter) ConvertNewOrderSingleMessage(rawMsg []byte, msgProps map[string]interface{}) (tradeOrder *schema.TradeOrder, de *domain_error.Error) {

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

	account, _, de1 := util.GetStringValueInField(msgProps, "account")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	clOrdID, _, de1 := util.GetStringValueInField(msgProps, "clOrdID")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

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

	symbol, _, de1 := util.GetStringValueInField(msgProps, "symbol")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	securityExchange, _, de1 := util.GetStringValueInField(msgProps, "securityExchange")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	side, _, de1 := util.GetStringValueInField(msgProps, "side")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	openClose, _, de1 := util.GetStringValueInField(msgProps, "openClose")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	transactTime, de1 := getTimeStamp(rawMsg, msgProps, "transactTime")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	orderQty, _, de1 := util.GetFloatValueInField(msgProps, "orderQty")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	cashOrderQty, _, de1 := util.GetFloatValueInField(msgProps, "cashOrderQty")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	ordType, _, de1 := util.GetStringValueInField(msgProps, "ordType")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	price, _, de1 := util.GetFloatValueInField(msgProps, "price")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}

	currency, _, de1 := util.GetStringValueInField(msgProps, "currency")
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
		if qfiiAccounts[account] {
			channelCode = "citic_qfii"
		} else {
			channelCode = "gfgfix"
			// 在titans目前的实现里，北向目前固定为1，南向为3
			handlInst = "3"
		}
	}

	// securityType, 有些问题，建议不要传了

	// 算法参数，集中处理
	targetStrategy, ok, de1 := util.GetStringValueInField(msgProps, "targetStrategy")
	if de1 != nil {
		return a.returnOfConvertNewOrderSingleMessageOnError(extendAttr, de1)
	}
	// 参考：http://wiki.gf.com.cn/display/ZYTZGL/New+Order+-+Single
	switch targetStrategy {
		case "1001": targetStrategy ="VWAP"
		case "1002": targetStrategy ="TWAP"
		case "1003": targetStrategy ="POV"
		case "1004": targetStrategy ="ICEBERG"
		case "1005": targetStrategy ="SNIPER"
		case "2001": targetStrategy ="VWAP"
		case "2002": targetStrategy ="TWAP"
		case "2003": targetStrategy ="POV"
	}
	var algParams string
	var rawParams interface{}
	if ok {
		rawParams, ok = msgProps["strategyParameters"]
		if ok {
			algParamsJson, err := json.Marshal(rawParams)
			if err != nil {
				de = domain_error.Build(domain_error.CONVERT_NEW_ORDER_SINGLE_MSG_ERR_CODE, err, rawMsg)
				return
			}
			algParams = string(algParamsJson)
		}
	}
	var algParamsMap map[string]interface{}
	if ok {
		algParamsMap, _ = rawParams.(map[string]interface{})
	}
	tradeOrder = &schema.TradeOrder{
		ID:                 int64(id),
		Account:            account,
		AppOrdID:           clOrdID,
		HandlInst:          handlInst,
		Symbol:             symbol,
		SecurityExchange:   securityExchange,
		Side:               side,
		OpenClose:          openClose,
		TransactTime:       transactTime,
		OrderQty:           orderQty,
		CashOrderQty:       cashOrderQty,
		OrdType:            ordType,
		Price:              price,
		Currency:           currency,
		AlgName:            targetStrategy,
		AlgParams:          algParams,
		ExtendAttr:         string(extendAttr),
		IsDirectOrd:        true,
		ChannelCode:        channelCode,
		ExtendAttrMap:      msgProps,
		AlgParamsMap:       algParamsMap,
		OrdDraftUpdateUser: draftUpdateUser,
		OrdExecUser:        execUser,
	}

	return
}

func (a *TitansCrossBorderAPIAdapter) RefineAndValidate(tradeOrder *schema.TradeOrder, trade bool) *domain_error.Error {
	// 卖出时持券数量检查
	a.checkPosition(tradeOrder)
	return nil
}

func (a *TitansCrossBorderAPIAdapter) processExecUserInfo(extendAttrMap map[string]interface{}, tradeOrder *schema.TradeOrder) {
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
func (a *TitansCrossBorderAPIAdapter) ConvertTradeOrderParams(tradeOrder *schema.TradeOrder) (appOrd interface{}, de *domain_error.Error) {
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

func (a *TitansCrossBorderAPIAdapter) ProcessNewOrderSingleError(tradeOrder *schema.TradeOrder, de *domain_error.Error) (rejectMsg []byte) {
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

	// Todo：在生成错误回报时，必须设置
	msg[constant.ReqMsgSeq] = tradeOrder.MsgSeq

	//msg["text"] = fmt.Sprintf("交易委托已经被拒绝，原因：%s\n", de.ErrorString())
	rejectMsg, _ = json.Marshal(msg)
	return
}

// ExecutionReport 和 OrderCancelReject其实都在这里处理了。区别可以根据 CxlRejResponseTo 字段是否有值
func (a *TitansCrossBorderAPIAdapter) ConvertTradeResponseMessage(tradeResp *types.TradeActionRespReturn) (tradeRespMsg []byte, de *domain_error.Error) {

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
		msg["transactTime"] = timeutil.ConvertMillisecondsToTime(resp.TransactTime).In(timeutil.CnTimeLocation).Format("20060102-15:04:05")
	} else { // OrderCancelReject
		//常规回报，一定要支持system_guid字段
		msg[constant.TradeActionRespGuidField] = resp.GetCacheKey()

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

func (a *TitansCrossBorderAPIAdapter) ConvertOrderCancelRequestMessage(rawMsg []byte, msgProps map[string]interface{}) (ordCxlReq *types.ApplicationOrderCancelRequest, de *domain_error.Error) {

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

func (a *TitansCrossBorderAPIAdapter) ProcesssOrderCancelRequestError(ordCxlReq *types.ApplicationOrderCancelRequest, tradeOrder *schema.TradeOrder, de *domain_error.Error) (rejectMsg []byte) {
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

func (a *TitansCrossBorderAPIAdapter) ErrorCouldBeIgnoreAfterReview(de *domain_error.Error) (ignoreAfterReview bool) {
	return
}

func (a *TitansCrossBorderAPIAdapter) ExtractAppOrderID(rawMsg []byte, msgProps map[string]interface{})(appOrdID string, de *domain_error.Error){
	return
}
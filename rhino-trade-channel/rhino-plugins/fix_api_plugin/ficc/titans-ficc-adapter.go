package ficc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"rhino-common/domain_error"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/timeutil"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/types"
	"rhino-plugins/fix_api_plugin/fixutil"
	"strconv"
	"strings"
	"time"

	attrenum "rhino-common/enum"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
)

type TitansFiccFixApiAdapter struct {
	applicationCfg *domain_cfg.ApplicationCfg
	configMap      map[string]*schema.ApplicationCfgItem
	configPath     string
	appSettings    *quickfix.Settings
	phoneNumMap    map[quickfix.SessionID]string
}

func NewTitansFiccFixApiAdapter(applicationCfg *domain_cfg.ApplicationCfg) (adapter *TitansFiccFixApiAdapter, de *domain_error.Error) {
	configMap := applicationCfg.GetApplicationCfgItemMap()
	cfgItem := configMap["PathFixServer"]
	if cfgItem == nil {
		de = domain_error.Build(domain_error.GENERIC_ERR_CODE, errors.New("PathFixServer not config"))
		return
	}
	configPath := cfgItem.ConfigItemValue
	_, err := os.Open(configPath)
	if err != nil {
		de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return
	}

	adapter = &TitansFiccFixApiAdapter{applicationCfg: applicationCfg, configMap: configMap, configPath: configPath}

	configFileContent, err := os.ReadFile(adapter.GetConfigPath())
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, fmt.Sprintf("fail to read fixserver config from %s\n", adapter.GetConfigPath()))
	}
	appSettings, err := quickfix.ParseSettings(bytes.NewReader(configFileContent))
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, fmt.Sprintf("fail to parse fixserver setting from %s\n", adapter.GetConfigPath()))
	}

	adapter.appSettings = appSettings
	adapter.phoneNumMap = map[quickfix.SessionID]string{}
	for sessionID, sessionSetting := range appSettings.SessionSettings() {
		phoneNum, _ := sessionSetting.Setting("PhoneNum")
		if phoneNum != "" {
			adapter.phoneNumMap[sessionID] = phoneNum
		}
	}

	return
}

func (a *TitansFiccFixApiAdapter) GetConfigPath() (configPath string) {
	return a.configPath
}

func (a *TitansFiccFixApiAdapter) LoginValidate(username, password string, sessionID quickfix.SessionID) (ok bool) {

	// 从配置取电话号码
	sessionSetting, ok := a.appSettings.SessionSettings()[sessionID]
	if !ok {
		log.Printf("session not found : %s\n", sessionID)
		return false
	}
	phoneNum, err := sessionSetting.Setting("PhoneNum")
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to get phone number")
		return false
	}
	loginServiceUrl, err := sessionSetting.Setting("LoginServiceUrl")
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to get loginServiceUrl")
		return false
	}
	loginClientID, err := sessionSetting.Setting("LoginClientID")
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to get loginClientID")
		return false
	}
	loginClientKey, err := sessionSetting.Setting("LoginClientKey")
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to get loginClientKey")
		return false
	}

	log.Printf("======> sessionID:%v, phoneNum:%s, loginServiceUrl:%s, loginClientID:%s, loginClientKey:%s\n", sessionID, phoneNum, loginServiceUrl, loginClientID, loginClientKey)

	ok, _, err = Authenticate(loginServiceUrl, loginClientID, loginClientKey, phoneNum, password)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to invoke Authenticate service")
		return false
	}

	return
}

func (a *TitansFiccFixApiAdapter) DecodeForNewOrderSingle(message *quickfix.Message, sessionID quickfix.SessionID) (msgProps map[string]interface{}, rejErr quickfix.MessageRejectError) {

	senderCompID, _ := message.Header.GetString(tag.TargetCompID)
	targetCompID, _ := message.Header.GetString(tag.SenderCompID)

	msgProps = map[string]interface{}{
		// Todo: ExtendAttrs增加senderCompID、targetCompID的定义
		"senderCompID":     senderCompID,
		"targetCompID":     targetCompID,
		"businessType":     "收益互换",
		"ordSource":        "FIX",
		"securityExchange": "8",
		"handlInst":        "1",
		"customer":         targetCompID,
		"execUser":         targetCompID,
		"phoneNum":         a.phoneNumMap[sessionID],
	}

	de := fixutil.SetValueByTag(msgProps, message, "clOrdID", tag.ClOrdID, message.Body.FieldMap, attrenum.AttrValueType_STRING)
	if rejErr = fixutil.ConvertDomainErrToFixRejection(de, tag.ClOrdID); rejErr != nil {
		return
	}
	de = fixutil.SetValueByTag(msgProps, message, "account", tag.Account, message.Body.FieldMap, attrenum.AttrValueType_INT)
	if rejErr = fixutil.ConvertDomainErrToFixRejection(de, tag.Account); rejErr != nil {
		return
	}
	de = fixutil.SetValueByTag(msgProps, message, "planCode", 20000, message.Body.FieldMap, attrenum.AttrValueType_STRING)
	if rejErr = fixutil.ConvertDomainErrToFixRejection(de, quickfix.Tag(20000)); rejErr != nil {
		return
	}
	de = fixutil.SetValueByTag(msgProps, message, "price", tag.Price, message.Body.FieldMap, attrenum.AttrValueType_FLOAT)
	if rejErr = fixutil.ConvertDomainErrToFixRejection(de, tag.Price); rejErr != nil {
		return
	}
	de = fixutil.SetValueByTag(msgProps, message, "quantity", tag.OrderQty, message.Body.FieldMap, attrenum.AttrValueType_INT)
	if rejErr = fixutil.ConvertDomainErrToFixRejection(de, tag.OrderQty); rejErr != nil {
		return
	}
	de = fixutil.SetValueByTag(msgProps, message, "settlType", 20001, message.Body.FieldMap, attrenum.AttrValueType_STRING)
	if rejErr = fixutil.ConvertDomainErrToFixRejection(de, quickfix.Tag(20001)); rejErr != nil {
		return
	}
	de = fixutil.SetValueByTag(msgProps, message, "side", tag.Side, message.Body.FieldMap, attrenum.AttrValueType_STRING)
	if rejErr = fixutil.ConvertDomainErrToFixRejection(de, tag.Side); rejErr != nil {
		return
	}
	de = fixutil.SetValueByTag(msgProps, message, "symbol", tag.Symbol, message.Body.FieldMap, attrenum.AttrValueType_STRING)
	if rejErr = fixutil.ConvertDomainErrToFixRejection(de, tag.Symbol); rejErr != nil {
		return
	}
	de = fixutil.SetValueByTag(msgProps, message, "transactTime", tag.TransactTime, message.Body.FieldMap, attrenum.AttrValueType_STRING)
	if rejErr = fixutil.ConvertDomainErrToFixRejection(de, tag.TransactTime); rejErr != nil {
		return
	}
	de = fixutil.SetValueByTag(msgProps, message, "ordType", tag.OrdType, message.Body.FieldMap, attrenum.AttrValueType_STRING)
	if rejErr = fixutil.ConvertDomainErrToFixRejection(de, tag.OrdType); rejErr != nil {
		return
	}

	handlInst, _ := message.Body.FieldMap.GetString(tag.HandlInst)
	// 兼容性设计，确保当handlInst不传参时也可以处理
	if len(handlInst) > 0 {
		de = fixutil.SetValueByTag(msgProps, message, "handlInst", tag.HandlInst, message.Body.FieldMap, attrenum.AttrValueType_STRING)
		if rejErr = fixutil.ConvertDomainErrToFixRejection(de, tag.HandlInst); rejErr != nil {
			return
		}
	}

	return
}

func (a *TitansFiccFixApiAdapter) AutoTurnToReviewForErrors(order *schema.TradeOrder, de *domain_error.Error) (turnToReview bool, displayErr *domain_error.Error) {
	if de.Code == domain_error.CAPITAL_AMOUNT_NOT_ENOUGH_ERR_CODE {

		if order.ExtendAttrMap != nil {
			order.ExtendAttrMap["remark"] = de.Msg
		}

		turnToReview = true
		displayErr = domain_error.Build(domain_error.CAPITAL_AMOUNT_NEED_REVIEW_ERR_CODE, nil)
		displayErr.Refine(domain_error.WARNING, order)
		return
	}
	if de.Code == domain_error.QUOTA_LIMIT_EXCEEDED_ERR_CODE {

		if order.ExtendAttrMap != nil {
			order.ExtendAttrMap["remark"] = de.Msg
		}
		
		turnToReview = true
		displayErr = domain_error.Build(domain_error.QUOTA_LIMIT_EXCEEDED_NEED_REVIEW_ERR_CODE, nil, "")
		displayErr.Msg = de.Msg + "，订单已转为待审批，请联系业务人员"
		displayErr.Refine(domain_error.WARNING, order)
		return
	}
	return
}

func (a *TitansFiccFixApiAdapter) DecodeForOrderCancelRequest(message *quickfix.Message, sessionID quickfix.SessionID) (applicationOrderCancelRequest *types.ApplicationOrderCancelRequest, rejErr quickfix.MessageRejectError) {

	account, de := fixutil.GetValueByTag(message, "account", tag.Account, message.Body.FieldMap, attrenum.AttrValueType_STRING)
	if rejErr = fixutil.ConvertDomainErrToFixRejection(de, tag.Account); rejErr != nil {
		return
	}
	origClOrdID, de := fixutil.GetValueByTag(message, "origClOrdID", tag.OrigClOrdID, message.Body.FieldMap, attrenum.AttrValueType_STRING)
	if rejErr = fixutil.ConvertDomainErrToFixRejection(de, tag.OrigClOrdID); rejErr != nil {
		return
	}
	transactTimeStr, de := fixutil.GetValueByTag(message, "transactTime", tag.TransactTime, message.Body.FieldMap, attrenum.AttrValueType_STRING)
	if rejErr = fixutil.ConvertDomainErrToFixRejection(de, tag.TransactTime); rejErr != nil {
		return
	}
	var transactTime int64
	if transactTimeStr != "" {

		t, err := timeutil.ParseTimeStrToTimeByTimeLocation(timeutil.TransactTimeLayout, transactTimeStr.(string), time.UTC) //20241107-09:46:10
		if err != nil {
			de = domain_error.Build(domain_error.ILLEGAL_FIT_TAG_ERR_CODE, rejErr, "TransactTime", tag.TransactTime)
			if rejErr = fixutil.ConvertDomainErrToFixRejection(de, tag.TransactTime); rejErr != nil {
				return
			}
		}
		if time.Since(t) > 60*time.Second {
			err = fmt.Errorf("expire transactTime:%s", transactTimeStr)
			de = domain_error.Build(domain_error.TX_TIME_TOO_EARLY_ERR_CODE, err, transactTimeStr)
			if rejErr = fixutil.ConvertDomainErrToFixRejection(de, tag.TransactTime); rejErr != nil {
				return
			}
		}
		transactTime = timeutil.ConvertTimeToMilliseconds(t)
	}

	applicationOrderCancelRequest = &types.ApplicationOrderCancelRequest{
		ActionUser: account.(string),
		AppOrdID:   origClOrdID.(string),
		ActionTime: transactTime,
	}

	return
}

func (a *TitansFiccFixApiAdapter) ConvertTradeResponseMessage(tradeResp *types.TradeActionRespReturn) (message *quickfix.Message, de *domain_error.Error) {

	resp := tradeResp.CurrentTradeActionResp
	msg := make(map[string]interface{})
	tradeOrder := tradeResp.GetTradeOrder()
	if len(tradeOrder.ExtendAttrMap) > 0 {
		msg = tradeOrder.ExtendAttrMap
	} else if len(tradeOrder.ExtendAttr) > 0 {
		err := json.Unmarshal([]byte(tradeOrder.ExtendAttr), &msg)
		if err != nil {
			de = domain_error.Build(domain_error.GENERIC_ERR_CODE, fmt.Errorf("ConvertTradeResponseMessage:: fail to parse tradeOrder.ExtendAttr:%s", tradeOrder.ExtendAttr))
			return
		}
	}
	if len(msg) == 0 {
		de = domain_error.Build(domain_error.GENERIC_ERR_CODE, fmt.Errorf("trade order ExtendAttrMap is empty"))
		return
	}

	if resp.CxlRejResponseTo == "" { // ExecutionReport
		senderCompID, _, _ := attrutil.GetAttrValue(msg, "senderCompID", attrenum.AttrValueType_STRING)
		if senderCompID == "" {
			return
		}
		targetCompID, _, _ := attrutil.GetAttrValue(msg, "targetCompID", attrenum.AttrValueType_STRING)
		if targetCompID == "" {
			return
		}
		message = quickfix.NewMessage()
		message.Header.SetString(tag.MsgType, string(enum.MsgType_EXECUTION_REPORT))
		message.Header.SetString(tag.BeginString, "FIX.4.2")
		message.Header.SetString(tag.SenderCompID, senderCompID.(string))
		message.Header.SetString(tag.TargetCompID, targetCompID.(string))
		message.Body.SetString(tag.OrdType, tradeOrder.OrdType)
		message.Body.SetString(tag.SecurityExchange, tradeOrder.SecurityExchange)
		message.Body.SetString(tag.HandlInst, tradeOrder.HandlInst)
		message.Body.SetString(tag.ClOrdID, tradeOrder.AppOrdID)
		message.Body.SetString(tag.OrderID, tradeOrder.ClOrdID)
		message.Body.SetString(tag.ExecID, resp.ExecID)

		// 20251217改，直接取resp的ExecType
		//execType := resp.OrdStatus
		//message.Body.SetString(tag.ExecType, execType)
		execType := resp.ExecType
		message.Body.SetString(tag.ExecType, execType)

		// 20251217改，直接取resp的ExecTransType
		// if execType == "6" || execType == "4" {
		// 	message.Body.SetString(tag.ExecTransType, "1")
		// } else {
		// 	message.Body.SetString(tag.ExecTransType, "0")
		// }
		message.Body.SetString(tag.ExecTransType, resp.ExecTransType)

		// 20251217改，直接取resp的OrdStatus
		// if resp.OrdStatus == "1" && resp.PendingCancel {
		// 	message.Body.SetString(tag.OrdStatus, "6")
		// } else {
		// 	message.Body.SetString(tag.OrdStatus, resp.OrdStatus)
		// }
		message.Body.SetString(tag.OrdStatus, resp.OrdStatus)

		if resp.OrdRejReason != "" {
			message.Body.SetString(tag.Text, resp.OrdRejReason)
		}
		message.Body.SetString(tag.Account, tradeOrder.Account)
		message.Body.SetString(tag.Symbol, tradeOrder.Symbol)
		message.Body.SetString(tag.Side, tradeOrder.Side)
		message.Body.SetInt(tag.OrderQty, int(tradeOrder.OrderQty))
		price := strconv.FormatFloat(tradeOrder.Price, 'f', -1, 64)
		message.Body.SetString(tag.Price, price)
		if resp.LastShares > 0 {
			message.Body.SetInt(tag.LastQty, int(resp.LastShares))

			// “成交收益率 20002”、“成交全价 20003”分别取订单的意向收益率和意向全价，“费后全价 20004”则根据订单实际处理效率取相应档位的佣金率实时计算
			//quickfix.Tag(20001)
			respPrice, _, _ := attrutil.GetAttrValue(resp.ExtendAttrMap, "respPrice", attrenum.AttrValueType_FLOAT)
			respDirtyPrice, _, _ := attrutil.GetAttrValue(resp.ExtendAttrMap, "respDirtyPrice", attrenum.AttrValueType_FLOAT)
			respDirtyPriceWithFee, _, _ := attrutil.GetAttrValue(resp.ExtendAttrMap, "respDirtyPriceWithFee", attrenum.AttrValueType_FLOAT)
			respYtm, _, _ := attrutil.GetAttrValue(resp.ExtendAttrMap, "respYtm", attrenum.AttrValueType_FLOAT)

			message.Body.SetString(tag.LastPx, strconv.FormatFloat(respPrice.(float64), 'f', -1, 64))
			message.Body.SetString(quickfix.Tag(20003), strconv.FormatFloat(respDirtyPrice.(float64), 'f', -1, 64))
			message.Body.SetString(quickfix.Tag(20004), strconv.FormatFloat(respDirtyPriceWithFee.(float64), 'f', -1, 64))
			message.Body.SetString(quickfix.Tag(20002), strconv.FormatFloat(respYtm.(float64), 'f', -1, 64))
		}
		message.Body.SetInt(tag.LeavesQty, int(resp.LeavesQty))
		if resp.CumQty > 0 {
			message.Body.SetInt(tag.CumQty, int(resp.CumQty))
			
			avgPrice, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "avgPrice", attrenum.AttrValueType_FLOAT)
			message.Body.SetString(tag.AvgPx, strconv.FormatFloat(avgPrice.(float64), 'f', -1, 64))
			
		} else {
			message.Body.SetInt(tag.CumQty, 0)
			message.Body.SetInt(tag.AvgPx, 0)
		}
		message.Body.SetString(tag.TransactTime, timeutil.ConvertMillisecondsToTime(resp.TransactTime).In(time.UTC).Format(timeutil.TransactTimeLayout))
	} else { // OrderCancelReject
		//tradeAction := tradeResp.GetTradeActionLatestResp()
		senderCompID, _, _ := attrutil.GetAttrValue(msg, "senderCompID", attrenum.AttrValueType_STRING)
		if senderCompID == "" {
			return
		}
		targetCompID, _, _ := attrutil.GetAttrValue(msg, "targetCompID", attrenum.AttrValueType_STRING)
		if targetCompID == "" {
			return
		}
		message = quickfix.NewMessage()
		message.Header.SetString(tag.BeginString, "FIX.4.2")
		message.Header.SetString(tag.MsgType, string(enum.MsgType_ORDER_CANCEL_REJECT))
		message.Header.SetString(tag.SenderCompID, senderCompID.(string))
		message.Header.SetString(tag.TargetCompID, targetCompID.(string))
		// if tradeAction != nil {
		// 	message.Body.SetString(tag.ClOrdID, tradeAction.ClOrdID)
		// }
		message.Body.SetString(tag.OrigClOrdID, tradeOrder.AppOrdID)
		message.Body.SetString(tag.OrderID, tradeOrder.ClOrdID)
		message.Body.SetString(tag.OrdStatus, tradeOrder.OrdStatus)
		message.Body.SetString(tag.CxlRejResponseTo, resp.CxlRejResponseTo)
		message.Body.SetString(tag.Text, resp.OrdRejReason)
		message.Body.SetString(tag.Account, tradeOrder.Account)
		message.Body.SetString(tag.TransactTime, timeutil.ConvertMillisecondsToTime(resp.TransactTime).In(time.UTC).Format(timeutil.TransactTimeLayout))
	}

	return
}

func (a *TitansFiccFixApiAdapter) GetFixPortOpenTimeRange() (begin string, end string, layout string, local *time.Location) {

	strs := strings.Split(a.configMap["TradingTime"].ConfigItemValue, "-")
	local = timeutil.CnTimeLocation
	dateStr := time.Now().In(local).Format(time.DateOnly) + " "
	beginTime, err := time.ParseInLocation(time.DateOnly+" 15:04", dateStr+strs[0], local)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "error in GetFixPortOpenTimeRange")
	}
	endTime, err := time.ParseInLocation(time.DateOnly+" 15:04", dateStr+strs[1], local)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "error in GetFixPortOpenTimeRange")
	}

	beginTime = beginTime.Add(-1 * time.Hour)
	endTime = endTime.Add(time.Hour)

	layout = "15:04"
	begin = beginTime.In(local).Format(layout)
	end = endTime.In(local).Format(layout)

	return
}

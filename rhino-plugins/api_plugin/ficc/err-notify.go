package ficc

import (
	"bytes"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/domain_error/notify_provider"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	"rhino-plugins/api_plugin/util"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// 企业微信消息结构体
type WechatMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content       string   `json:"content"`
		MentionedList []string `json:"mentioned_list"`
		//MentionedMobileList []string `json:"mentioned_mobile_list"`
	} `json:"text"`
}

func (a *TitansFiccAPIAdapter) initErrNotifier(configMap map[string]*schema.ApplicationCfgItem) {
	var err error
	a.errNotifier, err = notify_provider.NewWechatErrorNotifyProvider(configMap["WechatServiceUrl"].ConfigItemValue, configMap["DomainErrLogPath"].ConfigItemValue, time.NewTicker(24*time.Hour), time.NewTicker(5*time.Second), time.Sunday, func(logLine string) (interface{}, error) {

		if len(logLine) == 0 {
			return nil, nil
		}

		de := &domain_error.Error{}

		err := json.Unmarshal([]byte(logLine), de)
		if err != nil {
			return nil, err
		}

		if de.Order == nil || len(de.Order.ExtendAttrMap) == 0 {
			return nil, nil
		}

		extendAttrMap := de.Order.ExtendAttrMap

		rawCounterparty, ok := extendAttrMap["account"]
		if !ok {
			rawCounterparty = "-"
		}
		rawCustomer := "-"
		counterpartyID, ok, _ := attrutil.GetAttrValue(extendAttrMap, "account", enum.AttrValueType_INT)
		if ok {
			valList, ok, _ := a.autoSyncRepo.Get("Counterparty", strconv.Itoa(counterpartyID.(int)))
			if ok && len(valList) > 0 {
				ctpyShortName, ok, _ := attrutil.GetAttrValue(valList[0], "Counterparty", enum.AttrValueType_STRING)
				if ok && len(ctpyShortName.(string)) > 0 {
					rawCounterparty = ctpyShortName.(string)
				}

				companyName, ok, _ := attrutil.GetAttrValue(valList[0], "CounterpartyCompanyName", enum.AttrValueType_STRING)
				if ok && len(companyName.(string)) > 0 {
					rawCustomer = companyName.(string)
				}
			}
		}

		var orderTime string
		var realTimeStamp int64
		if de.Order.TransactTime > 0 {
			realTimeStamp = de.Order.TransactTime
			orderTime = timeutil.ConvertMillisecondsToTime(de.Order.TransactTime).In(timeutil.CnTimeLocation).Format(timeutil.TransactTimeLayout)
		}
		if orderTime == "" {
			transactTime, de := getTimeStamp(nil, extendAttrMap, "transactTime")
			if de == nil && transactTime > 0 {
				realTimeStamp = transactTime
				orderTime = timeutil.ConvertMillisecondsToTime(transactTime).In(timeutil.CnTimeLocation).Format(timeutil.TransactTimeLayout)
			}
		}

		if realTimeStamp > 0 {
			orderTimeStamp := timeutil.ConvertMillisecondsToTime(realTimeStamp)
			// 如果是撤单请求，需要取domainErr的timestamp
			if strings.Contains(de.Msg, "撤单请求待处理") && len(de.Timestamp) > 0 {
				t, err := time.ParseInLocation(timeutil.TransactTimeLayout, de.Timestamp, timeutil.CnTimeLocation)
				if err == nil {
					orderTimeStamp = t
				} else {
					log.Printf("fail to ParseInLocation for %s, error:%v\n", de.Timestamp, err)
				}
			}
			if time.Since(orderTimeStamp) > 3*time.Hour {
				log.Println("ignore expired alert message")
				return nil, nil
			}
		}

		if orderTime == "" {
			orderTime = de.Timestamp
		}

		appOrdID, ok := extendAttrMap["clOrdID"]
		if !ok {
			appOrdID = "-"
		}
		errMsg := de.Msg

		orderSource, ok := extendAttrMap["ordSource"]
		if !ok {
			orderSource = "-"
		}
		switch orderSource {
		case "goats":
			orderSource = "GOATS客户端"
		case "titans":
			orderSource = "TITANS"
		}

		tradeSide := "-"
		side, ok, _ := util.GetStringValueInField(extendAttrMap, "side")
		if ok {
			if side == "1" {
				tradeSide = "bid"
			} else {
				tradeSide = "ofr"
			}
		}

		symbol, ok := extendAttrMap["symbol"]
		if !ok {
			symbol = "-"
		} else {
			_symbol, ok := symbol.(string)
			if ok {
				idx := strings.Index(symbol.(string), ".")
				if idx > 0 {
					symbol = _symbol[:idx]
				}
			}
		}

		settleType, ok := extendAttrMap["settlType"]
		if !ok {
			settleType = "-"
		} else {
			_settleType, ok := settleType.(string)
			if ok {
				idx := strings.Index(_settleType, "+")
				if idx > 0 {
					settleType = _settleType[idx:]
				}
			}
		}

		orderQty, ok := extendAttrMap["quantity"]
		if !ok {
			orderQty = "-"
		} else {
			orderQty, _, _ = util.GetFloatValueInField(extendAttrMap, "quantity")
			orderQty = int(orderQty.(float64) / 10000)
		}

		price, ok := extendAttrMap["price"]
		if !ok {
			price = "-"
		}

		dirtyPrice, ok := extendAttrMap["dirtyPrice"]
		if !ok {
			dirtyPrice = "-"
		}

		ytm, ok := extendAttrMap["ytm"]
		if !ok {
			ytm = "-"
		}

		handlInst, ok := extendAttrMap["handlInst"]
		if !ok {
			handlInst = "-"
		} else {
			if handlInst == "1" {
				handlInst = "快速交易"
			} else if handlInst == "3" {
				handlInst = "普通交易"
			}
		}

		orderCancel := ""
		if strings.Contains(errMsg, "撤单请求待处理") {
			orderCancel = " ref"
		}

		var buf *bytes.Buffer
		if de.Level == domain_error.ERROR {
			buf = bytes.NewBufferString("【债券互换拒单提醒】\n")
		} else if de.Level == domain_error.WARNING {
			buf = bytes.NewBufferString("【债券互换业务提醒】\n")
		}

		buf.WriteString(fmt.Sprintf("下单时间：%v\n", orderTime))
		if de.Level == domain_error.ERROR {
			buf.WriteString(fmt.Sprintf("拒单原因：%v\n", errMsg))
		} else if de.Level == domain_error.WARNING {
			buf.WriteString(fmt.Sprintf("业务提示：%v\n", errMsg))
		}
		buf.WriteString(fmt.Sprintf("订单编号：%v\n", appOrdID))
		buf.WriteString(fmt.Sprintf("下单渠道：%v\n", orderSource))
		buf.WriteString(fmt.Sprintf("交易对手：%v\n", rawCounterparty))
		buf.WriteString(fmt.Sprintf("客户名称：%v\n", rawCustomer))
		// buf.WriteString("合约方向：多头\n")
		//buf.WriteString(fmt.Sprintf("交易方向：%v\n", tradeSide))
		//buf.WriteString(fmt.Sprintf("标的代码：%v\n", symbol))
		buf.WriteString(fmt.Sprintf("交易效率：%v\n", handlInst))
		//buf.WriteString(fmt.Sprintf("清算速度：%v\n", settleType))
		//buf.WriteString(fmt.Sprintf("券面总额(万元)：%v\n", orderQty))
		buf.WriteString(fmt.Sprintf("交易要素：%v %v %v %v %v%v\n", tradeSide, symbol, orderQty, ytm, settleType, orderCancel))
		buf.WriteString(fmt.Sprintf("意向净价(元)：%v\n", price))
		buf.WriteString(fmt.Sprintf("意向全价(元)：%v\n", dirtyPrice))
		// if ytm == "-" {
		// 	buf.WriteString(fmt.Sprintf("意向到期收益率：%v\n", ytm))
		// } else {
		// 	buf.WriteString(fmt.Sprintf("意向到期收益率：%v%%\n", ytm))
		// }

		// 构造消息
		message := &WechatMessage{
			MsgType: "text",
		}
		message.Text.Content = buf.String()
		message.Text.MentionedList = []string{"@all"}
		//message.Text.MentionedMobileList = []string{"@all"}
		return message, nil
	})
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to NewWechatErrorNotifyProvider")
	}

	// 绑定错误提示函数
	domain_error.ErrorNotifyFunction = a.errNotifier.NotifyError
}

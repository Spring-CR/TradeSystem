package citic_qfii

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
	"rhino-plugins/channel_plugin/fix"
	"rhino-trade-channel/adapter/data_convert/domain_order_validate"
	"rhino-trade-channel/adapter/data_convert/fix_convert"
	"strconv"
	"time"

	newordersingle42 "github.com/quickfixgo/fix42/newordersingle"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
)

const (
	datetimeLayout = "20060102-15:04:05"
)

type CiticFIXChannel struct {
	*fix.GenericFIXChannel
	useCiticUat bool
	quotesUrl   string
}

func NewCiticFIXChannel(cfg *domain_cfg.TradeChannelCfg, onTradeActionResp func(tradeActionResp *schema.TradeActionResp)) (channel *CiticFIXChannel, de *domain_error.Error) {

	// genericChannel, de1 := fix.NewGenericFIXChannel(cfg)
	// if de1 != nil {
	// 	return nil, de1
	// }

	// channel = &CiticFIXChannel{genericChannel}

	channel = &CiticFIXChannel{}

	genericChannel, de1 := fix.NewGenericFIXChannel(cfg, onTradeActionResp, channel.CleanTradeActionResp, channel.RefineOrderCancelID)
	if de1 != nil {
		return nil, de1
	}

	channel.GenericFIXChannel = genericChannel

	for _, item := range cfg.GetTradeChannelCfgItems() {
		if item.ConfigItemName == "UseCiticUat" && item.ConfigItemValue == "1" {
			channel.useCiticUat = true
			log.Println("UseCiticUat!")
			break
		}
	}
	for _, item := range cfg.GetTradeChannelCfgItems() {
		if item.ConfigItemName == "QuotesUrl" {
			channel.quotesUrl = item.ConfigItemValue
			log.Printf("QuotesUrl:%s", channel.quotesUrl)
			break
		}
	}

	return
}

func (c *CiticFIXChannel) ValidateOrder(order *schema.TradeOrder) (channelOrder interface{}, de *domain_error.Error) {

	var rawOrder newordersingle42.NewOrderSingle
	var msg *quickfix.Message

	rawOrder = fix_convert.DomainOrderToRawOrderFix42(order)

	// 检查业务类型是否为北上期权
	if order.ExtendAttrMap != nil && order.ExtendAttrMap["businessType"] != "N_CROSS_OPTION" {
		de = domain_error.Build(domain_error.TRADE_CHANNEL_ORDER_VALIDATE_ERR_CODE, nil, "当标的为A股、柜台为中信QFII时，业务类型应记录为：北上期权")
		return
	}

	// 参考wiki：http://wiki.gf.com.cn/pages/viewpage.action?pageId=241660445&preview=/241660445/243410398/Citics%20Algo%20Trading%20Service%EF%BC%88Thor%EF%BC%89FIX%20Specs%20_RootNet%20QFII%20DSA_20230731.xlsx
	// http://wiki.gf.com.cn/pages/viewpage.action?pageId=243402961
	// Todo wiki里可能漏了北交所
	// 设置交易所
	isShangHai := false
	switch order.SymbolSfx {
	case "SZ":
		rawOrder.SetString(tag.ExDestination, "SZ")
	case "SH":
		rawOrder.SetString(tag.ExDestination, "SS")
		isShangHai = true
	default:
		de = domain_error.Build(domain_error.DATA_CONVERT_UNSUPPORT_EXCHANGE_ERR_CODE, nil, order.SecurityExchange)
		return
		// case "BJ":
		// 	rawOrder.SetString(tag.ExDestination, "BJ") // Todo wiki里可能漏了北交所，经过确认，目前不支持北交所
	}

	// 对市价委托进行特殊设置，完全参考北向trademodule的现有逻辑
	if order.OrdType == string(enum.OrdType_MARKET) {
		if order.AlgName == "" { // 非算法
			rawOrder.SetString(tag.TimeInForce, "3")
			rawOrder.SetString(quickfix.Tag(6183), "5")
		}
		if isShangHai && order.Price == 0 {
			//rawOrder.SetStopPx(decimal.NewFromFloat(order.Price), 4)
			protectPrice, err := calcShProtectPrice(c.quotesUrl, order.Symbol2, order.Side)
			if err != nil {
				de = domain_error.Build(domain_error.DATA_CONVERT_GET_PROECTED_PRICE_ERR_CODE, err, order.Symbol2, order.Side)
				return
			}
			rawOrder.SetString(quickfix.Tag(99), fmt.Sprintf("%.4f", float64(protectPrice)/100.0))
		}
	}

	// 设置算法参数
	var algParams map[string]interface{}
	var err error
	if order.AlgName != "" {
		if order.AlgParams == "" {
			de = domain_error.Build(domain_error.DATA_CONVERT_ALG_PARAMS_NOT_PROVIDED_ERR_CODE, nil, order.AlgName)
			return
		}

		// 在接口ConvertNewOrderSingleMessage中，已经设置了ExtendAttr和AlgParam的map形式，可以减少重复的计算
		if len(order.AlgParamsMap) > 0 {
			algParams = order.AlgParamsMap
		} else {
			algParams, err = attrutil.ParseAttrMapString(order.AlgParams)
			if err != nil {
				de = domain_error.Build(domain_error.DATA_CONVERT_ALG_PARAMS_EXTRACT_ERR_CODE, err, order.AlgName, order.AlgParams)
				return
			}
		}

		rawOrder.SetString(quickfix.Tag(6065), "3")

		switch order.AlgName {
		case "POV":
			rawOrder.SetString(6061, "VOLINLINE")
			de = setAlgTimeParam(rawOrder, 6062, algParams, "startTime", timeutil.CnTimeLocation)
			if de != nil {
				return
			}
			de = setAlgTimeParam(rawOrder, 6063, algParams, "endTime", timeutil.CnTimeLocation)
			if de != nil {
				return
			}
			de = setAlgRatioParam(rawOrder, 6064, algParams, "povRatio")
			if de != nil {
				return
			}
		case "TWAP":
			rawOrder.SetString(6061, "TWAP")
			de = setAlgTimeParam(rawOrder, 6062, algParams, "startTime", timeutil.CnTimeLocation)
			if de != nil {
				return
			}
			de = setAlgTimeParam(rawOrder, 6063, algParams, "endTime", timeutil.CnTimeLocation)
			if de != nil {
				return
			}
			rawOrder.SetString(quickfix.Tag(6064), "0")
		case "VWAP":
			rawOrder.SetString(6061, "VWAP")
			de = setAlgTimeParam(rawOrder, 6062, algParams, "startTime", timeutil.CnTimeLocation)
			if de != nil {
				return
			}
			de = setAlgTimeParam(rawOrder, 6063, algParams, "endTime", timeutil.CnTimeLocation)
			if de != nil {
				return
			}
			rawOrder.SetString(quickfix.Tag(6064), "0")
		}
	}

	// 使用无后缀的symbol
	rawOrder.SetSymbol(order.Symbol2)

	// DSA+限价且价格不输入时，后端将OrdType转为市价单报盘（已与产品沟通）
	if order.Price <= 0 && order.OrdType == string(enum.OrdType_LIMIT) {
		rawOrder.SetString(tag.OrdType, string(enum.OrdType_MARKET))
	}

	// 如果是中信uat环境转换一下账号
	if c.useCiticUat {
		rawOrder.SetAccount("322")
	}

	msg = rawOrder.ToMessage()

	c.QueryHeader(&msg.Header)

	return msg, nil
}

func (c *CiticFIXChannel) NewOrderSingle(order *schema.TradeOrder, afterTradeOrderUpsert func(tradeOrder *schema.TradeOrder), syncTradeOrderAfterError func(tradeOrder *schema.TradeOrder)) (duplicatedOrder bool, de *domain_error.Error) {

	// 先检查session是否就绪
	if !c.IsSessionReady() {
		return duplicatedOrder, domain_error.Build(domain_error.FIX_SESSION_NOT_READY_ERR_CODE, nil)
	}
	var channelOrder interface{}
	duplicatedOrder, channelOrder, de = domain_order_validate.CheckAndBaseSetDomainOrder(order, c.GetChannelConfig(), c.ValidateOrder, afterTradeOrderUpsert)
	if de != nil {
		return duplicatedOrder, de
	}
	// do nothing!
	if duplicatedOrder {
		return duplicatedOrder, nil
	}
	msg := channelOrder.(*quickfix.Message)
	// 更新ClOrdID
	msg.Body.SetString(quickfix.Tag(11), order.ClOrdID)
	err := quickfix.Send(msg)
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

// func setAlgParam(rawOrder newordersingle42.NewOrderSingle, tag int, algParams map[string]interface{}, field string) *domain_error.Error {
// 	val, ok, err := attrutil.GetAttrValue(algParams, field, enum.AttrValueType_STRING)
// 	if err != nil {
// 		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
// 		return de
// 	}
// 	if ok {
// 		rawOrder.SetString(quickfix.Tag(tag), val.(string))
// 	}
// 	return nil
// }

func setAlgTimeParam(rawOrder newordersingle42.NewOrderSingle, tag int, algParams map[string]interface{}, field string, timeLoc *time.Location) *domain_error.Error {
	val, ok, err := attrutil.GetAttrValue(algParams, field, enum.AttrValueType_STRING)
	if err != nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return de
	}
	// 目标格式：20060609-09:30:00, GMT
	if ok && val != "" {
		strVal := val.(string)
		// 输入格式和目标格式不一样
		t, err := timeutil.ParseTimeStrToTimeByTimeLocation(time.DateTime, strVal, timeLoc)
		if err != nil {
			de := domain_error.Build(domain_error.DATA_CONVERT_ALG_PARAMS_EXTRACT_TIME_FIELD_ERR_CODE, err, strVal)
			return de
		}
		// 转为GMT时间
		strVal = t.In(time.UTC).Format(datetimeLayout)
		rawOrder.SetString(quickfix.Tag(tag), strVal)
	}
	return nil
}

func setAlgRatioParam(rawOrder newordersingle42.NewOrderSingle, tag int, algParams map[string]interface{}, field string) *domain_error.Error {
	val, ok, err := attrutil.GetAttrValue(algParams, field, enum.AttrValueType_FLOAT)
	if err != nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return de
	}
	if ok {
		floatVal := val.(float64)
		// 根据中信的协议要转成百分比
		rawOrder.SetString(quickfix.Tag(tag), strconv.Itoa(int(floatVal*100)))
	}
	return nil
}

func (c *CiticFIXChannel) CleanTradeActionResp(tradeActionResp *schema.TradeActionResp) (ignore bool) {
	return
}

func (c *CiticFIXChannel) RefineOrderCancelID(rawTargetClOrdID string, rawClOrdID string, order *schema.TradeOrder) (targetClOrdID, clOrdID string) {
	// 不改
	return rawTargetClOrdID, rawClOrdID
}

func (c *CiticFIXChannel) OnMarketClosed(orderCache *order_cache.OrderCache) {
	
}
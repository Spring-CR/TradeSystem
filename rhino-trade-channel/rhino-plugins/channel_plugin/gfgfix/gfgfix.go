package gfgfix

import (
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-core/domain_cfg"
	"rhino-core/order_domain/order_cache"
	"rhino-core/schema"
	"rhino-core/store/app_store"
	"rhino-plugins/channel_plugin/fix"
	"rhino-trade-channel/adapter/data_convert/domain_order_validate"
	"rhino-trade-channel/adapter/data_convert/fix_convert"
	"strings"

	newordersingle42 "github.com/quickfixgo/fix42/newordersingle"
	"github.com/quickfixgo/quickfix"
	"github.com/shopspring/decimal"
)

type GFGFIXChannel struct {
	*fix.GenericFIXChannel
	useHkUat bool
}

func NewGFGFIXChannel(cfg *domain_cfg.TradeChannelCfg, onTradeActionResp func(tradeActionResp *schema.TradeActionResp)) (channel *GFGFIXChannel, de *domain_error.Error) {

	channel = &GFGFIXChannel{}

	genericChannel, de1 := fix.NewGenericFIXChannel(cfg, onTradeActionResp, channel.CleanTradeActionResp, channel.RefineOrderCancelID)
	if de1 != nil {
		return nil, de1
	}

	channel.GenericFIXChannel = genericChannel

	for _, item := range cfg.GetTradeChannelCfgItems() {
		if item.ConfigItemName == "UseHkUat" && item.ConfigItemValue == "1" {
			channel.useHkUat = true
			log.Println("UseHkUat!")
			break
		}
	}

	return
}

var trimingPrefixes = []string{
	"GFZQ-HK2GO-", // 较长前缀优先
	"HKHEDGE-",
}

func (c *GFGFIXChannel) trimClOrdIDPrefix(s string) string {
	// Todo：如果gfgfix改成128长度限制，以下这两行代码就可以去掉
	s = strings.Replace(s, "-t-", "-titans-", 1)
	s = strings.Replace(s, "-c-", "-crossborder-", 1)
	for _, prefix := range trimingPrefixes {
		if strings.HasPrefix(s, prefix) {
			return s[len(prefix):]
		}
	}
	return s
}

func (c *GFGFIXChannel) CleanTradeActionResp(tradeActionResp *schema.TradeActionResp) (ignore bool) {
	log.Printf("start to CleanTradeActionResp, old clOrdID:%s\n", tradeActionResp.ClOrdID)
	tradeActionResp.ClOrdID = c.trimClOrdIDPrefix(tradeActionResp.ClOrdID)
	if tradeActionResp.OrigClOrdID != "" {
		tradeActionResp.OrigClOrdID = c.trimClOrdIDPrefix(tradeActionResp.OrigClOrdID)
	}
	return
}

func (c *GFGFIXChannel) RefineOrderCancelID(rawTargetClOrdID string, rawClOrdID string, order *schema.TradeOrder) (targetClOrdID, clOrdID string) {

	businessSource, err := getBusinessSource(order.ExtendAttrMap)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to RefineOrderCancelID")
		return rawTargetClOrdID, rawClOrdID
	}

	targetClOrdID = strings.Replace(rawTargetClOrdID, "-titans-", "-t-", 1)
	targetClOrdID = strings.Replace(targetClOrdID, "-crossborder-", "-c-", 1)

	clOrdID = strings.Replace(rawClOrdID, "-titans-", "-t-", 1)
	clOrdID = strings.Replace(clOrdID, "-crossborder-", "-c-", 1)

	return businessSource + "-" + targetClOrdID, businessSource + "-" + clOrdID
}

func (c *GFGFIXChannel) ValidateOrder(order *schema.TradeOrder) (channelOrder interface{}, de *domain_error.Error) {

	var rawOrder newordersingle42.NewOrderSingle
	var msg *quickfix.Message

	rawOrder = fix_convert.DomainOrderToRawOrderFix42(order)

	if strings.HasPrefix(order.SecurityExchange, "HK") {
		rawOrder.SetSecurityExchange("HK")
		order.SecurityExchangeRegion="HK"
	} else {
		switch order.SecurityExchange {
		case "TSE": // 日股
			rawOrder.SetSecurityExchange("TSE")
		default: // 美股
			rawOrder.SetSecurityExchange("US")
		}
	}

	// gfgfix根据特定的逻辑来设置TEXT Tag 58
	text := ""

	var extendAttrMap map[string]interface{}
	var algParamMap map[string]interface{}
	var err error

	if len(order.ExtendAttrMap) > 0 {
		extendAttrMap = order.ExtendAttrMap
	} else {
		log.Printf("here will cost much time, attrutil.ParseAttrMapString(order.ExtendAttr)")
		extendAttrMap, err = attrutil.ParseAttrMapString(order.ExtendAttr)
		if err != nil {
			de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
			return
		}
		order.ExtendAttrMap = extendAttrMap
	}

	if len(order.AlgParamsMap) > 0 {
		algParamMap = order.AlgParamsMap
	} else {
		algParamMap, err = attrutil.ParseAttrMapString(order.AlgParams)
		if err != nil {
			de = domain_error.Build(domain_error.DATA_CONVERT_ALG_PARAMS_EXTRACT_ERR_CODE, err, order.AlgName, order.AlgParams)
			return
		}
	}

	// // 兼容gfgfix关于tag11前缀的需求
	// businessSource, err := getBusinessSource(extendAttrMap)
	// if err != nil {
	// 	de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
	// 	return
	// }
	// rawOrder.SetClOrdID(businessSource + "-" + order.ClOrdID)

	tradeMode, err := getTradeMode(order, extendAttrMap)
	if err != nil {
		de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return
	}

	clientName, err := getClientShortName(extendAttrMap)
	if err != nil {
		de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return
	}

	// 匹配香港uat测试环境，强制增加SSCSKING才能通过
	if clientName == "" && c.useHkUat {
		clientName = "SSCSKING"
	}

	if order.OrderQty > 0 {
		rawOrder.SetOrderQty(decimal.NewFromFloat(order.OrderQty), 0)
	}

	var algStartTime, algEndTime, percent, maxVol, displayQty, displayQtyUnit, displayQtyVariance interface{}
	if order.AlgName != "" {
		algStartTime, err = getAlgStartTime(algParamMap)
		if err != nil {
			de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
			return
		}
		algEndTime, err = getAlgEndTime(algParamMap)
		if err != nil {
			de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
			return
		}
		percent, de = getPercent(algParamMap)
		if de != nil {
			return
		}
		maxVol, de = getMaxVol(algParamMap)
		if de != nil {
			return
		}
		displayQty, err = getDisplayQty(algParamMap)
		if err != nil {
			de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
			return
		}
		displayQtyUnit, err = getDisplayQtyUnit(algParamMap)
		if err != nil {
			de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
			return
		}
		displayQtyVariance, err = getDisplayQtyVariance(algParamMap)
		if err != nil {
			de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
			return
		}
	}

	if order.OrdType == string(enum.OrdType_MARKET) {
		switch order.AlgName {
		case "": //非算法市价单
			text = fmt.Sprintf("TRADE_MODE:%s CLIENT_SHORTNAME:%s ORDER_TYPE:%s QUANTITY:%.2f",
				tradeMode, clientName, "STRAIGHT_MARKET_ORDER", order.OrderQty)
		case "POV": //POV
			text = fmt.Sprintf("TRADE_MODE:%s CLIENT_SHORTNAME:%s ORDER_TYPE:%s QUANTITY:%.2f START_DATETIME:%s END_DATETIME:%s PERCENT:%.4f%%",
				tradeMode, clientName, order.AlgName, order.OrderQty, algStartTime, algEndTime, percent.(float64)*100.0)
		case "VWAP": //VWAP
			text = fmt.Sprintf("TRADE_MODE:%s CLIENT_SHORTNAME:%s ORDER_TYPE:%s QUANTITY:%.2f START_DATETIME:%s END_DATETIME:%s MAX_VOL:%.4f%%",
				tradeMode, clientName, order.AlgName, order.OrderQty, algStartTime, algEndTime, maxVol.(float64)*100.0)
		case "TWAP": //TWAP
			text = fmt.Sprintf("TRADE_MODE:%s CLIENT_SHORTNAME:%s ORDER_TYPE:%s QUANTITY:%.2f START_DATETIME:%s END_DATETIME:%s MAX_VOL:%.4f%%",
				tradeMode, clientName, order.AlgName, order.OrderQty, algStartTime, algEndTime, maxVol.(float64)*100.0)
		}
	} else if order.OrdType == string(enum.OrdType_LIMIT) {
		switch order.AlgName {
		case "": //非算法限价单
			text = fmt.Sprintf("TRADE_MODE:%s CLIENT_SHORTNAME:%s ORDER_TYPE:%s PRICE:%.4f QUANTITY:%.2f",
				tradeMode, clientName, "STRAIGHT_LIMIT_ORDER", order.Price, order.OrderQty)
		case "POV": //POV
			text = fmt.Sprintf("TRADE_MODE:%s CLIENT_SHORTNAME:%s ORDER_TYPE:%s PRICE:%.4f QUANTITY:%.2f START_DATETIME:%s END_DATETIME:%s PERCENT:%.4f%%",
				tradeMode, clientName, order.AlgName, order.Price, order.OrderQty, algStartTime, algEndTime, percent.(float64)*100.0)
		case "VWAP": //VWAP
			text = fmt.Sprintf("TRADE_MODE:%s CLIENT_SHORTNAME:%s ORDER_TYPE:%s PRICE:%.4f QUANTITY:%.2f START_DATETIME:%s END_DATETIME:%s MAX_VOL:%.4f%%",
				tradeMode, clientName, order.AlgName, order.Price, order.OrderQty, algStartTime, algEndTime, maxVol.(float64)*100.0)
		case "TWAP": //TWAP
			text = fmt.Sprintf("TRADE_MODE:%s CLIENT_SHORTNAME:%s ORDER_TYPE:%s PRICE:%.4f QUANTITY:%.2f START_DATETIME:%s END_DATETIME:%s MAX_VOL:%.4f%%",
				tradeMode, clientName, order.AlgName, order.Price, order.OrderQty, algStartTime, algEndTime, maxVol.(float64)*100.0)
		case "SNIPER": //SNIPER
			text = fmt.Sprintf("TRADE_MODE:%s CLIENT_SHORTNAME:%s ORDER_TYPE:%s PRICE:%.4f QUANTITY:%.2f START_DATETIME:%s END_DATETIME:%s",
				tradeMode, clientName, order.AlgName, order.Price, order.OrderQty, algStartTime, algEndTime)
		case "ICEBERG": //ICEBERG
			if strings.Contains(strings.ToUpper(order.SecurityExchangeRegion), "HK") { //港股单
				text = fmt.Sprintf("TRADE_MODE:%s CLIENT_SHORTNAME:%s ORDER_TYPE:%s PRICE:%.4f QUANTITY:%.2f START_DATETIME:%s END_DATETIME:%s DISPLAY_QTY:%d DISPLAY_QTY_UNIT:%s DISPLAY_QTY_VARIANCE:%d",
					tradeMode, clientName, order.AlgName, order.Price, order.OrderQty, algStartTime, algEndTime, displayQty, displayQtyUnit, displayQtyVariance)
			} else { //美股单
				text = fmt.Sprintf("TRADE_MODE:%s CLIENT_SHORTNAME:%s ORDER_TYPE:%s PRICE:%.4f QUANTITY:%.2f START_DATETIME:%s END_DATETIME:%s DISPLAY_QTY:%d",
					tradeMode, clientName, order.AlgName, order.Price, order.OrderQty, algStartTime, algEndTime, displayQty)
			}
		}
	}

	premarket, err := isPremarket(extendAttrMap)
	if err != nil {
		de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return
	}
	if premarket {
		text = "PREMARKET " + text
	}

	rawOrder.SetText(text)

	// 使用无后缀的symbol
	rawOrder.SetSymbol(order.Symbol2)

	msg = rawOrder.ToMessage()

	c.QueryHeader(&msg.Header)

	return msg, nil
}

func (c *GFGFIXChannel) NewOrderSingle(order *schema.TradeOrder, afterTradeOrderUpsert func(tradeOrder *schema.TradeOrder), syncTradeOrderAfterError func(tradeOrder *schema.TradeOrder)) (duplicatedOrder bool, de *domain_error.Error) {

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
	// 兼容gfgfix关于tag11前缀的需求
	// businessSource, err := getBusinessSource(order.ExtendAttrMap)
	// if err != nil {
	// 	de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
	// 	return
	// }
	// realClOrdID := businessSource + "-" + order.ClOrdID
	realClOrdID, err := c.getRealClOrdID(order)
	if err != nil {
		de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return
	}
	msg.Body.SetString(quickfix.Tag(11), realClOrdID)

	err = quickfix.Send(msg)
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

func (c *GFGFIXChannel) getRealClOrdID(order *schema.TradeOrder) (string, error) {
	// 更新ClOrdID
	// 兼容gfgfix关于tag11前缀的需求
	businessSource, err := getBusinessSource(order.ExtendAttrMap)
	if err != nil {
		return "", err
	}

	// Todo：如果gfgfix改成128长度限制，以下这两行代码就可以去掉
	realClOrdID := order.ClOrdID
	switch order.SystemCode {
	case "titans":
		realClOrdID = strings.Replace(realClOrdID, "-titans-", "-t-", 1)
	}
	switch order.BusinessCode {
	case "crossborder":
		realClOrdID = strings.Replace(realClOrdID, "-crossborder-", "-c-", 1)
	}

	realClOrdID = businessSource + "-" + realClOrdID

	return realClOrdID, nil
}

var (
	sellVal  = string(enum.Side_Sell)
	buyVal   = string(enum.Side_Buy)
	closeVal = string(enum.OpenClose_CLOSE)
	openVal  = string(enum.OpenClose_OPEN)
)

func getTradeMode(order *schema.TradeOrder, valMap map[string]interface{}) (val interface{}, err error) {

	openClose := order.OpenClose
	if len(openClose) == 0 {
		return "LOW_TOUCH", nil
	}

	side := order.Side

	if openClose == openVal && side == sellVal || openClose == closeVal && side == buyVal {
		return "HIGH_TOUCH", nil
	}

	return "LOW_TOUCH", nil
}

func getClientShortName(valMap map[string]interface{}) (val interface{}, err error) {
	val, _, err = attrutil.GetAttrValue(valMap, "clientShortName", enum.AttrValueType_STRING)
	return
}

func getAlgStartTime(valMap map[string]interface{}) (val interface{}, err error) {
	val, _, err = attrutil.GetAttrValue(valMap, "startTime", enum.AttrValueType_STRING)
	return
}

func getAlgEndTime(valMap map[string]interface{}) (val interface{}, err error) {
	val, _, err = attrutil.GetAttrValue(valMap, "endTime", enum.AttrValueType_STRING)
	return
}

func getAlgRatioParam(algParams map[string]interface{}, field string) (float64, *domain_error.Error) {
	val, ok, err := attrutil.GetAttrValue(algParams, field, enum.AttrValueType_FLOAT)
	if err != nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return 0, de
	}
	if ok {
		return val.(float64), nil
	}
	return 0, nil
}

func getPercent(valMap map[string]interface{}) (val interface{}, de *domain_error.Error) {
	val, de = getAlgRatioParam(valMap, "povRatio")
	return
}

func getMaxVol(valMap map[string]interface{}) (val interface{}, de *domain_error.Error) {
	val, de = getAlgRatioParam(valMap, "maxVol")
	return
}

func getDisplayQty(valMap map[string]interface{}) (val interface{}, err error) {
	val, _, err = attrutil.GetAttrValue(valMap, "displayQty", enum.AttrValueType_INT)
	return
}

func getDisplayQtyUnit(valMap map[string]interface{}) (val interface{}, err error) {
	val, _, err = attrutil.GetAttrValue(valMap, "displayQtyUnit", enum.AttrValueType_STRING)
	return
}

func getDisplayQtyVariance(valMap map[string]interface{}) (val interface{}, err error) {
	val, _, err = attrutil.GetAttrValue(valMap, "displayQtyVariance", enum.AttrValueType_INT)
	return
}

func getBusinessSource(valMap map[string]interface{}) (val string, err error) {
	var v interface{}
	v, _, err = attrutil.GetAttrValue(valMap, "businessSource", enum.AttrValueType_STRING)
	if err != nil {
		return
	}
	val = v.(string)
	if val == "" {
		return "HKHEDGE", nil
	}
	return
}

func isPremarket(valMap map[string]interface{}) (bool, error) {
	val, _, err := attrutil.GetAttrValue(valMap, "premarket", enum.AttrValueType_BOOL)
	if err != nil {
		return false, err
	}
	return val.(bool), err
}

func (c *GFGFIXChannel) OnMarketClosed(orderCache *order_cache.OrderCache) {
	
}
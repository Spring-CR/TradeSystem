package atp_cbf

import (
	"errors"
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
	"time"

	newordersingle44 "github.com/quickfixgo/fix44/newordersingle"
	"github.com/quickfixgo/quickfix"
	"github.com/shopspring/decimal"
)

type AtpFIXChannel struct {
	*fix.GenericFIXChannel
}

func NewAtpFIXChannel(cfg *domain_cfg.TradeChannelCfg, onTradeActionResp func(tradeActionResp *schema.TradeActionResp)) (channel *AtpFIXChannel, de *domain_error.Error) {

	channel = &AtpFIXChannel{}

	genericChannel, de1 := fix.NewGenericFIXChannel(cfg, onTradeActionResp, channel.CleanTradeActionResp, channel.RefineOrderCancelID)
	if de1 != nil {
		return nil, de1
	}

	channel.GenericFIXChannel = genericChannel

	return
}

func (c *AtpFIXChannel) ValidateOrder(order *schema.TradeOrder) (channelOrder interface{}, de *domain_error.Error) {

	var rawOrder newordersingle44.NewOrderSingle
	var msg *quickfix.Message

	rawOrder = fix_convert.DomainOrderToRawOrderFix44(order)

	// 使用原始的symbol
	rawOrder.SetSymbol(order.Symbol)

	// 要取整数
	if order.OrderQty > 0 {
		rawOrder.SetOrderQty(decimal.NewFromFloat(order.OrderQty), 0)
	}

	// 设置算法参数
	switch order.AlgName {
	case "POV":
		rawOrder.SetString(quickfix.Tag(40), "p")

		de = setParticipationPercent(rawOrder, order.AlgParamsMap)
		if de != nil {
			return
		}

		de = setCumulateUnmatchedQty(rawOrder, order.AlgParamsMap)
		if de != nil {
			return
		}

		de = c.setCommonParamsForPovTwap(rawOrder, order.AlgParamsMap)
		if de != nil {
			return
		}

	case "TWAP":
		rawOrder.SetString(quickfix.Tag(40), "t")

		de = c.setCommonParamsForPovTwap(rawOrder, order.AlgParamsMap)
		if de != nil {
			return
		}

	}

	// 显示设置expirtTime
	// if order.ExpireTime <= 0 && order.AlgName != "" {
	// 	tNow := time.Now()
	// 	_, _, _, endTime := c.GetChannelConfig().GetTradeChannelDetails().GetExchangeDate(tNow)
	// 	endTime = endTime.Add(-5 * time.Minute)
	// 	if endTime.Before(tNow) {
	// 		endTime = tNow
	// 	}
	// 	rawOrder.SetExpireTime(endTime)
	// }

	msg = rawOrder.ToMessage()

	c.QueryHeader(&msg.Header)

	return msg, nil
}

func (c *AtpFIXChannel) setCommonParamsForPovTwap(rawOrder newordersingle44.NewOrderSingle, algParams map[string]interface{}) (de *domain_error.Error) {

	de = setLimitPrice(rawOrder, algParams)
	if de != nil {
		return
	}

	de = setPriceType(rawOrder, algParams)
	if de != nil {
		return
	}

	de = c.setTime(rawOrder, algParams)
	if de != nil {
		return
	}

	de = setSubOrderQty(rawOrder, algParams)
	if de != nil {
		return
	}

	de = setPriceDiffInTick(rawOrder, algParams)
	if de != nil {
		return
	}

	de = setTriggerTimeInterval(rawOrder, algParams)
	if de != nil {
		return
	}

	de = setCancelTimeInterval(rawOrder, algParams)
	if de != nil {
		return
	}

	de = setMaxCancelOrderTimes(rawOrder, algParams)
	if de != nil {
		return
	}
	return
}

func setParticipationPercent(rawOrder newordersingle44.NewOrderSingle, algParams map[string]interface{}) *domain_error.Error {
	val, ok, err := attrutil.GetAttrValue(algParams, "participationPercent", enum.AttrValueType_FLOAT)
	if err != nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return de
	}
	if ok {
		percent := val.(float64)
		if percent <= 0 || percent >= 100 {
			return domain_error.Build(domain_error.GENERIC_ERR_CODE, errors.New("participationPercent must be greater than 0 and less than 100"))
		}
		percent = percent / 100.0
		rawOrder.SetString(quickfix.Tag(20035), fmt.Sprintf("%.4f", percent))
	}
	return nil
}

func setCumulateUnmatchedQty(rawOrder newordersingle44.NewOrderSingle, algParams map[string]interface{}) *domain_error.Error {
	val, ok, err := attrutil.GetAttrValue(algParams, "cumulateUnmatchedQty", enum.AttrValueType_BOOL)
	if err != nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return de
	}
	if ok {
		cumulateUnmatchedQty := val.(bool)
		if cumulateUnmatchedQty {
			rawOrder.SetString(quickfix.Tag(20026), "Y")
		} else {
			rawOrder.SetString(quickfix.Tag(20026), "N")
		}
	}
	return nil
}

func getEnableLimitPrice(algParams map[string]interface{}) (bool, *domain_error.Error) {
	val, ok, err := attrutil.GetAttrValue(algParams, "enableLimitPrice", enum.AttrValueType_BOOL)
	if err != nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return false, de
	}
	if ok {
		return val.(bool), nil
	}
	return false, nil
}

func setLimitPrice(rawOrder newordersingle44.NewOrderSingle, algParams map[string]interface{}) *domain_error.Error {
	enableLimitPrice, de := getEnableLimitPrice(algParams)
	if de != nil {
		return de
	}
	if enableLimitPrice {
		rawOrder.SetString(quickfix.Tag(20013), "Y")
		val, ok, err := attrutil.GetAttrValue(algParams, "limitPrice", enum.AttrValueType_FLOAT)
		if err != nil {
			de := domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
			return de
		}
		if ok {
			limitPrice := val.(float64)
			if limitPrice <= 0 {
				return domain_error.Build(domain_error.GENERIC_ERR_CODE, errors.New("limitPrice must be greater than 0"))
			}
			rawOrder.SetString(quickfix.Tag(20014), fmt.Sprintf("%.4f", limitPrice))
		}
	} else {
		rawOrder.SetString(quickfix.Tag(20013), "N")
	}
	return nil
}

var priceTypeMap = map[string]bool{
	"L": true, //Last Trade
	"B": true, //Market Buy
	"S": true, //Market Sell
	"P": true, //Best Bid/Offer
}

func setPriceType(rawOrder newordersingle44.NewOrderSingle, algParams map[string]interface{}) *domain_error.Error {
	val, ok, err := attrutil.GetAttrValue(algParams, "priceType", enum.AttrValueType_STRING)
	if err != nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return de
	}
	if ok {
		priceType := val.(string)
		if !priceTypeMap[priceType] {
			de := domain_error.Build(domain_error.GENERIC_ERR_CODE, errors.New("pirce type for ALGO can only be L, B, S or P"))
			return de
		}
		rawOrder.SetString(quickfix.Tag(20011), priceType)
	}
	return nil
}

func (c *AtpFIXChannel) setTime(rawOrder newordersingle44.NewOrderSingle, algParams map[string]interface{}) *domain_error.Error {
	val, ok, err := attrutil.GetAttrValue(algParams, "beginTime", enum.AttrValueType_INT)
	if err != nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return de
	}
	if ok && val.(int) > 0 {
		rawOrder.SetEffectiveTime(timeutil.ConvertMillisecondsToTime(int64(val.(int))))
	}

	val, ok, err = attrutil.GetAttrValue(algParams, "endTime", enum.AttrValueType_INT)
	if err != nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return de
	}
	if ok && val.(int) > 0 {
		rawOrder.SetExpireTime(timeutil.ConvertMillisecondsToTime(int64(val.(int))))
	} else {
		tNow := time.Now()
		_, _, _, endTime := c.GetChannelConfig().GetTradeChannelDetails().GetExchangeDate(tNow)
		endTime = endTime.Add(-5 * time.Minute)
		if endTime.Before(tNow) {
			endTime = tNow
		}
		rawOrder.SetExpireTime(endTime)
	}

	return nil
}

func setSubOrderQty(rawOrder newordersingle44.NewOrderSingle, algParams map[string]interface{}) *domain_error.Error {
	val, ok, err := attrutil.GetAttrValue(algParams, "subOrderQty", enum.AttrValueType_INT)
	if err != nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return de
	}
	if ok && val.(int) > 0 {
		rawOrder.SetInt(quickfix.Tag(20015), val.(int))
	}
	return nil
}

func setPriceDiffInTick(rawOrder newordersingle44.NewOrderSingle, algParams map[string]interface{}) *domain_error.Error {
	val, ok, err := attrutil.GetAttrValue(algParams, "priceDiffInTick", enum.AttrValueType_INT)
	if err != nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return de
	}
	if ok && val.(int) >= 0 {
		rawOrder.SetInt(quickfix.Tag(20012), val.(int))
	}
	return nil
}

func setTriggerTimeInterval(rawOrder newordersingle44.NewOrderSingle, algParams map[string]interface{}) *domain_error.Error {
	val, ok, err := attrutil.GetAttrValue(algParams, "triggerTimeIntervalSeconds", enum.AttrValueType_FLOAT)
	if err != nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return de
	}
	if ok && val.(float64) > 0 {
		rawOrder.SetInt(quickfix.Tag(20018), int(val.(float64)*1000))
	}
	return nil
}

func setCancelTimeInterval(rawOrder newordersingle44.NewOrderSingle, algParams map[string]interface{}) *domain_error.Error {
	val, ok, err := attrutil.GetAttrValue(algParams, "cancelTimeIntervalSeconds", enum.AttrValueType_FLOAT)
	if err != nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return de
	}
	if ok && val.(float64) > 0 {
		rawOrder.SetInt(quickfix.Tag(20019), int(val.(float64)*1000))
	}
	return nil
}

func setMaxCancelOrderTimes(rawOrder newordersingle44.NewOrderSingle, algParams map[string]interface{}) *domain_error.Error {
	val, ok, err := attrutil.GetAttrValue(algParams, "maxCancelOrderTimes", enum.AttrValueType_INT)
	if err != nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return de
	}
	if ok && val.(int) > 0 {
		rawOrder.SetInt(quickfix.Tag(20020), val.(int))
	}
	return nil
}

func (c *AtpFIXChannel) NewOrderSingle(order *schema.TradeOrder, afterTradeOrderUpsert func(tradeOrder *schema.TradeOrder), syncTradeOrderAfterError func(tradeOrder *schema.TradeOrder)) (duplicatedOrder bool, de *domain_error.Error) {

	log.Printf("======>NewOrderSingle, clOrdid:%s\n", order.ClOrdID)
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

func (c *AtpFIXChannel) CleanTradeActionResp(tradeActionResp *schema.TradeActionResp) (ignore bool) {
	ignore = len(tradeActionResp.ClOrdID) < 9 || tradeActionResp.ClOrdID[8] != '-'
	return
}

func (c *AtpFIXChannel) RefineOrderCancelID(rawTargetClOrdID string, rawClOrdID string, order *schema.TradeOrder) (targetClOrdID, clOrdID string) {
	// 不改
	return rawTargetClOrdID, rawClOrdID
}

func (c *AtpFIXChannel) OnMarketClosed(orderCache *order_cache.OrderCache) {
	
}

package hk_fut

import (
	"rhino-common/domain_error"
	"rhino-core/domain_cfg"
	"rhino-core/order_domain/order_cache"
	"rhino-core/schema"
	"rhino-core/store/app_store"
	"rhino-plugins/channel_plugin/fix"
	"rhino-trade-channel/adapter/data_convert/domain_order_validate"
	"rhino-trade-channel/adapter/data_convert/fix_convert"

	newordersingle42 "github.com/quickfixgo/fix42/newordersingle"
	"github.com/quickfixgo/quickfix"
)

type HKFutFIXChannel struct {
	*fix.GenericFIXChannel
}

func NewHKFutFIXChannel(cfg *domain_cfg.TradeChannelCfg, onTradeActionResp func(tradeActionResp *schema.TradeActionResp)) (channel *HKFutFIXChannel, de *domain_error.Error) {

	channel = &HKFutFIXChannel{}

	genericChannel, de1 := fix.NewGenericFIXChannel(cfg, onTradeActionResp, channel.CleanTradeActionResp, channel.RefineOrderCancelID)
	if de1 != nil {
		return nil, de1
	}

	channel.GenericFIXChannel = genericChannel

	return
}

func (c *HKFutFIXChannel) ValidateOrder(order *schema.TradeOrder) (channelOrder interface{}, de *domain_error.Error) {

	var rawOrder newordersingle42.NewOrderSingle
	var msg *quickfix.Message

	rawOrder = fix_convert.DomainOrderToRawOrderFix42(order)

	// 使用无后缀的symbol
	rawOrder.SetSymbol(order.Symbol2)
	if order.SecurityExchange != "" {
		rawOrder.SetSecurityExchange(order.SecurityExchange)
	}

	msg = rawOrder.ToMessage()

	c.QueryHeader(&msg.Header)

	return msg, nil
}

func (c *HKFutFIXChannel) NewOrderSingle(order *schema.TradeOrder, afterTradeOrderUpsert func(tradeOrder *schema.TradeOrder), syncTradeOrderAfterError func(tradeOrder *schema.TradeOrder)) (duplicatedOrder bool, de *domain_error.Error) {

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

func (c *HKFutFIXChannel) CleanTradeActionResp(tradeActionResp *schema.TradeActionResp) (ignore bool) {
	return
}

func (c *HKFutFIXChannel) RefineOrderCancelID(rawTargetClOrdID string, rawClOrdID string, order *schema.TradeOrder) (targetClOrdID, clOrdID string) {
	// 不改
	return rawTargetClOrdID, rawClOrdID
}

func (c *HKFutFIXChannel) OnMarketClosed(orderCache *order_cache.OrderCache) {
	
}
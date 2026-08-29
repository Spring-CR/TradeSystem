package order_position_manager

import "rhino-core/schema"

type QuotaLocker struct {
	tradeOrder          *schema.TradeOrder      // 交易订单
	freezeQty           float64                 // 订单冻结额度
	freezeQtyBuf        float64                 // T1_SHORT - T0_SHORT
	unfreezeQty         float64                 // 订单解冻额度
	qtyFreezeTime       int64                   // 额度冻结时间，毫秒时间戳
	qtyUnfreezeTime     int64                   // 额度冻结时间，毫秒时间戳
	lastTradeActionResp *schema.TradeActionResp // 最后一个交易回报
}

func newQuotaLocker(tradeOrder *schema.TradeOrder, freezeQty float64, freezeQtyBuf float64) *QuotaLocker {
	timeStamp := tradeOrder.OrdCreateTime
	if tradeOrder.OrdStatusUpdateTime > timeStamp {
		timeStamp = tradeOrder.OrdStatusUpdateTime
	}
	return &QuotaLocker{
		tradeOrder:    tradeOrder,
		freezeQty:     freezeQty,
		freezeQtyBuf:  freezeQtyBuf,
		qtyFreezeTime: timeStamp,
	}
}

func (l *QuotaLocker) unfreeze(unfreezeQty float64, tradeOrder *schema.TradeOrder, lastTradeActionResp *schema.TradeActionResp) {
	l.unfreezeQty = unfreezeQty
	l.lastTradeActionResp = lastTradeActionResp
	l.qtyUnfreezeTime = tradeOrder.OrdStatusUpdateTime
}

func (l *QuotaLocker) GetFreezeQty() float64 {
	return l.freezeQty
}

func (l *QuotaLocker) GetFreezeQtyBuf() float64 {
	return l.freezeQtyBuf
}

func (l *QuotaLocker) GetUnfreezeQty() float64 {
	return l.unfreezeQty
}

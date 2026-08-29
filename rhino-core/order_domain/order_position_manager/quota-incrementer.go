package order_position_manager

import "rhino-core/schema"

type FillReport struct {
	tradeActionResp *schema.TradeActionResp // 成交回报
	incQty          float64                 // 持仓增加数量
}

func newFillReport(tradeActionResp *schema.TradeActionResp, incQty float64) *FillReport {
	return &FillReport{tradeActionResp: tradeActionResp, incQty: incQty}
}

type QuotaIncrementer struct {
	tradeOrder  *schema.TradeOrder // 交易订单
	fillReports []*FillReport      // 成交回报列表
	totalIncQty float64            // 累计增加额度
}

func newQuotaIncrementer(tradeOrder *schema.TradeOrder) *QuotaIncrementer {
	return &QuotaIncrementer{tradeOrder: tradeOrder}
}

func (i *QuotaIncrementer) increase(increaseQty float64, tradeActionResp *schema.TradeActionResp) {
	i.totalIncQty += increaseQty
	i.fillReports = append(i.fillReports, newFillReport(tradeActionResp, increaseQty))
}

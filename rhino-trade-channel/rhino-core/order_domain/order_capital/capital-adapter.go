package order_capital

import (
	"rhino-common/domain_error"
	"rhino-core/schema"
)

type CapitalAdapter interface {
	WouldOrderCapitalIncrease(order *schema.TradeOrder) bool
	WouldOrderCapitalDecrease(order *schema.TradeOrder) bool
	FreezeOrderCapitalOnTradeBegin(force bool, order *schema.TradeOrder) (freezedAmount float64, turnToReviewOnErr bool, de *domain_error.Error)
	RollbackFreezedOrderCapitalForTradeError(order *schema.TradeOrder)
}

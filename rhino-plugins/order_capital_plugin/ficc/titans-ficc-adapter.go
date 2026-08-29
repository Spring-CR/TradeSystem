package ficc

import (
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	positionplugin "rhino-plugins/order_position_plugin/ficc"
)

type TitansFiccOrderCapitalAdapter struct {
	applicationCfg *domain_cfg.ApplicationCfg
}

func NewTitansFiccOrderCapitalAdapter(applicationCfg *domain_cfg.ApplicationCfg) (adapter *TitansFiccOrderCapitalAdapter, de *domain_error.Error) {
	log.Printf("construct TitansFiccOrderCapitalAdapter...")
	adapter = &TitansFiccOrderCapitalAdapter{applicationCfg: applicationCfg}
	return
}

func (a *TitansFiccOrderCapitalAdapter) WouldOrderCapitalDecrease(order *schema.TradeOrder) bool {
	return order.Side == "O"
}

func (a *TitansFiccOrderCapitalAdapter) WouldOrderCapitalIncrease(order *schema.TradeOrder) bool {
	return order.OpenClose == "C"
}

func (a *TitansFiccOrderCapitalAdapter) FreezeOrderCapitalOnTradeBegin(force bool, order *schema.TradeOrder) (freezedAmount float64, turnToReviewOnErr bool, de *domain_error.Error) {
	var availableAmt float64
	var ok bool
	var err error
	log.Printf("start to FreezeCapitalForTradeStart, appOrdID=%s\n", order.AppOrdID)
	availableAmt, freezedAmount, ok, err = positionplugin.CapitalController.FreezeCapitalForTradeStart(force, order)
	if err != nil {
		de = domain_error.Build(domain_error.CANNOT_FREEZE_CAPITAL_ON_TRADE_START_ERR_CODE, err, order.AppOrdID)
		return
	}

	// 除了!ok，还要对比金额才知道是否可用余额不足引起的问题
	if !force && !ok && freezedAmount > availableAmt {
		de = domain_error.Build(domain_error.CAPITAL_AMOUNT_NOT_ENOUGH_ERR_CODE, fmt.Errorf("可用余额不足[可用=%v, 需冻结=%v]", availableAmt, freezedAmount))
		turnToReviewOnErr = true
		return
	}

	return
}

func (a *TitansFiccOrderCapitalAdapter) RollbackFreezedOrderCapitalForTradeError(order *schema.TradeOrder) {
	positionplugin.CapitalController.RollbackFreezedOrderCapitalForTradeError(order)
}

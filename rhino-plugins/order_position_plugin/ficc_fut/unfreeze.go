package ficc_fut

import (
	"encoding/json"
	"errors"
	"math"
	"rhino-common/domain_error"
	"rhino-core/order_domain/order_position_manager"
	"rhino-core/schema"
	statusficc "rhino-plugins/order_status_plugin/ficc"
)

func (a *FiccFutOrderPositionAdapter) CalculateUnfreezeQuotaInPositionUnit(positionUnit *order_position_manager.PositionUnit, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, metadata interface{}, quotaLocker *order_position_manager.QuotaLocker) (unfreezeQuota float64, ok bool) {

	if !a.isOrderFinished(tradeActionResp) {
		return
	}

	if len(tradeOrder.ExtendAttrMap) == 0 {
		return
	}

	positionRecord, ok := metadata.(*PositionRecord)
	if !ok {
		errMsg := "metadata is not type of *PositionRecord"
		domain_error.ProcessSevereError(true, 5, nil, errors.New(errMsg), errMsg)
	}

	freezeQty := quotaLocker.GetFreezeQty()
	cumQty := float64(tradeOrder.CumQty)

	a.orderLog.Printf(tradeOrder, tradeActionResp, "[CalculateUnfreezeQuotaInPositionUnit=%v], OrderQty=%v, FreezeQty=%v, CumQty=%v", positionUnit.GetName(), tradeOrder.OrderQty, freezeQty, cumQty)

	defer func() {
		ok = true
	}()

	// 按持仓类型，计算解冻的持仓量
	switch positionUnit.GetName() {

	case positionUnitNameLong:

		if tradeOrder.Side != sideSell {
			return
		}

		// Ord_[SELL][END].UFRZQTY_TMP = Ord_[SELL][RUN].FRZQTY - MIN(Ord_[SELL][RUN].FRZQTY, Ord_[SELL][END].QTYREAL)
		unfreezeQuotaTmp := freezeQty - math.Min(freezeQty, cumQty)

		// Ord_[SELL][END].UFRZQTY = NAV > 0 ? Ord_[SELL][END].UFRZQTY_TMP: ( NAV + Ord_[SELL][END].UFRZQTY_TMP <= 0 ? 0 : NAV + Ord_[SELL][END].UFRZQTY_TMP )
		if positionRecord.NetPosition > 0 {
			unfreezeQuota = unfreezeQuotaTmp
		} else {
			unfreezeQuota = positionRecord.NetPosition + unfreezeQuotaTmp
			if unfreezeQuota <= 0 {
				unfreezeQuota = 0
			}
		}

		ok = true
		return

	case positionUnitNameShort:

		if tradeOrder.Side != sideBuy {
			return
		}

		// Ord_[BUY][END].UFRZQTY_TMP = Ord_[BUY][RUN].FRZQTY - MIN(Ord_[BUY][RUN].FRZQTY, Ord_[BUY][END].QTYREAL)
		unfreezeQuotaTmp := freezeQty - math.Min(freezeQty, cumQty)

		// Ord_[BUY][END].UFRZQTY = NAV<0 ? Ord_[BUY][END].UFRZQTY_TMP : ( NAV >= Ord_[BUY][END].UFRZQTY_TMP ? 0 : Ord_[BUY][END].UFRZQTY_TMP - NAV )
		if positionRecord.NetPosition < 0 {
			unfreezeQuota = unfreezeQuotaTmp
		} else {
			unfreezeQuota = -positionRecord.NetPosition + unfreezeQuotaTmp
			if unfreezeQuota <= 0 {
				unfreezeQuota = 0
			}
		}

		ok = true
		return
	}

	return
}

func (a *FiccFutOrderPositionAdapter) AfterUnfreezeQuotaInPositionUnit(positionUnit *order_position_manager.PositionUnit, unfreezeQuota float64, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, metadata interface{}, quotaLocker *order_position_manager.QuotaLocker, lastPositionUnit bool, rollback bool) {

	positionRecord, ok := metadata.(*PositionRecord)
	if !ok {
		errMsg := "metadata is not type of *PositionRecord"
		domain_error.ProcessSevereError(true, 5, nil, errors.New(errMsg), errMsg)
	}

	switch positionUnit.GetName() {

	case positionUnitNameLong:

		positionRecord.LongAvailablePosition += unfreezeQuota

	case positionUnitNameShort:

		positionRecord.ShortAvailablePosition += unfreezeQuota
	}

	if lastPositionUnit {
		a.marginExposureManager.DeleteOrder(tradeOrder)
		a.marginExposureManager.CalculateMarginExposure(positionRecord, tradeOrder)

		js, _ := json.MarshalIndent(positionRecord, "", "  ")
		a.orderLog.Printf(tradeOrder, nil, "[AfterUnfreezeQuotaInPositionUnit] PositionRecord=%s", js)
	}
}

func (a *FiccFutOrderPositionAdapter) isOrderFinished(resp *schema.TradeActionResp) bool {
	return statusficc.EndStatus[resp.OrdStatus]
}

package ficc_v2

import (
	"encoding/json"
	"errors"
	"math"
	"rhino-common/domain_error"
	"rhino-core/order_domain/order_position_manager"
	"rhino-core/schema"
	statusficc "rhino-plugins/order_status_plugin/ficc"
)

func (a *TitansFiccOrderPositionAdapter) CalculateUnfreezeQuotaInPositionUnit(positionUnit *order_position_manager.PositionUnit, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, metadata interface{}, quotaLocker *order_position_manager.QuotaLocker) (unfreezeQuota float64, ok bool) {

	if !a.isOrderFinished(tradeActionResp) {
		return
	}

	if len(tradeOrder.ExtendAttrMap) == 0 {
		return
	}

	val := tradeOrder.ExtendAttrMap["settlType"]
	settleType, _ := val.(string)

	freezeQty := quotaLocker.GetFreezeQty()
	cumQty := float64(tradeOrder.CumQty)

	a.orderLog.Printf(tradeOrder, tradeActionResp, "[CalculateUnfreezeQuotaInPositionUnit=%v] SettlType=%v, OrderQty=%v, FreezeQty=%v, CumQty=%v", positionUnit.GetName(), settleType, tradeOrder.OrderQty, freezeQty, cumQty)

	defer func() {
		ok = true
	}()

	// 按持仓类型，计算解冻的持仓量
	switch positionUnit.GetName() {

	case positionUnitNameLongT0:

		if tradeOrder.Side != sideSell {
			return
		}

		switch settleType {
		case settleTypeT0:
			unfreezeQuota = freezeQty - math.Min(freezeQty, cumQty)
			ok = true
			return

		case settleTypeT1:
			unfreezeQuota = freezeQty - math.Min(freezeQty, cumQty)
			ok = true
			return
		}

	case positionUnitNameLongT1:

		if tradeOrder.Side != sideSell {
			return
		}

		switch settleType {
		case settleTypeT0:
			unfreezeQuota = freezeQty - math.Min(freezeQty, cumQty)
			ok = true
			return

		case settleTypeT1:
			unfreezeQuota = freezeQty - math.Min(freezeQty, cumQty)
			ok = true
			return
		}

	case positionUnitNameShortT0:

		if tradeOrder.Side != sideBuy {
			return
		}

		switch settleType {
		case settleTypeT0:
			unfreezeQuota = freezeQty - math.Min(freezeQty, cumQty)
			ok = true
			return

		case settleTypeT1:
			unfreezeQuota = freezeQty - math.Min(freezeQty, math.Max(0, cumQty - quotaLocker.GetFreezeQtyBuf()))
			ok = true
			return
		}

	case positionUnitNameShortT1:

		if tradeOrder.Side != sideBuy {
			return
		}

		switch settleType {
		case settleTypeT0:
			unfreezeQuota = freezeQty - math.Min(freezeQty, cumQty)
			ok = true
			return

		case settleTypeT1:
			unfreezeQuota = freezeQty - math.Min(freezeQty, cumQty)
			ok = true
			return
		}
	}

	return
}

func (a *TitansFiccOrderPositionAdapter) AfterUnfreezeQuotaInPositionUnit(positionUnit *order_position_manager.PositionUnit, unfreezeQuota float64, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, metadata interface{}, quotaLocker *order_position_manager.QuotaLocker, lastPositionUnit bool, rollback bool) {

	positionRecord, ok := metadata.(*PositionRecord)
	if !ok {
		errMsg := "metadata is not type of *PositionRecord"
		domain_error.ProcessSevereError(true, 5, nil, errors.New(errMsg), errMsg)
	}

	switch positionUnit.GetName() {

	case positionUnitNameLongT0:

		positionRecord.LongAvailablePositionT0 += unfreezeQuota

		/*if positionRecord.NetPositionT1 >= 0 {

			maxLongAvailablePositionT0 := positionRecord.NetPositionT0
			if positionRecord.LongAvailablePositionT0 > maxLongAvailablePositionT0 {
				positionRecord.LongAvailablePositionT0 = maxLongAvailablePositionT0
			}

		} else {
			positionRecord.LongAvailablePositionT0 = 0
		}*/

	case positionUnitNameLongT1:

		positionRecord.LongAvailablePositionT1 += unfreezeQuota

		/*if positionRecord.NetPositionT1 >= 0 {

			maxLongAvailablePositionT1 := positionRecord.NetPositionT1
			if positionRecord.LongAvailablePositionT1 > maxLongAvailablePositionT1 {
				positionRecord.LongAvailablePositionT1 = maxLongAvailablePositionT1
			}

		} else {
			positionRecord.LongAvailablePositionT1 = 0
		}*/

	case positionUnitNameShortT0:

		positionRecord.ShortAvailablePositionT0 += unfreezeQuota

		/*if positionRecord.NetPositionT1 <= 0 {

			maxShortAvailablePositionT0 := -positionRecord.NetPositionT0
			if positionRecord.ShortAvailablePositionT0 > maxShortAvailablePositionT0 {
				positionRecord.ShortAvailablePositionT0 = maxShortAvailablePositionT0
			}

		} else {
			positionRecord.ShortAvailablePositionT0 = 0
		}*/

	case positionUnitNameShortT1:

		positionRecord.ShortAvailablePositionT1 += unfreezeQuota

		/*if positionRecord.NetPositionT1 <= 0 {

			maxShortAvailablePositionT1 := -positionRecord.NetPositionT1
			if positionRecord.ShortAvailablePositionT1 > maxShortAvailablePositionT1 {
				positionRecord.ShortAvailablePositionT1 = maxShortAvailablePositionT1
			}

		} else {
			positionRecord.ShortAvailablePositionT1 = 0
		}*/
	}

	//if lastPositionUnit && rollback {
	if lastPositionUnit {
		a.marginExposureManager.DeleteOrder(tradeOrder)
		a.marginExposureManager.CalculateMarginExposure(positionRecord, tradeOrder)

		js, _ := json.MarshalIndent(positionRecord, "", "  ")
		a.orderLog.Printf(tradeOrder, nil, "[AfterUnfreezeQuotaInPositionUnit] PositionRecord=%s", js)
	}
}

func (a *TitansFiccOrderPositionAdapter) isOrderFinished(resp *schema.TradeActionResp) bool {
	return statusficc.EndStatus[resp.OrdStatus]
}

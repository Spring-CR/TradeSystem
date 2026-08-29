package ficc_v2

import (
	"encoding/json"
	"errors"
	"rhino-common/domain_error"
	"rhino-core/order_domain/order_position_manager"
	"rhino-core/schema"
)

func (a *TitansFiccOrderPositionAdapter) CalculateIncreaseQuotaInPositionUnit(positionUnit *order_position_manager.PositionUnit, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, metadata interface{}) (increaseQuota float64, ok bool) {

	if tradeActionResp.LastShares <= 0 {
		return
	}

	positionRecord, ok := metadata.(*PositionRecord)
	if !ok {
		errMsg := "metadata is not type of *PositionRecord"
		domain_error.ProcessSevereError(true, 5, nil, errors.New(errMsg), errMsg)
	}

	if len(tradeOrder.ExtendAttrMap) == 0 {
		return
	}

	val := tradeOrder.ExtendAttrMap["settlType"]
	settleType, _ := val.(string)

	a.orderLog.Printf(tradeOrder, tradeActionResp, "[CalculateIncreaseQuotaInPositionUnit=%v] SettlType=%v, OrderQty=%v, LastShares=%v, NetPositionT0=%v, NetPositionT1=%v", positionUnit.GetName(), settleType, tradeOrder.OrderQty,tradeOrder.LastShares, positionRecord.NetPositionT0, positionRecord.NetPositionT1)

	// 按持仓类型，计算增加持仓量
	switch positionUnit.GetName() {

	// 计算T+0多头可用持仓
	case positionUnitNameLongT0:

		if tradeOrder.Side != sideBuy {
			return
		}

		switch settleType {
		case settleTypeT0:
			if positionRecord.NetPositionT1 > 0 {
				increaseQuota = float64(tradeActionResp.LastShares)
				ok = true
				return
			} else {
				increaseQuota = positionRecord.NetPositionT1 + float64(tradeActionResp.LastShares)
				if increaseQuota <= 0 {
					increaseQuota = 0
				}
				ok = true
				return
			}
		case settleTypeT1: // T+1 买单成交回报不能提升T+0持仓
			return
		}

	// 计算T+1多头可用持仓
	case positionUnitNameLongT1:

		if tradeOrder.Side != sideBuy {
			return
		}

		switch settleType {
		case settleTypeT0:
			if positionRecord.NetPositionT1 > 0 {
				increaseQuota = float64(tradeActionResp.LastShares)
				ok = true
				return
			} else {
				increaseQuota = positionRecord.NetPositionT1 + float64(tradeActionResp.LastShares)
				if increaseQuota <= 0 {
					increaseQuota = 0
				}
				ok = true
				return
			}
		case settleTypeT1:
			if positionRecord.NetPositionT1 > 0 {
				increaseQuota = float64(tradeActionResp.LastShares)
				ok = true
				return
			} else {
				increaseQuota = positionRecord.NetPositionT1 + float64(tradeActionResp.LastShares)
				if increaseQuota <= 0 {
					increaseQuota = 0
				}
				ok = true
				return
			}
		}

	// 计算T+0空头可用持仓
	case positionUnitNameShortT0:
		// T+0空头可用持仓计算
		if tradeOrder.Side != sideSell {
			return
		}

		switch settleType {
		case settleTypeT0:
			if positionRecord.NetPositionT1 < 0 {
				increaseQuota = float64(tradeActionResp.LastShares)
				ok = true
				return
			} else {
				increaseQuota = -positionRecord.NetPositionT1 + float64(tradeActionResp.LastShares)
				if increaseQuota <= 0 {
					increaseQuota = 0
				}
				ok = true
				return
			}
		case settleTypeT1: // T+1 买单成交回报不能提升T+0持仓
			return
		}

	// 计算T+1空头可用持仓
	case positionUnitNameShortT1:

		if tradeOrder.Side != sideSell {
			return
		}

		switch settleType {
		case settleTypeT0:
			if positionRecord.NetPositionT1 < 0 {
				increaseQuota = float64(tradeActionResp.LastShares)
				ok = true
				return
			} else {
				increaseQuota = -positionRecord.NetPositionT1 + float64(tradeActionResp.LastShares)
				if increaseQuota <= 0 {
					increaseQuota = 0
				}
				ok = true
				return
			}
		case settleTypeT1:
			if positionRecord.NetPositionT1 < 0 {
				increaseQuota = float64(tradeActionResp.LastShares)
				ok = true
				return
			} else {
				increaseQuota = -positionRecord.NetPositionT1 + float64(tradeActionResp.LastShares)
				if increaseQuota <= 0 {
					increaseQuota = 0
				}
				ok = true
				return
			}
		}
	}

	return
}

func (a *TitansFiccOrderPositionAdapter) AfterIncreaseQuotaInPositionUnit(positionUnit *order_position_manager.PositionUnit, increaseQty float64, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, metadata interface{}, quotaIncrementer *order_position_manager.QuotaIncrementer, lastPositionUnit bool) {

	positionRecord, ok := metadata.(*PositionRecord)
	if !ok {
		errMsg := "metadata is not type of *PositionRecord"
		domain_error.ProcessSevereError(true, 5, nil, errors.New(errMsg), errMsg)
	}

	switch positionUnit.GetName() {

	case positionUnitNameLongT0:

		positionRecord.LongAvailablePositionT0 += increaseQty

	case positionUnitNameLongT1:

		positionRecord.LongAvailablePositionT1 += increaseQty

	case positionUnitNameShortT0:

		positionRecord.ShortAvailablePositionT0 += increaseQty

	case positionUnitNameShortT1:

		positionRecord.ShortAvailablePositionT1 += increaseQty
	}

	if lastPositionUnit {
		a.marginExposureManager.AddOrder(tradeOrder)
		a.marginExposureManager.CalculateMarginExposure(positionRecord, tradeOrder)

		js, _ := json.MarshalIndent(positionRecord, "", "  ")
		a.orderLog.Printf(tradeOrder, nil, "[AfterIncreaseQuotaInPositionUnit] PositionRecord=%s", js)
	}
}

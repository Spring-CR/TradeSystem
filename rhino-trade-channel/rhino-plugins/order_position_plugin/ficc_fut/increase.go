package ficc_fut

import (
	"encoding/json"
	"errors"
	"log"
	"rhino-common/domain_error"
	"rhino-core/order_domain/order_position_manager"
	"rhino-core/schema"
)

func (a *FiccFutOrderPositionAdapter) CalculateIncreaseQuotaInPositionUnit(positionUnit *order_position_manager.PositionUnit, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, metadata interface{}) (increaseQuota float64, ok bool) {

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

	a.orderLog.Printf(tradeOrder, tradeActionResp, "[CalculateIncreaseQuotaInPositionUnit=%v], OrderQty=%v, LastShares=%v, NetPosition=%v", positionUnit.GetName(), tradeOrder.OrderQty, tradeOrder.LastShares, positionRecord.NetPosition)

	// 按持仓类型，计算增加持仓量
	switch positionUnit.GetName() {

	// 计算多头可用持仓
	case positionUnitNameLong:

		if tradeOrder.Side != sideBuy {
			return
		}

		if positionRecord.NetPosition > 0 {
			increaseQuota = float64(tradeActionResp.LastShares)
			ok = true
			return
		} else {
			increaseQuota = positionRecord.NetPosition + float64(tradeActionResp.LastShares)
			if increaseQuota <= 0 {
				increaseQuota = 0
			}
			ok = true
			return
		}

	// 计算空头可用持仓
	case positionUnitNameShort:
		// T+0空头可用持仓计算
		if tradeOrder.Side != sideSell {
			return
		}

		if positionRecord.NetPosition < 0 {
			increaseQuota = float64(tradeActionResp.LastShares)
			ok = true
			return
		} else {
			increaseQuota = -positionRecord.NetPosition + float64(tradeActionResp.LastShares)
			if increaseQuota <= 0 {
				increaseQuota = 0
			}
			ok = true
			return
		}
	}

	return
}

func (a *FiccFutOrderPositionAdapter) AfterIncreaseQuotaInPositionUnit(positionUnit *order_position_manager.PositionUnit, increaseQty float64, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, metadata interface{}, quotaIncrementer *order_position_manager.QuotaIncrementer, lastPositionUnit bool) {

	positionRecord, ok := metadata.(*PositionRecord)
	if !ok {
		errMsg := "metadata is not type of *PositionRecord"
		domain_error.ProcessSevereError(true, 5, nil, errors.New(errMsg), errMsg)
	}

	switch positionUnit.GetName() {

	case positionUnitNameLong:

		positionRecord.LongAvailablePosition += increaseQty

	case positionUnitNameShort:

		positionRecord.ShortAvailablePosition += increaseQty
	}

	if lastPositionUnit {

		log.Printf("AfterIncreaseQuotaInPositionUnit, AddOrder:%s\n", tradeOrder.AppOrdID)

		a.marginExposureManager.AddOrder(tradeOrder)
		a.marginExposureManager.CalculateMarginExposure(positionRecord, tradeOrder)

		js, _ := json.MarshalIndent(positionRecord, "", "  ")
		a.orderLog.Printf(tradeOrder, nil, "[AfterIncreaseQuotaInPositionUnit] PositionRecord=%s", js)
	}
}

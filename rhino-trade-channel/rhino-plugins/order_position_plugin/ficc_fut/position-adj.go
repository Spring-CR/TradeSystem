package ficc_fut

import (
	"fmt"
	"rhino-common/context/constant"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	apiplugin "rhino-plugins/api_plugin/ficc_fut"
	"time"
)

// TradeOrder:
// account  --- 必须
// side     --- 必须
// price    --- 必须
// orderQty --- 必须
// symbol2  --- 必须
// planCode --- 必须
// priceCNY           --- 派生
// priceCNYWithFee    --- 派生
// avgPriceWithFee    --- 派生
// contractMultiplier --- 派生
// exchangeRateCNY    --- 派生
// commissionType     --- 派生
// commissionValue    --- 派生

// TradeActionResp:
// respPriceWithFee
func (a *FiccFutOrderPositionAdapter) PreparePositionAdjustmentParams(tradeOrder *schema.TradeOrder) (mockTradeOrder *schema.TradeOrder, mockTradeActionResp *schema.TradeActionResp, de *domain_error.Error) {

	mockTradeOrder = tradeOrder

	if mockTradeOrder.AppOrdID == "" {
		appOrdID := fmt.Sprintf("%s%v", constant.MockPositionOrdIDPrefix, timeutil.ConvertTimeToMicroseconds(time.Now()))
		mockTradeOrder.AppOrdID = appOrdID
	}

	appOrdID := mockTradeOrder.AppOrdID

	mockTradeOrder.OrdStatus = "2"
	mockTradeOrder.OrdStatus2 = "2"

	mockTradeActionResp = &schema.TradeActionResp{
		ExecID:        appOrdID,
		AppOrdID:      appOrdID,
		OrdStatus:     "2",
		ExecType:      "2",
		TransactTime:  mockTradeOrder.TransactTime,
		ExtendAttrMap: make(map[string]interface{}),
	}

	mockTradeActionResp.Account = mockTradeOrder.Account
	mockTradeActionResp.Side = mockTradeOrder.Side

	// 设置费前价格
	mockTradeOrder.AvgPx = mockTradeOrder.Price
	mockTradeOrder.LastPx = mockTradeOrder.Price
	mockTradeActionResp.Price = mockTradeOrder.Price
	mockTradeActionResp.AvgPx = mockTradeOrder.Price
	mockTradeActionResp.LastPx = mockTradeOrder.Price

	// 设置数量
	mockTradeOrder.CumQty = int64(mockTradeOrder.OrderQty)
	mockTradeOrder.LastShares = int64(mockTradeOrder.OrderQty)
	mockTradeActionResp.OrderQty = mockTradeOrder.OrderQty
	mockTradeActionResp.CumQty = int64(mockTradeOrder.OrderQty)
	mockTradeActionResp.LastShares = int64(mockTradeOrder.OrderQty)

	// 设置标的
	mockTradeActionResp.Symbol = mockTradeOrder.Symbol

	commissionType, _, _ := attrutil.GetAttrValue(mockTradeOrder.ExtendAttrMap, "commissionType", enum.AttrValueType_STRING)
	commissionValue, _, _ := attrutil.GetAttrValue(mockTradeOrder.ExtendAttrMap, "commissionValue", enum.AttrValueType_FLOAT)

	// 设置含费均价
	avgPriceWithFee := apiplugin.GetPriceWithFee(mockTradeOrder.AvgPx, commissionType.(string), commissionValue.(float64), mockTradeOrder.Side)
	mockTradeOrder.ExtendAttrMap["avgPriceWithFee"] = avgPriceWithFee
	mockTradeActionResp.ExtendAttrMap["respPriceWithFee"] = avgPriceWithFee

	return
}
/*
func (a *FiccFutOrderPositionAdapter) PreparePositionAdjustmentParamsByPositionBaseDiff(_positionRecordBase, _positionRecordCurr interface{}) (mockTradeOrder *schema.TradeOrder, mockTradeActionResp *schema.TradeActionResp, de *domain_error.Error) {

	positionRecordBase, ok := _positionRecordBase.(PositionRecord)
	if !ok {
		return
	}

	positionRecordCurr, ok := _positionRecordCurr.(PositionRecord)
	if !ok {
		return
	}

	netPositionBase := math.Round(positionRecordBase.NetPosition)
	netPositionCurr := math.Round(positionRecordCurr.NetPosition)

	contractMultiplier := positionRecordBase.ContractMultiplier
	if contractMultiplier <= 0 {
		contractMultiplier = 1
	}

	// 如果净持仓不变，就不调了
	if netPositionBase == netPositionCurr {
		return
	}

	// 计算方向
	side := "1"
	if netPositionBase < netPositionCurr {
		side = "2"
	}

	// 计算数量
	ordQty := math.Round(netPositionBase - netPositionCurr)

	var price float64
	// 计算价格
	if netPositionBase == 0 {

		if netPositionCurr > 0 {
			price = positionRecordCurr.InitLongPriceCost / netPositionCurr / contractMultiplier
		} else { // netPositionCurr < 0
			price = positionRecordCurr.InitShortPriceCost / -netPositionCurr / contractMultiplier
		}

	} else {

		if netPositionBase > 0 { // 现在是多头

			if netPositionCurr >= 0 {

				if netPositionBase > netPositionCurr { // 多头的数量增加了
					if positionRecordBase.LongPriceCost > positionRecordCurr.LongPriceCost { // 名本也是增加的
						price = (positionRecordBase.LongPriceCost - positionRecordCurr.LongPriceCost) / (netPositionBase - netPositionCurr) / contractMultiplier
					} else { // 多头的数量增加了，但是名本反而减少了
						// 让名本不变，即不能减少
						price = 0.0
					}
				} else {
					// 多头减少，均价其实不变
					price = 0.0
				}

			} else {
				// 原来空头，现在要变成多头
				price = positionRecordBase.LongPriceCost / netPositionBase / contractMultiplier
			}

		} else { // 现在是空头， netPositionBase < 0 

			if netPositionCurr >= 0 {
				// 原来多头，现在要变成空头
				price = positionRecordBase.ShortPriceCost / -netPositionBase / contractMultiplier
			} else {

			}
		}
	}

	return
}*/

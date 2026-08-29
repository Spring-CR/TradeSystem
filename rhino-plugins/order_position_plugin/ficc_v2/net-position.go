package ficc_v2

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-core/schema"
	"rhino-core/types"
	"strings"
)

func (a *TitansFiccOrderPositionAdapter) resetPostions(tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, positionRecord *PositionRecord) {
	if positionRecord.NetPositionT1 >= 0 {
		maxLongAvailablePositionT0 := positionRecord.NetPositionT0
		if positionRecord.LongAvailablePositionT0 > maxLongAvailablePositionT0 {
			positionRecord.LongAvailablePositionT0 = maxLongAvailablePositionT0
		}
	} else {
		positionRecord.LongAvailablePositionT0 = 0
	}

	if positionRecord.NetPositionT1 >= 0 {
		maxLongAvailablePositionT1 := positionRecord.NetPositionT1
		if positionRecord.LongAvailablePositionT1 > maxLongAvailablePositionT1 {
			positionRecord.LongAvailablePositionT1 = maxLongAvailablePositionT1
		}
	} else {
		positionRecord.LongAvailablePositionT1 = 0
	}

	if positionRecord.NetPositionT1 <= 0 {
		maxShortAvailablePositionT0 := -positionRecord.NetPositionT0
		if positionRecord.ShortAvailablePositionT0 > maxShortAvailablePositionT0 {
			positionRecord.ShortAvailablePositionT0 = maxShortAvailablePositionT0
		}
	} else {
		positionRecord.ShortAvailablePositionT0 = 0
	}

	if positionRecord.NetPositionT1 <= 0 {
		maxShortAvailablePositionT1 := -positionRecord.NetPositionT1
		if positionRecord.ShortAvailablePositionT1 > maxShortAvailablePositionT1 {
			positionRecord.ShortAvailablePositionT1 = maxShortAvailablePositionT1
		}
	} else {
		positionRecord.ShortAvailablePositionT1 = 0
	}

	js, _ := json.MarshalIndent(positionRecord, "", "  ")
	a.orderLog.Printf(tradeOrder, tradeActionResp, "[ResetPostions] PositionRecord=%s", js)
}

// 计算净持仓和持仓均价
func (a *TitansFiccOrderPositionAdapter) AfterUpdateQuota(tradeResp *types.TradeActionRespReturn, metadata interface{}) {

	tradeOrder := tradeResp.GetTradeOrder()
	tradeActionResp := tradeResp.CurrentTradeActionResp
	positionRecord, ok := metadata.(*PositionRecord)
	if !ok {
		errMsg := "metadata is not type of *PositionRecord"
		domain_error.ProcessSevereError(true, 5, nil, errors.New(errMsg), errMsg)
	}

	if tradeActionResp.LastShares <= 0 {
		a.resetPostions(tradeOrder, tradeActionResp, positionRecord)
		return
	}

	if len(tradeOrder.ExtendAttrMap) == 0 {
		return
	}
	val := tradeOrder.ExtendAttrMap["settlType"]
	settleType, _ := val.(string)

	T0 := positionRecord.NetPositionT0
	T1 := positionRecord.NetPositionT1
	N := float64(tradeActionResp.LastShares)

	a.orderLog.Printf(tradeOrder, tradeActionResp, "[CalculateNetPosition] settleType=%v, T0=%v, T1=%v, N=%v", settleType, T0, T1, N)

	switch tradeOrder.Side {
	case sideSell:
		a.processSellTradeResp(T0, T1, N, settleType, tradeOrder, tradeActionResp, positionRecord)
	case sideBuy:
		a.processBuyTradeResp(T0, T1, N, settleType, tradeOrder, tradeActionResp, positionRecord)
	}

	a.resetPostions(tradeOrder, tradeActionResp, positionRecord)

	a.marginExposureManager.UpdateOrder(tradeResp)
	a.marginExposureManager.CalculateMarginExposure(positionRecord, tradeOrder)

	js, _ := json.MarshalIndent(positionRecord, "", "  ")
	a.orderLog.Printf(tradeOrder, tradeActionResp, "[AfterCalculateNetPosition] PositionRecord=%s", js)
}

func (a *TitansFiccOrderPositionAdapter) processSellTradeResp(T0 float64, T1 float64, N float64, settleType string, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, positionRecord *PositionRecord) {

	log.Printf("processSellTradeResp, T0:%v, T1:%v, N:%v, settleType:%v\n", T0, T1, N, settleType)

	switch settleType {

	case settleTypeT0:

		if T0 >= 0 && T1 >= T0 {

			var newT0, newT1 float64

			if T1-N >= 0 {
				newT0 = math.Max(0, T0-N)
				newT1 = T1 - N
			} else {
				newT1 = T1 - N
				newT0 = newT1
			}

			positionRecord.NetPositionT0 = newT0
			positionRecord.NetPositionT1 = newT1

			a.calculateAvgPrice(T1, newT1, tradeOrder, tradeActionResp, positionRecord)

			return
		}

		if T0 <= 0 && T1 < 0 {

			newT0 := T0 - N
			newT1 := T1 - N

			positionRecord.NetPositionT0 = newT0
			positionRecord.NetPositionT1 = newT1

			a.calculateAvgPrice(T1, newT1, tradeOrder, tradeActionResp, positionRecord)

			return
		}

	case settleTypeT1:

		if T0 >= 0 && T1 >= T0 {

			newT0 := math.Max(0, T0-N)
			newT1 := T1 - N

			positionRecord.NetPositionT0 = newT0
			positionRecord.NetPositionT1 = newT1

			a.calculateAvgPrice(T1, newT1, tradeOrder, tradeActionResp, positionRecord)

			return
		}

		if T0 <= 0 && T1 < 0 {

			newT0 := T0
			newT1 := T1 - N

			positionRecord.NetPositionT0 = newT0
			positionRecord.NetPositionT1 = newT1

			a.calculateAvgPrice(T1, newT1, tradeOrder, tradeActionResp, positionRecord)

			return
		}
	}
}

func (a *TitansFiccOrderPositionAdapter) processBuyTradeResp(T0 float64, T1 float64, N float64, settleType string, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, positionRecord *PositionRecord) {

	log.Printf("processBuyTradeResp, T0:%v, T1:%v, N:%v, settleType:%v\n", T0, T1, N, settleType)

	switch settleType {

	case settleTypeT0:

		if T0 >= 0 && T1 >= T0 {

			newT0 := T0 + N
			newT1 := T1 + N

			positionRecord.NetPositionT0 = newT0
			positionRecord.NetPositionT1 = newT1

			a.calculateAvgPrice(T1, newT1, tradeOrder, tradeActionResp, positionRecord)

			return
		}

		if T0 <= 0 && T1 < 0 {

			var newT0, newT1 float64

			if T1+N <= 0 {
				newT0 = math.Min(0, T0+N)
				newT1 = T1 + N
			} else {
				newT0 = T1 + N
				newT1 = T1 + N
			}

			positionRecord.NetPositionT0 = newT0
			positionRecord.NetPositionT1 = newT1

			a.calculateAvgPrice(T1, newT1, tradeOrder, tradeActionResp, positionRecord)

			return
		}

	case settleTypeT1:

		if T0 >= 0 && T1 >= T0 {

			newT0 := T0
			newT1 := T1 + N

			positionRecord.NetPositionT0 = newT0
			positionRecord.NetPositionT1 = newT1

			a.calculateAvgPrice(T1, newT1, tradeOrder, tradeActionResp, positionRecord)

			return
		}

		if T0 <= 0 && T1 < 0 {

			newT1 := T1 + N
			newT0 := math.Max(T0, math.Min(0, newT1))

			positionRecord.NetPositionT0 = newT0
			positionRecord.NetPositionT1 = newT1

			a.calculateAvgPrice(T1, newT1, tradeOrder, tradeActionResp, positionRecord)

			return
		}
	}
}

func (a *TitansFiccOrderPositionAdapter) initStrictCtpyMap() {
	strictCtpyMap := make(map[string]bool)
	configMap := a.applicationCfg.GetApplicationCfgItemMap()
	cfgItem, ok := configMap["StrictCtpy"]
	if ok {
		strictCtpys := strings.Split(cfgItem.ConfigItemValue, ",")
		for _, strictCtpy := range strictCtpys {
			strictCtpy = strings.TrimSpace(strictCtpy)
			if strictCtpy == "" {
				continue
			}
			strictCtpyMap[strictCtpy] = true
		}
	}
	a.strictCtpyMap = strictCtpyMap
}

func (a *TitansFiccOrderPositionAdapter) calculateAvgPrice(T1 float64, newT1 float64, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, positionRecord *PositionRecord) {

	log.Printf("calculateAvgPrice, T1:%v, newT1:%v, appOrdID:%v, symbol:%v\n", T1, newT1, tradeOrder.AppOrdID, tradeOrder.Symbol)

	if newT1 == 0 {

		positionRecord.LongCleanPriceCost = 0.0
		positionRecord.LongDirtyPriceCost = 0.0
		positionRecord.LongDirtyPriceWithFeeCost = 0.0
		positionRecord.ShortCleanPriceCost = 0.0
		positionRecord.ShortDirtyPriceCost = 0.0
		positionRecord.ShortDirtyPriceWithFeeCost = 0.0
		positionRecord.LongAvgCleanPrice = 0.0
		positionRecord.LongAvgDirtyPrice = 0.0
		positionRecord.LongAvgDirtyPriceWithFee = 0.0
		positionRecord.ShortAvgCleanPrice = 0.0
		positionRecord.ShortAvgDirtyPrice = 0.0
		positionRecord.ShortAvgDirtyPriceWithFee = 0.0

		return
	}

	cleanPrice, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "price", enum.AttrValueType_FLOAT)
	dirtyPrice, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "dirtyPrice", enum.AttrValueType_FLOAT)
	dirtyPriceWithFee, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "dirtyPriceWithFee", enum.AttrValueType_FLOAT)
	parValue, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "parValue", enum.AttrValueType_FLOAT)
	if parValue == 0 {
		parValue = 100.0
	}

	diff := newT1 - T1

	if newT1 > 0 {

		totalAmount := newT1 / parValue.(float64)

		if T1 >= 0 {

			diffAmount := diff / parValue.(float64)

			if diff >= 0 { // 多头增多

				if a.strictCtpyMap[tradeOrder.Account] {
					// 当多头增多时，价格需要取真实成交价格
					cleanPrice, dirtyPrice, dirtyPriceWithFee = a.getRealPrice(cleanPrice, dirtyPrice, dirtyPriceWithFee, tradeActionResp)
					log.Printf("calculateAvgPrice1, cleanPrice:%v, dirtyPrice:%v, dirtyPriceWithFee:%v, appOrdID:%v, symbol:%v\n", cleanPrice, dirtyPrice, dirtyPriceWithFee, tradeOrder.AppOrdID, tradeOrder.Symbol)
				}

				positionRecord.LongCleanPriceCost += diffAmount * cleanPrice.(float64)
				positionRecord.LongDirtyPriceCost += diffAmount * dirtyPrice.(float64)
				positionRecord.LongDirtyPriceWithFeeCost += diffAmount * dirtyPriceWithFee.(float64)
				positionRecord.LongAvgCleanPrice = positionRecord.LongCleanPriceCost / totalAmount
				positionRecord.LongAvgDirtyPrice = positionRecord.LongDirtyPriceCost / totalAmount
				positionRecord.LongAvgDirtyPriceWithFee = positionRecord.LongDirtyPriceWithFeeCost / totalAmount

				return

			} else { // 多头减少
				// 多头减少，即使要按真实成交价算，也要用之前的均价，这样可以确保卖出不影响均价
				positionRecord.LongCleanPriceCost += diffAmount * positionRecord.LongAvgCleanPrice
				positionRecord.LongDirtyPriceCost += diffAmount * positionRecord.LongAvgDirtyPrice
				positionRecord.LongDirtyPriceWithFeeCost += diffAmount * positionRecord.LongAvgDirtyPriceWithFee
				positionRecord.LongAvgCleanPrice = positionRecord.LongCleanPriceCost / totalAmount
				positionRecord.LongAvgDirtyPrice = positionRecord.LongDirtyPriceCost / totalAmount
				positionRecord.LongAvgDirtyPriceWithFee = positionRecord.LongDirtyPriceWithFeeCost / totalAmount

				return
			}

		} else { // 空头耗尽，多头产生

			if a.strictCtpyMap[tradeOrder.Account] {
				// 当多头增多时，价格需要取真实成交价格
				cleanPrice, dirtyPrice, dirtyPriceWithFee = a.getRealPrice(cleanPrice, dirtyPrice, dirtyPriceWithFee, tradeActionResp)
				log.Printf("calculateAvgPrice2, cleanPrice:%v, dirtyPrice:%v, dirtyPriceWithFee:%v, appOrdID:%v, symbol:%v\n", cleanPrice, dirtyPrice, dirtyPriceWithFee, tradeOrder.AppOrdID, tradeOrder.Symbol)
			}

			positionRecord.ShortCleanPriceCost = 0.0
			positionRecord.ShortDirtyPriceCost = 0.0
			positionRecord.ShortDirtyPriceWithFeeCost = 0.0
			positionRecord.ShortAvgCleanPrice = 0.0
			positionRecord.ShortAvgDirtyPrice = 0.0
			positionRecord.ShortAvgDirtyPriceWithFee = 0.0

			diffAmount := newT1 / parValue.(float64)

			positionRecord.LongCleanPriceCost += diffAmount * cleanPrice.(float64)
			positionRecord.LongDirtyPriceCost += diffAmount * dirtyPrice.(float64)
			positionRecord.LongDirtyPriceWithFeeCost += diffAmount * dirtyPriceWithFee.(float64)
			positionRecord.LongAvgCleanPrice = positionRecord.LongCleanPriceCost / totalAmount
			positionRecord.LongAvgDirtyPrice = positionRecord.LongDirtyPriceCost / totalAmount
			positionRecord.LongAvgDirtyPriceWithFee = positionRecord.LongDirtyPriceWithFeeCost / totalAmount
		}

		return
	}

	if newT1 < 0 {

		totalAmount := -newT1 / parValue.(float64)

		if T1 <= 0 {

			diffAmount := -diff / parValue.(float64)

			if diff <= 0 { // 空头增多

				if a.strictCtpyMap[tradeOrder.Account] {
					// 当空头增多时，价格需要取真实成交价格
					cleanPrice, dirtyPrice, dirtyPriceWithFee = a.getRealPrice(cleanPrice, dirtyPrice, dirtyPriceWithFee, tradeActionResp)
					log.Printf("calculateAvgPrice3, cleanPrice:%v, dirtyPrice:%v, dirtyPriceWithFee:%v, appOrdID:%v, symbol:%v\n", cleanPrice, dirtyPrice, dirtyPriceWithFee, tradeOrder.AppOrdID, tradeOrder.Symbol)
				}

				positionRecord.ShortCleanPriceCost += diffAmount * cleanPrice.(float64)
				positionRecord.ShortDirtyPriceCost += diffAmount * dirtyPrice.(float64)
				positionRecord.ShortDirtyPriceWithFeeCost += diffAmount * dirtyPriceWithFee.(float64)
				positionRecord.ShortAvgCleanPrice = positionRecord.ShortCleanPriceCost / totalAmount
				positionRecord.ShortAvgDirtyPrice = positionRecord.ShortDirtyPriceCost / totalAmount
				positionRecord.ShortAvgDirtyPriceWithFee = positionRecord.ShortDirtyPriceWithFeeCost / totalAmount

				return

			} else { // 空头减少
				// 空头减少，即使要按真实成交价算，也要用之前的均价，这样可以确保买入不影响均价
				positionRecord.ShortCleanPriceCost += diffAmount * positionRecord.ShortAvgCleanPrice
				positionRecord.ShortDirtyPriceCost += diffAmount * positionRecord.ShortAvgDirtyPrice
				positionRecord.ShortDirtyPriceWithFeeCost += diffAmount * positionRecord.ShortAvgDirtyPriceWithFee
				positionRecord.ShortAvgCleanPrice = positionRecord.ShortCleanPriceCost / totalAmount
				positionRecord.ShortAvgDirtyPrice = positionRecord.ShortDirtyPriceCost / totalAmount
				positionRecord.ShortAvgDirtyPriceWithFee = positionRecord.ShortDirtyPriceWithFeeCost / totalAmount

				return
			}

		} else { // 多头耗尽，空头产生

			if a.strictCtpyMap[tradeOrder.Account] {
				// 当空头增多时，价格需要取真实成交价格
				cleanPrice, dirtyPrice, dirtyPriceWithFee = a.getRealPrice(cleanPrice, dirtyPrice, dirtyPriceWithFee, tradeActionResp)
				log.Printf("calculateAvgPrice4, cleanPrice:%v, dirtyPrice:%v, dirtyPriceWithFee:%v, appOrdID:%v, symbol:%v\n", cleanPrice, dirtyPrice, dirtyPriceWithFee, tradeOrder.AppOrdID, tradeOrder.Symbol)
			}

			positionRecord.LongCleanPriceCost = 0.0
			positionRecord.LongDirtyPriceCost = 0.0
			positionRecord.LongDirtyPriceWithFeeCost = 0.0
			positionRecord.LongAvgCleanPrice = 0.0
			positionRecord.LongAvgDirtyPrice = 0.0
			positionRecord.LongAvgDirtyPriceWithFee = 0.0

			diffAmount := -newT1 / parValue.(float64)

			positionRecord.ShortCleanPriceCost += diffAmount * cleanPrice.(float64)
			positionRecord.ShortDirtyPriceCost += diffAmount * dirtyPrice.(float64)
			positionRecord.ShortDirtyPriceWithFeeCost += diffAmount * dirtyPriceWithFee.(float64)
			positionRecord.ShortAvgCleanPrice = positionRecord.ShortCleanPriceCost / totalAmount
			positionRecord.ShortAvgDirtyPrice = positionRecord.ShortDirtyPriceCost / totalAmount
			positionRecord.ShortAvgDirtyPriceWithFee = positionRecord.ShortDirtyPriceWithFeeCost / totalAmount
		}

		return
	}
}

func (a *TitansFiccOrderPositionAdapter) getRealPrice(cleanPrice, dirtyPrice, dirtyPriceWithFee interface{}, tradeActionResp *schema.TradeActionResp) (realCleanPrice, realDirtyPrice, realDirtyPriceWithFee interface{}) {

	respPrice, ok1, _ := attrutil.GetAttrValue(tradeActionResp.ExtendAttrMap, "respPrice", enum.AttrValueType_FLOAT)
	respDirtyPrice, ok2, _ := attrutil.GetAttrValue(tradeActionResp.ExtendAttrMap, "respDirtyPrice", enum.AttrValueType_FLOAT)
	respDirtyPriceWithFee, ok3, _ := attrutil.GetAttrValue(tradeActionResp.ExtendAttrMap, "respDirtyPriceWithFee", enum.AttrValueType_FLOAT)

	if ok1 && ok2 && ok3 {
		realCleanPrice = respPrice
		realDirtyPrice = respDirtyPrice
		realDirtyPriceWithFee = respDirtyPriceWithFee
	} else {
		realCleanPrice = cleanPrice
		realDirtyPrice = dirtyPrice
		realDirtyPriceWithFee = dirtyPriceWithFee
	}

	return
}
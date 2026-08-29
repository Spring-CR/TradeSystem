package ficc_fut

import (
	"encoding/json"
	"errors"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-core/schema"
	"rhino-core/types"
)

func (a *FiccFutOrderPositionAdapter) resetPostions(tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, positionRecord *PositionRecord) {
	if positionRecord.NetPosition >= 0 {
		maxLongAvailablePosition := positionRecord.NetPosition
		if positionRecord.LongAvailablePosition > maxLongAvailablePosition {
			positionRecord.LongAvailablePosition = maxLongAvailablePosition
		}
	} else {
		positionRecord.LongAvailablePosition = 0
	}

	if positionRecord.NetPosition <= 0 {
		maxShortAvailablePosition := -positionRecord.NetPosition
		if positionRecord.ShortAvailablePosition > maxShortAvailablePosition {
			positionRecord.ShortAvailablePosition = maxShortAvailablePosition
		}
	} else {
		positionRecord.ShortAvailablePosition = 0
	}

	js, _ := json.MarshalIndent(positionRecord, "", "  ")
	a.orderLog.Printf(tradeOrder, tradeActionResp, "[ResetPostions] PositionRecord=%s", js)
}

// 计算净持仓和持仓均价
func (a *FiccFutOrderPositionAdapter) AfterUpdateQuota(tradeResp *types.TradeActionRespReturn, metadata interface{}) {

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

	T := positionRecord.NetPosition
	N := float64(tradeActionResp.LastShares)

	a.orderLog.Printf(tradeOrder, tradeActionResp, "[CalculateNetPosition], T=%v, N=%v", T, N)

	switch tradeOrder.Side {
	case sideSell:
		a.processSellTradeResp(T, N, tradeOrder, tradeActionResp, positionRecord)
	case sideBuy:
		a.processBuyTradeResp(T, N, tradeOrder, tradeActionResp, positionRecord)
	}

	a.resetPostions(tradeOrder, tradeActionResp, positionRecord)

	a.marginExposureManager.UpdateOrder(tradeResp)
	a.marginExposureManager.CalculateMarginExposure(positionRecord, tradeOrder)

	js, _ := json.MarshalIndent(positionRecord, "", "  ")
	a.orderLog.Printf(tradeOrder, tradeActionResp, "[AfterCalculateNetPosition] PositionRecord=%s", js)
}

func (a *FiccFutOrderPositionAdapter) processSellTradeResp(T float64, N float64, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, positionRecord *PositionRecord) {

	log.Printf("processSellTradeResp, T:%v, N:%v\n", T, N)

	var newT float64

	newT = T - N

	positionRecord.NetPosition = newT

	a.calculateAvgPrice(T, newT, tradeOrder, tradeActionResp, positionRecord)

	return

}

func (a *FiccFutOrderPositionAdapter) processBuyTradeResp(T float64, N float64, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, positionRecord *PositionRecord) {

	log.Printf("processBuyTradeResp, T:%v, N:%v\n", T, N)

	newT := T + N

	positionRecord.NetPosition = newT

	a.calculateAvgPrice(T, newT, tradeOrder, tradeActionResp, positionRecord)

	return

}

func (a *FiccFutOrderPositionAdapter) calculateAvgPrice(T float64, newT float64, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, positionRecord *PositionRecord) {

	log.Printf("calculateAvgPrice, T:%v, newT:%v, appOrdID:%v, symbol:%v\n", T, newT, tradeOrder.AppOrdID, tradeOrder.Symbol)

	if newT == 0 {

		positionRecord.LongPriceCost = 0.0
		positionRecord.LongPriceWithFeeCost = 0.0
		positionRecord.LongPriceCNYWithFeeCost = 0.0
		positionRecord.ShortPriceCost = 0.0
		positionRecord.ShortPriceWithFeeCost = 0.0
		positionRecord.ShortPriceCNYWithFeeCost = 0.0
		positionRecord.LongAvgPrice = 0.0
		positionRecord.LongAvgPriceWithFee = 0.0
		positionRecord.ShortAvgPrice = 0.0
		positionRecord.ShortAvgPriceWithFee = 0.0

		return
	}

	// FICC期货业务要求使用真实的对冲价格
	// cleanPrice, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "price", enum.AttrValueType_FLOAT)
	// dirtyPrice, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "dirtyPrice", enum.AttrValueType_FLOAT)
	// dirtyPriceWithFee, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "dirtyPriceWithFee", enum.AttrValueType_FLOAT)
	// parValue, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "parValue", enum.AttrValueType_FLOAT)
	// if parValue == 0 {
	// 	parValue = 100.0
	// }

	diff := newT - T

	if newT > 0 {

		totalAmount := newT * positionRecord.ContractMultiplier

		if T >= 0 {

			diffAmount := diff * positionRecord.ContractMultiplier

			if diff >= 0 { // 多头增多

				// 当多头增多时，价格需要取真实成交价格
				price, priceWithFee, priceCNYWithFee := a.getRealPrice(tradeOrder, tradeActionResp)
				log.Printf("calculateAvgPrice1, price:%v, priceWithFee:%v, appOrdID:%v, symbol:%v\n", price, priceWithFee, tradeOrder.AppOrdID, tradeOrder.Symbol)

				positionRecord.LongPriceCost += diffAmount * price
				positionRecord.LongPriceWithFeeCost += diffAmount * priceWithFee
				positionRecord.LongPriceCNYWithFeeCost += diffAmount * priceCNYWithFee
				positionRecord.LongAvgPrice = positionRecord.LongPriceCost / totalAmount
				positionRecord.LongAvgPriceWithFee = positionRecord.LongPriceWithFeeCost / totalAmount

				return

			} else { // 多头减少
				// 多头减少，即使要按真实成交价算，也要用之前的均价，这样可以确保卖出不影响均价
				positionRecord.LongPriceCost += diffAmount * positionRecord.LongAvgPrice
				positionRecord.LongPriceWithFeeCost += diffAmount * positionRecord.LongAvgPriceWithFee
				positionRecord.LongPriceCNYWithFeeCost = positionRecord.LongPriceCNYWithFeeCost * newT / T
				positionRecord.LongAvgPrice = positionRecord.LongPriceCost / totalAmount
				positionRecord.LongAvgPriceWithFee = positionRecord.LongPriceWithFeeCost / totalAmount

				return
			}

		} else { // 空头耗尽，多头产生

			// 当多头增多时，价格需要取真实成交价格
			price, priceWithFee, priceCNYWithFee := a.getRealPrice(tradeOrder, tradeActionResp)
			log.Printf("calculateAvgPrice2, price:%v, priceWithFee:%v, appOrdID:%v, symbol:%v\n", price, priceWithFee, tradeOrder.AppOrdID, tradeOrder.Symbol)

			positionRecord.ShortPriceCost = 0.0
			positionRecord.ShortPriceWithFeeCost = 0.0
			positionRecord.ShortPriceCNYWithFeeCost = 0.0
			positionRecord.ShortAvgPrice = 0.0
			positionRecord.ShortAvgPriceWithFee = 0.0

			diffAmount := newT * positionRecord.ContractMultiplier

			positionRecord.LongPriceCost += diffAmount * price
			positionRecord.LongPriceWithFeeCost += diffAmount * priceWithFee
			positionRecord.LongPriceCNYWithFeeCost += diffAmount * priceCNYWithFee
			positionRecord.LongAvgPrice = positionRecord.LongPriceCost / totalAmount
			positionRecord.LongAvgPriceWithFee = positionRecord.LongPriceWithFeeCost / totalAmount
		}

		return
	}

	if newT < 0 {

		totalAmount := -newT * positionRecord.ContractMultiplier

		if T <= 0 {

			diffAmount := -diff * positionRecord.ContractMultiplier

			if diff <= 0 { // 空头增多

				// 当空头增多时，价格需要取真实成交价格
				price, priceWithFee, priceCNYWithFee := a.getRealPrice(tradeOrder, tradeActionResp)
				log.Printf("calculateAvgPrice3, price:%v, priceWithFee:%v, appOrdID:%v, symbol:%v\n", price, priceWithFee, tradeOrder.AppOrdID, tradeOrder.Symbol)

				positionRecord.ShortPriceCost += diffAmount * price
				positionRecord.ShortPriceWithFeeCost += diffAmount * priceWithFee
				positionRecord.ShortPriceCNYWithFeeCost += diffAmount * priceCNYWithFee
				positionRecord.ShortAvgPrice = positionRecord.ShortPriceCost / totalAmount
				positionRecord.ShortAvgPriceWithFee = positionRecord.ShortPriceWithFeeCost / totalAmount

				return

			} else { // 空头减少
				// 空头减少，即使要按真实成交价算，也要用之前的均价，这样可以确保买入不影响均价
				positionRecord.ShortPriceCost += diffAmount * positionRecord.ShortAvgPrice
				positionRecord.ShortPriceWithFeeCost += diffAmount * positionRecord.ShortAvgPriceWithFee
				positionRecord.ShortPriceCNYWithFeeCost = positionRecord.ShortPriceCNYWithFeeCost * newT / T
				positionRecord.ShortAvgPrice = positionRecord.ShortPriceCost / totalAmount
				positionRecord.ShortAvgPriceWithFee = positionRecord.ShortPriceWithFeeCost / totalAmount

				return
			}

		} else { // 多头耗尽，空头产生

			// 当空头增多时，价格需要取真实成交价格
			price, priceWithFee, priceCNYWithFee := a.getRealPrice(tradeOrder, tradeActionResp)
			log.Printf("calculateAvgPrice4, price:%v, priceWithFee:%v, appOrdID:%v, symbol:%v\n", price, priceWithFee, tradeOrder.AppOrdID, tradeOrder.Symbol)

			positionRecord.LongPriceCost = 0.0
			positionRecord.LongPriceWithFeeCost = 0.0
			positionRecord.LongPriceCNYWithFeeCost = 0.0
			positionRecord.LongAvgPrice = 0.0
			positionRecord.LongAvgPriceWithFee = 0.0

			diffAmount := -newT * positionRecord.ContractMultiplier

			positionRecord.ShortPriceCost += diffAmount * price
			positionRecord.ShortPriceWithFeeCost += diffAmount * priceWithFee
			positionRecord.ShortPriceCNYWithFeeCost += diffAmount * priceCNYWithFee
			positionRecord.ShortAvgPrice = positionRecord.ShortPriceCost / totalAmount
			positionRecord.ShortAvgPriceWithFee = positionRecord.ShortPriceWithFeeCost / totalAmount
		}

		return
	}
}

func (a *FiccFutOrderPositionAdapter) getRealPrice(tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp) (price, priceWithFee, priceCNYWithFee float64) {

	price = tradeActionResp.LastPx
	_priceWithFee, _, _ := attrutil.GetAttrValue(tradeActionResp.ExtendAttrMap, "respPriceWithFee", enum.AttrValueType_FLOAT)
	priceWithFee = _priceWithFee.(float64)

	if priceWithFee <= 0 {

		priceWithFee = price
	}

	// 转为本币币种
	exchangeRateCNY, ok, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "exchangeRateCNY", enum.AttrValueType_FLOAT)
	if !ok {
		exchangeRateCNY = 1.0
	}

	priceCNYWithFee = priceWithFee * exchangeRateCNY.(float64)

	return
}


// 废弃取本币币种价格
func (a *FiccFutOrderPositionAdapter) getRealPrice2(tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp) (price, priceWithFee float64) {

	price = tradeActionResp.LastPx
	_priceWithFee, _, _ := attrutil.GetAttrValue(tradeActionResp.ExtendAttrMap, "respPriceWithFee", enum.AttrValueType_FLOAT)
	priceWithFee = _priceWithFee.(float64)

	if priceWithFee <= 0 {

		priceWithFee = price
	}

	// 转为本币币种
	exchangeRateCNY, ok, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "exchangeRateCNY", enum.AttrValueType_FLOAT)
	if !ok {
		exchangeRateCNY = 1.0
	}

	price *= exchangeRateCNY.(float64)
	priceWithFee *= exchangeRateCNY.(float64)

	return
}

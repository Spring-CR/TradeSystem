package ficc

import (
	"encoding/json"
	"errors"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-core/schema"
	"rhino-core/types"
	"strings"

	"github.com/manucorporat/try"
)

func (a *TitansFiccOrderStatusAdapter) initStrictCtpyMap() {
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

func (a *TitansFiccOrderStatusAdapter) setRespPriceInfo(tradeActionResp *schema.TradeActionResp, tradeActionLatestResp *schema.TradeActionLatestResp, tradeActionRespList []*schema.TradeActionResp, tradeOrder *schema.TradeOrder, traceableTradeOrder *types.TraceableTradeOrder) {

	tradeActionResp.ExtendAttrMap["respPrice"] = tradeOrder.Price
	tradeActionResp.ExtendAttrMap["respDirtyPrice"] = tradeOrder.ExtendAttrMap["dirtyPrice"]
	tradeActionResp.ExtendAttrMap["respYtm"] = tradeOrder.ExtendAttrMap["ytm"]

	if !a.strictCtpyMap[tradeOrder.Account] {
		return
	}

	settlType, ok, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "settlType", enum.AttrValueType_STRING)
	if !ok {
		return
	}

	parValue, ok, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "parValue", enum.AttrValueType_INT)
	if !ok {
		return
	}

	valList, ok, _ := a.applicationCfg.GetAutoSyncRepo().Get("Security", tradeOrder.Symbol)
	if !ok || len(valList) == 0 {
		return
	}
	symbolData := valList[len(valList)-1]

	respPrice := tradeActionResp.LastPx
	respDirtyPrice, respYtm, err := a.ytmutil.ComputeDirtyPrice(a.applicationCfg, settlType.(string), tradeOrder, parValue.(int), respPrice, symbolData, true)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to ComputeDirtyPrice")
		return
	}

	tradeActionResp.ExtendAttrMap["respPrice"] = respPrice
	tradeActionResp.ExtendAttrMap["respDirtyPrice"] = respDirtyPrice
	tradeActionResp.ExtendAttrMap["respYtm"] = respYtm
}

// 成交的时候才调用的
func (a *TitansFiccOrderStatusAdapter) getOrderPrice(tradeOrder *schema.TradeOrder) (avgPrice, avgDirtyPrice, avgDirtyPriceWithFee float64) {

	avgPrice = tradeOrder.Price
	_avgDirtyPrice, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "dirtyPrice", enum.AttrValueType_FLOAT)
	avgDirtyPrice = _avgDirtyPrice.(float64)
	_dirtyPriceWithFee, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "dirtyPriceWithFee", enum.AttrValueType_FLOAT)
	avgDirtyPriceWithFee = _dirtyPriceWithFee.(float64)

	settlType, ok, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "settlType", enum.AttrValueType_STRING)
	if !ok {
		return
	}

	parValue, ok, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "parValue", enum.AttrValueType_INT)
	if !ok {
		return
	}

	commissionRate, ok, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "commissionRate", enum.AttrValueType_FLOAT)
	if !ok {
		return
	}

	valList, ok, _ := a.applicationCfg.GetAutoSyncRepo().Get("Security", tradeOrder.Symbol)
	if !ok || len(valList) == 0 {
		return
	}
	symbolData := valList[len(valList)-1]

	log.Printf("ComputeDirtyPrice, settlType:%v, parValue:%v, tradeOrder.AvgPx:%v, symbolData:%v\n", settlType, parValue, tradeOrder.AvgPx, symbolData)
	var err error
	avgDirtyPrice, _, err = a.ytmutil.ComputeDirtyPrice(a.applicationCfg, settlType.(string), tradeOrder, parValue.(int), tradeOrder.AvgPx, symbolData, false)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to ComputeDirtyPrice")
		return
	}

	avgPrice = tradeOrder.AvgPx

	if tradeOrder.Side == "1" { // 买单
		avgDirtyPriceWithFee = avgDirtyPrice * (1 + commissionRate.(float64))
	} else { // 卖单
		avgDirtyPriceWithFee = avgDirtyPrice * (1 - commissionRate.(float64))
	}

	return
}

func (a *TitansFiccOrderStatusAdapter) setOrderPriceInfo(tradeActionResp *schema.TradeActionResp, tradeActionLatestResp *schema.TradeActionLatestResp, tradeActionRespList []*schema.TradeActionResp, tradeOrder *schema.TradeOrder, traceableTradeOrder *types.TraceableTradeOrder, orderUpdateAttributes map[string]interface{}) {

	// 发生交易时设置订单维度的成交信息
	if fillExecType[tradeActionResp.ExecType] && a.strictCtpyMap[tradeOrder.Account] {
		avgPrice, avgDirtyPrice, avgDirtyPriceWithFee := a.getOrderPrice(tradeOrder)
		tradeOrder.ExtendAttrMap["avgPrice"] = avgPrice
		tradeOrder.ExtendAttrMap["avgDirtyPrice"] = avgDirtyPrice
		tradeOrder.ExtendAttrMap["avgDirtyPriceWithFee"] = avgDirtyPriceWithFee

		orderUpdateAttributes["avgPrice"] = avgPrice
		orderUpdateAttributes["avgDirtyPrice"] = avgDirtyPrice
		orderUpdateAttributes["avgDirtyPriceWithFee"] = avgDirtyPriceWithFee
	}

	_, ok, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "commissionRate1", enum.AttrValueType_FLOAT)
	if EndStatus[tradeActionResp.OrdStatus] && ok {
		cumQty := float64(tradeOrder.CumQty)
		if cumQty == 0 {
			return
		}

		var tradeAction *types.TraceableTradeActionResp
		try.This(func() {
			tradeAction = traceableTradeOrder.GetTraceableTradeActionRespByRootClOrdID(tradeActionLatestResp.RootClOrdID)
		}).Catch(func(err try.E) {})
		if tradeAction == nil {
			errMsg := "fail to get tradeAction by RootClOrdID:" + tradeActionLatestResp.RootClOrdID
			domain_error.ProcessSevereError(false, 0, nil, errors.New(errMsg), errMsg)
			return
		}

		var tradeActionRespListForTrade []*schema.TradeActionResp
		try.This(func() {
			tradeActionRespListForTrade = tradeAction.GetTradeActionRespListWithoutLock()
		}).Catch(func(err try.E) {})
		if tradeActionRespListForTrade == nil {
			errMsg := "fail to get tradeActionRespListForTrade"
			domain_error.ProcessSevereError(false, 0, nil, errors.New(errMsg), errMsg)
			return
		}

		commissionRate := 0.0
		for _, tradeActionResp := range tradeActionRespListForTrade {
			js, _ := json.Marshal(tradeActionResp.ExtendAttrMap)
			log.Printf("===>tradeActionResp:%s, tradeActionResp.LastShares:%v, tradeActionResp.ExtendAttrMap:%s\n", tradeActionResp.ExecID, tradeActionResp.LastShares, js)
			if tradeActionResp.LastShares <= 0 {
				continue
			}
			respCommissionRate, _, _ := attrutil.GetAttrValue(tradeActionResp.ExtendAttrMap, "respCommissionRate", enum.AttrValueType_FLOAT)
			commissionRate += respCommissionRate.(float64) * float64(tradeActionResp.LastShares)
		}
		commissionRate /= cumQty

		log.Printf("1.reconfig tradeOrder ExtendAttrMap for tradeOrder %s when tradeOrder is ended! commissionRate:%v\n", tradeOrder.ClOrdID, commissionRate)

		avgDirtyPrice, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "avgDirtyPrice", enum.AttrValueType_FLOAT)
		if tradeOrder.Side == "1" { // 买单
			tradeOrder.ExtendAttrMap["commissionRate"] = commissionRate
			tradeOrder.ExtendAttrMap["avgDirtyPriceWithFee"] = avgDirtyPrice.(float64) * (1 + commissionRate)
		} else { // 卖单
			tradeOrder.ExtendAttrMap["commissionRate"] = commissionRate
			tradeOrder.ExtendAttrMap["avgDirtyPriceWithFee"] = avgDirtyPrice.(float64) * (1 - commissionRate)
		}

		log.Printf("2.reconfig tradeOrder ExtendAttrMap for tradeOrder %s when tradeOrder is ended! commissionRate:%v\n", tradeOrder.ClOrdID, tradeOrder.ExtendAttrMap["commissionRate"])

		orderUpdateAttributes["commissionRate"] = tradeOrder.ExtendAttrMap["commissionRate"]
		orderUpdateAttributes["avgDirtyPriceWithFee"] = tradeOrder.ExtendAttrMap["avgDirtyPriceWithFee"]
	}
}

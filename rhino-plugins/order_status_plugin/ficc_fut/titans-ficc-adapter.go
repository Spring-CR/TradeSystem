package ficc

import (
	"encoding/json"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/types"
	"time"

	"github.com/manucorporat/try"
)

var (
	filledStatus          = string(enum.OrdStatus_Filled)
	partiallyFilledStatus = string(enum.OrdFillStatus_PartiallyFilled)
)

type FiccFutOrderStatusAdapter struct {
	applicationCfg *domain_cfg.ApplicationCfg
}

func NewFiccFutOrderStatusAdapter(c *domain_cfg.ApplicationCfg) (adapter *FiccFutOrderStatusAdapter, de *domain_error.Error) {
	log.Printf("construct FiccFutOrderStatusAdapter...")

	adapter = &FiccFutOrderStatusAdapter{applicationCfg: c}

	return
}

func (a *FiccFutOrderStatusAdapter) LookUpTraceableTradeActionResp(directOrderMap map[string]*types.TraceableTradeOrder, tradeActionResp *schema.TradeActionResp, tradeActionRespMap map[string]*types.TraceableTradeActionResp) (traceableTradeActionResp *types.TraceableTradeActionResp, ok bool) {

	data, _ := json.Marshal(tradeActionResp)
	log.Printf("LookUpTraceableTradeActionResp, clOrdID=%s, tradeActionResp=%s, appOrdID=%s\n", tradeActionResp.ClOrdID, data, tradeActionResp.AppOrdID)

	clOrdID := tradeActionResp.ClOrdID
	// 如果是撤单相关的执行回报，或者撤单拒绝回报
	if (!fillExecType[tradeActionResp.ExecType] && tradeActionResp.OrdStatus == string(enum.OrdStatus_PendingCancel)) || tradeActionResp.OrdStatus == string(enum.OrdStatus_Canceled) || len(tradeActionResp.CxlRejResponseTo) > 0 {

		log.Printf("receved cancel response, clOrdID:%s\n", clOrdID)
		if clOrdID == "" {
			clOrdID = tradeActionResp.OrigClOrdID
			log.Printf("original clOrdID is empty, reset as %s\n", clOrdID)
		}

		traceableTradeAction := tradeActionRespMap[clOrdID]
		log.Printf("finish get traceableTradeAction, traceableTradeAction == nil :%v\n", traceableTradeAction == nil)
		if traceableTradeAction == nil {
			return
		}

		traceableTradeOrder := getOrder(traceableTradeAction.GetTradeOrder().AppOrdID, directOrderMap)
		log.Printf("finish get traceableTradeOrder, traceableTradeOrder == nil :%v\n", traceableTradeOrder == nil)
		if traceableTradeOrder == nil {
			return
		}

		traceableTradeOrder.GetLock().Lock()
		defer traceableTradeOrder.GetLock().Unlock()

		log.Printf("success found order, clOrdID:%s, ordStatus:%s\n", traceableTradeOrder.GetBasicInfo().ClOrdID, traceableTradeOrder.GetBasicInfo().OrdStatus)

		actions := traceableTradeOrder.GetTradeActionsWithoutLock()
		for i := len(actions) - 1; i >= 0; i-- {
			action := actions[i]
			if action.GetTradeActionLatestResp().ActionType == string(enum.ActionType_Withdraw) {

				log.Printf("found order cancel action:%s\n", tradeActionResp.ClOrdID)

				traceableTradeActionResp = action
				ok = true
				if tradeActionResp.ClOrdID == "" {
					// 重新设置ClOrdID
					tradeActionResp.ClOrdID = traceableTradeActionResp.GetTradeActionLatestResp().ClOrdID
					log.Printf("original ClOrdID is empty, reset as %s\n", tradeActionResp.ClOrdID)
				}
				if tradeActionResp.OrdStatus == "" {
					tradeActionResp.OrdStatus = traceableTradeActionResp.GetTradeOrder().OrdStatus2
					log.Printf("original OrdStatus is empty, reset as %s\n", tradeActionResp.OrdStatus)
				}

				return
			}
		}

		log.Printf("cannot not found order cancel action for %s, check traceableTradeAction!=nil: %v\n", clOrdID, traceableTradeAction != nil)
		if traceableTradeAction != nil {
			return traceableTradeAction, true
		}

		return
	}

	// log.Printf("LookUpTraceableTradeActionResp, tradeActionRespMap.Len=%d\n", len(tradeActionRespMap))
	// for _, v := range tradeActionRespMap {
	// 	log.Printf("clOrdID=%s\n", v.GetTradeActionLatestResp().ClOrdID)
	// }
	if clOrdID == "" && tradeActionResp.AppOrdID != "" {
		// 待审核等草稿态时会走到这一步
		traceableTradeOrder := getOrder(tradeActionResp.AppOrdID, directOrderMap)
		if traceableTradeOrder == nil {
			return
		}
		traceableTradeActionResp = traceableTradeOrder.GetLatestTraceableTradeActionRespByAppOrdID(tradeActionResp.AppOrdID)
		if traceableTradeActionResp != nil {
			ok = true
		}
		return
	}

	if clOrdID == "" && tradeActionResp.AppOrdID == "" {
		return
	}

	traceableTradeActionResp, ok = tradeActionRespMap[clOrdID]
	return
}

var (
	fillExecType = map[string]bool{
		"F": true,
		"1": true,
		"2": true,
	}
)

func (a *FiccFutOrderStatusAdapter) UpdateTradeActionLatestResp(tradeActionResp *schema.TradeActionResp, traceableTradeActionResp *types.TraceableTradeActionResp) {

	// 这里不需要了，因为star会正常的返回pandingcancel了
	//if tradeActionResp.OrdStatus == string(enum.OrdStatus_New) && (traceableTradeActionResp.GetTradeOrder().OrdStatus != string(enum.OrdStatus_PendingNew) && traceableTradeActionResp.GetTradeOrder().OrdStatus != "") {
	//	tradeActionResp.OrdStatus = string(enum.OrdStatus_PendingCancel)
	//}

	// 更新latestresp
	traceableTradeActionResp.GetTradeActionLatestResp().OrderID = tradeActionResp.OrderID
	traceableTradeActionResp.GetTradeActionLatestResp().OrigClOrdID = tradeActionResp.OrigClOrdID
	traceableTradeActionResp.GetTradeActionLatestResp().ExecID = tradeActionResp.ExecID
	traceableTradeActionResp.GetTradeActionLatestResp().ExecType = tradeActionResp.ExecType
	traceableTradeActionResp.GetTradeActionLatestResp().OrdStatus = tradeActionResp.OrdStatus
	traceableTradeActionResp.GetTradeActionLatestResp().OrdRejReason = tradeActionResp.OrdRejReason
	traceableTradeActionResp.GetTradeActionLatestResp().CxlRejResponseTo = tradeActionResp.CxlRejResponseTo
	traceableTradeActionResp.GetTradeActionLatestResp().ExecRestatementReason = tradeActionResp.ExecRestatementReason

	// 来自Order的基础信息
	traceableTradeActionResp.GetTradeActionLatestResp().Account = traceableTradeActionResp.GetTradeOrder().Account
	traceableTradeActionResp.GetTradeActionLatestResp().Symbol = traceableTradeActionResp.GetTradeOrder().Symbol
	traceableTradeActionResp.GetTradeActionLatestResp().SymbolSfx = traceableTradeActionResp.GetTradeOrder().SymbolSfx
	traceableTradeActionResp.GetTradeActionLatestResp().SecurityID = traceableTradeActionResp.GetTradeOrder().SecurityID
	traceableTradeActionResp.GetTradeActionLatestResp().IDSource = traceableTradeActionResp.GetTradeOrder().IDSource
	traceableTradeActionResp.GetTradeActionLatestResp().SecurityType = traceableTradeActionResp.GetTradeOrder().SecurityType
	traceableTradeActionResp.GetTradeActionLatestResp().Side = traceableTradeActionResp.GetTradeOrder().Side
	traceableTradeActionResp.GetTradeActionLatestResp().OrderQty = traceableTradeActionResp.GetTradeOrder().OrderQty
	traceableTradeActionResp.GetTradeActionLatestResp().CashOrderQty = traceableTradeActionResp.GetTradeOrder().CashOrderQty
	traceableTradeActionResp.GetTradeActionLatestResp().OrdType = traceableTradeActionResp.GetTradeOrder().OrdType
	traceableTradeActionResp.GetTradeActionLatestResp().Price = traceableTradeActionResp.GetTradeOrder().Price
	traceableTradeActionResp.GetTradeActionLatestResp().Currency = traceableTradeActionResp.GetTradeOrder().Currency // order里不一定有币种信息
	traceableTradeActionResp.GetTradeActionLatestResp().EffectiveTime = tradeActionResp.EffectiveTime
	traceableTradeActionResp.GetTradeActionLatestResp().ExpireTime = tradeActionResp.ExpireTime

	// 默认的处理成交信息的逻辑
	if tradeActionResp.LastShares > 0 {
		traceableTradeActionResp.GetTradeActionLatestResp().LastShares = tradeActionResp.LastShares
	}
	if tradeActionResp.LastPx > 0 {
		traceableTradeActionResp.GetTradeActionLatestResp().LastPx = tradeActionResp.LastPx
	}

	// 时序信息
	traceableTradeActionResp.GetTradeActionLatestResp().TransactTime = tradeActionResp.TransactTime
	traceableTradeActionResp.GetTradeActionLatestResp().MsgTime = tradeActionResp.MsgTime
	traceableTradeActionResp.GetTradeActionLatestResp().MsgSeq = tradeActionResp.MsgSeq

	tradeActionLatestResp := traceableTradeActionResp.GetTradeActionLatestResp()
	order := traceableTradeActionResp.GetTradeOrder()

	// 处理状态和交易信息
	if fillExecType[tradeActionResp.ExecType] {

		var contractMultiplier float64 = 1.0
		extendAttrMap := traceableTradeActionResp.GetTradeOrder().ExtendAttrMap
		val, ok, _ := attrutil.GetAttrValue(extendAttrMap, "contractMultiplier", enum.AttrValueType_FLOAT)
		if ok {
			contractMultiplier = val.(float64)
		}

		ordQty := int64(order.OrderQty)
		log.Printf("===> ordQty:%v\n", ordQty)
		log.Printf("===> ord.CumQty:%v\n", order.CumQty)
		log.Printf("===> ord.LeavesQty:%v\n", order.LeavesQty)
		log.Printf("===> ord.LastPx:%v\n", order.LastPx)
		log.Printf("===> ord.LastShares:%v\n", order.LastShares)
		log.Printf("===> ord.AvgPx:%v\n", order.AvgPx)

		log.Printf("===> tradeActionLatestResp.CumQty:%v\n", tradeActionLatestResp.CumQty)
		log.Printf("===> tradeActionLatestResp.LeavesQty:%v\n", tradeActionLatestResp.LeavesQty)
		log.Printf("===> tradeActionLatestResp.LastPx:%v\n", tradeActionLatestResp.LastPx)
		log.Printf("===> tradeActionLatestResp.LastShares:%v\n", tradeActionLatestResp.LastShares)
		log.Printf("===> tradeActionLatestResp.AvgPx:%v\n", tradeActionLatestResp.AvgPx)

		tradeActionLatestResp.CumQty += tradeActionResp.LastShares
		tradeActionLatestResp.LeavesQty = ordQty - tradeActionLatestResp.CumQty

		totalAmount := order.AvgPx*float64(order.CumQty)*contractMultiplier + float64(tradeActionResp.LastShares)*contractMultiplier*tradeActionResp.LastPx
		tradeActionLatestResp.AvgPx = totalAmount / (float64(tradeActionLatestResp.CumQty) * contractMultiplier)

		if tradeActionLatestResp.LeavesQty <= 0 {
			tradeActionLatestResp.OrdStatus = filledStatus
			tradeActionResp.OrdStatus = filledStatus
		} else {
			tradeActionLatestResp.OrdStatus = partiallyFilledStatus
			tradeActionResp.OrdStatus = partiallyFilledStatus
		}
		tradeActionResp.CumQty = tradeActionLatestResp.CumQty
		tradeActionResp.AvgPx = tradeActionLatestResp.AvgPx
		tradeActionResp.LeavesQty = tradeActionLatestResp.LeavesQty
	} else {
		tradeActionLatestResp.CumQty = order.CumQty
		tradeActionLatestResp.AvgPx = order.AvgPx
		switch tradeActionResp.OrdStatus {
		case string(enum.OrdStatus_PendingNew), string(enum.OrdStatus_New):
			tradeActionLatestResp.LeavesQty = int64(order.OrderQty)
		case string(enum.OrdStatus_PendingCancel), string(enum.OrdStatus_PendingReplace):
			log.Printf("===>order.clOrdID:%v, order.LeavesQty:%v\n", order.ClOrdID, order.LeavesQty)
			tradeActionLatestResp.LeavesQty = order.LeavesQty
		case string(enum.OrdStatus_Replaced):
			tradeActionLatestResp.LeavesQty = int64(order.OrderQty - float64(order.CumQty))
			if tradeActionLatestResp.LeavesQty < 0 {
				tradeActionLatestResp.LeavesQty = 0
			}
		default:
			tradeActionLatestResp.LeavesQty = 0
		}
		tradeActionResp.CumQty = tradeActionLatestResp.CumQty
		tradeActionResp.AvgPx = tradeActionLatestResp.AvgPx
		tradeActionResp.LeavesQty = tradeActionLatestResp.LeavesQty
	}

	log.Printf("===> set tradeActionLatestResp, OrdStatus:%v, LastShares:%v, LastPx:%v, LeavesQty:%v, CumQty:%v, AvgPx:%v, tradeActionResp.ClOrdID:%v, tradeActionResp.OrigClOrdID:%v\n", tradeActionLatestResp.OrdStatus, tradeActionLatestResp.LastShares, tradeActionLatestResp.LastPx, tradeActionLatestResp.LeavesQty, tradeActionLatestResp.CumQty, tradeActionLatestResp.AvgPx, tradeActionResp.ClOrdID, tradeActionResp.OrigClOrdID)
}

func (a *FiccFutOrderStatusAdapter) UpdateOrderStatusForOrderCancelReject(tradeActionResp *schema.TradeActionResp, tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder) {
	order.OrdStatus = tradeActionLatestResp.OrdStatus
}

var (
	EndStatus = map[string]bool{
		string(enum.OrdStatus_Filled):     true,
		string(enum.OrdStatus_DoneForDay): true,
		string(enum.OrdStatus_Canceled):   true,
		string(enum.OrdStatus_Stopped):    true,
		string(enum.OrdStatus_Rejected):   true,
		string(enum.OrdStatus_Suspended):  true,
		string(enum.OrdStatus_Expired):    true,
	}
	DuringTradingStatus = map[string]bool{
		string(enum.OrdStatus_New):                true,
		string(enum.OrdStatus_PartiallyFilled):    true,
		string(enum.OrdStatus_Replaced):           true,
		string(enum.OrdStatus_PendingCancel):      true,
		string(enum.OrdStatus_PendingNew):         true,
		string(enum.OrdStatus_Calculated):         true,
		string(enum.OrdStatus_AcceptedForBidding): true,
		string(enum.OrdStatus_PendingReplace):     true,
	}
	TradedStatus = map[string]bool{
		string(enum.OrdStatus_PartiallyFilled): true,
		string(enum.OrdStatus_Filled):          true,
	}
)

// func (a *FiccFutOrderStatusAdapter) setRespExtendAttrs(tradeActionResp *schema.TradeActionResp, tradeOrder *schema.TradeOrder) {
// 	commissionRate, _, _ := attrutil.GetAttrValue(tradeActionResp.ExtendAttrMap, "respCommissionRate", enum.AttrValueType_FLOAT)
// 	dirtyPrice, _, _ := attrutil.GetAttrValue(tradeActionResp.ExtendAttrMap, "respDirtyPrice", enum.AttrValueType_FLOAT)

//		log.Printf("resetRespCommission, side:%v, dirtyPrice:%v, commissionRate:%v, appOrdID:%v\n", tradeOrder.Side, dirtyPrice, commissionRate, tradeOrder.AppOrdID)
//		if tradeOrder.Side == "1" { // 买单
//			tradeActionResp.ExtendAttrMap["respDirtyPriceWithFee"] = dirtyPrice.(float64) * (1 + commissionRate.(float64))
//		} else { // 卖单
//			tradeActionResp.ExtendAttrMap["respDirtyPriceWithFee"] = dirtyPrice.(float64) * (1 - commissionRate.(float64))
//		}
//	}

func (a *FiccFutOrderStatusAdapter) getPriceWithFee(price float64, tradeOrder *schema.TradeOrder) float64 {

	commissionValue, ok1, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "commissionValue", enum.AttrValueType_FLOAT)
	commissionType, ok2, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "commissionType", enum.AttrValueType_STRING)

	if !ok1 || !ok2 {
		return price
	}

	oldPrice := price

	if tradeOrder.Side == "1" { // 买单

		switch commissionType {
		case "FEE_BY_RATE":
			price = price * (1 + commissionValue.(float64))
		case "FEE_BY_PER_SHARE":
			price = price + commissionValue.(float64)
		}

	} else { // 卖单

		switch commissionType {
		case "FEE_BY_RATE":
			price = price * (1 - commissionValue.(float64))
		case "FEE_BY_PER_SHARE":
			price = price - commissionValue.(float64)
		}
	}

	if price < 0 {
		price = oldPrice
	} 

	return price
}

func (a *FiccFutOrderStatusAdapter) setRespExtendAttrs(tradeActionResp *schema.TradeActionResp, tradeOrder *schema.TradeOrder) {

	respPriceWithFee := a.getPriceWithFee(tradeActionResp.LastPx, tradeOrder)

	log.Printf("setRespExtendAttrs, side:%v, appOrdID:%v, rawPrice:%v, newPrice:%v, execID:%v\n", tradeOrder.Side, tradeOrder.AppOrdID, tradeActionResp.LastPx, respPriceWithFee, tradeActionResp.ExecID)

	tradeActionResp.ExtendAttrMap["respPriceWithFee"] = respPriceWithFee
}

func (a *FiccFutOrderStatusAdapter) UpdateOrderStatusForExecutionReport(tradeActionResp *schema.TradeActionResp, tradeActionLatestResp *schema.TradeActionLatestResp, tradeActionRespList []*schema.TradeActionResp, order *schema.TradeOrder, traceableTradeOrder *types.TraceableTradeOrder) (orderUpdateAttributes map[string]interface{}) {

	orderUpdateAttributes = map[string]interface{}{}
	// stars的原生回报里，成交不带订单状态
	// 当订单状态处于pending，并且收到部分成交回报的时候，不能更改订单信息，其他情况需要更改订单信息
	//log.Printf("tradeActionResp, ExecID=%s, order.OrdStatus=%s, tradeActionResp.ExecType=%s, order.LeavesQty=%v\n", tradeActionResp.ExecID, order.OrdStatus, tradeActionResp.ExecType, order.LeavesQty)
	//if order.OrdStatus != string(enum.OrdStatus_PendingCancel) || !fillExecType[tradeActionResp.ExecType] || order.LeavesQty <= 0 {

	// bugfix: 剩余成交量，要用tradeActionLatestResp的值
	log.Printf("tradeActionResp, ExecID=%s, order.OrdStatus=%s, tradeActionResp.ExecType=%s, tradeActionResp.LeavesQty=%v\n", tradeActionResp.ExecID, order.OrdStatus, tradeActionResp.ExecType, tradeActionResp.LeavesQty)
	if order.OrdStatus != string(enum.OrdStatus_PendingCancel) || !fillExecType[tradeActionResp.ExecType] || tradeActionResp.LeavesQty <= 0 {
		order.OrdStatus = tradeActionLatestResp.OrdStatus
	}
	log.Printf("==> new order.OrdStatus=%s\n", order.OrdStatus)
	//tradeActionResp.PendingCancel = order.OrdStatus == "6"

	if fillExecType[tradeActionResp.ExecType] {
		order.LastShares = tradeActionLatestResp.LastShares
		order.LastPx = tradeActionLatestResp.LastPx
		order.LeavesQty = tradeActionLatestResp.LeavesQty
		order.CumQty = tradeActionLatestResp.CumQty
		order.AvgPx = tradeActionLatestResp.AvgPx

		if tradeActionResp.ExtendAttrMap != nil {
			// 对于成交回报，设置其扩展属性
			a.setRespExtendAttrs(tradeActionResp, order)
		}

		// 设置费后订单均价
		avgPriceWithFee := a.getPriceWithFee(order.AvgPx, order)
		orderUpdateAttributes["avgPriceWithFee"] = avgPriceWithFee

	} else {

		if tradeActionResp.OrdStatus == "A" || tradeActionResp.OrdStatus == "0" {
			log.Printf("set order.LeavesQty=%v for order:%v\n", order.OrderQty, order.ClOrdID)
			order.LeavesQty = int64(order.OrderQty)
		}

		tradeActionLatestResp.LastShares = 0
		tradeActionLatestResp.LastPx = 0
		tradeActionLatestResp.LeavesQty = order.LeavesQty
		tradeActionLatestResp.CumQty = order.CumQty
		tradeActionLatestResp.AvgPx = order.AvgPx

		if EndStatus[order.OrdStatus] {
			tradeActionLatestResp.LeavesQty = 0
		}

		tradeActionResp.LastShares = tradeActionLatestResp.LastShares
		tradeActionResp.LastPx = tradeActionLatestResp.LastPx
		tradeActionResp.LeavesQty = tradeActionLatestResp.LeavesQty
		tradeActionResp.CumQty = tradeActionLatestResp.CumQty
		tradeActionResp.AvgPx = tradeActionLatestResp.AvgPx
	}

	order.OrdRejReason = tradeActionLatestResp.OrdRejReason

	// 20251217改，直接在这里设置tradeActionResp ExecType、OrdStatus、ExecTransType
	execType := tradeActionResp.OrdStatus
	tradeActionResp.ExecType = execType
	if order.OrdStatus == "6" && execType == "1" {
		tradeActionResp.OrdStatus = "6"
	}
	if execType == "6" || execType == "4" {
		tradeActionResp.ExecTransType = "1"
	} else {
		tradeActionResp.ExecTransType = "0"
	}

	if order.OrdStatus == "8" {
		domain_error.BuildWithDetails(domain_error.ERROR, order, domain_error.ORDER_REJECT_ERR_CODE, nil, order.AppOrdID, order.OrdRejReason)
	}

	return
}

func getOrder(clOrdID string, directOrderMap map[string]*types.TraceableTradeOrder) *types.TraceableTradeOrder {

	var order *types.TraceableTradeOrder

	errHappen := false
	i := 0

	for (errHappen || i == 0) && i < 10 {
		try.This(func() {
			order = directOrderMap[clOrdID]
		}).Catch(func(err try.E) {
			log.Printf("error occur while look up direct order, error:%v\n", err)
			errHappen = true
		})

		if !errHappen {
			return order
		}

		time.Sleep(15 * time.Millisecond)
		i++
	}

	domain_error.ProcessSevereError(false, 0, nil, fmt.Errorf("fail to get order after 10 times"), fmt.Sprintf("fail to get order after 10 times, clOrdID:%s\n", clOrdID))

	return nil
}

package ficc_v2

import (
	"math"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-core/schema"
	"rhino-core/types"
	"rhino-plugins/order_status_plugin/ficc"
)

type OrderLeftQtyWrapper struct {
	orderID    string
	price      float64
	leftQty    int64
	arrayIndex int
}

func NewOrderLeftQtyWrapper(tradeOrder *schema.TradeOrder) *OrderLeftQtyWrapper {
	inst := &OrderLeftQtyWrapper{}
	inst.orderID = tradeOrder.AppOrdID
	price, ok, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "dirtyPriceWithFee", enum.AttrValueType_FLOAT)
	if !ok {
		price = tradeOrder.Price
	}

	inst.price = price.(float64)
	inst.leftQty = int64(tradeOrder.OrderQty)
	return inst
}

func (m *OrderLeftQtyWrapper) GetOrderPrice() float64    { return m.price }
func (m *OrderLeftQtyWrapper) GetOrderID() string        { return m.orderID }
func (m *OrderLeftQtyWrapper) GetOrderLeftQty() int64    { return m.leftQty }
func (m *OrderLeftQtyWrapper) SetOrderLeftQty(qty int64) { m.leftQty = qty }
func (m *OrderLeftQtyWrapper) GetArrayIndex() int        { return m.arrayIndex }
func (m *OrderLeftQtyWrapper) SetArrayIndex(idx int)     { m.arrayIndex = idx }

type OrderLists struct {
	parValue      float64
	orderBuyList  *OrderLeftQtyManager
	orderSellList *OrderLeftQtyManager
}

func NewOrderLists(parValue float64) *OrderLists {
	orderLists := &OrderLists{parValue: parValue, orderBuyList: NewOrderLeftQtyManager(sideBuy, parValue), orderSellList: NewOrderLeftQtyManager(sideSell, parValue)}
	return orderLists
}

func (l *OrderLists) AddOrder(tradeOrder *schema.TradeOrder) {

	orderLeftQty := NewOrderLeftQtyWrapper(tradeOrder)

	switch tradeOrder.Side {
	case sideBuy:
		l.orderBuyList.AddQty(orderLeftQty)
	case sideSell:
		l.orderSellList.AddQty(orderLeftQty)
	}
}

func (l *OrderLists) UpdateOrder(tradeResp *types.TradeActionRespReturn) {

	var orderList *OrderLeftQtyManager
	tradeOrder := tradeResp.GetTradeOrder()
	switch tradeOrder.Side {
	case sideBuy:
		orderList = l.orderBuyList
	case sideSell:
		orderList = l.orderSellList
	}

	tradeActionResp := tradeResp.CurrentTradeActionResp

	if ficc.EndStatus[tradeActionResp.OrdStatus] {
		orderList.DeleteQty(tradeOrder.AppOrdID)
		return
	}

	if fillExecType[tradeActionResp.ExecType] && tradeActionResp.LastShares > 0 {
		orderList.UpdateQty(tradeOrder.AppOrdID, tradeActionResp.LeavesQty)
		return
	}
}

// 强制删除，比如校验不通过、下单失败时，直接删除占用额度
func (l *OrderLists) DeleteOrder(tradeOrder *schema.TradeOrder) {
	switch tradeOrder.Side {
	case sideBuy:
		l.orderBuyList.DeleteQty(tradeOrder.AppOrdID)
	case sideSell:
		l.orderSellList.DeleteQty(tradeOrder.AppOrdID)
	}
}

type MarginExposureManager struct {
	defaultLongMarginRatio  float64
	defaultShortMarginRatio float64
	orderListsMap           map[string]*OrderLists // key = 交易对手ID + appOrdID
}

func NewMarginExposureManager(longMarginRatio, shortMarginRatio float64) *MarginExposureManager {
	inst := &MarginExposureManager{defaultLongMarginRatio: longMarginRatio, defaultShortMarginRatio: shortMarginRatio, orderListsMap: make(map[string]*OrderLists)}
	return inst
}

func (m *MarginExposureManager) AddOrder(tradeOrder *schema.TradeOrder) {
	key := m.getKey(tradeOrder)
	orderLists, ok := m.orderListsMap[key]
	if !ok {
		orderLists = NewOrderLists(m.getParValue(tradeOrder))
		m.orderListsMap[key] = orderLists
	}
	orderLists.AddOrder(tradeOrder)
}

func (m *MarginExposureManager) DeleteOrder(tradeOrder *schema.TradeOrder) {
	key := m.getKey(tradeOrder)
	orderLists, ok := m.orderListsMap[key]
	if !ok {
		return
	}
	orderLists.DeleteOrder(tradeOrder)
}

func (m *MarginExposureManager) UpdateOrder(tradeResp *types.TradeActionRespReturn) {
	tradeOrder := tradeResp.GetTradeOrder()
	key := m.getKey(tradeOrder)
	orderLists, ok := m.orderListsMap[key]
	if !ok {
		m.AddOrder(tradeOrder)
		orderLists, ok = m.orderListsMap[key]
		if !ok {
			return
		}
	}
	orderLists.UpdateOrder(tradeResp)
}

func (m *MarginExposureManager) getKey(tradeOrder *schema.TradeOrder) string {
	return tradeOrder.Account + "_" + tradeOrder.Symbol
}

func (m *MarginExposureManager) getParValue(tradeOrder *schema.TradeOrder) float64 {
	parValue, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "parValue", enum.AttrValueType_FLOAT)
	if parValue.(float64) <= 0 {
		parValue = 100.0
	}
	return parValue.(float64)
}

// 计算初保占用
func (m *MarginExposureManager) CalculateMarginExposure(positionRecord *PositionRecord, tradeOrder *schema.TradeOrder) {

	var longMarginRatio, shortMarginRatio float64
	if positionRecord.LongMarginRatio > 0 {
		longMarginRatio = positionRecord.LongMarginRatio
	} else {
		longMarginRatio = m.defaultLongMarginRatio
	}
	if positionRecord.ShortMarginRatio > 0 {
		shortMarginRatio = positionRecord.ShortMarginRatio
	} else {
		shortMarginRatio = m.defaultShortMarginRatio
	}

	key := m.getKey(tradeOrder)
	orderLists, ok := m.orderListsMap[key]
	if !ok {
		return
	}

	// 获取净持仓（以T1净持仓为准）
	netPosition := positionRecord.NetPositionT1

	if netPosition >= 0 {

		positionRecord.MaxLongMarginOccupancy = positionRecord.LongDirtyPriceWithFeeCost*longMarginRatio + orderLists.orderBuyList.GetTotalLeftAmount()/orderLists.orderBuyList.GetParValue()*longMarginRatio
		if math.Abs(positionRecord.MaxLongMarginOccupancy) < 0.001 {
			positionRecord.MaxLongMarginOccupancy = 0.0
		}

		totalSellOrderLeftQty := orderLists.orderSellList.GetTotalLeftQty()
		if totalSellOrderLeftQty < int64(netPosition) {
			positionRecord.MaxShortMarginOccupancy = 0
		} else {
			positionRecord.MaxShortMarginOccupancy = orderLists.orderSellList.GetOversoldLeftAmount(netPosition) / orderLists.orderSellList.GetParValue() * shortMarginRatio
			if math.Abs(positionRecord.MaxShortMarginOccupancy) < 0.001 {
				positionRecord.MaxShortMarginOccupancy = 0.0
			}
		}

		positionRecord.MaxMarginOccupancy = math.Max(positionRecord.MaxLongMarginOccupancy, positionRecord.MaxShortMarginOccupancy)

		return

	} else {

		positionRecord.MaxShortMarginOccupancy = positionRecord.ShortDirtyPriceWithFeeCost*shortMarginRatio + orderLists.orderSellList.GetTotalLeftAmount()/orderLists.orderSellList.GetParValue()*shortMarginRatio
		if math.Abs(positionRecord.MaxShortMarginOccupancy) < 0.001 {
			positionRecord.MaxShortMarginOccupancy = 0.0
		}

		totalBuyOrderLeftQty := orderLists.orderBuyList.GetTotalLeftQty()
		if totalBuyOrderLeftQty < -int64(netPosition) {
			positionRecord.MaxLongMarginOccupancy = 0
		} else {
			positionRecord.MaxLongMarginOccupancy = orderLists.orderBuyList.GetOversoldLeftAmount(netPosition) / orderLists.orderBuyList.GetParValue() * longMarginRatio
			if math.Abs(positionRecord.MaxLongMarginOccupancy) < 0.001 {
				positionRecord.MaxLongMarginOccupancy = 0.0
			}
		}

		positionRecord.MaxMarginOccupancy = math.Max(positionRecord.MaxLongMarginOccupancy, positionRecord.MaxShortMarginOccupancy)

		return
	}
}

// 计算初保占用
func (m *MarginExposureManager) CalculateInitMarginExposure(positionRecord *PositionRecord) {

	var longMarginRatio, shortMarginRatio float64
	if positionRecord.LongMarginRatio > 0 {
		longMarginRatio = positionRecord.LongMarginRatio
	} else {
		longMarginRatio = m.defaultLongMarginRatio
	}
	if positionRecord.ShortMarginRatio > 0 {
		shortMarginRatio = positionRecord.ShortMarginRatio
	} else {
		shortMarginRatio = m.defaultShortMarginRatio
	}

	// 获取净持仓（以T1净持仓为准）
	netPosition := positionRecord.NetPositionT1

	if netPosition >= 0 {

		positionRecord.MaxLongMarginOccupancy = positionRecord.LongDirtyPriceWithFeeCost * longMarginRatio
		if math.Abs(positionRecord.MaxLongMarginOccupancy) < 0.001 {
			positionRecord.MaxLongMarginOccupancy = 0.0
		}

		positionRecord.MaxShortMarginOccupancy = 0

		positionRecord.MaxMarginOccupancy = positionRecord.MaxLongMarginOccupancy

		return

	} else {

		positionRecord.MaxShortMarginOccupancy = positionRecord.ShortDirtyPriceWithFeeCost * shortMarginRatio
		if math.Abs(positionRecord.MaxShortMarginOccupancy) < 0.001 {
			positionRecord.MaxShortMarginOccupancy = 0.0
		}

		positionRecord.MaxLongMarginOccupancy = 0

		positionRecord.MaxMarginOccupancy = positionRecord.MaxShortMarginOccupancy

		return
	}
}

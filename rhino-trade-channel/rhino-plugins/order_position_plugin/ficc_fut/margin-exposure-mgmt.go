package ficc_fut

import (
	"fmt"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/leftqty"
	"rhino-common/utils/logger"
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
	priceCNYWithFee, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "priceCNYWithFee", enum.AttrValueType_FLOAT)
	if priceCNYWithFee == 0 {
		priceCNYWithFee = tradeOrder.Price
	}
	inst.price = priceCNYWithFee.(float64)
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
	contractMultiplier float64
	orderBuyList       *leftqty.OrderLeftQtyManager
	orderSellList      *leftqty.OrderLeftQtyManager
}

func NewOrderLists(contractMultiplier float64) *OrderLists {
	orderLists := &OrderLists{contractMultiplier: contractMultiplier, orderBuyList: leftqty.NewOrderLeftQtyManager(sideBuy, contractMultiplier), orderSellList: leftqty.NewOrderLeftQtyManager(sideSell, contractMultiplier)}
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

	var orderList *leftqty.OrderLeftQtyManager
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
	defaultLongMarginRatio   float64
	defaultShortMarginRatio  float64
	orderListsMap            map[string]*OrderLists // key = 交易对手ID + 标的
	symbolQtyAggregrator     *SymbolQtyAggregrator
	accountAmountAggregrator *AccountAmountAggregrator
	orderLog                 *logger.OrderLog
}

func NewMarginExposureManager(longMarginRatio, shortMarginRatio float64, orderLog *logger.OrderLog) *MarginExposureManager {
	inst := &MarginExposureManager{defaultLongMarginRatio: longMarginRatio, defaultShortMarginRatio: shortMarginRatio, orderListsMap: make(map[string]*OrderLists), symbolQtyAggregrator: NewSymbolQtyAggregrator(orderLog), accountAmountAggregrator: NewAccountAmountAggregrator(orderLog)}
	return inst
}

func (m *MarginExposureManager) AddOrder(tradeOrder *schema.TradeOrder) {
	key := m.getKey(tradeOrder)
	orderLists, ok := m.orderListsMap[key]
	if !ok {
		orderLists = NewOrderLists(m.getContractMultiplier(tradeOrder))
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
	key := fmt.Sprintf("%v-%v", tradeOrder.Account, tradeOrder.ExtendAttrMap["symbol2"])
	return key
}

func (m *MarginExposureManager) getContractMultiplier(tradeOrder *schema.TradeOrder) float64 {
	contractMultiplier, ok, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "contractMultiplier", enum.AttrValueType_FLOAT)
	if !ok || contractMultiplier == 0 {
		contractMultiplier = 1
	}
	return contractMultiplier.(float64)
}

// 计算初保占用
func (m *MarginExposureManager) CalculateMarginExposure(positionRecord *PositionRecord, tradeOrder *schema.TradeOrder) {

	key := m.getKey(tradeOrder)
	orderLists, ok := m.orderListsMap[key]
	if !ok {
		return
	}

	buyOrderLeftQty := float64(orderLists.orderBuyList.GetTotalLeftQty())
	sellOrderLeftQty := float64(orderLists.orderSellList.GetTotalLeftQty())

	positionRecord.BuyOrderLeftQty = buyOrderLeftQty
	positionRecord.SellOrderLeftQty = sellOrderLeftQty
	positionRecord.BuyOrderLeftCost = orderLists.orderBuyList.GetTotalLeftAmount() * positionRecord.ContractMultiplier
	positionRecord.SellOrderLeftCost = orderLists.orderSellList.GetTotalLeftAmount() * positionRecord.ContractMultiplier

	m.symbolQtyAggregrator.CreateOrUpdatePositionRecord(positionRecord)
	m.accountAmountAggregrator.CreateOrUpdatePositionRecord(positionRecord)
}

// 计算初保占用
func (m *MarginExposureManager) CalculateInitMarginExposure(positionRecord *PositionRecord) {

	positionRecord.BuyOrderLeftQty = 0
	positionRecord.SellOrderLeftQty = 0
	positionRecord.BuyOrderLeftCost = 0
	positionRecord.SellOrderLeftCost = 0

	m.symbolQtyAggregrator.CreateOrUpdatePositionRecord(positionRecord)
	m.accountAmountAggregrator.CreateOrUpdatePositionRecord(positionRecord)
}

func (m *MarginExposureManager) GetQty(positionRecord *PositionRecord) (float64, float64, float64, float64) {
	aggregratorItem, ok := m.symbolQtyAggregrator.entry[positionRecord.Symbol2]
	if !ok {
		return 0, 0, 0, 0
	}
	return aggregratorItem.longQty, aggregratorItem.shortQty, aggregratorItem.buyOrderLeftQty, aggregratorItem.sellOrderLeftQty
}

func (m *MarginExposureManager) GetCost(positionRecord *PositionRecord) (float64, float64, float64, float64) {
	aggregratorItem, ok := m.accountAmountAggregrator.entry[positionRecord.Account]
	if !ok {
		return 0, 0, 0, 0
	}
	return aggregratorItem.longCost, aggregratorItem.shortCost, aggregratorItem.buyOrderLeftCost, aggregratorItem.sellOrderLeftCost
}

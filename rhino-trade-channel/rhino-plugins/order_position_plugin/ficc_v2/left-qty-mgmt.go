package ficc_v2

import (
	"math"
	"slices"
	"sort"
)

type OrderLeftQty interface {
	GetOrderPrice() float64
	GetOrderID() string
	GetOrderLeftQty() int64
	SetOrderLeftQty(leftQty int64)
	GetArrayIndex() int
	SetArrayIndex(arrayIndex int)
}

type OrderLeftQtyManager struct {
	tradeSide       string
	parValue        float64
	totalLeftQty    int64
	totalLeftAmount float64 // 未除券面价格
	qtyList         []OrderLeftQty
	qtyMap          map[string]OrderLeftQty
}

func NewOrderLeftQtyManager(tradeSide string, parValue float64) *OrderLeftQtyManager {
	return &OrderLeftQtyManager{tradeSide: tradeSide, parValue: parValue, qtyMap: map[string]OrderLeftQty{}}
}

func (m *OrderLeftQtyManager) AddQty(orderLeftQty OrderLeftQty) {

	if _, exists := m.qtyMap[orderLeftQty.GetOrderID()]; exists {
		return
	}

	if orderLeftQty.GetOrderLeftQty() <= 0 {
		return
	}

	n := len(m.qtyList)

	if n == 0 {
		orderLeftQty.SetArrayIndex(0)
		m.qtyList = append(m.qtyList, orderLeftQty)
		m.qtyMap[orderLeftQty.GetOrderID()] = orderLeftQty
		m.totalLeftQty = orderLeftQty.GetOrderLeftQty()
		m.totalLeftAmount = float64(orderLeftQty.GetOrderLeftQty()) * orderLeftQty.GetOrderPrice()
		return
	}

	idx := sort.Search(n, func(i int) bool {
		return m.qtyList[i].GetOrderPrice() < orderLeftQty.GetOrderPrice()
	})

	orderLeftQty.SetArrayIndex(idx)
	m.qtyList = slices.Insert(m.qtyList, idx, orderLeftQty)
	m.qtyMap[orderLeftQty.GetOrderID()] = orderLeftQty
	m.totalLeftQty += orderLeftQty.GetOrderLeftQty()
	m.totalLeftAmount += float64(orderLeftQty.GetOrderLeftQty()) * orderLeftQty.GetOrderPrice()

	n = len(m.qtyList)
	for i := idx + 1; i < n; i++ {
		qtyEl := m.qtyList[i]
		qtyEl.SetArrayIndex(qtyEl.GetArrayIndex() + 1)
	}
}

func (m *OrderLeftQtyManager) DeleteQty(orderID string) {

	orderLeftQty, ok := m.qtyMap[orderID]
	if !ok {
		return
	}
	idx := orderLeftQty.GetArrayIndex()

	m.qtyList = slices.Delete(m.qtyList, idx, idx+1)
	delete(m.qtyMap, orderID)
	m.totalLeftQty -= orderLeftQty.GetOrderLeftQty()
	m.totalLeftAmount -= float64(orderLeftQty.GetOrderLeftQty()) * orderLeftQty.GetOrderPrice()

	n := len(m.qtyList)
	if idx == n {
		return
	}

	for i := idx; i < n; i++ {
		el := m.qtyList[i]
		el.SetArrayIndex(el.GetArrayIndex() - 1)
	}
}

func (m *OrderLeftQtyManager) UpdateQty(orderID string, newLeftQty int64) {

	orderLeftQty, ok := m.qtyMap[orderID]
	if !ok {
		return
	}

	if newLeftQty <= 0 {
		m.DeleteQty(orderID)
		return
	}

	oldLeftQty := orderLeftQty.GetOrderLeftQty()
	if oldLeftQty == newLeftQty {
		return
	}

	orderLeftQty.SetOrderLeftQty(newLeftQty)
	m.totalLeftQty = m.totalLeftQty - oldLeftQty + newLeftQty
	m.totalLeftAmount = m.totalLeftAmount - float64(oldLeftQty)*orderLeftQty.GetOrderPrice() + float64(newLeftQty)*orderLeftQty.GetOrderPrice()
}

// 无论多空头，netPosition要求为正数，返回的值未除parValue
func (m *OrderLeftQtyManager) GetOversoldLeftAmount(netPosition float64) (oversoldLeftAmount float64) {

	absNetPosition := int64(math.Abs(netPosition))

	if m.totalLeftQty <= absNetPosition {
		return 0
	}

	J := 0
	var tmpQty int64
	var lastAdjQty int64
	n := len(m.qtyList)

	oversoldQty := m.totalLeftQty - absNetPosition

	for j := 0; j < n; j++ {
		tmpQty += m.qtyList[j].GetOrderLeftQty()
		if tmpQty >= oversoldQty {
			J = j
			lastAdjQty = tmpQty - oversoldQty
			break
		}
	}

	for i := 0; i <= J; i++ {
		var leftQty int64
		qtyEl := m.qtyList[i]
		if i < J {
			leftQty = qtyEl.GetOrderLeftQty()
		} else {
			leftQty = qtyEl.GetOrderLeftQty() - lastAdjQty
		}

		oversoldLeftAmount += qtyEl.GetOrderPrice() * float64(leftQty)
	}

	return
}

func (m *OrderLeftQtyManager) GetTotalLeftQty() int64 {
	return m.totalLeftQty
}

func (m *OrderLeftQtyManager) GetTotalLeftAmount() float64 {
	return m.totalLeftAmount
}

func (m *OrderLeftQtyManager) GetQtyList() []OrderLeftQty {
	return m.qtyList
}

func (m *OrderLeftQtyManager) GetQtyMap() map[string]OrderLeftQty {
	return m.qtyMap
}

func (m *OrderLeftQtyManager) GetParValue() float64 {
	return m.parValue
}

// func (m *OrderLeftQtyManager) ComputeMaxMarginExposure(netPosition float64, avgPrice float64) {

// }

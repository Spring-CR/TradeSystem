package leftqty_test

import (
	"math"
	"rhino-plugins/order_position_plugin/ficc_v2"
	"testing"
)

type MockOrder struct {
	orderID    string
	price      float64
	leftQty    int64
	arrayIndex int
}

func (m *MockOrder) GetOrderPrice() float64 { return m.price }
func (m *MockOrder) GetOrderID() string     { return m.orderID }
func (m *MockOrder) GetOrderLeftQty() int64 { return m.leftQty }
func (m *MockOrder) SetOrderLeftQty(qty int64) { m.leftQty = qty }
func (m *MockOrder) GetArrayIndex() int     { return m.arrayIndex }
func (m *MockOrder) SetArrayIndex(idx int)  { m.arrayIndex = idx }

func TestAddQty(t *testing.T) {
	mgr := ficc_v2.NewOrderLeftQtyManager("BUY", 100)
	
	// 测试添加第一个订单
	order1 := &MockOrder{orderID: "1", price: 100, leftQty: 10}
	mgr.AddQty(order1)
	
	if mgr.GetTotalLeftQty() != 10 {
		t.Errorf("添加后总数量错误，期望 10，实际 %d", mgr.GetTotalLeftQty())
	}
	if mgr.GetTotalLeftAmount() != 1000 {
		t.Errorf("添加后总金额错误，期望 1000，实际 %f", mgr.GetTotalLeftAmount())
	}
	if len(mgr.GetQtyList()) != 1 {
		t.Errorf("订单列表长度错误，期望 1，实际 %d", len(mgr.GetQtyList()))
	}
	if mgr.GetQtyList()[0].GetOrderID() != "1" {
		t.Errorf("订单ID错误，期望 '1'，实际 '%s'", mgr.GetQtyList()[0].GetOrderID())
	}
	if order1.arrayIndex != 0 {
		t.Errorf("订单索引错误，期望 0，实际 %d", order1.arrayIndex)
	}

	// 测试添加更高价格的订单（应插入到前面）
	order2 := &MockOrder{orderID: "2", price: 110, leftQty: 5}
	mgr.AddQty(order2)
	
	if mgr.GetTotalLeftQty() != 15 {
		t.Errorf("添加后总数量错误，期望 15，实际 %d", mgr.GetTotalLeftQty())
	}
	if mgr.GetTotalLeftAmount() != 1550 { // 100*10 + 110*5
		t.Errorf("添加后总金额错误，期望 1550，实际 %f", mgr.GetTotalLeftAmount())
	}
	if len(mgr.GetQtyList()) != 2 {
		t.Errorf("订单列表长度错误，期望 2，实际 %d", len(mgr.GetQtyList()))
	}
	if mgr.GetQtyList()[0].GetOrderID() != "2" {
		t.Errorf("订单顺序错误，期望 '2' 在位置0")
	}
	if mgr.GetQtyList()[1].GetOrderID() != "1" {
		t.Errorf("订单顺序错误，期望 '1' 在位置1")
	}
	if order2.arrayIndex != 0 {
		t.Errorf("订单2索引错误，期望 0，实际 %d", order2.arrayIndex)
	}
	if order1.arrayIndex != 1 {
		t.Errorf("订单1索引错误，期望 1，实际 %d", order1.arrayIndex)
	}

	// 测试添加重复订单
	mgr.AddQty(order1)
	if mgr.GetTotalLeftQty() != 15 || len(mgr.GetQtyList()) != 2 {
		t.Errorf("重复添加订单不应生效")
	}

	// 测试添加零数量订单
	order3 := &MockOrder{orderID: "3", price: 90, leftQty: 0}
	mgr.AddQty(order3)
	if mgr.GetTotalLeftQty() != 15 || len(mgr.GetQtyList()) != 2 {
		t.Errorf("零数量订单不应被添加")
	}
}

func TestDeleteQty(t *testing.T) {
	mgr := ficc_v2.NewOrderLeftQtyManager("SELL", 100)
	orders := []*MockOrder{
		{orderID: "1", price: 100, leftQty: 10},
		{orderID: "2", price: 110, leftQty: 5},
		{orderID: "3", price: 90, leftQty: 8},
	}

	for _, o := range orders {
		mgr.AddQty(o)
	}

	// 删除中间订单
	mgr.DeleteQty("1")
	
	expectedAmount := float64(110*5 + 90*8)
	
	if mgr.GetTotalLeftQty() != 13 {
		t.Errorf("删除后总数量错误，期望 13，实际 %d", mgr.GetTotalLeftQty())
	}
	if mgr.GetTotalLeftAmount() != expectedAmount {
		t.Errorf("删除后总金额错误，期望 %f，实际 %f", expectedAmount, mgr.GetTotalLeftAmount())
	}
	if len(mgr.GetQtyList()) != 2 {
		t.Errorf("订单列表长度错误，期望 2，实际 %d", len(mgr.GetQtyList()))
	}
	if mgr.GetQtyList()[0].GetOrderID() != "2" {
		t.Errorf("订单顺序错误，期望 '2' 在位置0")
	}
	if mgr.GetQtyList()[1].GetOrderID() != "3" {
		t.Errorf("订单顺序错误，期望 '3' 在位置1")
	}
	if mgr.GetQtyList()[0].GetArrayIndex() != 0 {
		t.Errorf("订单索引错误，期望 0，实际 %d", mgr.GetQtyList()[0].GetArrayIndex())
	}
	if mgr.GetQtyList()[1].GetArrayIndex() != 1 {
		t.Errorf("订单索引错误，期望 1，实际 %d", mgr.GetQtyList()[1].GetArrayIndex())
	}

	// 删除不存在的订单
	mgr.DeleteQty("999")
	if mgr.GetTotalLeftQty() != 13 || len(mgr.GetQtyList()) != 2 {
		t.Errorf("删除不存在的订单不应影响状态")
	}
}

func TestUpdateQty(t *testing.T) {
	mgr := ficc_v2.NewOrderLeftQtyManager("BUY", 100)
	order := &MockOrder{orderID: "1", price: 100, leftQty: 10}
	mgr.AddQty(order)

	// 更新数量
	mgr.UpdateQty("1", 15)
	if mgr.GetTotalLeftQty() != 15 {
		t.Errorf("更新后总数量错误，期望 15，实际 %d", mgr.GetTotalLeftQty())
	}
	if mgr.GetTotalLeftAmount() != 1500 {
		t.Errorf("更新后总金额错误，期望 1500，实际 %f", mgr.GetTotalLeftAmount())
	}
	if order.leftQty != 15 {
		t.Errorf("订单数量未更新，期望 15，实际 %d", order.leftQty)
	}

	// 更新为0（应删除）
	mgr.UpdateQty("1", 0)
	if mgr.GetTotalLeftQty() != 0 {
		t.Errorf("更新为0后总数量应为0，实际 %d", mgr.GetTotalLeftQty())
	}
	if mgr.GetTotalLeftAmount() != 0 {
		t.Errorf("更新为0后总金额应为0，实际 %f", mgr.GetTotalLeftAmount())
	}
	if len(mgr.GetQtyList()) != 0 {
		t.Errorf("订单列表应为空，实际长度 %d", len(mgr.GetQtyList()))
	}

	// 更新不存在的订单
	mgr.UpdateQty("999", 10)
	if mgr.GetTotalLeftQty() != 0 || len(mgr.GetQtyList()) != 0 {
		t.Errorf("更新不存在的订单不应影响状态")
	}
}

func TestGetOversoldLeftAmount(t *testing.T) {
	mgr := ficc_v2.NewOrderLeftQtyManager("SELL", 100)
	orders := []*MockOrder{
		{orderID: "1", price: 110, leftQty: 5}, // 最高价
		{orderID: "2", price: 100, leftQty: 10},
		{orderID: "3", price: 90, leftQty: 8},  // 最低价
	}

	for _, o := range orders {
		mgr.AddQty(o)
	}

	// 定义比较函数，处理浮点精度
	nearlyEqual := func(a, b float64) bool {
		return math.Abs(a-b) < 1e-5
	}

	tests := []struct {
		name         string
		netPosition  float64
		expected     float64
	}{
		{"净头寸大于总数量", 30, 0},
		{"部分超售1", 20, 110*3},         // 超售3张，取最高价订单3张
		{"部分超售2", 15, 110*5 + 100*3},  // 超售8张，取最高价5张+次高价3张
		{"完全超售", 0, 110*5 + 100*10 + 90*8}, // 超售23张，取所有订单
		{"负净头寸", -10, 110*5 + 100*8},  // 超售13张，取最高价5张+次高价8张
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mgr.GetOversoldLeftAmount(tt.netPosition)
			if !nearlyEqual(result, tt.expected) {
				t.Errorf("净头寸 %.0f 预期 %.1f，实际 %.1f", tt.netPosition, tt.expected, result)
			}
		})
	}
}

func TestOrderSorting(t *testing.T) {
	mgr := ficc_v2.NewOrderLeftQtyManager("BUY", 100)
	
	// 添加不同价格的订单
	orders := []*MockOrder{
		{orderID: "1", price: 105, leftQty: 5},
		{orderID: "2", price: 100, leftQty: 10},
		{orderID: "3", price: 110, leftQty: 8},
		{orderID: "4", price: 95, leftQty: 7},
		{orderID: "5", price: 102, leftQty: 6},
	}
	
	for _, o := range orders {
		mgr.AddQty(o)
	}
	
	// 验证排序顺序（价格从高到低）
	expectedOrder := []string{"3", "1", "5", "2", "4"}
	qtyList := mgr.GetQtyList()
	if len(qtyList) != len(expectedOrder) {
		t.Fatalf("订单数量错误，期望 %d，实际 %d", len(expectedOrder), len(qtyList))
	}
	
	for i, id := range expectedOrder {
		if qtyList[i].GetOrderID() != id {
			t.Errorf("位置 %d 订单错误，期望 %s，实际 %s", i, id, qtyList[i].GetOrderID())
		}
		if qtyList[i].GetArrayIndex() != i {
			t.Errorf("订单 %s 索引错误，期望 %d，实际 %d", id, i, qtyList[i].GetArrayIndex())
		}
	}
	
	// 验证总数量和总金额
	var totalQty int64
	var totalAmount float64
	for _, o := range orders {
		totalQty += o.leftQty
		totalAmount += float64(o.leftQty) * o.price
	}
	
	if mgr.GetTotalLeftQty() != totalQty {
		t.Errorf("总数量错误，期望 %d，实际 %d", totalQty, mgr.GetTotalLeftQty())
	}
	if mgr.GetTotalLeftAmount() != totalAmount {
		t.Errorf("总金额错误，期望 %f，实际 %f", totalAmount, mgr.GetTotalLeftAmount())
	}
}

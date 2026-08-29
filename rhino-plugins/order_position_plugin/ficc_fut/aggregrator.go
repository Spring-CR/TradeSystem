package ficc_fut

import (
	"encoding/json"
	"rhino-common/utils/logger"
)

type SymbolQtyAggregratorItem struct {
	id               string
	longQty          float64
	shortQty         float64
	buyOrderLeftQty  float64
	sellOrderLeftQty float64
	entry            map[string]*PositionRecord // key = {account}-{symbol2}
	orderLog         *logger.OrderLog
}

func NewSymbolQtyAggregratorItem(id string, orderLog *logger.OrderLog) *SymbolQtyAggregratorItem {

	inst := &SymbolQtyAggregratorItem{
		id:       id,
		entry:    map[string]*PositionRecord{},
		orderLog: orderLog,
	}

	return inst
}

func (a *SymbolQtyAggregratorItem) Refresh() {

	var longQty float64
	var shortQty float64
	var buyOrderLeftQty float64
	var sellOrderLeftQty float64

	for _, p := range a.entry {
		if p.NetPosition > 0 {
			longQty += p.NetPosition
		} else {
			shortQty += -p.NetPosition
		}
		buyOrderLeftQty += p.BuyOrderLeftQty
		sellOrderLeftQty += p.SellOrderLeftQty
	}

	a.longQty = longQty
	a.shortQty = shortQty
	a.buyOrderLeftQty = buyOrderLeftQty
	a.sellOrderLeftQty = sellOrderLeftQty
}

func (a *SymbolQtyAggregratorItem) CreateOrUpdatePositionRecord(positionRecord *PositionRecord) {

	key := positionRecord.Key
	_, ok := a.entry[key]
	if !ok {
		a.entry[key] = positionRecord
	}

	a.Refresh()

	jsData, _ := json.MarshalIndent(positionRecord, "", " ")
	a.orderLog.PrintfWithKey(a.id+"-qty-aggr", "[SymbolQtyAggregrate][CreateOrUpdatePositionRecord] - positionRecord:%s, key:%v, longQty:%v, shortQty:%v, buyOrderLeftQty:%v, sellOrderLeftQty:%v", jsData, a.id, a.longQty, a.shortQty, a.buyOrderLeftQty, a.sellOrderLeftQty)
}

func (a *SymbolQtyAggregratorItem) CreateOrUpdatePositionRecord2(positionRecord *PositionRecord) {

	key := positionRecord.Key
	jsData, _ := json.MarshalIndent(positionRecord, "", " ")
	oldPositionRecord, ok := a.entry[key]

	if !ok {

		a.entry[key] = positionRecord

		if positionRecord.NetPosition > 0 {
			a.longQty += positionRecord.NetPosition
		}

		if positionRecord.NetPosition < 0 {
			a.shortQty += -positionRecord.NetPosition
		}

		// ----------------------------------------------------

		a.buyOrderLeftQty += positionRecord.BuyOrderLeftQty
		a.sellOrderLeftQty += positionRecord.SellOrderLeftQty

		a.orderLog.PrintfWithKey(a.id+"-qty-aggr", "[SymbolQtyAggregrate][CreateOrUpdatePositionRecord] - positionRecord:%s\n key:%v, longQty:%v, shortQty:%v, buyOrderLeftQty:%v, sellOrderLeftQty:%v\n", jsData, a.id, a.longQty, a.shortQty, a.buyOrderLeftQty, a.sellOrderLeftQty)

		return
	}

	var newLongQty, newShortQty, oldLongQty, oldShortQty float64
	if positionRecord.NetPosition > 0 {
		newLongQty = positionRecord.NetPosition
		newShortQty = 0
	}
	if positionRecord.NetPosition < 0 {
		newLongQty = 0
		newShortQty = -positionRecord.NetPosition
	}
	if oldPositionRecord.NetPosition > 0 {
		oldLongQty = oldPositionRecord.NetPosition
		oldShortQty = 0
	}
	if oldPositionRecord.NetPosition < 0 {
		oldLongQty = 0
		oldShortQty = -oldPositionRecord.NetPosition
	}

	// 计算diff
	longQtyDiff := newLongQty - oldLongQty
	shortQtyDiff := newShortQty - oldShortQty

	a.longQty += longQtyDiff
	a.shortQty += shortQtyDiff

	if a.longQty < 0 {
		a.longQty = 0
	}

	if a.shortQty < 0 {
		a.shortQty = 0
	}

	// --------------------------------------------------------

	// 获取新的BuyOrderLeftQty、SellOrderLeftQty
	newBuyOrderLeftQty := positionRecord.BuyOrderLeftQty
	newSellOrderLeftQty := positionRecord.SellOrderLeftQty

	// 获取旧的BuyOrderLeftQty、SellOrderLeftQty
	oldBuyOrderLeftQty := oldPositionRecord.BuyOrderLeftQty
	oldSellOrderLeftQty := oldPositionRecord.SellOrderLeftQty

	// 计算diff
	buyOrderLeftQtyDiff := newBuyOrderLeftQty - oldBuyOrderLeftQty
	sellOrderLeftQtyDiff := newSellOrderLeftQty - oldSellOrderLeftQty

	a.buyOrderLeftQty += buyOrderLeftQtyDiff
	a.sellOrderLeftQty += sellOrderLeftQtyDiff

	if a.buyOrderLeftQty < 0 {
		a.buyOrderLeftQty = 0
	}

	if a.sellOrderLeftQty < 0 {
		a.sellOrderLeftQty = 0
	}

	return
}

type SymbolQtyAggregrator struct {
	entry    map[string]*SymbolQtyAggregratorItem // key = symbol2
	orderLog *logger.OrderLog
}

func NewSymbolQtyAggregrator(orderLog *logger.OrderLog) *SymbolQtyAggregrator {

	inst := &SymbolQtyAggregrator{entry: make(map[string]*SymbolQtyAggregratorItem), orderLog: orderLog}

	return inst
}

func (a *SymbolQtyAggregrator) CreateOrUpdatePositionRecord(positionRecord *PositionRecord) {

	key := positionRecord.Symbol2

	symbolQtyAggregratorItem, ok := a.entry[key]

	if !ok {
		symbolQtyAggregratorItem = NewSymbolQtyAggregratorItem(key, a.orderLog)
		a.entry[key] = symbolQtyAggregratorItem
	}

	symbolQtyAggregratorItem.CreateOrUpdatePositionRecord(positionRecord)
}

package ficc_fut

import (
	"encoding/json"
	"rhino-common/utils/logger"
	"strconv"
)

type AccountAmountAggregratorItem struct {
	id                string
	longCost          float64
	shortCost         float64
	buyOrderLeftCost  float64
	sellOrderLeftCost float64
	entry             map[string]*PositionRecord // key = {account}-{symbol2}
	orderLog          *logger.OrderLog
}

func NewAccountAmountAggregratorItem(id string, orderLog *logger.OrderLog) *AccountAmountAggregratorItem {

	inst := &AccountAmountAggregratorItem{
		id:       id,
		entry:    map[string]*PositionRecord{},
		orderLog: orderLog,
	}

	return inst
}

func (a *AccountAmountAggregratorItem) Refresh() {

	var longCost float64
	var shortCost float64
	var buyOrderLeftCost float64
	var sellOrderLeftCost float64

	for _, p := range a.entry {
		if p.NetPosition > 0 {
			longCost += p.LongPriceCNYWithFeeCost
		} else {
			shortCost += -p.ShortPriceCNYWithFeeCost
		}
		buyOrderLeftCost += p.BuyOrderLeftCost
		sellOrderLeftCost += p.SellOrderLeftCost
	}

	a.longCost = longCost
	a.shortCost = shortCost
	a.buyOrderLeftCost = buyOrderLeftCost
	a.sellOrderLeftCost = sellOrderLeftCost
}

func (a *AccountAmountAggregratorItem) CreateOrUpdatePositionRecord(positionRecord *PositionRecord) {

	key := positionRecord.Key

	_, ok := a.entry[key]

	defer func() {
		jsData, _ := json.MarshalIndent(positionRecord, "", " ")
		a.orderLog.PrintfWithKey(a.id+"-amt-aggr", "[AccountAmounAggregrate][CreateOrUpdatePositionRecord] - positionRecord:%s, key:%v, longCost:%v, shortCost:%v, buyOrderLeftCost:%v, sellOrderLeftCost:%v", jsData, a.id, a.longCost, a.shortCost, a.buyOrderLeftCost, a.sellOrderLeftCost)
	}()

	if !ok {
		a.entry[key] = positionRecord
	}

	a.Refresh()

	return
}

func (a *AccountAmountAggregratorItem) CreateOrUpdatePositionRecord2(positionRecord *PositionRecord) {

	key := positionRecord.Key

	oldPositionRecord, ok := a.entry[key]

	defer func() {
		jsData, _ := json.MarshalIndent(positionRecord, "", " ")
		a.orderLog.PrintfWithKey(a.id+"-amt-aggr", "[AccountAmounAggregrate][CreateOrUpdatePositionRecord] - positionRecord:%s\n key:%v, longCost:%v, shortCost:%v, buyOrderLeftCost:%v, sellOrderLeftCost:%v\n", jsData, a.id, a.longCost, a.shortCost, a.buyOrderLeftCost, a.sellOrderLeftCost)
	}()

	if !ok {

		a.entry[key] = positionRecord

		a.longCost += positionRecord.LongPriceCNYWithFeeCost
		a.shortCost += positionRecord.ShortPriceCNYWithFeeCost

		// ----------------------------------------------------

		a.buyOrderLeftCost += positionRecord.BuyOrderLeftCost * positionRecord.ExchangeRateCNY
		a.sellOrderLeftCost += positionRecord.SellOrderLeftCost * positionRecord.ExchangeRateCNY

		return
	}

	// 获取新的longCost、shortCost
	newLongCost := positionRecord.LongPriceCNYWithFeeCost
	newShortCost := positionRecord.ShortPriceCNYWithFeeCost

	// 获取旧的longCost、shortCost
	oldLongCost := oldPositionRecord.LongPriceCNYWithFeeCost
	oldShortCost := oldPositionRecord.ShortPriceCNYWithFeeCost

	// 计算diff
	longCostDiff := newLongCost - oldLongCost
	shortCostDiff := newShortCost - oldShortCost

	a.longCost += longCostDiff
	a.shortCost += shortCostDiff

	if a.longCost < 0 {
		a.longCost = 0
	}

	if a.shortCost < 0 {
		a.shortCost = 0
	}

	// ----------------------------------------------------

	// 获取新的BuyOrderLeftCost、SellOrderLeftCost
	newBuyOrderLeftCost := positionRecord.BuyOrderLeftCost
	newSellOrderLeftCost := positionRecord.SellOrderLeftCost

	// 获取旧的BuyOrderLeftCost、SellOrderLeftCost
	oldBuyOrderLeftCost := oldPositionRecord.BuyOrderLeftCost
	oldSellOrderLeftCost := oldPositionRecord.SellOrderLeftCost

	// 计算diff
	buyOrderLeftCostDiff := newBuyOrderLeftCost - oldBuyOrderLeftCost
	sellOrderLeftCostDiff := newSellOrderLeftCost - oldSellOrderLeftCost

	a.buyOrderLeftCost += buyOrderLeftCostDiff * positionRecord.ExchangeRateCNY
	a.sellOrderLeftCost += sellOrderLeftCostDiff * positionRecord.ExchangeRateCNY

	if a.buyOrderLeftCost < 0 {
		a.buyOrderLeftCost = 0
	}

	if a.sellOrderLeftCost < 0 {
		a.sellOrderLeftCost = 0
	}

	return
}

type AccountAmountAggregrator struct {
	entry    map[int]*AccountAmountAggregratorItem // key = symbol2
	orderLog *logger.OrderLog
}

func NewAccountAmountAggregrator(orderLog *logger.OrderLog) *AccountAmountAggregrator {

	inst := &AccountAmountAggregrator{entry: make(map[int]*AccountAmountAggregratorItem), orderLog: orderLog}

	return inst
}

func (a *AccountAmountAggregrator) CreateOrUpdatePositionRecord(positionRecord *PositionRecord) {

	key := positionRecord.Account

	accountAmountAggregratorItem, ok := a.entry[key]

	if !ok {
		accountAmountAggregratorItem = NewAccountAmountAggregratorItem(strconv.Itoa(key), a.orderLog)
		a.entry[key] = accountAmountAggregratorItem
	}

	accountAmountAggregratorItem.CreateOrUpdatePositionRecord(positionRecord)
}

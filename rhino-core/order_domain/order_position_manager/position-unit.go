package order_position_manager

import (
	"rhino-common/utils/logger"
	"rhino-core/schema"
	"sync"
)

type PositionUnit struct {
	unitName            string                       // 持仓单元名称
	lock                *sync.RWMutex                // 读写锁
	baseQuota           float64                      // 日初额度
	quota               float64                      // 可用额度
	totalFreezeQty      float64                      // 总冻结量
	totalUnfreezeQty    float64                      // 总解冻量
	quotaLockerMap      map[string]*QuotaLocker      // quotaLockerMap, key=AppOrdID
	quotaIncrementerMap map[string]*QuotaIncrementer // quotaIncrementerMap, key=AppOrdID
	orderLog            *logger.OrderLog
}

func NewPositionUnit(unitName string, baseQuota float64, orderLog *logger.OrderLog) *PositionUnit {

	pu := &PositionUnit{
		unitName:            unitName,
		lock:                &sync.RWMutex{},
		baseQuota:           baseQuota,
		quota:               baseQuota,
		totalFreezeQty:      0,
		totalUnfreezeQty:    0,
		quotaLockerMap:      map[string]*QuotaLocker{},
		quotaIncrementerMap: map[string]*QuotaIncrementer{},
		orderLog:            orderLog,
	}

	return pu
}

func (u *PositionUnit) freezeQuota(positionAdapter PositionAdapter, metadata interface{}, tradeOrder *schema.TradeOrder, lastPositionUnit bool) {

	freezeQty, freezeQtyBuf, ok := positionAdapter.CalculateOrderFreezeQuotaInPositionUnit(u, tradeOrder, metadata)
	u.orderLog.Printf(tradeOrder, nil, "[PositionUnit=%v] freezeQty=%v, ok=%v", u.unitName, freezeQty, ok)
	if !ok {
		return
	}

	quotaLocker := newQuotaLocker(tradeOrder, freezeQty, freezeQtyBuf)

	u.lock.Lock()
	defer u.lock.Unlock()

	// 确保幂等
	if _, ok := u.quotaLockerMap[tradeOrder.AppOrdID]; ok {
		return
	}

	u.quotaLockerMap[tradeOrder.AppOrdID] = quotaLocker

	u.totalFreezeQty += freezeQty

	positionAdapter.AfterFreezeQuotaInPositionUnit(u, freezeQty, tradeOrder, metadata, quotaLocker, lastPositionUnit)
}

func (u *PositionUnit) increaseQuota(positionAdapter PositionAdapter, metadata interface{}, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, lastPositionUnit bool) {

	increaseQty, ok := positionAdapter.CalculateIncreaseQuotaInPositionUnit(u, tradeOrder, tradeActionResp, metadata)
	if !ok {
		return
	}
	u.orderLog.Printf(tradeOrder, tradeActionResp, "[PositionUnit=%v] increaseQty=%v, ok=%v", u.unitName, increaseQty, ok)

	u.lock.Lock()
	defer u.lock.Unlock()

	quotaIncrementer, ok := u.quotaIncrementerMap[tradeOrder.AppOrdID]
	if !ok {
		quotaIncrementer = newQuotaIncrementer(tradeOrder)
		u.quotaIncrementerMap[tradeOrder.AppOrdID] = quotaIncrementer
	}

	quotaIncrementer.increase(increaseQty, tradeActionResp)

	u.quota += increaseQty

	positionAdapter.AfterIncreaseQuotaInPositionUnit(u, increaseQty, tradeOrder, tradeActionResp, metadata, quotaIncrementer, lastPositionUnit)
}

func (u *PositionUnit) unfreezeQuota(positionAdapter PositionAdapter, metadata interface{}, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, lastPositionUnit bool) {

	u.lock.Lock()
	defer u.lock.Unlock()

	quotaLocker, ok := u.quotaLockerMap[tradeOrder.AppOrdID]
	if !ok {
		return
	}

	unfreezeQty, ok := positionAdapter.CalculateUnfreezeQuotaInPositionUnit(u, tradeOrder, tradeActionResp, metadata, quotaLocker)
	if !ok {
		return
	}
	u.orderLog.Printf(tradeOrder, tradeActionResp, "[PositionUnit=%v] unfreezeQty=%v, ok=%v", u.unitName, unfreezeQty, ok)

	quotaLocker.unfreeze(unfreezeQty, tradeOrder, tradeActionResp)

	u.quota += unfreezeQty
	u.totalUnfreezeQty += unfreezeQty

	positionAdapter.AfterUnfreezeQuotaInPositionUnit(u, unfreezeQty, tradeOrder, tradeActionResp, metadata, quotaLocker, lastPositionUnit, false)
}

func (u *PositionUnit) rollbackFreezeQuota(positionAdapter PositionAdapter, metadata interface{}, tradeOrder *schema.TradeOrder, lastPositionUnit bool) {

	u.lock.Lock()
	defer u.lock.Unlock()

	quotaLocker, ok := u.quotaLockerMap[tradeOrder.AppOrdID]
	if !ok {
		return
	}

	unfreezeQty := quotaLocker.GetFreezeQty()
	quotaLocker.unfreeze(unfreezeQty, tradeOrder, nil)

	u.orderLog.Printf(tradeOrder, nil, "[PositionUnit=%v] rollbackFreezeQty=%v, ok=%v", unfreezeQty, unfreezeQty, ok)

	u.quota += unfreezeQty
	u.totalUnfreezeQty += unfreezeQty

	positionAdapter.AfterUnfreezeQuotaInPositionUnit(u, unfreezeQty, tradeOrder, nil, metadata, quotaLocker, lastPositionUnit, true)

	// rollback后删除map entry
	delete(u.quotaLockerMap, tradeOrder.AppOrdID)
}

func (u *PositionUnit) GetName() string {
	return u.unitName
}

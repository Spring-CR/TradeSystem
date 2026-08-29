package order_domain

import (
	"log"
	"rhino-common/domain_error"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/timeutil"
	"rhino-core/store/app_store"
	"rhino-core/types"
	"time"
)

// 1. 按AppOrdID取出订单。去掉-【找不到订单、订单当前状态不是草稿态的都要报错】。
// 2. 将订单的状态设置为删除态，保存订单时，需要在一个事务中同时插入TradeActionLatestResp对象，该对象将保存改动前的订单内容。
// 3. 通过TradeActionLatestResp的ActionKey去重。
// 4. OrderCache同步更新。
func (o *OrderOrchestrator) UpdateOrderAttributes(attrUpdateReq *types.ApplicationOrderAttributeUpdateRequest) (de *domain_error.Error) {

	// 开启数据库事务
	tx, de1 := dbutil.BeginTx(o.applicationCfg.GetAppDB())
	if de1 != nil {
		return de1
	}

	// 在内存模型和数据库中，都需要存在该订单才能执行更新逻辑
	traceableTradeOrder, ok := o.orderCache.GetOrderByAppOrdID(attrUpdateReq.AppOrdID)
	if !ok {
		de = domain_error.Build(domain_error.CANNOT_FIND_ORDER_BY_APP_ORD_ID_ERR_CODE, nil, attrUpdateReq.AppOrdID)
		dbutil.RollbackTx(tx)
		return
	}

	log.Printf("get traceableTradeOrder with AppOrdID=%s in OrderCache\n", traceableTradeOrder.GetBasicInfo().AppOrdID)

	// 加锁，直到最后
	traceableTradeOrder.GetLock().Lock()
	defer traceableTradeOrder.GetLock().Unlock()

	log.Printf("lock to UpdateOrderAttributes\n")

	// 状态以OrderCache的为准
	orderDraft := traceableTradeOrder.GetBasicInfo()

	if attrUpdateReq.ActionTime <= 0 {
		attrUpdateReq.ActionTime = timeutil.ConvertTimeToMilliseconds(time.Now())
	}

	oldAttrs := make(map[string]interface{})
	for k, newVal := range attrUpdateReq.UpdateAttributes {
		v, ok := orderDraft.ExtendAttrMap[k]
		if ok {
			oldAttrs[k] = v
			orderDraft.ExtendAttrMap[k] = newVal
		}
	}
	js, _ := json.Marshal(orderDraft.ExtendAttrMap)
	orderDraft.ExtendAttr = string(js)

	// 数据库处理失败时，需要进行状态的回滚
	rollBackStatus := func() {
		for k, oldVal := range oldAttrs {
			orderDraft.ExtendAttrMap[k] = oldVal
		}
		js, _ := json.Marshal(orderDraft.ExtendAttrMap)
		orderDraft.ExtendAttr = string(js)
	}

	// 从数据库更新order
	err := app_store.UpdateTradeOrderExtendAttr(tx, orderDraft)
	if err != nil {
		de = domain_error.Build(domain_error.CANNOT_UPDATE_ORDER_DRAFT_ERR_CODE, err)
		dbutil.RollbackTx(tx)
		rollBackStatus()
		return
	}

	de = dbutil.CommitTx(tx)
	if de != nil {
		rollBackStatus()
		return
	}

	log.Printf("CommitTx for UpdateTradeOrderExtendAttr")

	o.orderCache.UpdateTradeOrderAttributes(attrUpdateReq.AppOrdID, attrUpdateReq.UpdateAttributes)

	return
}

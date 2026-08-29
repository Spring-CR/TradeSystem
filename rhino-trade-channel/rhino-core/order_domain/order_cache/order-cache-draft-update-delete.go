package order_cache

import (
	"darkpool-common/bean"
	"log"
	"rhino-common/domain_error"
	"rhino-core/schema"
	"rhino-core/types"
)

func (c *OrderCache) UpdateTradeOrderDraft(oldOrder *types.TraceableTradeOrder, tradeActionLatestResp *schema.TradeActionLatestResp, newOrder *schema.TradeOrder) {

	if c.IsSlave() {
		var ok bool
		oldOrder, ok = c.GetOrderByAppOrdID(newOrder.AppOrdID)
		if !ok {
			de := domain_error.Build(domain_error.CANNOT_FIND_ORDER_BY_APP_ORD_ID_ERR_CODE, nil, newOrder.AppOrdID)
			domain_error.ProcessSevereError(false, 0, de, nil, "error occurs in UpdateTradeOrderDraft")
			return
		}

		// 幂等性处理，判断是否重复处理
		// 如果是slave，检查是否已经操作过，如果是重复的动作，则因该忽略
		// 否则，更新tradeActionLatestRespKeyMap
		key := tradeActionLatestResp.ActionKey
		c.tradeActionLatestRespKeyMapLock.Lock()
		if c.tradeActionLatestRespKeyMap[key] { // 重复输入
			c.tradeActionLatestRespKeyMapLock.Unlock()
			log.Printf("===> skip for slave, UpdateTradeOrderDraft, return for duplicate action key:%s\n", key)
			return
		}
		c.tradeActionLatestRespKeyMap[key] = true
		c.tradeActionLatestRespKeyMapLock.Unlock()
	}

	// 更新内存模型
	// 对于多次草稿保存的情况，用拷贝的方式更新已存在的order记录，并添加TradeOrderLatestResp记录
	bean.Copy(oldOrder.GetBasicInfo()).From(newOrder)
	// 从order找出TraceableTradeOrder对象
	// 改成用InitTradeAction，因为操作draft前面旧的draft都已经加锁，使用不带锁的AddTradeAction才能避免死锁问题
	oldOrder.InitTradeAction(types.NewTraceableTradeActionResp(oldOrder, oldOrder.GetBasicInfo(), tradeActionLatestResp, c.IsSlave()))

	c.afterUpdateTradeOrderDraft(tradeActionLatestResp, oldOrder.GetBasicInfo())
}

func (c *OrderCache) DeleteTradeOrderDraft(oldOrder *types.TraceableTradeOrder, tradeActionLatestResp *schema.TradeActionLatestResp, deletedOrderDraft *schema.TradeOrder) {

	if c.IsSlave() {
		var ok bool
		oldOrder, ok = c.GetOrderByAppOrdID(tradeActionLatestResp.AppOrdID)
		if !ok {
			de := domain_error.Build(domain_error.CANNOT_FIND_ORDER_BY_APP_ORD_ID_ERR_CODE, nil, tradeActionLatestResp.AppOrdID)
			domain_error.ProcessSevereError(false, 0, de, nil, "error occurs in DeleteTradeOrderDraft")
			return
		}

		// 幂等性处理，判断是否重复处理
		// 如果是slave，检查是否已经操作过，如果是重复的动作，则因该忽略
		// 否则，更新tradeActionLatestRespKeyMap
		key := tradeActionLatestResp.ActionKey
		c.tradeActionLatestRespKeyMapLock.Lock()
		if c.tradeActionLatestRespKeyMap[key] { // 重复输入
			c.tradeActionLatestRespKeyMapLock.Unlock()
			log.Printf("===> skip for slave, DeleteTradeOrderDraft, return for duplicate action key:%s\n", key)
			return
		}
		c.tradeActionLatestRespKeyMap[key] = true
		c.tradeActionLatestRespKeyMapLock.Unlock()

		// 更新内存模型（只有slave角色需要更新，因为对于master，已经直接在OrderOrchestrator里更新了订单的内存模型）
		bean.Copy(oldOrder.GetBasicInfo()).From(deletedOrderDraft)
	}

	// 从order找出TraceableTradeOrder对象
	// InitTradeAction是不加锁的，否则，因前面oldOrder已经加锁，就会出现死锁的问题
	oldOrder.InitTradeAction(types.NewTraceableTradeActionResp(oldOrder, oldOrder.GetBasicInfo(), tradeActionLatestResp, c.IsSlave()))
	log.Println("add trade action for draft delection")
	c.afterDeleteTradeOrderDraft(tradeActionLatestResp, oldOrder.GetBasicInfo())
}

func (c *OrderCache) UpdateTradeOrderAttributes(appOrdID string, updateAttributes map[string]interface{}) {

	if c.IsSlave() {
		var ok bool
		oldOrder, ok := c.GetOrderByAppOrdID(appOrdID)
		if !ok {
			de := domain_error.Build(domain_error.CANNOT_FIND_ORDER_BY_APP_ORD_ID_ERR_CODE, nil, oldOrder.GetBasicInfo().AppOrdID)
			domain_error.ProcessSevereError(false, 0, de, nil, "error occurs in UpdateTradeOrderAttributes")
			return
		}

		// 更新内存模型
		for k, v := range updateAttributes {
			oldOrder.GetBasicInfo().ExtendAttrMap[k] = v
		}
		js, _ := json.Marshal(oldOrder.GetBasicInfo().ExtendAttrMap)
		oldOrder.GetBasicInfo().ExtendAttr = string(js)
	}

	c.afterUpdateTradeOrderAttributes(appOrdID, updateAttributes)
}

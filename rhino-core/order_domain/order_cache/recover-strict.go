package order_cache

import (
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/store/app_store"
	"rhino-core/types"
	"sort"
	"time"
)

var (
	// Todo：状态定义发生变化时这里得改
	statusWithClOrdID = map[string]bool{
		string(enum.OrdStatus_Ready):                true,
		string(enum.OrdStatus_Submit):               true,
		string(enum.OrdStatus_InternalSubmitFailed): true,

		string(enum.OrdStatus_New):                true,
		string(enum.OrdStatus_PartiallyFilled):    true,
		string(enum.OrdStatus_Filled):             true,
		string(enum.OrdStatus_DoneForDay):         true,
		string(enum.OrdStatus_Canceled):           true,
		string(enum.OrdStatus_Replaced):           true,
		string(enum.OrdStatus_PendingCancel):      true,
		string(enum.OrdStatus_Stopped):            true,
		string(enum.OrdStatus_Rejected):           true,
		string(enum.OrdStatus_Suspended):          true,
		string(enum.OrdStatus_PendingNew):         true,
		string(enum.OrdStatus_Calculated):         true,
		string(enum.OrdStatus_Expired):            true,
		string(enum.OrdStatus_AcceptedForBidding): true,
		string(enum.OrdStatus_PendingReplace):     true,
	}
)

func (c *OrderCache) strictRecover() (int64, int64, []*OrderCacheSyncMessage) {

	begin := time.Now()
	log.Println("Start to recover data in strict mode...")

	// 第1步：获取状态恢复所需的数据源
	tradeOrders, tradeActionLatestResps, tradeActionResps, newestOffest, oldestOffset, syncMsgs := c.prepairDataForStrictRecover()
	log.Printf("===> Complete step#1: prepairDataForStrictRecover, tradeOrders.Len=%d, tradeActionLatestResps.Len=%d, tradeActionResps.Len=%d, syncMsgs.Len=%d\n", len(tradeOrders), len(tradeActionLatestResps), len(tradeActionResps), len(syncMsgs))

	// 第2步：构建订单拓扑
	var orderMap map[string]*types.TraceableTradeOrder
	c.rootOrders, c.directOrderMap, orderMap = c.constructTradeOrderTopology(tradeOrders)
	log.Printf("===> Complete step#2: constructTradeOrderTopology， rootOrders.Len=%d, directOrderMap.Len=%d, orderMap.Len=%d\n", len(c.rootOrders), len(c.directOrderMap), len(orderMap))

	if c.IsSlave() {
		for _, rootOrder := range c.rootOrders {
			order := rootOrder.GetBasicInfo()
			// 对于先保存，再执行的订单，仅用[order.AppOrdID+order.ClOrdID]会漏掉一些key
			c.rootOrderKeyMap[order.AppOrdID+order.ClOrdID] = true
			c.rootOrderKeyMap[order.AppOrdID] = true
		}
		for _, tradeActionLatestResp := range tradeActionLatestResps {
			key := tradeActionLatestResp.ActionKey
			c.tradeActionLatestRespKeyMap[key] = true
		}
	}

	// 第3步：添加交易动作
	c.tradeActionRespMap = c.attachTradeAction(tradeActionLatestResps, orderMap)
	log.Printf("===> Complete step#3: attachTradeAction tradeActionRespMap.Len=%d\n", len(c.tradeActionRespMap))

	// 第4步：更新交易订单
	// 把发过的成交回报找出来
	syncSendMap := map[string]bool{}
	for _, syncMsg := range syncMsgs {
		if syncMsg.MessageType == OrderCacheSyncMessageType_UpdateByTradeActionResp && syncMsg.TradeActionResp != nil {
			syncSendMap[syncMsg.TradeActionResp.GetCacheKey()] = true
		}
	}
	log.Printf("===> syncSendMap.Len=%d\n", len(syncSendMap))
	c.updateOrderStatusByTradeActionResp(syncSendMap, tradeActionResps, syncMsgs, orderMap, c.directOrderMap)
	log.Printf("===> Complete step#4: updateOrderStatusByTradeActionResp\n")

	log.Printf("Finish Recovering in strict mode, time cost: %v\n", time.Since(begin))

	var directOrders []*types.TraceableTradeOrder
	for _, order := range c.directOrderMap {
		directOrders = append(directOrders, order)
	}
	sort.Slice(directOrders, func(i, j int) bool {
		return directOrders[i].GetBasicInfo().ID < directOrders[j].GetBasicInfo().ID
	})
	// for _, order := range directOrders {
	// 	c.initPositionForNewOrder(order.GetBasicInfo())
	// }

	return newestOffest, oldestOffset, syncMsgs
}

func (c *OrderCache) updateOrderStatusByTradeActionResp(syncSendMap map[string]bool, tradeActionResps []*schema.TradeActionResp, syncMsgs []*OrderCacheSyncMessage, orderMap map[string]*types.TraceableTradeOrder, directOrderMap map[string]*types.TraceableTradeOrder) {

	processTradeActionResp := func(tradeActionResp *schema.TradeActionResp) {

		tradeActionResp.RecoverExtendAttrMap()

		clOrdID := tradeActionResp.ClOrdID
		tradeAction, ok := c.tradeActionRespMap[clOrdID]
		if !ok { // 如果找不到，就构造一个
			traceableOrder := orderMap[clOrdID]
			if traceableOrder != nil {
				log.Printf("cannot find trade action for clOrdId:%s, will create one\n", clOrdID)
				order := traceableOrder.GetBasicInfo()
				// 创建交易动作跟踪对象
				tradeActionLatestResp := &schema.TradeActionLatestResp{
					ActionUser:        order.OrdExecUser,
					ActionTime:        order.OrdStatusUpdateTime,
					ActionMsgTime:     order.OrdStatusUpdateTime,
					ActionType:        order.LatestActionType,
					ActionKey:         order.AppOrdID, // 对于下单委托，action key直接取ClOrdID
					AppOrdID:          order.AppOrdID,
					RootClOrdID:       order.ClOrdID,
					ClOrdID:           order.ClOrdID,
					ChannelCode:       order.ChannelCode,
					StreamInputMsgSeq: order.MsgSeq,
				}

				tradeAction = types.NewTraceableTradeActionResp(traceableOrder, order, tradeActionLatestResp, c.IsSlave())
				c.tradeActionRespMap[tradeActionLatestResp.ClOrdID] = tradeAction
				// 订单内追加交易动作
				traceableOrder.AddTradeAction(tradeAction)
			}
		}

		if tradeAction == nil {
			log.Printf("there's no trade action for %s\n", clOrdID)
			return
		}

		// msgTime, msgSeq := tradeAction.GetLatestMsgTimeAndMsgSeq()
		// // 忽略重复的交易回报
		// if tradeActionResp.MsgTime <= msgTime-100 || tradeActionResp.MsgSeq <= msgSeq {
		// 	return
		// }

		// tradeAction.UpdateByTradeActionResp(tradeActionResp)

		// 直接用tradeAction的update，会没法返回成交回报给应用，并且没法同步状态到OrderReport节点
		// if !tradeAction.UpdateByTradeActionResp(tradeActionResp) {
		// 	log.Printf("===> skip for tradeActionResp, key=%s\n", tradeActionResp.GetCacheKey())
		// 	return
		// }

		key := tradeActionResp.GetCacheKey()
		sendSyncMsg := !syncSendMap[key]
		wg, _ := c.doUpdateByTradeActionResp(tradeActionResp, sendSyncMsg, true)
		if sendSyncMsg {
			syncSendMap[key] = true
		}
		
		if wg!= nil {
			wg.Wait()
		}
	}

	// 1、原来的方案是分别用数据库、kafka的成交回报分别过一遍来更新订单状态
	// // 先用数据库的回报更新状态
	// for _, tradeActionResp := range tradeActionResps {
	// 	processTradeActionResp(tradeActionResp)
	// }

	// // 用syncMsg的，再过一次
	// for _, syncMsg := range syncMsgs {
	// 	if syncMsg.MessageType == OrderCacheSyncMessageType_UpdateByTradeActionResp {
	// 		processTradeActionResp(syncMsg.TradeActionResp)
	// 	}
	// }

	// 2、最新方案是先融合、排序，用合并后的成交回报来更新订单状态
	tradeActionResps = c.mergeTradeActionTradeActionResp(tradeActionResps, syncMsgs)
	log.Printf("===> total tradeActionResps.Len=%d\n", len(tradeActionResps))
	begin := time.Now()

	orderOrResps := c.sortOrderAndResps(directOrderMap, tradeActionResps)
	for _, orderOrResp := range orderOrResps {
		if orderOrResp.tradeOrder != nil {
			c.initPositionForNewOrder(orderOrResp.tradeOrder)
		} else {
			processTradeActionResp(orderOrResp.tradeActionResp)
		}
	}

	// for _, tradeActionResp := range tradeActionResps {
	// 	processTradeActionResp(tradeActionResp)
	// }
	log.Printf("===> finish replay all tradeActionResps, time cost:%v\n", time.Since(begin))
}

func (c *OrderCache) mergeTradeActionTradeActionResp(tradeActionResps []*schema.TradeActionResp, syncMsgs []*OrderCacheSyncMessage) []*schema.TradeActionResp {

	m := make(map[string]*schema.TradeActionResp)

	for _, tradeActionResp := range tradeActionResps {
		m[tradeActionResp.GetCacheKey()] = tradeActionResp
	}

	for _, syncMsg := range syncMsgs {
		if syncMsg.MessageType == OrderCacheSyncMessageType_UpdateByTradeActionResp {
			kafkaTradeActionResp := syncMsg.TradeActionResp
			if _, ok := m[kafkaTradeActionResp.GetCacheKey()]; !ok {
				tradeActionResps = append(tradeActionResps, kafkaTradeActionResp)
				jsData, _ := json.Marshal(kafkaTradeActionResp)
				log.Printf("======> Add kafkaTradeActionResp:%s\n", jsData)
			}
		}
	}

	//排序
	sortTradeActionResps(tradeActionResps)

	return tradeActionResps
}

func (c *OrderCache) attachTradeAction(tradeActionLatestResps []*schema.TradeActionLatestResp, orderMap map[string]*types.TraceableTradeOrder) (tradeActionRespMap map[string]*types.TraceableTradeActionResp) {
	tradeActionRespMap = make(map[string]*types.TraceableTradeActionResp)
	for _, tradeActionLatestResp := range tradeActionLatestResps {

		//log.Printf("======>tradeActionLatestResp: AppOrdID=%s, ClOrdID=%s\n", tradeActionLatestResp.AppOrdID, tradeActionLatestResp.ClOrdID)
		/*order, ok := orderMap[tradeActionLatestResp.ClOrdID]
		if ok {
			tradeAction := types.NewTraceableTradeActionResp(order.GetBasicInfo(), tradeActionLatestResp, c.IsSlave())
			tradeActionRespMap[tradeActionLatestResp.ClOrdID] = tradeAction
			// 订单内追加交易动作
			order.AddTradeAction(tradeAction)
		} else {
			order, ok = orderMap[tradeActionLatestResp.AppOrdID]
			if ok {
				// 未提交的订单，也需要增补交易动作
				tradeAction := types.NewTraceableTradeActionResp(order.GetBasicInfo(), tradeActionLatestResp, c.IsSlave())
				// 订单内追加交易动作
				order.AddTradeAction(tradeAction)
			}
		}*/

		hasClOrdID := false
		order, ok := orderMap[tradeActionLatestResp.ClOrdID]
		if ok {
			hasClOrdID = true
		} else {
			order, ok = orderMap[tradeActionLatestResp.AppOrdID]
		}

		if ok {
			// 未提交的订单，也需要增补交易动作
			tradeAction := types.NewTraceableTradeActionResp(order, order.GetBasicInfo(), tradeActionLatestResp, c.IsSlave())
			// 订单内追加交易动作
			order.AddTradeAction(tradeAction)
			if hasClOrdID {
				//	tradeActionRespMap[tradeActionLatestResp.ClOrdID] = tradeAction
				log.Printf("found order by tradeActionLatestResp.ClOrdID=%s\n", tradeActionLatestResp.ClOrdID)
			}
			tradeActionRespMap[tradeActionLatestResp.ClOrdID] = tradeAction
		} else {
			domain_error.ProcessSevereError(false, 0, nil, nil, fmt.Sprintf("cannot find order by appOrdID %s in strict recover process", tradeActionLatestResp.AppOrdID))
		}
	}
	return
}

func strToMap(str string) map[string]interface{} {
	m := make(map[string]interface{})
	if str == "" {
		return m
	}
	err := json.Unmarshal([]byte(str), &m)
	domain_error.ProcessSevereError(false, 0, nil, err, "fail to conver str to map")
	return m
}

func (c *OrderCache) constructTradeOrderTopology(tradeOrders []*schema.TradeOrder) (rootOrders []*types.TraceableTradeOrder, directOrderMap, orderMap map[string]*types.TraceableTradeOrder) {

	directOrderMap = make(map[string]*types.TraceableTradeOrder)
	orderMap = make(map[string]*types.TraceableTradeOrder)

	// channelTimezoneMap := make(map[string]int)
	// for _, c := range c.applicationCfg.GetTradeChannels() {
	// 	channelTimezoneMap[c.TradeChannel.ChannelCode] = c.TradeChannel.TimeZone
	// }

	channelMap := make(map[string]*domain_cfg.TradeChannelDetails)
	for _, c := range c.applicationCfg.GetTradeChannels() {
		channelMap[c.TradeChannel.ChannelCode] = c
	}

	log.Printf("======> constructTradeOrderTopology, tradeOrders.Len=%d\n", len(tradeOrders))

	for i, tradeOrder := range tradeOrders {

		tradeOrder.ExtendAttrMap = strToMap(tradeOrder.ExtendAttr)
		tradeOrder.AlgParamsMap = strToMap(tradeOrder.AlgParams)

		order := types.NewTraceableTradeOrder(tradeOrder)

		// 判断是否为母单
		if len(tradeOrder.ParentClOrdID) == 0 {
			rootOrders = append(rootOrders, order)
			if tradeOrder.IsDirectOrd && len(tradeOrder.AppOrdID) > 0 {
				directOrderMap[tradeOrder.AppOrdID] = order
			}
		}

		if len(tradeOrder.ClOrdID) == 0 { // 尝试创建ClOrdID
			if statusWithClOrdID[tradeOrder.OrdStatus] || len(tradeOrder.OrdStatus) > 0 && (tradeOrder.OrdStatus[0] >= '0' && tradeOrder.OrdStatus[0] <= '9' || tradeOrder.OrdStatus[0] >= 'A' && tradeOrder.OrdStatus[0] <= 'Z') {
				// tz, ok := channelTimezoneMap[tradeOrder.ChannelCode]
				// if ok {
				// 	dateStr, err := timeutil.GetDateStrForTimeZone(tradeOrder.TransactTime, tz)
				// 	if err != nil {
				// 		domain_error.ProcessSevereError(true, 5, nil, err, "fail to get date string for order")
				// 	}
				// 	tradeOrder.TradeDate = timeutil.Parse8BitDateStrToNum(dateStr)
				// 	tradeOrder.ClOrdID = fmt.Sprintf("%s-%s-%s-%d", dateStr, tradeOrder.SystemCode, tradeOrder.BusinessCode, tradeOrder.ID)
				// } else {
				// 	domain_error.ProcessSevereError(false, 0, nil, fmt.Errorf("fail to get time zone for channel:%s, order:%s", tradeOrder.ChannelCode, tradeOrder.AppOrdID), "fail to get timezone")
				// }
				ch, ok := channelMap[tradeOrder.ChannelCode]
				if ok {
					dateNum := ch.GetCurrentExchangeDate()
					tradeOrder.TradeDate = dateNum
					tradeOrder.ClOrdID = fmt.Sprintf("%d-%s-%s-%d", dateNum, tradeOrder.SystemCode, tradeOrder.BusinessCode, tradeOrder.ID)
				} else {
					domain_error.ProcessSevereError(false, 0, nil, fmt.Errorf("fail to get time zone for channel:%s, order:%s", tradeOrder.ChannelCode, tradeOrder.AppOrdID), "fail to get timezone")
				}
			}
		}

		if len(tradeOrder.ClOrdID) > 0 {
			orderMap[tradeOrder.ClOrdID] = order
		}

		if len(tradeOrder.AppOrdID) > 0 {
			orderMap[tradeOrder.AppOrdID] = order
		}

		log.Printf("======> order#%d, status:%s, clOrdID:%s, AppOrdID:%s, ChannelCode:%s\n", i, tradeOrder.OrdStatus, tradeOrder.ClOrdID, tradeOrder.AppOrdID, tradeOrder.ChannelCode)
	}

	for _, tradeOrder := range tradeOrders {
		if len(tradeOrder.ParentClOrdID) > 0 {
			parentOrder, ok := orderMap[tradeOrder.ParentClOrdID]
			if ok {
				order := orderMap[tradeOrder.ClOrdID]
				if order == nil {
					order = orderMap[tradeOrder.AppOrdID]
				}
				if order != nil {
					parentOrder.AddSubOrder(order)
				}
			}
		}
	}

	return
}

func (c *OrderCache) prepairDataForStrictRecover() (tradeOrders []*schema.TradeOrder, tradeActionLatestResps []*schema.TradeActionLatestResp, tradeActionResps []*schema.TradeActionResp, newestOffest, oldestOffset int64, syncMsgs []*OrderCacheSyncMessage) {
	log.Printf("======> In prepairDataForStrictRecover, c.applicationCfg:%v\n", c)
	var pageSize int64 = 100000
	var err error
	tradeOrders, err = dbutil.PaginateQueryAll[schema.TradeOrder](c.applicationCfg.GetAppDB(), pageSize, app_store.FindAllTradeOrdersInRangeOrderById)
	c.processDataError(err)
	sort.Slice(tradeOrders, func(i, j int) bool {
		return tradeOrders[i].OrdCreateTime < tradeOrders[j].OrdCreateTime
	})

	tradeActionLatestResps, err = dbutil.PaginateQueryAll[schema.TradeActionLatestResp](c.applicationCfg.GetAppDB(), pageSize, app_store.FindAllTradeActionLatestRespsInRangeOrderById)
	c.processDataError(err)
	sort.Slice(tradeActionLatestResps, func(i, j int) bool {
		return tradeActionLatestResps[i].ActionTime < tradeActionLatestResps[j].ActionTime
	})

	tradeActionResps, err = dbutil.PaginateQueryAll[schema.TradeActionResp](c.applicationCfg.GetAppDB(), pageSize, app_store.FindAllTradeActionRespsInRangeOrderById)
	c.processDataError(err)
	/*
		sort.Slice(tradeActionResps, func(i, j int) bool {
			if tradeActionResps[i].MsgTime == tradeActionResps[j].MsgTime {
				log.Printf("MsgTime is equal for tradeActionResps with [ClOrdID=%s ID=%d Status=%s MsgSeq=%d] and  [ClOrdID=%s ID=%d Status=%s MsgSeq=%d] \n", tradeActionResps[i].ClOrdID, tradeActionResps[i].ID, tradeActionResps[i].OrdStatus, tradeActionResps[i].MsgSeq, tradeActionResps[j].ClOrdID, tradeActionResps[j].ID, tradeActionResps[j].OrdStatus, tradeActionResps[j].MsgSeq)
				return tradeActionResps[i].MsgSeq < tradeActionResps[j].MsgSeq
			}
			return tradeActionResps[i].MsgTime < tradeActionResps[j].MsgTime
		})
	*/
	sortTradeActionResps(tradeActionResps)

	for i, tradeActionResp := range tradeActionResps {
		log.Printf("======>#%d tradeActionResp: ID=%d, ClOrdID=%s, Status=%s, TransactTime=%d, MsgTime=%d, MsgSeq=%d\n", i+1, tradeActionResp.ID, tradeActionResp.ClOrdID, tradeActionResp.OrdStatus, tradeActionResp.TransactTime, tradeActionResp.MsgTime, tradeActionResp.MsgSeq)
	}

	newestOffest, oldestOffset, syncMsgs, err = c.getHistoricalOrderCacheSyncMessages()
	c.processDataError(err)

	return
}

func sortTradeActionResps(tradeActionResps []*schema.TradeActionResp) {
	sort.Slice(tradeActionResps, func(i, j int) bool {
		if tradeActionResps[i].MsgTime == tradeActionResps[j].MsgTime {
			log.Printf("MsgTime is equal for tradeActionResps with [ClOrdID=%s ID=%d Status=%s MsgSeq=%d] and  [ClOrdID=%s ID=%d Status=%s MsgSeq=%d] \n", tradeActionResps[i].ClOrdID, tradeActionResps[i].ID, tradeActionResps[i].OrdStatus, tradeActionResps[i].MsgSeq, tradeActionResps[j].ClOrdID, tradeActionResps[j].ID, tradeActionResps[j].OrdStatus, tradeActionResps[j].MsgSeq)
			return tradeActionResps[i].MsgSeq < tradeActionResps[j].MsgSeq
		}
		return tradeActionResps[i].MsgTime < tradeActionResps[j].MsgTime
	})
}

func (c *OrderCache) processDataError(err error) {
	if err != nil && dbutil.IsDbRecordEmptyError(err) {
		err = nil
	}
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to recover data in strict mode!")
	}
}

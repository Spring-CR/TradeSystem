package order_domain

import (
	"fmt"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-core/adapter_registry"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/types"
	"rhino-trade-channel/channel"
)

type OrderExecutor struct {
	channelMap        map[string]channel.TradeChannelInterface
	orderOrchestrator *OrderOrchestrator
	applicationCfg    *domain_cfg.ApplicationCfg
	executorAdapter   OrderExecutorAdapter
}

func NewOrderExecutor(applicationCfg *domain_cfg.ApplicationCfg, channelMap map[string]channel.TradeChannelInterface, orderOrchestrator *OrderOrchestrator) *OrderExecutor {
	inst := &OrderExecutor{applicationCfg: applicationCfg, channelMap: channelMap, orderOrchestrator: orderOrchestrator}
	adapterPath := applicationCfg.GetOrdExecutorAdapterPath()
	if adapterPath != "" {
		_executorAdapter, de, err := adapter_registry.CallAdapterFunction(adapterPath, applicationCfg)
		if err != nil {
			domain_error.ProcessSevereError(true, 5, nil, err, fmt.Sprintf("fail to construct OrderExecutorAdapter for path:%s", adapterPath))
		}

		if de != domain_error.NilDomainError {
			domain_error.ProcessSevereError(true, 5, de.(*domain_error.Error), err, fmt.Sprintf("fail to construct OrderExecutorAdapter for path:%s", adapterPath))
		}

		inst.executorAdapter = _executorAdapter.(OrderExecutorAdapter)
	}
	return inst
}

func (e *OrderExecutor) NewOrderSingle(order *schema.TradeOrder) (duplicatedOrder bool, de *domain_error.Error) {

	channelCode := order.ChannelCode
	tradeChannel, ok := e.channelMap[channelCode]
	if !ok {
		de = domain_error.Build(domain_error.ORDER_EXECUTOR_TRADE_CHANNEL_NOT_FOUND_ERR_CODE, nil, channelCode)
		return
	}

	// 注意，如果是orderID > 0 , 说明是从草稿态编辑，然后再发单，因此，有些属性需要继续保留原始值
	if order.ID > 0 {
		cacheDraft, ok := e.orderOrchestrator.orderCache.GetOrderByAppOrdID(order.AppOrdID)
		if ok {
			// Todo：当TradeOrder数据结果变化时，这里可能需要调整字段
			cacheDraftOrder := cacheDraft.GetBasicInfo()
			order.ClGroupOrdID = cacheDraftOrder.ClGroupOrdID
			order.ClOrdID = cacheDraftOrder.ClOrdID
			order.OrdID = cacheDraftOrder.OrdID
			order.ParentClOrdID = cacheDraftOrder.ParentClOrdID
			order.OrdCreator = cacheDraftOrder.OrdCreator
			order.OrdCreateTime = cacheDraftOrder.OrdCreateTime
			order.OrdDraftDelFlag = cacheDraftOrder.OrdDraftDelFlag
			order.OrdDraftDelUser = cacheDraftOrder.OrdDraftDelUser
			order.OrdDraftDelTime = cacheDraftOrder.OrdDraftDelTime
			order.ReviewFlag = cacheDraftOrder.ReviewFlag
			order.ReviewerScope = cacheDraftOrder.ReviewerScope
			order.Reviewer = cacheDraftOrder.Reviewer
			order.ApproveStatus = cacheDraftOrder.ApproveStatus
			order.ReviewTime = cacheDraftOrder.ReviewTime
			order.WorkerAffinity = cacheDraftOrder.WorkerAffinity
		}
	}

	duplicatedOrder, de = tradeChannel.NewOrderSingle(order, e.orderOrchestrator.orchestrateOrder, e.orderOrchestrator.syncOrder)

	if e.executorAdapter != nil {
		e.executorAdapter.AfterNewOrderSingle(order, duplicatedOrder, de)
	}

	return
}

var (
	statusCouldBeCancel = map[string]bool{
		string(enum.OrdStatus_PendingNew):         true,
		string(enum.OrdStatus_New):                true,
		string(enum.OrdStatus_PartiallyFilled):    true,
		string(enum.OrdStatus_PendingReplace):     true,
		string(enum.OrdStatus_Replaced):           true,
		string(enum.OrdStatus_Suspended):          true,
		string(enum.OrdStatus_Calculated):         true,
		string(enum.OrdStatus_AcceptedForBidding): true,
	}
)

func (e *OrderExecutor) OrderCancelRequest(orderCancelRequest *types.ApplicationOrderCancelRequest) (order *schema.TradeOrder, de *domain_error.Error) {
	appOrdID := orderCancelRequest.AppOrdID
	traceableTradeOrder, ok := e.orderOrchestrator.orderCache.GetOrderByAppOrdID(appOrdID)
	if !ok {
		de = domain_error.Build(domain_error.CANNOT_FIND_ORDER_BY_APP_ORD_ID_ERR_CODE, nil, appOrdID)
		return
	}
	// Todo 当前只能处理直通订单，后续其他类型订单也要支持。
	order = traceableTradeOrder.GetBasicInfo()
	if order.IsDirectOrd {

		// 为了测试，先注释
		supportCancel := statusCouldBeCancel[order.OrdStatus]
		if !supportCancel {
			de = domain_error.Build(domain_error.ORDER_WITH_STATUS_WHICH_CANNOT_BE_CANCEL_ERR_CODE, nil, appOrdID, enum.GetOrdStatusDisplayName(order.OrdStatus))
			return
		}
		// 根据ClOrdID链式设置规律，找到目标ClOrdID
		targetClOrdID := traceableTradeOrder.GetActivedClOrdID()
		channelCode := order.ChannelCode
		tradeChannel, ok := e.channelMap[channelCode]
		if !ok {
			de = domain_error.Build(domain_error.ORDER_EXECUTOR_TRADE_CHANNEL_NOT_FOUND_ERR_CODE, nil, channelCode)
			return
		}
		de = tradeChannel.OrderCancelRequest(orderCancelRequest.ActionUser, orderCancelRequest.ActionTime, orderCancelRequest.ActionKey, targetClOrdID, orderCancelRequest.StreamInputMsgSeq, order, e.orderOrchestrator.cancelDirectOrder, e.orderOrchestrator.syncTradeActionLatestResp)

		if e.executorAdapter != nil {
			e.executorAdapter.AfterOrderCancelRequest(orderCancelRequest, order, de)
		}

		return
	}

	return
}

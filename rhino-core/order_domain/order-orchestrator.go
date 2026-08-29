package order_domain

import (
	"darkpool-common/timeutil"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-core/adapter_registry"
	"rhino-core/domain_cfg"
	"rhino-core/order_domain/order_archive"
	"rhino-core/order_domain/order_cache"
	"rhino-core/order_domain/order_capital"
	"rhino-core/order_domain/order_manager"
	"rhino-core/order_domain/order_position_manager"
	"rhino-core/order_domain/order_purge"
	"rhino-core/order_domain/schedule"
	"rhino-core/schema"
	"rhino-core/store/app_store"
	"rhino-core/types"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type orderOrResp struct {
	isOrder bool
	order   *schema.TradeOrder
	resp    *schema.TradeActionResp
}

type OrderOrchestrator struct {
	applicationCfg             *domain_cfg.ApplicationCfg
	buffer                     chan *orderOrResp
	orderCache                 *order_cache.OrderCache
	orderArchiver              *order_archive.OrderArchiver
	orderPurger                *order_purge.OrderPurger
	directTradeOrderManager    OrderManagerInterface
	algTradeOrderManager       OrderManagerInterface
	instrTradeOrderManager     OrderManagerInterface
	crossDateTradeOrderManager OrderManagerInterface
	capitalCalculator          *order_capital.CapitalCalculator
	positionManager            *order_position_manager.PositionManager
	//positionCalculator         *order_position.PositionCalculator
}

func NewOrderOrchestrator(applicationCfg *domain_cfg.ApplicationCfg, tradeRespCh chan *types.TradeActionRespReturn) *OrderOrchestrator {

	tradeRespChs := []chan *types.TradeActionRespReturn{tradeRespCh}

	/*var positionCalculator *order_position.PositionCalculator
	// 初始化持仓计算器
	if applicationCfg.GetOrdPositionAdapterPath() != "" {
		var de *domain_error.Error
		positionCalculator, de = order_position.NewPositionCalculator(applicationCfg, true, nil)
		if de != nil {
			domain_error.ProcessSevereError(true, 5, de, nil, "fail to create positionCalculator")
		}
		// 增加持仓计算器的订单执行回报channel
		tradeRespChs = append(tradeRespChs, positionCalculator.GetTradeRespCh())

		log.Printf("finish NewPositionCalculator")
	}*/
	var positionManager *order_position_manager.PositionManager
	// 初始化持仓计算器
	if applicationCfg.GetOrdPositionAdapterPath() != "" {
		var de *domain_error.Error
		positionManager, de = order_position_manager.NewPositionManager(applicationCfg, true, nil)
		if de != nil {
			domain_error.ProcessSevereError(true, 5, de, nil, "fail to create positionCalculator")
		}
		// 增加持仓计算器的订单执行回报channel
		//tradeRespChs = append(tradeRespChs, positionCalculator.GetTradeRespCh())
		tradeRespChs = append(tradeRespChs, positionManager.GetTradeRespCh())

		log.Printf("finish NewPositionCalculator")
	}

	var capitalCalculator *order_capital.CapitalCalculator
	// 初始化资金计算器
	if applicationCfg.GetOrdCapitalAdapterPath() != "" {
		var de *domain_error.Error
		capitalCalculator, de = order_capital.NewCapitalCalculator(applicationCfg)
		if de != nil {
			domain_error.ProcessSevereError(true, 5, de, nil, "fail to create capitalCalculator")
		}
		log.Printf("finish NewCapitalCalculator")
	}

	var scheduleAdapter schedule.ScheduleAdapter
	scheduleAdapterPath := applicationCfg.GetScheduleAdapterPath()
	// 初始化调度插件
	if scheduleAdapterPath != "" {

		// 从注册表获取适配器的构造函数（目前，对于apiAdpater，是无参数函数，有参函数需要根据特殊情况来处理了）
		_scheduleAdapter, de, err := adapter_registry.CallAdapterFunction(scheduleAdapterPath, applicationCfg)
		if err != nil {
			domain_error.ProcessSevereError(true, 5, nil, err, "fail to construct ScheduleAdapter from "+scheduleAdapterPath)
		}
		if de != domain_error.NilDomainError {
			domain_error.ProcessSevereError(true, 5, nil, err, "fail to construct ScheduleAdapter from "+scheduleAdapterPath)
		}
		// 获取apiAdapter
		scheduleAdapter = _scheduleAdapter.(schedule.ScheduleAdapter)

		log.Printf("finish create scheduleAdapter")
	}

	inst := &OrderOrchestrator{
		applicationCfg:          applicationCfg,
		buffer:                  make(chan *orderOrResp, 4096),
		orderCache:              order_cache.NewOrderCache(true, applicationCfg, positionManager, tradeRespChs, nil, nil, nil, nil, nil, nil, nil, nil, nil),
		directTradeOrderManager: order_manager.NewDirectTradeOrderManager(),
		capitalCalculator:       capitalCalculator,
		positionManager:         positionManager,
		//positionCalculator:      positionCalculator,
	}

	log.Printf("finish NewOrderOrchestrator")

	// 初始化归档器
	inst.orderArchiver = order_archive.NewOrderArchiver(applicationCfg, inst.orderCache, scheduleAdapter)
	inst.orderPurger = order_purge.NewOrderPurger(applicationCfg, inst.orderCache, inst.orderArchiver, scheduleAdapter, positionManager)

	inst.start()
	return inst
}

func (o *OrderOrchestrator) start() {
	go func() {
		for {
			item := <-o.buffer
			if item.isOrder {
				o.doOrchestrateOrder(item.order)
			} else {
				o.doUpdateOrderStatus(item.resp)
			}
		}
	}()

	o.orderArchiver.Start()
	o.orderPurger.Start()
}

func (o *OrderOrchestrator) orchestrateOrder(tradeOrder *schema.TradeOrder) {
	//o.buffer <- &orderOrResp{isOrder: true, order: tradeOrder}
	o.doOrchestrateOrder(tradeOrder)
}

func (o *OrderOrchestrator) doOrchestrateOrder(tradeOrder *schema.TradeOrder) {
	if tradeOrder.IsDirectOrd {
		o.directTradeOrderManager.OrchestrateOrder(tradeOrder, o.orderCache)
	} else if tradeOrder.IsAlgOrd {
		o.algTradeOrderManager.OrchestrateOrder(tradeOrder, o.orderCache)
	} else if tradeOrder.IsInstrOrd {
		o.instrTradeOrderManager.OrchestrateOrder(tradeOrder, o.orderCache)
	} else if tradeOrder.IsCrossDateOrd {
		o.crossDateTradeOrderManager.OrchestrateOrder(tradeOrder, o.orderCache)
	}
}

func (o *OrderOrchestrator) syncOrder(tradeOrder *schema.TradeOrder) {
	o.orderCache.SyncOrder(tradeOrder)
}

func (o *OrderOrchestrator) doUpdateOrderStatus(resp *schema.TradeActionResp) {
	//time.Sleep(time.Second) // 验证是否处理缓慢时，会引发quickfix异常！Yes！当quickfix对消息处理很慢的时候，会引起异常！
	orderUpdateAttributes := o.orderCache.UpdateByTradeActionResp(resp)
	if len(orderUpdateAttributes) > 0 {
		o.UpdateOrderAttributes(&types.ApplicationOrderAttributeUpdateRequest{
			AppOrdID : resp.AppOrdID,
			UpdateAttributes: orderUpdateAttributes,
		})
	}
}

func (o *OrderOrchestrator) updateOrderStatus(resp *schema.TradeActionResp) {
	o.buffer <- &orderOrResp{isOrder: false, resp: resp}
	//o.doUpdateOrderStatus(resp)
}

func (o *OrderOrchestrator) cancelDirectOrder(tradeActionLatestResp *schema.TradeActionLatestResp) {
	o.directTradeOrderManager.CancelOrder(tradeActionLatestResp, o.orderCache)
}

func (o *OrderOrchestrator) syncTradeActionLatestResp(tradeActionLatestResp *schema.TradeActionLatestResp) {
	o.orderCache.SyncTradeActionLatestResp(tradeActionLatestResp)
}

func (o *OrderOrchestrator) GetOrderByAppOrdID(appOrdID string) (order *schema.TradeOrder, ok bool) {
	v, ok := o.orderCache.GetOrderByAppOrdID(appOrdID)
	if !ok {
		return nil, false
	}
	return v.GetBasicInfo(), true
}

func (o *OrderOrchestrator) GetTraceableOrderByAppOrdID(appOrdID string) (order *types.TraceableTradeOrder, ok bool) {
	v, ok := o.orderCache.GetOrderByAppOrdID(appOrdID)
	if !ok {
		return nil, false
	}
	return v, true
}

// SaveOrderDraft：保存订单草稿。基本逻辑：
//  1. 判断order的id是否为0，为0则是插入操作；插入时，在一个事务中同时插入TradeActionLatestResp对象。
//  2. 如果不是0，按AppOrdID取出订单。找不到订单、订单当前状态不是草稿态的都要报错。保存订单时，需要在一个事务中同时插入TradeActionLatestResp对象，该对象将保存改动前的订单内容。
//  3. 只能用于http，同步API。
//  4. Todo 在apiengine，需要调用适配器接口，提前设置好WorkerAffinity，因为后续的订单处理，可能走stream api。
//     实际上，如果对于一个订单的全部操作都使用http同步API，可以不考虑WorkerAffinity的问题。因为http意味着客户端要等前一步完成，才能继续执行下一步操作，已经是天然串行的模式。
func (o *OrderOrchestrator) SaveOrderDraft(order *schema.TradeOrder, actionType enum.ActionType) (de *domain_error.Error) {
	// 在插入数据库前设置时间
	nowTime := timeutil.ConvertTimeToMilliseconds(time.Now())
	order.OrdDraftUpdateTime = nowTime
	order.OrdStatusUpdateTime = nowTime
	switch actionType {
	case enum.ActionType_Draft:
		order.OrdStatus = string(enum.OrdStatus_Draft)
	case enum.ActionType_SubmitForReview:
		order.OrdStatus = string(enum.OrdStatus_PendingReview)
		order.ApproveStatus = int(enum.ApproveStatus_PendingReview)
		order.Reviewer = ""
	case enum.ActionType_WithdrawForReview:
		order.OrdStatus = string(enum.OrdStatus_Draft)
		order.ApproveStatus = int(enum.ApproveStatus_NotSubmit)
	case enum.ActionType_CancelForReview:
		order.OrdStatus = string(enum.OrdStatus_ReviewCanceled)
		order.ApproveStatus = int(enum.ApproveStatus_NotSubmit)
	case enum.ActionType_ReviewCompleted:
		if order.ApproveStatus == int(enum.ApproveStatus_Approved) {
			order.OrdStatus = string(enum.OrdStatus_ReviewApproved)
		} else if order.ApproveStatus == int(enum.ApproveStatus_Rejected) {
			order.OrdStatus = string(enum.OrdStatus_ReviewRejected)
		}
	}

	order.LatestActionType = string(actionType)
	tx, de1 := dbutil.BeginTx(o.applicationCfg.GetAppDB())
	if de1 != nil {
		return de1
	}
	order.DBInsertTime = timeutil.ConvertTimeToMilliseconds(time.Now())

	// 构建TradeActionLatestResp对象
	tradeActionLatestResp := &schema.TradeActionLatestResp{
		ActionUser:        order.OrdDraftUpdateUser,
		ActionTime:        order.OrdStatusUpdateTime,
		ActionMsgTime:     order.OrdStatusUpdateTime,
		ActionType:        order.LatestActionType,
		AppOrdID:          order.AppOrdID,
		RootClOrdID:       order.ClOrdID,
		ClOrdID:           order.ClOrdID,
		ChannelCode:       order.ChannelCode,
		StreamInputMsgSeq: order.MsgSeq,
	}
	if enum.ActionType_ReviewCompleted == actionType {
		tradeActionLatestResp.ActionUser = order.Reviewer
		order.ReviewTime = tradeActionLatestResp.ActionTime
	}

	tradeActionLatestResp.ActionKey = tradeActionLatestResp.GetCacheKey() // 该actionkey带有时间戳，可以满足多次更新的幂等操作
	doInsert := false
	var traceableTradeOrder *types.TraceableTradeOrder
	var ok bool
	if order.ID == 0 { // 插入订单

		// 订单创建时间
		order.OrdCreateTime = nowTime
		order.OrdCreator = order.OrdDraftUpdateUser

		doInsert = true

		err := app_store.InsertTradeOrder(tx, order)
		if err != nil {
			if dbutil.IsMysqlDuplicateEntryError(err) {
				de = domain_error.Build(domain_error.DUPLICATE_ORDER_ERR_CODE, err, order.AppOrdID)
				dbutil.RollbackTx(tx)
				return
			} else {
				de = domain_error.Build(domain_error.CANNOT_INSERT_ORDER_DRAFT_ERR_CODE, err)
				dbutil.RollbackTx(tx)
				return
			}
		}

		// 插入TradeActionLatestResp记录
		err = app_store.InsertTradeActionLatestResp(tx, tradeActionLatestResp)
		//if err != nil && !dbutil.IsMysqlDuplicateEntryError(err) { // 排除重复插入的错误
		if err != nil { // 重复插入需要报错
			de = domain_error.Build(domain_error.CANNOT_INSERT_ORDER_DRAFT_ERR_CODE, err)
			dbutil.RollbackTx(tx)
			return
		}

	} else { // 保存订单

		// 在内存模型和数据库中，都需要存在该订单才能执行更新逻辑
		traceableTradeOrder, ok = o.orderCache.GetOrderByAppOrdID(order.AppOrdID)
		if !ok {
			de = domain_error.Build(domain_error.CANNOT_FIND_ORDER_BY_APP_ORD_ID_ERR_CODE, nil, order.AppOrdID)
			dbutil.RollbackTx(tx)
			return
		}

		// 加锁，直到最后
		traceableTradeOrder.GetLock().Lock()
		defer traceableTradeOrder.GetLock().Unlock()

		// 状态以OrderCache的为准
		draftBeforeUpdate := traceableTradeOrder.GetBasicInfo()
		// 如果当前不是草稿状态或待审核状态，就不能多次保存
		if draftBeforeUpdate.OrdStatus != string(enum.OrdStatus_Draft) && draftBeforeUpdate.OrdStatus != string(enum.OrdStatus_PendingReview) {
			de = domain_error.Build(domain_error.CANNOT_UPDATE_ORDER_DRAFT_FOR_STATUS_RESON_ERR_CODE, nil)
			log.Printf("current draftBeforeUpdate.OrdStatus:%v\n", draftBeforeUpdate.OrdStatus)
			dbutil.RollbackTx(tx)
			return
		}

		// 保持订单创建等基础信息不变
		order.OrdCreateTime = draftBeforeUpdate.OrdCreateTime
		order.OrdCreator = draftBeforeUpdate.OrdCreator
		order.ChannelCode = draftBeforeUpdate.ChannelCode
		order.DBInsertTime = draftBeforeUpdate.DBInsertTime

		// 取出历史记录保存到当前的TradeOrderLatestResp记录中
		// 当前的 traceableTradeOrder.GetBasicInfo() 就是更新之前的数据，因此，可以不必从数据库获取该信息
		// draftBeforeUpdate, err := app_store.GetTradeOrderByAppOrdId(tx, order.AppOrdID)
		// if err != nil {
		// 	dbutil.RollbackTx(tx)
		// 	de = domain_error.Build(domain_error.CANNOT_GET_TRADE_ORDER_BY_APP_ORD_ID_ERR_CODE, err, order.AppOrdID)
		// 	return
		// }

		// 构建TradeActionLatestResp对象。其中，修改前的内容需要存入。
		draftBeforeUpdateData, err := json.Marshal(draftBeforeUpdate)
		if err != nil {
			de = domain_error.Build(domain_error.CANNOT_UPDATE_ORDER_DRAFT_ERR_CODE, err, order.AppOrdID)
			dbutil.RollbackTx(tx)
			return
		}
		// 保存改动前的记录
		tradeActionLatestResp.OrdDraftBeforeUpdate = string(draftBeforeUpdateData)

		err = app_store.UpdateTradeOrderByAppOrdId(tx, order)
		if err != nil {
			de = domain_error.Build(domain_error.CANNOT_UPDATE_ORDER_DRAFT_ERR_CODE, err)
			dbutil.RollbackTx(tx)
			return
		}

		// 插入TradeActionLatestResp记录
		err = app_store.InsertTradeActionLatestResp(tx, tradeActionLatestResp)
		if err != nil {
			if dbutil.IsMysqlDuplicateEntryError(err) { // 重复更新，直接返回
				dbutil.RollbackTx(tx)
				return
			} else {
				de = domain_error.Build(domain_error.CANNOT_UPDATE_ORDER_DRAFT_ERR_CODE, err)
				dbutil.RollbackTx(tx)
				return
			}
		}
	}

	de = dbutil.CommitTx(tx)
	if de != nil {
		return
	}

	// 最后，需要更新内存模型。需要区分添加；还是更新order的内容
	if doInsert {
		// 需要设置ID，但是预存指令因为没有执行，还没有transacTime，因此也无法构建ClOrdID
		// 在内存模型中新增order记录
		o.orchestrateOrder(order)
	} else {
		// 对于多次草稿保存的情况，用拷贝的方式更新已存在的order记录，并添加TradeOrderLatestResp记录
		//bean.Copy(traceableTradeOrder.GetBasicInfo()).From(order)
		// 从order找出TraceableTradeOrder对象
		//traceableTradeOrder.AddTradeAction(types.NewTraceableTradeActionResp(traceableTradeOrder.GetBasicInfo(), tradeActionLatestResp, o.orderCache.IsSlave()))

		o.orderCache.UpdateTradeOrderDraft(traceableTradeOrder, tradeActionLatestResp, order)
	}

	return
}

// 1. 按AppOrdID取出订单。找不到订单、订单当前状态不是草稿态的都要报错。
// 2. 将订单的状态设置为删除态，保存订单时，需要在一个事务中同时插入TradeActionLatestResp对象，该对象将保存改动前的订单内容。
// 3. 通过TradeActionLatestResp的ActionKey去重。
// 4. OrderCache同步更新。
func (o *OrderOrchestrator) DeleteOrderDraft(orderDraftDeletion *types.ApplicationOrderDraftDeleteRequest) (de *domain_error.Error) {

	// 开启数据库事务
	tx, de1 := dbutil.BeginTx(o.applicationCfg.GetAppDB())
	if de1 != nil {
		return de1
	}

	// 在内存模型和数据库中，都需要存在该订单才能执行更新逻辑
	traceableTradeOrder, ok := o.orderCache.GetOrderByAppOrdID(orderDraftDeletion.AppOrdID)
	if !ok {
		de = domain_error.Build(domain_error.CANNOT_FIND_ORDER_BY_APP_ORD_ID_ERR_CODE, nil, orderDraftDeletion.AppOrdID)
		dbutil.RollbackTx(tx)
		return
	}

	log.Printf("get traceableTradeOrder with AppOrdID=%s in OrderCache\n", traceableTradeOrder.GetBasicInfo().AppOrdID)

	// 加锁，直到最后
	traceableTradeOrder.GetLock().Lock()
	defer traceableTradeOrder.GetLock().Unlock()

	log.Printf("lock to DeleteOrderDraft\n")

	// 状态以OrderCache的为准
	orderDraft := traceableTradeOrder.GetBasicInfo()
	// 如果当前不是草稿状态，就不能删除
	if orderDraft.OrdStatus != string(enum.OrdStatus_Draft) {
		de = domain_error.Build(domain_error.CANNOT_DELETE_ORDER_DRAFT_FOR_STATUS_RESON_ERR_CODE, nil)
		dbutil.RollbackTx(tx)
		return
	}

	// 构建TradeActionLatestResp对象
	tradeActionLatestResp := &schema.TradeActionLatestResp{
		ActionUser:        orderDraftDeletion.ActionUser,
		ActionTime:        orderDraftDeletion.ActionTime,
		ActionMsgTime:     orderDraftDeletion.ActionTime,
		ActionType:        string(enum.ActionType_Delete),
		AppOrdID:          orderDraft.AppOrdID,
		RootClOrdID:       orderDraft.ClOrdID,
		ClOrdID:           orderDraft.ClOrdID,
		ChannelCode:       orderDraft.ChannelCode,
		StreamInputMsgSeq: orderDraft.MsgSeq,
	}
	tradeActionLatestResp.ActionKey = tradeActionLatestResp.GetCacheKey() // 该actionkey带有时间戳，可以满足多次更新的幂等操作

	// 先插入tradeActionLatestResp，利用ActionKey实现幂等
	err := app_store.InsertTradeActionLatestResp(tx, tradeActionLatestResp)

	log.Printf("InsertTradeActionLatestResp for DeleteOrderDraft\n")

	if err != nil {
		if dbutil.IsMysqlDuplicateEntryError(err) { // 重复删除，直接返回
			dbutil.RollbackTx(tx)
			return
		} else {
			de = domain_error.Build(domain_error.CANNOT_DELETE_ORDER_DRAFT_ERR_CODE, err)
			dbutil.RollbackTx(tx)
			return
		}
	}

	// 数据库处理失败时，需要进行状态的回滚
	oldOrdDraftDelFlag := orderDraft.OrdDraftDelFlag
	oldOrdDraftDelTime := orderDraft.OrdDraftDelTime
	oldOrdDraftDelUser := orderDraft.OrdDraftDelUser
	oldOrdStatusUpdateTime := orderDraft.OrdStatusUpdateTime
	oldOrdStatus := orderDraft.OrdStatus
	oldLatestActionType := orderDraft.LatestActionType
	rollBackStatus := func() {
		orderDraft.OrdDraftDelFlag = oldOrdDraftDelFlag
		orderDraft.OrdDraftDelTime = oldOrdDraftDelTime
		orderDraft.OrdDraftDelUser = oldOrdDraftDelUser
		orderDraft.OrdStatusUpdateTime = oldOrdStatusUpdateTime
		orderDraft.OrdStatus = oldOrdStatus
		orderDraft.LatestActionType = oldLatestActionType
	}

	// 设置order草稿删除的相关信息
	orderDraft.OrdDraftDelFlag = 1
	orderDraft.OrdDraftDelTime = orderDraftDeletion.ActionTime
	orderDraft.OrdDraftDelUser = orderDraftDeletion.ActionUser
	orderDraft.OrdStatusUpdateTime = orderDraftDeletion.ActionTime
	orderDraft.OrdStatus = string(enum.OrdStatus_Draft_Deleted)
	orderDraft.LatestActionType = string(enum.ActionType_Delete)

	// 从数据库更新order
	err = app_store.UpdateTradeOrderByAppOrdId(tx, orderDraft)
	if err != nil {
		de = domain_error.Build(domain_error.CANNOT_UPDATE_ORDER_DRAFT_ERR_CODE, err)
		dbutil.RollbackTx(tx)
		rollBackStatus()
		return
	}

	log.Println("UpdateTradeOrderByAppOrdId")

	de = dbutil.CommitTx(tx)
	if de != nil {
		return
	}

	log.Printf("CommitTx for UpdateTradeOrderByAppOrdId")

	// 最后，需要更新内存模型
	o.orderCache.DeleteTradeOrderDraft(traceableTradeOrder, tradeActionLatestResp, orderDraft)

	log.Printf("update  orderCache after DeleteOrderDraft\n")

	return
}

// 强制归档
func (o *OrderOrchestrator) ForceArchiving() {
	o.orderArchiver.ForceArchiving()
}

func (o *OrderOrchestrator) ForcePurging() {
	o.orderPurger.ForcePurging()
}

func (o *OrderOrchestrator) GetOrderCache() *order_cache.OrderCache {
	return o.orderCache
}

// func (o *OrderOrchestrator) GetPositionCalculator() *order_position.PositionCalculator {
// 	return o.positionCalculator
// }

func (o *OrderOrchestrator) GetPositionManager() *order_position_manager.PositionManager {
	return o.positionManager
}

func (o *OrderOrchestrator) GetCapitalCalculator() *order_capital.CapitalCalculator {
	return o.capitalCalculator
}

package order_cache

import (
	"darkpool-common/bean"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/kafka"
	"rhino-core/adapter_registry"
	"rhino-core/domain_cfg"
	"rhino-core/order_domain/order_position_manager"
	"rhino-core/order_domain/order_status_adapter"
	"rhino-core/schema"
	"rhino-core/types"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// 这里是用于分析延迟
	totalAddRootOrderDuration         int64
	totalAddRootOrderWithLockDuration int64
	totalCreateOrderObjectDuration    int64
)

type OrderCache struct {
	//lock               *sync.RWMutex
	master                 bool
	applicationCfg         *domain_cfg.ApplicationCfg
	rootOrdersLock         *sync.RWMutex
	directOrderMapLock     *sync.RWMutex
	tradeActionRespMapLock *sync.RWMutex
	rootOrders             []*types.TraceableTradeOrder               // 根订单，即母单
	directOrderMap         map[string]*types.TraceableTradeOrder      // 直接订单(为直通订单类型的母单，子单也有直通订单，这里特指的是母单)map，以ClOrdID为key(改成以 AppOrderID 为key)
	tradeActionRespMap     map[string]*types.TraceableTradeActionResp // 交易回报map，以ClOrdID为key。注意，对于改单操作，这里的ClOrdID就不一定是母单，而是由系统指定的用于将交易动作和后续交易回报进行关联的id。
	tradeRespCh            []chan *types.TradeActionRespReturn
	producer               *kafka.TenaciousProducerWithBuffer //*kafka.TenaciousProducer
	consumer               *kafka.TenaciousConsumerWithBuffer //*kafka.BlocklessConsumer
	positionManager        *order_position_manager.PositionManager

	// 对于master，在初始化时从数据库和channel加载数据初始化；对于slave，从kafka最老的一条数据开始加载数据初始化
	tradeActionLatestRespKeyMapLock *sync.RWMutex
	tradeActionLatestRespKeyMap     map[string]bool
	rootOrderKeyMap                 map[string]bool

	// Todo 还需要补充关于order draft update、order draft delete的相关handler
	// 这几个是slave角色下的handler
	_afterAddRootTradeOrder            func(order *types.TraceableTradeOrder)
	_afterUpdateByTradeActionResp      func(tradeActionResp *schema.TradeActionResp, traceableTradeActionResp *types.TraceableTradeActionResp)
	_afterAddTradeActionForDirectOrder func(tradeAction *types.TraceableTradeActionResp)
	_afterUpdateTradeOrderDraft        func(tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder)
	_afterDeleteTradeOrderDraft        func(tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder)
	_afterUpdateTradeOrderAttributes   func(appOrdID string, updateAttrs map[string]interface{})
	_afterReset                        func()
	_afterSyncOrder                    func(tradeOrder *schema.TradeOrder)
	_afterSyncTradeActionLatestResp    func(tradeActionLatestResp *schema.TradeActionLatestResp)

	lookUpTraceableTradeActionResp        func(directOrderMap map[string]*types.TraceableTradeOrder, tradeActionResp *schema.TradeActionResp, tradeActionRespMap map[string]*types.TraceableTradeActionResp) (traceableTradeActionResp *types.TraceableTradeActionResp, ok bool)
	updateTradeActionLatestResp           func(tradeActionResp *schema.TradeActionResp, traceableTradeActionResp *types.TraceableTradeActionResp)
	updateOrderStatusForOrderCancelReject func(tradeActionResp *schema.TradeActionResp, tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder)
	updateOrderStatusForExecutionReport   func(tradeActionResp *schema.TradeActionResp, tradeActionLatestResp *schema.TradeActionLatestResp, tradeActionRespList []*schema.TradeActionResp, order *schema.TradeOrder, traceableTradeOrder *types.TraceableTradeOrder) (orderUpdateAttributes map[string]interface{})
}

func NewOrderCache(master bool, applicationCfg *domain_cfg.ApplicationCfg, positionManager *order_position_manager.PositionManager, tradeRespCh []chan *types.TradeActionRespReturn,
	afterAddRootTradeOrder func(order *types.TraceableTradeOrder),
	afterUpdateByTradeActionResp func(tradeActionResp *schema.TradeActionResp, traceableTradeActionResp *types.TraceableTradeActionResp),
	afterAddTradeActionForDirectOrder func(tradeAction *types.TraceableTradeActionResp),
	afterUpdateTradeOrderDraft func(tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder),
	afterDeleteTradeOrderDraft func(tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder),
	afterUpdateTradeOrderAttributes func(appOrdID string, updateAttrs map[string]interface{}),
	afterReset func(),
	afterSyncOrder func(tradeOrder *schema.TradeOrder),
	afterSyncTradeActionLatestResp func(tradeActionLatestResp *schema.TradeActionLatestResp)) *OrderCache {

	inst := &OrderCache{
		master:                             master,
		applicationCfg:                     applicationCfg,
		positionManager:                    positionManager,
		rootOrdersLock:                     &sync.RWMutex{},
		directOrderMapLock:                 &sync.RWMutex{},
		tradeActionRespMapLock:             &sync.RWMutex{},
		directOrderMap:                     map[string]*types.TraceableTradeOrder{},
		tradeActionRespMap:                 map[string]*types.TraceableTradeActionResp{},
		tradeRespCh:                        tradeRespCh,
		_afterAddRootTradeOrder:            afterAddRootTradeOrder,
		_afterUpdateByTradeActionResp:      afterUpdateByTradeActionResp,
		_afterAddTradeActionForDirectOrder: afterAddTradeActionForDirectOrder,
		_afterUpdateTradeOrderDraft:        afterUpdateTradeOrderDraft,
		_afterDeleteTradeOrderDraft:        afterDeleteTradeOrderDraft,
		_afterUpdateTradeOrderAttributes:   afterUpdateTradeOrderAttributes,
		_afterReset:                        afterReset,
		_afterSyncOrder:                    afterSyncOrder,
		_afterSyncTradeActionLatestResp:    afterSyncTradeActionLatestResp,
		lookUpTraceableTradeActionResp:     defaultLookUpTraceableTradeActionResp,
	}

	adapterPath := applicationCfg.GetOrdStatusAdapterPath()
	if adapterPath != "" {

		// 从注册表获取适配器的构造函数（目前，对于apiAdpater，是无参数函数，有参函数需要根据特殊情况来处理了）
		_orderStatusAdapter, de, err := adapter_registry.CallAdapterFunction(adapterPath, applicationCfg)
		if err != nil {
			domain_error.ProcessSevereError(true, 5, nil, err, fmt.Sprintf("fail to construct OrderStatusAdapter for path:%s", adapterPath))
		}

		if de != domain_error.NilDomainError {
			log.Printf("de:%+v\n", de)
			return nil
		}

		// 获取apiAdapter
		orderStatusAdapter := _orderStatusAdapter.(order_status_adapter.OrderStatusAdapter)
		log.Println("finish get orderStatusAdapter")
		inst.lookUpTraceableTradeActionResp = orderStatusAdapter.LookUpTraceableTradeActionResp
		inst.updateTradeActionLatestResp = orderStatusAdapter.UpdateTradeActionLatestResp
		inst.updateOrderStatusForOrderCancelReject = orderStatusAdapter.UpdateOrderStatusForOrderCancelReject
		inst.updateOrderStatusForExecutionReport = orderStatusAdapter.UpdateOrderStatusForExecutionReport
	}

	if master {
		inst.configForMaster()
	} else {
		inst.configForSlave()
	}

	go func() {
		for {
			inst.tradeActionRespMapLock.RLock()
			n := 0
			for _, resp := range inst.tradeActionRespMap {
				n += len(resp.GetTradeActionRespList())
			}
			log.Printf("IsMaster: %v, Order slice size: %d, Order map size: %d, Order trade map size: %d, Order trade resp size: %d\n", inst.master, len(inst.rootOrders), len(inst.directOrderMap), len(inst.tradeActionRespMap), n)
			inst.tradeActionRespMapLock.RUnlock()
			time.Sleep(15 * time.Second)
		}
	}()

	log.Println("finish NewOrderCache")

	return inst
}

var (
	ordStatusDraft          = string(enum.OrdStatus_Draft)
	ordStatusPendingReview  = string(enum.OrdStatus_PendingReview)
	ordStatusReviewApproved = string(enum.OrdStatus_ReviewApproved)
	ordStatusReviewRejected = string(enum.OrdStatus_ReviewRejected)
	ordStatusReviewCanceled = string(enum.OrdStatus_ReviewCanceled)
)

// draft == false && insert == false时，代表的是是执行order draft
func (c *OrderCache) AddRootTradeOrder(order *schema.TradeOrder) {

	if !c.master { // 只对slave角色有效
		c.rootOrdersLock.RLock()
		if c.rootOrderKeyMap[order.AppOrdID+order.ClOrdID] { // 因订单草稿态时ClOrdID是空字符串，如果仅用AppOrdID，当订单从草稿态变成执行态时，会漏掉一层转换逻辑
			// 重复处理，需要立即返回
			c.rootOrdersLock.RUnlock()
			log.Printf("===> skip for slave, AddRootTradeOrder, return for duplicate [order.AppOrdID+order.ClOrdID] key:%s\n", order.AppOrdID+order.ClOrdID)
			return
		}
		c.rootOrdersLock.RUnlock()
	}

	beginCreateOrderObject := time.Now()
	// 判断是否草稿
	draft := false
	switch order.OrdStatus {
	case ordStatusDraft, ordStatusPendingReview, ordStatusReviewApproved, ordStatusReviewRejected, ordStatusReviewCanceled:
		draft = true
	}
	log.Printf("AddRootTradeOrder, order.OrdStatus=%s, draft=%v, order.AppOrdID=%s\n", order.OrdStatus, draft, order.AppOrdID)
	var newOrder *types.TraceableTradeOrder
	// 如果是draft或者直接下单的情况，需要创建母单
	if order.DBInsertOnOrdExec || draft {

		newOrder = types.NewTraceableTradeOrder(order)

		c.rootOrdersLock.Lock()
		// 追加母单移到这里
		c.rootOrders = append(c.rootOrders, newOrder)
		if !c.master { // 只有slave角色需要设置
			c.rootOrderKeyMap[order.AppOrdID+order.ClOrdID] = true
		}
		c.rootOrdersLock.Unlock()
	} else {
		// 表明是从草稿订单下单
		// 从draft order里面，找到订单并更新其order基本信息
		var ok bool
		newOrder, ok = c.GetOrderByAppOrdID(order.AppOrdID)
		if ok {
			// 更新基本信息
			bean.Copy(newOrder.GetBasicInfo()).From(order)
		} else {
			domain_error.ProcessSevereError(false, 0, nil, fmt.Errorf("fail to get draft order by AppOrdId=%s", order.AppOrdID), "fail to get draft order by AppOrdId")
			newOrder = types.NewTraceableTradeOrder(order)
		}
	}

	// 否则，更新母单中的order信息
	var actionUser string
	if draft {
		actionUser = order.OrdDraftUpdateUser
		if order.LatestActionType == string(enum.ActionType_ReviewCompleted) {
			actionUser = order.Reviewer
		}
	} else {
		actionUser = order.OrdExecUser
	}
	// 创建交易动作跟踪对象
	tradeActionLatestResp := &schema.TradeActionLatestResp{
		ActionUser:        actionUser,
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
	if draft { // 对于draft，需要重新设置ActionKey
		tradeActionLatestResp.ActionKey = tradeActionLatestResp.GetCacheKey()
	}
	traceableTradeActionResp := types.NewTraceableTradeActionResp(newOrder, newOrder.GetBasicInfo(), tradeActionLatestResp, c.IsSlave())

	//newOrder.tradeActions = append(newOrder.tradeActions, traceableTradeActionResp)
	newOrder.InitTradeAction(traceableTradeActionResp)

	atomic.AddInt64(&totalCreateOrderObjectDuration, int64(time.Since(beginCreateOrderObject)))

	beginAddRootOrderWithLockDuration := time.Now()

	//c.lock.Lock()
	beginAddRootOrderDuration := time.Now()

	// 追加母单
	// 需要将追加母单的逻辑往上迁移，如果订单是从草稿态转执行态，此步是不需要增加的
	// c.rootOrdersLock.Lock()
	// c.rootOrders = append(c.rootOrders, newOrder)
	// if !c.master { // 只有slave角色需要设置
	// 	c.rootOrderKeyMap[order.AppOrdID] = true
	// }
	// c.rootOrdersLock.Unlock()

	c.directOrderMapLock.Lock()
	c.directOrderMap[order.AppOrdID] = newOrder // 改成以 AppOrdID 为key，更利于应用层使用
	c.directOrderMapLock.Unlock()

	//log.Printf("AddRootTradeOrder, DBInsertOnOrdExec=%v, ClOrdID=%s, AppOrdID=%s, draft=%v\n", order.DBInsertOnOrdExec, order.ClOrdID, order.AppOrdID, draft)

	if !draft { // 如果不是草稿，表明是已经提交的订单，是有ClOrdID的，使用ClOrdID作为key
		// 追加交易动作
		c.tradeActionRespMapLock.Lock()
		c.tradeActionRespMap[order.ClOrdID] = traceableTradeActionResp
		c.tradeActionRespMapLock.Unlock()
	} // 对于草稿态，因为tradeActionRespMap是用于跟踪回报的，由于草稿订单不可能存在回报，所以没有必要往tradeActionRespMap增加记录了。

	atomic.AddInt64(&totalAddRootOrderDuration, int64(time.Since(beginAddRootOrderDuration)))

	//c.lock.Unlock()

	atomic.AddInt64(&totalAddRootOrderWithLockDuration, int64(time.Since(beginAddRootOrderWithLockDuration)))

	c.afterAddRootTradeOrder(newOrder)
}

func (c *OrderCache) UpdateByTradeActionResp(tradeActionResp *schema.TradeActionResp)(orderUpdateAttributes map[string]interface{}) {
	_, orderUpdateAttributes = c.doUpdateByTradeActionResp(tradeActionResp, true, false)
	return
}

func (c *OrderCache) doUpdateByTradeActionResp(tradeActionResp *schema.TradeActionResp, sendSyncMsg bool, duringRecover bool) (wg *sync.WaitGroup, orderUpdateAttributes map[string]interface{}) {

	c.tradeActionRespMapLock.RLock()
	//traceableTradeActionResp, ok := c.tradeActionRespMap[tradeActionResp.ClOrdID]
	traceableTradeActionResp, ok := c.lookUpTraceableTradeActionResp(c.directOrderMap, tradeActionResp, c.tradeActionRespMap)
	c.tradeActionRespMapLock.RUnlock()

	if !ok {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, fmt.Errorf("cannot get traceableTradeActionResp from tradeActionRespMap by ClOrdID:%s", tradeActionResp.ClOrdID))
		domain_error.ProcessSevereError(false, 0, de, nil, "error on UpdateByTradeActionResp")
		return
	}

	// 确保对tradeActionResp的处理是幂等的，但是 MsgSeq 不应只跟MsgSeq有关，因为 fix server可能会重置session，需要同时考虑ClOrdID和ExecID。
	// 结合消息时间戳msgTime一起判断会更好，如果消息差大于100毫秒，表示可以明确前后顺序了（我司局域网内NTP的精度小于0.5ms）
	// 实际上，中信QFII可以确保execId每天不一样，其加入了日期；fidessa则是通过持续递增，保证execId唯一
	// 对于刚插入的tradeActionLatestResp，因 MsgTime 和 MsgSeq 都为0，不会进入一下if分支内
	// msgTime, msgSeq := traceableTradeActionResp.GetLatestMsgTimeAndMsgSeq()
	// // 忽略重复的交易回报
	// if traceableTradeActionResp == nil || tradeActionResp.MsgTime <= msgTime-100 || tradeActionResp.MsgSeq <= msgSeq {
	// 	if traceableTradeActionResp != nil {
	// 		log.Printf("===> skip for tradeActionResp, tradeActionResp.MsgSeq: %d, traceableTradeActionResp.tradeActionLatestResp.MsgSeq: %d\n", tradeActionResp.MsgSeq, msgSeq)
	// 	}
	// 	return
	// }

	// traceableTradeActionResp.UpdateByTradeActionResp(tradeActionResp)
	var added bool
	if orderUpdateAttributes, added = traceableTradeActionResp.UpdateByTradeActionResp(tradeActionResp, c.updateTradeActionLatestResp, c.updateOrderStatusForOrderCancelReject, c.updateOrderStatusForExecutionReport); !added {
		log.Printf("===> skip for tradeActionResp, key=%s\n", tradeActionResp.GetCacheKey())
		return
	}

	if c.tradeRespCh != nil {
		log.Printf("===> push traceableTradeActionResp to channel, tradeActionResp.ClOrdID:%s, tradeActionResp.AppOrdID:%s, tradeActionResp.ExecID:%s, OrigClOrdID:%s, OrdStatus:%s, LastShares:%v\n", tradeActionResp.ClOrdID, tradeActionResp.AppOrdID, tradeActionResp.ExecID, tradeActionResp.OrigClOrdID, tradeActionResp.OrdStatus, tradeActionResp.LastShares)
		// 发送到channel
		if duringRecover {
			wg = &sync.WaitGroup{}
			wg.Add(c.GetTradeRespChSize())
			log.Printf("create WaitGroup, size=%d\n", c.GetTradeRespChSize())
		}
		for _, tradeRespCh := range c.tradeRespCh {
			tradeActionRespReturn := types.NewTradeActionRespReturn(traceableTradeActionResp, tradeActionResp)
			if duringRecover {
				tradeActionRespReturn.WaitGroup = wg
			}
			tradeRespCh <- tradeActionRespReturn
		}
	}

	if sendSyncMsg {
		c.afterUpdateByTradeActionResp(tradeActionResp, traceableTradeActionResp)
	}

	return
}

func (c *OrderCache) GetTradeRespChSize() int {
	return len(c.tradeRespCh)
}

// Todo 目前只是从directOrder里返回，后续也要查找其他类型的Order
func (c *OrderCache) GetOrderByAppOrdID(appOrdID string) (traceableTradeOrder *types.TraceableTradeOrder, ok bool) {
	c.directOrderMapLock.RLock()
	traceableTradeOrder, ok = c.directOrderMap[appOrdID]
	c.directOrderMapLock.RUnlock()
	if ok {
		return
	}
	return
}

// 给直通订单添加交易动作，对于新建委托，在InitTradeAction已经添加了交易动作，这里则通常用于Cancel、Update的场景，请注意加锁
func (c *OrderCache) AddTradeActionForDirectOrder(tradeActionLatestResp *schema.TradeActionLatestResp) {

	var key string
	// 如果是slave，检查是否已经操作过，如果是重复的动作，则因该忽略
	if !c.master {
		key = tradeActionLatestResp.ActionKey
		c.tradeActionLatestRespKeyMapLock.RLock()
		if c.tradeActionLatestRespKeyMap[key] { // 重复输入
			c.tradeActionLatestRespKeyMapLock.RUnlock()
			log.Printf("===> skip for slave, AddTradeActionForDirectOrder, return for duplicate action key:%s\n", key)
			return
		}
		c.tradeActionLatestRespKeyMapLock.RUnlock()
	}

	c.directOrderMapLock.RLock()
	traceableTradeOrder := c.directOrderMap[tradeActionLatestResp.AppOrdID]
	c.directOrderMapLock.RUnlock()

	if traceableTradeOrder == nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, fmt.Errorf("cannot get traceableTradeOrder from directOrderMap by AppOrdID:%s", tradeActionLatestResp.AppOrdID))
		domain_error.ProcessSevereError(false, 0, de, nil, "error on AddTradeActionForDirectOrder")
		return
	}

	// 如果是slave，更新tradeActionLatestRespKeyMap
	if !c.master {
		c.tradeActionLatestRespKeyMapLock.Lock()
		if c.tradeActionLatestRespKeyMap[key] { // 重复输入(RUnlock之后，还需要在写锁内再判断一次才严谨)
			c.tradeActionLatestRespKeyMapLock.Unlock()
			return
		}
		c.tradeActionLatestRespKeyMap[key] = true
		c.tradeActionLatestRespKeyMapLock.Unlock()
	}

	newTradeAction := types.NewTraceableTradeActionResp(traceableTradeOrder, traceableTradeOrder.GetBasicInfo(), tradeActionLatestResp, c.IsSlave())
	traceableTradeOrder.AddTradeAction(newTradeAction)

	if tradeActionLatestResp.ClOrdID != "" {
		c.tradeActionRespMapLock.Lock()
		c.tradeActionRespMap[tradeActionLatestResp.ClOrdID] = newTradeAction
		c.tradeActionRespMapLock.Unlock()
	}

	c.afterAddTradeActionForDirectOrder(newTradeAction)
}

func GetAllDurations() (int64, int64, int64) {
	return totalAddRootOrderDuration, totalAddRootOrderWithLockDuration, totalCreateOrderObjectDuration
}

func (c *OrderCache) IsMaster() bool {
	return c.master
}

func (c *OrderCache) IsSlave() bool {
	return !c.master
}

func (c *OrderCache) IsDuplicateOrder(order *schema.TradeOrder) bool {
	appOrdID := order.AppOrdID
	if appOrdID == "" {
		return false
	}
	traceableTradeOrder, ok := c.GetOrderByAppOrdID(appOrdID)
	return ok && traceableTradeOrder.GetBasicInfo().ClOrdID != ""
}

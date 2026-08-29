package types

import (
	"log"
	"rhino-common/enum"
	"rhino-core/schema"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	OrderCancelRejectPrefix = "OrderCancelReject:"
)

type TraceableTradeActionResp struct {
	lock                  *sync.RWMutex
	traceableTradeOrder   *TraceableTradeOrder
	parent                *schema.TradeOrder
	tradeActionLatestResp *schema.TradeActionLatestResp
	tradeActionRespList   []*schema.TradeActionResp
	tradeActionRespMap    map[string]*schema.TradeActionResp
	slave                 bool
	lastTradeActionResp   unsafe.Pointer
}

func NewTraceableTradeActionResp(traceableTradeOrder *TraceableTradeOrder, parentOrder *schema.TradeOrder, tradeActionLatestResp *schema.TradeActionLatestResp, slave bool) *TraceableTradeActionResp {
	return &TraceableTradeActionResp{
		lock:                  &sync.RWMutex{},
		traceableTradeOrder:   traceableTradeOrder,
		parent:                parentOrder,
		tradeActionLatestResp: tradeActionLatestResp,
		tradeActionRespMap:    make(map[string]*schema.TradeActionResp),
		slave:                 slave,
	}
}

// 确保高效GC
func (r *TraceableTradeActionResp) Dispose() {
	/*
		2025/06/19 05:32:12 FromAdmin: 8=FIX.4.2|9=60|35=0|49=GFGFIX|56=GFHQTMP|34=21389|52=20250618-21:32:12.913|10=120|
		panic: runtime error: invalid memory address or nil pointer dereference
		[signal SIGSEGV: segmentation violation code=0x1 addr=0x10 pc=0x862cd5]

		goroutine 109 [running]:
		sync.(*RWMutex).RLock(...)
			/usr/local/go/src/sync/rwmutex.go:69
		rhino-core/types.(*TraceableTradeActionResp).GetTradeActionRespList(0x37e11d600?)
			/Users/spring/Documents/workspace_go/src/rhino-core/types/traceable-order.go:76 +0x35
		rhino-core/order_domain/order_cache.NewOrderCache.func1()
			/Users/spring/Documents/workspace_go/src/rhino-core/order_domain/order_cache/order-cache.go:96 +0x245
		created by rhino-core/order_domain/order_cache.NewOrderCache in goroutine 1
			/Users/spring/Documents/workspace_go/src/rhino-core/order_domain/order_cache/order-cache.go:91 +0x2f6

		遇到上述异常，因此不能将r.lock设置为空
	*/

	r.lock.Lock()
	defer r.lock.Unlock()

	//r.lock = nil
	r.parent = nil
	r.tradeActionLatestResp = nil
	for i := range r.tradeActionRespList {
		r.tradeActionRespList[i] = nil
	}
	r.tradeActionRespList = nil
	r.lastTradeActionResp = nil

	// 清理map资源
	for k := range r.tradeActionRespMap {
		delete(r.tradeActionRespMap, k)
	}
	r.tradeActionRespMap = make(map[string]*schema.TradeActionResp)
}

func (r *TraceableTradeActionResp) GetTradeOrder() *schema.TradeOrder {
	return r.parent
}

func (r *TraceableTradeActionResp) GetLastTradeActionResp() *schema.TradeActionResp {
	var resp *schema.TradeActionResp
	r.lock.RLock()
	l := r.tradeActionRespList
	if len(l) > 0 {
		resp = l[len(l)-1]
	}
	r.lock.RUnlock()
	return resp
}

// 只适用于slave角色
func (r *TraceableTradeActionResp) GetLastTradeActionRespByAtomicForSlaveOnly() *schema.TradeActionResp {
	return (*schema.TradeActionResp)(atomic.LoadPointer(&r.lastTradeActionResp)) // 原子地获取指针值并转换为正确的类型
}

func (r *TraceableTradeActionResp) GetTradeActionRespList() []*schema.TradeActionResp {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.tradeActionRespList
}

func (r *TraceableTradeActionResp) GetTradeActionRespListWithoutLock() []*schema.TradeActionResp {
	return r.tradeActionRespList
}

func (r *TraceableTradeActionResp) GetLatestMsgTimeAndMsgSeq() (msgTime int64, msgSeq int64) {
	msgTime = atomic.LoadInt64(&r.tradeActionLatestResp.MsgTime)
	msgSeq = atomic.LoadInt64(&r.tradeActionLatestResp.MsgSeq)
	return
}

// 需要分发状态：TradeActionResp
func (traceableTradeActionResp *TraceableTradeActionResp) UpdateByTradeActionResp(tradeActionResp *schema.TradeActionResp,
	updateTradeActionLatestResp func(tradeActionResp *schema.TradeActionResp, traceableTradeActionResp *TraceableTradeActionResp),
	updateOrderStatusForOrderCancelReject func(tradeActionResp *schema.TradeActionResp, tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder),
	updateOrderStatusForExecutionReport func(tradeActionResp *schema.TradeActionResp, tradeActionLatestResp *schema.TradeActionLatestResp, tradeActionRespList []*schema.TradeActionResp, order *schema.TradeOrder, traceableTradeOrder *TraceableTradeOrder) (orderUpdateAttributes map[string]interface{})) (orderUpdateAttributes map[string]interface{}, added bool) {

	traceableTradeActionResp.lock.Lock()

	// 如果存在，就不要更新了
	key := tradeActionResp.GetCacheKey()
	_, ok := traceableTradeActionResp.tradeActionRespMap[key]
	// 如果回报已经存在，则不再进行处理
	if ok {
		traceableTradeActionResp.lock.Unlock()
		added = false
		return
	}

	// 设置AppOrdID
	tradeActionResp.AppOrdID = traceableTradeActionResp.tradeActionLatestResp.AppOrdID

	if traceableTradeActionResp.slave { // 只适用于slave角色
		atomic.StorePointer(&traceableTradeActionResp.lastTradeActionResp, unsafe.Pointer(tradeActionResp)) // 原子地设置指针值
	}

	// 追加回报
	traceableTradeActionResp.tradeActionRespList = append(traceableTradeActionResp.tradeActionRespList, tradeActionResp)
	traceableTradeActionResp.tradeActionRespMap[key] = tradeActionResp
	added = true
	/*
		// 更新latestresp
		traceableTradeActionResp.tradeActionLatestResp.OrderID = tradeActionResp.OrderID
		traceableTradeActionResp.tradeActionLatestResp.OrigClOrdID = tradeActionResp.OrigClOrdID
		traceableTradeActionResp.tradeActionLatestResp.ExecID = tradeActionResp.ExecID
		traceableTradeActionResp.tradeActionLatestResp.ExecType = tradeActionResp.ExecType
		traceableTradeActionResp.tradeActionLatestResp.OrdStatus = tradeActionResp.OrdStatus
		traceableTradeActionResp.tradeActionLatestResp.OrdRejReason = tradeActionResp.OrdRejReason
		traceableTradeActionResp.tradeActionLatestResp.CxlRejResponseTo = tradeActionResp.CxlRejResponseTo
		traceableTradeActionResp.tradeActionLatestResp.ExecRestatementReason = tradeActionResp.ExecRestatementReason

		// 对于master，这部分其实只需要更新一次，或者并不需要更新
		if traceableTradeActionResp.slave && traceableTradeActionResp.tradeActionLatestResp.Symbol == "" {
			traceableTradeActionResp.tradeActionLatestResp.Account = tradeActionResp.Account
			traceableTradeActionResp.tradeActionLatestResp.Symbol = tradeActionResp.Symbol
			traceableTradeActionResp.tradeActionLatestResp.SymbolSfx = tradeActionResp.SymbolSfx
			traceableTradeActionResp.tradeActionLatestResp.SecurityID = tradeActionResp.SecurityID
			traceableTradeActionResp.tradeActionLatestResp.IDSource = tradeActionResp.IDSource
			traceableTradeActionResp.tradeActionLatestResp.SecurityType = tradeActionResp.SecurityType
			traceableTradeActionResp.tradeActionLatestResp.Side = tradeActionResp.Side
			traceableTradeActionResp.tradeActionLatestResp.OpenClose = tradeActionResp.OpenClose
			traceableTradeActionResp.tradeActionLatestResp.OrderQty = tradeActionResp.OrderQty
			traceableTradeActionResp.tradeActionLatestResp.CashOrderQty = tradeActionResp.CashOrderQty
			traceableTradeActionResp.tradeActionLatestResp.OrdType = tradeActionResp.OrdType
			traceableTradeActionResp.tradeActionLatestResp.Price = tradeActionResp.Price
			traceableTradeActionResp.tradeActionLatestResp.Currency = tradeActionResp.Currency
			traceableTradeActionResp.tradeActionLatestResp.EffectiveTime = tradeActionResp.EffectiveTime
			traceableTradeActionResp.tradeActionLatestResp.ExpireTime = tradeActionResp.ExpireTime
		}

		traceableTradeActionResp.tradeActionLatestResp.LastShares = tradeActionResp.LastShares
		traceableTradeActionResp.tradeActionLatestResp.LastPx = tradeActionResp.LastPx
		traceableTradeActionResp.tradeActionLatestResp.LeavesQty = tradeActionResp.LeavesQty
		traceableTradeActionResp.tradeActionLatestResp.CumQty = tradeActionResp.CumQty
		traceableTradeActionResp.tradeActionLatestResp.AvgPx = tradeActionResp.AvgPx
		traceableTradeActionResp.tradeActionLatestResp.TransactTime = tradeActionResp.TransactTime
		traceableTradeActionResp.tradeActionLatestResp.MsgTime = tradeActionResp.MsgTime
		traceableTradeActionResp.tradeActionLatestResp.MsgSeq = tradeActionResp.MsgSeq
	*/
	if updateTradeActionLatestResp == nil {
		updateTradeActionLatestResp = defaultUpdateTradeActionLatestResp
	}
	updateTradeActionLatestResp(tradeActionResp, traceableTradeActionResp)

	// 更新 tradeOrder 的状态，要求对于适配器，所有的属性设置都是有效的
	// 需要特别注意拒绝状态：
	// 1、如果是拒绝状态，且为Update操作的CancelReject，不应更新订单订单的状态
	// 2、其他情况下应更新订单状态
	//if traceableTradeActionResp.tradeActionLatestResp.ActionType == string(enum.ActionType_Update) && strings.HasPrefix(tradeActionResp.OrdRejReason, OrderCancelRejectPrefix) {
	// 【不对】实际上，当遇到OrderCancelReject的回包时，都不应该更新订单的成交状态了。
	// 实际上，当遇到OrderCancelReject的回包时，需要仅更新订单的状态，否则，订单会一直处于pending cancel的状态，而不能恢复到撤单之前的状态
	if len(tradeActionResp.CxlRejResponseTo) > 0 {

		log.Printf("should recover order status for order cancel reject happen!, ClOrdID=%s, OrdStatus=%s\n", tradeActionResp.ClOrdID, tradeActionResp.OrdStatus)
		/*
			traceableTradeActionResp.parent.OrdStatus = tradeActionResp.OrdStatus
		*/

		if updateOrderStatusForOrderCancelReject == nil {
			updateOrderStatusForOrderCancelReject = defaultUpdateOrderStatusForOrderCancelReject
		}
		updateOrderStatusForOrderCancelReject(tradeActionResp, traceableTradeActionResp.tradeActionLatestResp, traceableTradeActionResp.parent)

		// 回滚状态更新时间
		traceableTradeActionResp.parent.OrdStatusUpdateTime = traceableTradeActionResp.parent.OrdStatusUpdateTime2
	} else if tradeActionResp.CumQty >= traceableTradeActionResp.parent.CumQty { // 增加一个约束，只有resp的累积量大于等于当前的订单累积量时，才去更新订单的状态（拒单时也是等于）
		//traceableTradeActionResp.parent.OrdStatus = tradeActionResp.OrdStatus
		traceableTradeActionResp.parent.OrdStatusUpdateTime = tradeActionResp.TransactTime

		//if traceableTradeActionResp.slave { // 因为要每天都归档，所以master节点也是需要同步order的交易信息的
		// traceableTradeActionResp.parent.Account = tradeActionResp.Account
		// traceableTradeActionResp.parent.Symbol = tradeActionResp.Symbol
		// traceableTradeActionResp.parent.SymbolSfx = tradeActionResp.SymbolSfx
		// traceableTradeActionResp.parent.SecurityID = tradeActionResp.SecurityID
		// traceableTradeActionResp.parent.IDSource = tradeActionResp.IDSource
		// traceableTradeActionResp.parent.SecurityType = tradeActionResp.SecurityType
		// traceableTradeActionResp.parent.Side = tradeActionResp.Side
		// traceableTradeActionResp.parent.OpenClose = tradeActionResp.OpenClose
		// traceableTradeActionResp.parent.OrderQty = tradeActionResp.OrderQty
		// traceableTradeActionResp.parent.CashOrderQty = tradeActionResp.CashOrderQty
		// traceableTradeActionResp.parent.OrdType = tradeActionResp.OrdType
		// traceableTradeActionResp.parent.Price = tradeActionResp.Price
		// traceableTradeActionResp.parent.Currency = tradeActionResp.Currency

		/*
			traceableTradeActionResp.parent.LastShares = tradeActionResp.LastShares
			traceableTradeActionResp.parent.LastPx = tradeActionResp.LastPx
			traceableTradeActionResp.parent.LeavesQty = tradeActionResp.LeavesQty
			traceableTradeActionResp.parent.CumQty = tradeActionResp.CumQty
			traceableTradeActionResp.parent.AvgPx = tradeActionResp.AvgPx
			traceableTradeActionResp.parent.OrdRejReason = tradeActionResp.OrdRejReason
		*/

		// 20251217改，移到最后
		//updateOrderStatusForExecutionReport(tradeActionResp, traceableTradeActionResp.tradeActionLatestResp, traceableTradeActionResp.parent)

		// OrdStatusUpdateTime2 保存非pending cancel状态的更新时间
		if tradeActionResp.OrdStatus != "" && tradeActionResp.OrdStatus != string(enum.OrdStatus_PendingCancel) {
			traceableTradeActionResp.parent.OrdStatusUpdateTime2 = tradeActionResp.TransactTime
			traceableTradeActionResp.parent.OrdStatus2 = tradeActionResp.OrdStatus
		}

		if updateOrderStatusForExecutionReport == nil {
			updateOrderStatusForExecutionReport = defaultUpdateOrderStatusForExecutionReport
		}
		// 20251217，updateOrderStatusForExecutionReport，移到这里
		orderUpdateAttributes = updateOrderStatusForExecutionReport(tradeActionResp, traceableTradeActionResp.tradeActionLatestResp, traceableTradeActionResp.tradeActionRespList, traceableTradeActionResp.parent, traceableTradeActionResp.traceableTradeOrder)
		//}
	}

	traceableTradeActionResp.lock.Unlock()

	return
}

func (r *TraceableTradeActionResp) GetLock() *sync.RWMutex {
	return r.lock
}

func (r *TraceableTradeActionResp) GetTradeActionLatestResp() *schema.TradeActionLatestResp {
	return r.tradeActionLatestResp
}

type TraceableTradeOrder struct {
	lock         *sync.RWMutex
	tradeOrder   *schema.TradeOrder
	tradeActions []*TraceableTradeActionResp
	subOrders    []*TraceableTradeOrder // 子单
}

func NewTraceableTradeOrder(order *schema.TradeOrder) *TraceableTradeOrder {
	return &TraceableTradeOrder{lock: &sync.RWMutex{}, tradeOrder: order}
}

func (o *TraceableTradeOrder) Dispose() {
	o.lock = nil
	o.tradeOrder = nil
	for _, tradeAction := range o.tradeActions {
		if tradeAction != nil {
			tradeAction.Dispose()
		}
	}
	for _, subOrder := range o.subOrders {
		if subOrder != nil {
			subOrder.Dispose()
		}
	}
	for i := range o.tradeActions {
		o.tradeActions[i] = nil
	}
	for i := range o.subOrders {
		o.subOrders[i] = nil
	}
	o.tradeActions = nil
	o.subOrders = nil
}

func (o *TraceableTradeOrder) GetLatestTraceableTradeActionRespByAppOrdID(appOrdID string) *TraceableTradeActionResp {
	o.lock.RLock()
	defer o.lock.RUnlock()
	n := len(o.tradeActions)
	log.Printf("GetLatestTraceableTradeActionRespByAppOrdID, n=%d\n", n)
	for i := n - 1; i >= 0; i-- {
		tradeAction := o.tradeActions[i]
		log.Printf("===>tradeAction.parent.AppOrdID:%s\n", tradeAction.parent.AppOrdID)
		if tradeAction.parent.AppOrdID == appOrdID {
			return tradeAction
		}
	}
	return nil
}

func (o *TraceableTradeOrder) GetTraceableTradeActionRespByRootClOrdID(id string) *TraceableTradeActionResp {
	n := len(o.tradeActions)
	log.Printf("GetTraceableTradeActionRespByRootClOrdID, n=%d\n", n)
	for i := 0; i < n; i++ {
		tradeAction := o.tradeActions[i]
		rootClOrdID := tradeAction.tradeActionLatestResp.RootClOrdID
		log.Printf("===>tradeAction.tradeActionLatestResp.RootClOrdID:%s\n", rootClOrdID)
		if rootClOrdID == id {
			return tradeAction
		}
	}
	return nil
}

// 初始化不用加锁
// 需要分发状态：TradeOrder，TradeActionLatestResp
func (o *TraceableTradeOrder) InitTradeAction(firstTradeAction *TraceableTradeActionResp) {
	o.tradeActions = append(o.tradeActions, firstTradeAction)
}

// 需要分发状态：TradeActionLatestResp
func (o *TraceableTradeOrder) AddTradeAction(newTradeAction *TraceableTradeActionResp) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.tradeActions = append(o.tradeActions, newTradeAction)
}

// 获取基本信息
func (o *TraceableTradeOrder) GetBasicInfo() *schema.TradeOrder {
	return o.tradeOrder
}

func (o *TraceableTradeOrder) GetLock() *sync.RWMutex {
	return o.lock
}

// 请参考wiki：http://wiki.gf.com.cn/pages/viewpage.action?pageId=267709155
// 【算法2】
// （1）自右向左遍历 TraceableTradeOrder 的成员数组tradeActions。
// （2）对于每个tradeAction，检查成员对象 tradeActionLatestResp，只考虑 ActionType 等于 New 或 Update，并且 OrdRejReason不能以 OrderCancelReject 开头（以OrderCancelReject开头表明是更新订单被拒绝, 或者之前取消操作被直接拒绝的订单）。
// （3）自右向左找到第一个tradeAction，取ClOrdID返回。
func (o *TraceableTradeOrder) GetActivedClOrdID() (activeID string) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	n := len(o.tradeActions)
	for i := n - 1; i >= 0; i-- {
		tradeAction := o.tradeActions[i]
		actionType := tradeAction.tradeActionLatestResp.ActionType
		if actionType != string(enum.ActionType_New) && actionType != string(enum.ActionType_Update) {
			continue
		}
		if strings.HasPrefix(tradeAction.tradeActionLatestResp.OrdRejReason, OrderCancelRejectPrefix) {
			continue
		}
		return tradeAction.tradeActionLatestResp.ClOrdID
	}
	return
}

func (o *TraceableTradeOrder) GetTradeActions() []*TraceableTradeActionResp {
	o.lock.RLock()
	defer o.lock.RUnlock()

	return o.tradeActions
}

func (o *TraceableTradeOrder) GetTradeActionsWithoutLock() []*TraceableTradeActionResp {
	return o.tradeActions
}

func (o *TraceableTradeOrder) GetSubOrders() []*TraceableTradeOrder {
	o.lock.RLock()
	defer o.lock.RUnlock()

	return o.subOrders
}

func (o *TraceableTradeOrder) AddSubOrder(subOrder *TraceableTradeOrder) []*TraceableTradeOrder {
	o.lock.Lock()
	defer o.lock.Unlock()

	o.subOrders = append(o.subOrders, subOrder)

	return o.subOrders
}

// func (o *TraceableTradeOrder) GetCumCancelCount() (cumCancelCount int) {
// 	o.lock.RLock()
// 	defer o.lock.RUnlock()
// 	n := len(o.tradeActions)
// 	for i := n - 1; i >= 0; i-- {
// 		tradeAction := o.tradeActions[i]
// 		actionType := tradeAction.tradeActionLatestResp.ActionType
// 		if actionType == string(enum.ActionType_Withdraw) {
// 			cumCancelCount++
// 		}
// 	}
// 	return
// }

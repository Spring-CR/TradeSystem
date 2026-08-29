package order_status

import (
	"database/sql"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-core/domain_cfg"
	"rhino-core/order_domain/order_cache"
	"rhino-core/order_domain/order_position_manager"
	"rhino-core/schema"
	"rhino-core/types"
	"rhino-order-report/position"
)

type OrderStatusReplica struct {
	systemCode              string
	businessCode            string
	orderTableName          string
	orderRespTableName      string
	hisOrderTableName       string
	hisOrderRespTableName   string
	extendAttrItems         []*schema.ExtendAttrItem
	extendAttrMap           map[string]*schema.ExtendAttrItem
	tradeRespAttrItems      []*schema.TradeActionRespAttrItem
	dbConfig                string
	createOrderTableInitSql string
	insertOrderSql          string
	updateOrderSql          string
	fullUpdateOrderSql      string
	insertOrderRespSql      string
	applicationCfg          *domain_cfg.ApplicationCfg
	orderCache              *order_cache.OrderCache
	dbWrite                 *sql.DB
	dbRead                  *sql.DB
	memPosition             *position.MemPosition
	positionManager         *order_position_manager.PositionManager
}

func NewOrderStatusReplica(applicationCfg *domain_cfg.ApplicationCfg) *OrderStatusReplica {
	systemCode, businessCode := applicationCfg.GetSystemAndBusinessCodes()
	inst := &OrderStatusReplica{
		applicationCfg: applicationCfg,
		systemCode:     systemCode,
		businessCode:   businessCode,
		// hisOrderTableName     : fmt.Sprintf("trade_orders_extend_%s_%s", systemCode, businessCode),
		// hisOrderRespTableName : fmt.Sprintf("trade_action_resps_extend_%s_%s", systemCode, businessCode),
	}
	inst.initMemDb()
	log.Printf("======>GetApiKafkaBrokers:%v\n", applicationCfg.GetApiKafkaBrokers())

	var tradeRespChs []chan *types.TradeActionRespReturn

	var positionManager *order_position_manager.PositionManager
	// 初始化持仓计算器
	if applicationCfg.GetOrdPositionAdapterPath() != "" {
		inst.memPosition = position.NewMemPosition(applicationCfg)
		var de *domain_error.Error
		positionManager, de = order_position_manager.NewPositionManager(applicationCfg, false, inst.memPosition.ProcessPositionChangeEvent)
		if de != nil {
			domain_error.ProcessSevereError(true, 5, de, nil, "fail to create positionCalculator")
		}
		// 增加持仓计算器的订单执行回报channel
		tradeRespChs = append(tradeRespChs, positionManager.GetTradeRespCh())
		inst.positionManager = positionManager
	}

	inst.orderCache = order_cache.NewOrderCache(false, applicationCfg, positionManager, tradeRespChs, inst.afterAddRootTradeOrder, inst.afterUpdateByTradeActionResp, inst.afterAddTradeActionForDirectOrder, inst.afterUpdateTradeOrderDraft, inst.afterDeleteTradeOrderDraft, inst.afterUpdateTradeOrderAttributes, inst.afterReset, inst.afterSyncOrder, inst.afterSyncTradeActionLatestResp, inst.adjustPosition)

	// 原来注释的理由：只有reset的时候才需要调用。orderCache创建时，已经自动恢复数据
	// 解除注释的理由：orderCache创建时，自动恢复数据模型数据，要更新至sqllite，需要syncMsg驱动。但在stricRecover过程中，如果db的数据本身是完整的，则处理syncMsg消息时，将过滤掉重复的数据，导致sqlite数据无法恢复。
	//              所以，最佳方式是强制再根据模型拓扑，插入一次数据。
	inst.memInitByOrderTopology()

	return inst
}

var (
	ordStatusDraft          = string(enum.OrdStatus_Draft)
	ordStatusPendingReview  = string(enum.OrdStatus_PendingReview)
	ordStatusReviewApproved = string(enum.OrdStatus_ReviewApproved)
	ordStatusReviewRejected = string(enum.OrdStatus_ReviewRejected)
	ordStatusReviewCanceled = string(enum.OrdStatus_ReviewCanceled)
)

// var ordMap = map[string]bool{}
func (s *OrderStatusReplica) afterAddRootTradeOrder(order *types.TraceableTradeOrder) {
	log.Printf("afterAddRootTradeOrder, clOrdID=%s, DBInsertOnOrdExec=%v, appOrdID=%v\n", order.GetBasicInfo().ClOrdID, order.GetBasicInfo().DBInsertOnOrdExec, order.GetBasicInfo().AppOrdID)
	extendAttrMap := order.GetBasicInfo().ExtendAttrMap
	if extendAttrMap == nil {
		return
	}
	orderBasic := order.GetBasicInfo()
	draft := false
	switch orderBasic.OrdStatus {
	case ordStatusDraft, ordStatusPendingReview, ordStatusReviewApproved, ordStatusReviewRejected, ordStatusReviewCanceled:
		draft = true
	}
	if orderBasic.DBInsertOnOrdExec || draft { // insert

		err := s.insertOrder(s.dbWrite, orderBasic)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("error occurs while insert order, appOrdIDp=%s, error=%v\n", orderBasic.AppOrdID, err))
			return
		}

	} else { // full update
		err := s.fullUpdateOrder(s.dbWrite, orderBasic)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("error occurs while update order, appOrdIDp=%s, error=%v\n", orderBasic.AppOrdID, err))
		}
	}
}

func (s *OrderStatusReplica) afterUpdateByTradeActionResp(tradeActionResp *schema.TradeActionResp, traceableTradeActionResp *types.TraceableTradeActionResp) {
	log.Printf("afterUpdateByTradeActionResp, clOrdID=%s\n", tradeActionResp.ClOrdID)
	tradeOrder := traceableTradeActionResp.GetTradeOrder()
	if tradeOrder == nil {
		return
	}
	extendAttrMap := tradeOrder.ExtendAttrMap
	if extendAttrMap == nil {
		return
	}

	// lastTradeActionResp := traceableTradeActionResp.GetLastTradeActionRespByAtomicForSlaveOnly()
	// if lastTradeActionResp == nil {
	// 	lastTradeActionResp = traceableTradeActionResp.GetLastTradeActionResp()
	// }

	// tx, de := dbutil.BeginTx(s.dbWrite)
	// if de != nil {
	// 	domain_error.ProcessSevereError(false, 0, de, nil, "fail to begin tx in afterUpdateByTradeActionResp func")
	// }

	// err := s.insertOrderResp(tx, tradeOrder, lastTradeActionResp)
	// if err != nil {
	// 	domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("error occurs while insert order resp, appOrdIDp=%s, error=%v\n", tradeOrder.AppOrdID, err))
	// 	dbutil.RollbackTx(tx)
	// }

	// err = s.updateOrder(tx, tradeOrder)
	// if err != nil {
	// 	domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("error occurs while update order, appOrdIDp=%s, error=%v\n", tradeOrder.AppOrdID, err))
	// 	dbutil.RollbackTx(tx)
	// }

	// de = dbutil.CommitTx(tx)
	// if de != nil {
	// 	domain_error.ProcessSevereError(false, 0, de, nil, "fail to commit tx in afterUpdateByTradeActionResp func")
	// }

	// 不用tx方案，速度更快
	//err := s.insertOrderResp(s.dbWrite, tradeOrder, lastTradeActionResp) // tradeActionResp 就是当前要存储的tradeActionResp，不需要用lastTradeActionResp
	err := s.insertOrderResp(s.dbWrite, tradeOrder, tradeActionResp)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("error occurs while insert order resp, appOrdIDp=%s, error=%v\n", tradeOrder.AppOrdID, err))
	}

	err = s.updateOrder(s.dbWrite, tradeOrder)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("error occurs while update order, appOrdIDp=%s, error=%v\n", tradeOrder.AppOrdID, err))
	}
}

func (s *OrderStatusReplica) afterAddTradeActionForDirectOrder(tradeAction *types.TraceableTradeActionResp) {

}

func (s *OrderStatusReplica) afterUpdateTradeOrderDraft(tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder) {
	err := s.fullUpdateOrder(s.dbWrite, order)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("error occurs while update order, appOrdIDp=%s, error=%v\n", order.AppOrdID, err))
	}
}

func (s *OrderStatusReplica) afterDeleteTradeOrderDraft(tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder) {
	err := s.updateOrder(s.dbWrite, order)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("error occurs while update order, appOrdIDp=%s, error=%v\n", order.AppOrdID, err))
	}
}

func (s *OrderStatusReplica) afterUpdateTradeOrderAttributes(appOrdID string, updateAttrs map[string]interface{}) {
	err := s.updateOrderAttributes(s.dbWrite, appOrdID, updateAttrs)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("error occurs while update order, appOrdIDp=%s, error=%v\n", appOrdID, err))
	}
}

func (s *OrderStatusReplica) afterSyncOrder(order *schema.TradeOrder) {
	err := s.updateOrder(s.dbWrite, order)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("error occurs while update order, appOrdIDp=%s, error=%v\n", order.AppOrdID, err))
	}
}

func (s *OrderStatusReplica) afterSyncTradeActionLatestResp(tradeActionLatestResp *schema.TradeActionLatestResp) {

}

func (s *OrderStatusReplica) GetOrderCache() *order_cache.OrderCache {
	return s.orderCache
}

func (s *OrderStatusReplica) GetMemPosition() *position.MemPosition {
	return s.memPosition
}

func (s *OrderStatusReplica) adjustPosition(mockTradeOrder *schema.TradeOrder, mockTradeActionResp *schema.TradeActionResp) {
	s.positionManager.AdjustPosition(mockTradeOrder, mockTradeActionResp)
}

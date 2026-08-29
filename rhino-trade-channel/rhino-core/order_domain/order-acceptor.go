package order_domain

// 作者：林春泉

import (
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-core/schema"
	"rhino-core/types"
	"sync/atomic"
)

type OrderAcceptor struct {
	concurrencyLimit  int64
	concurrency       int64
	orderExecutor     OrderExecutorInterface
	orderOrchestrator *OrderOrchestrator
}

func NewOrderAcceptor(concurrencyLimit int64, orderExecutor OrderExecutorInterface, orderOrchestrator *OrderOrchestrator) *OrderAcceptor {
	inst := &OrderAcceptor{concurrencyLimit: concurrencyLimit, orderExecutor: orderExecutor, orderOrchestrator: orderOrchestrator}
	return inst
}

func (a *OrderAcceptor) AcceptNewOrderSingleRequest(order *schema.TradeOrder) (duplicatedOrder bool, de *domain_error.Error) {

	atomic.AddInt64(&a.concurrency, 1)
	defer atomic.AddInt64(&a.concurrency, -1)

	if order.IsDirectOrd {

		jsData, _ := json.Marshal(order)
		log.Printf("AcceptNewOrderSingleRequest, Order=%s\n", jsData)

		// 重复订单的检查需要提前，不然很多问题
		duplicatedOrder = a.orderOrchestrator.orderCache.IsDuplicateOrder(order)
		if duplicatedOrder {
			log.Printf("detect duplicated order:%v\n", order.AppOrdID)
			domain_error.BuildWithDetails(domain_error.ERROR, order, domain_error.DUPLICATE_ORDER_ERR_CODE, nil, order.AppOrdID)
			return
		}

		pm := a.orderOrchestrator.GetPositionManager()
		if pm != nil {
			// 申请持仓额度
			// 注意：一开始是没有clOrdID的，因此，只能使用AppOrdID了
			//de = pc.AcquireOrderQuota(false, order)
			// 确保幂等，即同一个订单并发多次冻结也仅首次有效

			// _, de = pm.FreezeQuota(false, order)
			// if de != nil {
			// 	de.Refine(domain_error.ERROR, order)
			// 	return
			// }

			// force与当前是否审批通过的状态有关，如果是限额检查未通过的异常，需要保存一下
			_, de = pm.FreezeQuota(order.ApproveStatus == int(enum.ApproveStatus_Approved), order)
			if de != nil {
				de.Refine(domain_error.ERROR, order)
				return
			}
		}

		/*
		var cc *order_capital.CapitalCalculator
		if order.ApproveStatus != int(enum.ApproveStatus_Approved) {
			cc = a.orderOrchestrator.GetCapitalCalculator()
			if cc != nil {
				// 申请资金额度
				_, _, de = cc.AcquireOrderCapital(false, order)
				if de != nil {
					de.Refine(domain_error.ERROR, order)
					return
				}
			}
		} else {
			cc = a.orderOrchestrator.GetCapitalCalculator()
			if cc != nil {
				// 申请资金额度
				_, _, de = cc.AcquireOrderCapital(true, order)
				if de != nil {
					de.Refine(domain_error.ERROR, order)
					return
				}
			}
		}
		*/

		order.ApproveStatus = int(enum.ApproveStatus_Approved)
		duplicatedOrder, de = a.orderExecutor.NewOrderSingle(order)
		if de != nil || duplicatedOrder {
			if pm != nil { 
				// 释放持仓额度
				// 确保幂等，即同一个订单并发多次解冻也仅首次有效
				pm.RollbackFreezeQuota(order)
			}
			/*
			if cc != nil {
				// 回滚冻结资金
				cc.RollbackFreezedOrderCapitalForTradeError(order)
			}
			*/
		}
		if de != nil {
			return
		}
	}

	return
}

func (a *OrderAcceptor) AcceptOrderCancelRequest(orderCxlReq *types.ApplicationOrderCancelRequest) (tradeOrder *schema.TradeOrder, de *domain_error.Error) {
	jsData, _ := json.Marshal(orderCxlReq)
	log.Printf("AcceptOrderCancelRequest, orderCxlReq=%s\n", jsData)
	tradeOrder, de = a.orderExecutor.OrderCancelRequest(orderCxlReq)
	return
}

// AcceptOrderDraft: 保存订单草稿。
func (a *OrderAcceptor) AcceptOrderDraft(order *schema.TradeOrder, actionType enum.ActionType) (de *domain_error.Error) {

	// 需要加锁，从OrderCache获取实时数据，所以不可以在这里做订单状态的检查
	// if order.OrdStatus != "" && order.OrdStatus != string(enum.OrdStatus_Draft) {
	// 	return domain_error.Build(domain_error.CANNOT_UPDATE_ORDER_DRAFT_FOR_STATUS_RESON_ERR_CODE, nil)
	// }

	return a.orderOrchestrator.SaveOrderDraft(order, actionType)
}

// AcceptOrderDraftDeletion: 删除订单草稿（假删除，即仅设置删除标识，而并非从数据库中移除记录）
func (a *OrderAcceptor) AcceptOrderDraftDeletion(orderDraftDeletion *types.ApplicationOrderDraftDeleteRequest) *domain_error.Error {
	return a.orderOrchestrator.DeleteOrderDraft(orderDraftDeletion)
}

func (a *OrderAcceptor) AcceptOrderAttributeUpdateRequest(attrUpdateReq *types.ApplicationOrderAttributeUpdateRequest) (de *domain_error.Error) {
	jsData, _ := json.Marshal(attrUpdateReq)
	log.Printf("AcceptOrderAttributeUpdateRequest, attrUpdateReq=%s\n", jsData)
	de = a.orderOrchestrator.UpdateOrderAttributes(attrUpdateReq)
	return
}

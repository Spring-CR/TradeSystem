package ficc

import (
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/types"
)

type TitansFiccExecutorAdapter struct {
}

func NewTitansFiccExecutorAdapter(c *domain_cfg.ApplicationCfg) (adapter *TitansFiccExecutorAdapter, de *domain_error.Error) {
	log.Printf("construct NewTitansFiccExecutorAdapter...")
	adapter = &TitansFiccExecutorAdapter{}
	return
}

func (a *TitansFiccExecutorAdapter) AfterNewOrderSingle(order *schema.TradeOrder, duplicatedOrder bool, de *domain_error.Error) {
	if !duplicatedOrder && de == nil {
		val, _, _ := attrutil.GetAttrValue(order.ExtendAttrMap, "handlInst", enum.AttrValueType_STRING)
		if val.(string) == "4" {
			domain_error.BuildWithDetails(domain_error.WARNING, order, domain_error.GENERIC_WARNING_CODE, nil, "有一笔人工补录订单的下单请求待处理")
		}
	}
}

func (a *TitansFiccExecutorAdapter) AfterOrderCancelRequest(orderCancelRequest *types.ApplicationOrderCancelRequest, order *schema.TradeOrder, de *domain_error.Error) {
	if de == nil {
		val, _, _ := attrutil.GetAttrValue(order.ExtendAttrMap, "handlInst", enum.AttrValueType_STRING)
		if val.(string) == "4" {
			domain_error.BuildWithDetails(domain_error.WARNING, order, domain_error.GENERIC_WARNING_CODE, nil, "有一笔人工补录订单的撤单请求待处理")
		}
	}
}

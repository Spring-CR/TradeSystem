package handlers

import (
	"errors"
	"log"
	"rhino-common/domain_error"
	"rhino-common/server/middleware"

	"github.com/gin-gonic/gin"
)

func (h *OrderHandler) AdjustPosition(c *gin.Context) {

	pm := h.engine.GetOrderOrchestrator().GetPositionManager()
	if pm == nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, errors.New("未配置持仓计算插件"), nil)
		if middleware.ProcessDomainError(de, c) {
			return
		}
	}

	msgProps := map[string]interface{}{}
	if !middleware.BindInputOption(c, &msgProps) {
		return
	}

	msgPropsJs, err := json.Marshal(msgProps)
	log.Printf("======>msgPropsJs:%s, err:%v\n", msgPropsJs, err)

	order, de := h.apiAdapter.ConvertNewOrderSingleMessage(nil, msgProps)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	order.SystemCode = h.systemCode
	order.BusinessCode = h.businessCode

	// 订单属性精华和校验
	de = h.apiAdapter.RefineAndValidate(order, false)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	mockTradeOrder, mockTradeActionResp, de := pm.PreparePositionAdjustmentParams(order)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	if mockTradeOrder == nil || mockTradeActionResp == nil {
		de := domain_error.Build(domain_error.GENERIC_ERR_CODE, errors.New("未正确实现持仓计算插件"))
		if middleware.ProcessDomainError(de, c) {
			return
		}
	}

	pm.AdjustPosition(mockTradeOrder, mockTradeActionResp)

	h.engine.GetOrderOrchestrator().GetOrderCache().AfterAdjustPosition(mockTradeOrder, mockTradeActionResp)
}

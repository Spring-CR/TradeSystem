package handlers

import (
	"rhino-common/server/middleware"
	"rhino-common/utils/dbutil"
	"rhino-core/domain_cfg"
	"rhino-order-report/api/options"
	"rhino-order-report/order_status"

	"github.com/gin-gonic/gin"
)

// var (
// 	json = jsoniter.ConfigCompatibleWithStandardLibrary
// )

type OrderReportHandler struct {
	orderStatus *order_status.OrderStatusReplica
}

func NewCurrentReportHandler(applicationCfg *domain_cfg.ApplicationCfg) *OrderReportHandler {
	inst := &OrderReportHandler{}
	inst.orderStatus = order_status.NewOrderStatusReplica(applicationCfg)
	return inst
}

func (h *OrderReportHandler) QueryTradeOrder(c *gin.Context) {

	query := &dbutil.StructuralQuery{}
	if !middleware.BindInputOption(c, query) {
		return
	}

	data, total, de := h.orderStatus.QueryTradeOrder(query, false)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	result := options.PagingRecord{Total: total, Data: data}

	middleware.ResponseJson(c, result)
}

func (h *OrderReportHandler) QueryTradeActionResp(c *gin.Context) {

	query := &dbutil.StructuralQuery{}
	if !middleware.BindInputOption(c, query) {
		return
	}

	data, total, de := h.orderStatus.QueryTradeActionResp(query, false)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	result := options.PagingRecord{Total: total, Data: data}

	middleware.ResponseJson(c, result)
}

func (h *OrderReportHandler) Dump(c *gin.Context) {

	h.orderStatus.GetOrderCache().Dump()
}

func (h *OrderReportHandler) QueryHisTradeOrder(c *gin.Context) {

	query := &dbutil.StructuralQuery{}
	if !middleware.BindInputOption(c, query) {
		return
	}

	data, total, de := h.orderStatus.QueryTradeOrder(query, true)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	result := options.PagingRecord{Total: total, Data: data}

	middleware.ResponseJson(c, result)
}

func (h *OrderReportHandler) QueryHisTradeActionResp(c *gin.Context) {

	query := &dbutil.StructuralQuery{}
	if !middleware.BindInputOption(c, query) {
		return
	}

	data, total, de := h.orderStatus.QueryTradeActionResp(query, true)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	result := options.PagingRecord{Total: total, Data: data}

	middleware.ResponseJson(c, result)
}

func (h* OrderReportHandler) QueryPosition(c *gin.Context) {

	query := &dbutil.StructuralQuery{}
	if !middleware.BindInputOption(c, query) {
		return
	}

	data, total, de := h.orderStatus.GetMemPosition().QueryPosition(query)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	result := options.PagingRecord{Total: total, Data: data}

	middleware.ResponseJson(c, result)
}
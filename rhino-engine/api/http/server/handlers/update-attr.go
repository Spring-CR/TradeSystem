package handlers

import (
	"log"
	"rhino-common/domain_error"
	"rhino-common/server/middleware"
	"rhino-core/types"

	"github.com/gin-gonic/gin"
)

func (h *OrderHandler) UpdateOrderAttribute(c *gin.Context) {

	attrUpdateReq := &types.ApplicationOrderAttributeUpdateRequest{}
	if !middleware.BindInputOption(c, &attrUpdateReq) {
		return
	}

	if len(attrUpdateReq.UpdateAttributes) == 0 {
		js, _ := json.Marshal(attrUpdateReq)
		log.Printf("UpdateAttributes is empty, UpdateAttributes=%s\n", js)
		return
	}

	// appOrdID, de := h.apiAdapter.ExtractApplicationOrderID(nil, attrUpdateReq.UpdateAttributes)
	// if middleware.ProcessDomainError(de, c) {
	// 	return
	// }

	// attrUpdateReq.AppOrdID = appOrdID
	var de *domain_error.Error
	_, ok := h.engine.GetOrderByAppOrdID(attrUpdateReq.AppOrdID)
	if !ok {
		de = domain_error.Build(domain_error.CANNOT_FIND_ORDER_BY_APP_ORD_ID_ERR_CODE, nil, attrUpdateReq.AppOrdID)
		middleware.ProcessDomainError(de, c)
		return
	}

	de = h.engine.GetOrderAcceptor().AcceptOrderAttributeUpdateRequest(attrUpdateReq)
	if middleware.ProcessDomainError(de, c) {
		return
	}
}
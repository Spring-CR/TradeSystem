package router

import (
	"rhino-order-report/api/api_const"
	"rhino-order-report/api/server/handlers"

	"github.com/gin-gonic/gin"
)

func setCurrentReportHandler(e *gin.Engine, handler *handlers.OrderReportHandler) {

	r := e.Group(api_const.RootRouteReport)
	{
		r.POST(api_const.SubCurrentTradeOrders, handler.QueryTradeOrder)
		r.POST(api_const.SubCurrentTradeActionResps, handler.QueryTradeActionResp)
		r.GET(api_const.SubCurrentDataDump, handler.Dump)
		r.POST(api_const.SubHistoryTradeOrders, handler.QueryHisTradeOrder)
		r.POST(api_const.SubHistoryTradeActionResps, handler.QueryHisTradeActionResp)
		r.POST(api_const.SubPosition, handler.QueryPosition)
	}
}

package router

import (
	"rhino-api/api/api_const"
	"rhino-api/api/server/handlers"

	"github.com/gin-gonic/gin"
)

func setTradeInstrRouter(e*gin.Engine) {

	r := e.Group(api_const.RouteTradeInstr)
	{
		r.POST(api_const.SubRouteSingleTrade, handlers.ExecuteSingleTradeInstr)
		r.GET(api_const.SubRouteTradeDeskOrderIds, handlers.FindTradeDeskOrderIdsByParentKey)
	}
}
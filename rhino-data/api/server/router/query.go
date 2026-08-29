package router

import (
	"rhino-data/api/api_const"
	"rhino-data/api/server/handlers"

	"github.com/gin-gonic/gin"
)

func setDataQueryHandler(e *gin.Engine, handler *handlers.DataQueryHandler) {

	r := e.Group(api_const.RootRouteDataQry)
	{
		r.GET("", handler.Query)
	}
}

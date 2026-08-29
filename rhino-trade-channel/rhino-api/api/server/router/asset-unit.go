package router

import (
	"rhino-api/api/api_const"
	"rhino-api/api/server/handlers"

	"github.com/gin-gonic/gin"
)

func setAssetUnitRouter(e*gin.Engine) {

	r := e.Group(api_const.RouteAssetUnit)
	{
		r.POST("", handlers.CreateAssetUnit)
		r.GET(api_const.SubRouteFindAll, handlers.FindAllAssetUnits)
		r.DELETE(api_const.SubRouteDelete, handlers.DeleteAssetUnit)
	}
}
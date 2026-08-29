package router

import (
	"rhino-data-sync/api/api_const"
	"rhino-data-sync/api/server/handlers"

	"github.com/gin-gonic/gin"
)

func setDataSyncNotifyHandler(e *gin.Engine, handler *handlers.DataSyncNotifyHandler) {

	r := e.Group(api_const.RootRouteDataSync)
	{
		r.POST(api_const.SubRouteNotify, handler.Notify)
		r.GET(api_const.SubRouteSyncLog, handler.GetSyncLog)
	}
}

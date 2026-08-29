package router

import (
	"rhino-data/api/api_const"
	"rhino-data/api/server/handlers"

	"github.com/gin-gonic/gin"
)

func setSyncLogHandler(e *gin.Engine, handler *handlers.SyncLogHandler) {

	r := e.Group(api_const.RootRouteSyncLogs)
	{
		r.GET("", handler.GetSyncLogs)
	}
}

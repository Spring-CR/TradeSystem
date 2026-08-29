package router

import (
	"net/http"
	"rhino-common/server/middleware"
	"rhino-data-sync/api/api_const"
	"rhino-data-sync/api/server/handlers"

	"github.com/gin-gonic/gin"
)

func GetApiRouters(dataSyncNotifyHandler *handlers.DataSyncNotifyHandler, apiToken string) http.Handler {
	
	e := gin.New()
	e.Use(gin.Recovery())

	if apiToken != "" {
		e.Use(middleware.SetCredentials(apiToken))
	}

	setPingRouter(e)
	setDataSyncNotifyHandler(e, dataSyncNotifyHandler)

	return e
}

func setPingRouter(e *gin.Engine) {

	r := e.Group(api_const.RoutePing)
	{
		r.GET("", func (c *gin.Context)  {
			c.String(200, "pong")
		})
	}
}
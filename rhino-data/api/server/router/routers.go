package router

import (
	"net/http"
	"rhino-common/server/middleware"
	"rhino-data/api/api_const"
	"rhino-data/api/server/handlers"

	"github.com/gin-gonic/gin"
)

func GetApiRouters(dataQueryHandler *handlers.DataQueryHandler, syncLogsHandler *handlers.SyncLogHandler, apiToken string) http.Handler  {

	e := gin.New()
	e.Use(gin.Recovery())

	if apiToken != "" {
		e.Use(middleware.SetCredentials(apiToken))
	}

	setPingRouter(e)
	setDataQueryHandler(e, dataQueryHandler)
	setSyncLogHandler(e, syncLogsHandler)

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
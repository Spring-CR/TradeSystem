package router

import (
	"net/http"
	"rhino-common/server/middleware"
	"rhino-order-report/api/api_const"
	"rhino-order-report/api/server/handlers"

	"github.com/gin-gonic/gin"
)

func GetApiRouters(currentReportHandler *handlers.OrderReportHandler, apiToken string) http.Handler  {

	e := gin.New()
	e.Use(gin.Recovery())

	if apiToken != "" {
		e.Use(middleware.SetCredentials(apiToken))
	}

	setPingRouter(e)
	setCurrentReportHandler(e, currentReportHandler)

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
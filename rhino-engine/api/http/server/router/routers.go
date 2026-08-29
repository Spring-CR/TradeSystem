package router

import (
	"net/http"
	"rhino-common/server/middleware"
	"rhino-engine/api/http/api_const"
	"rhino-engine/api/http/server/handlers"

	"github.com/gin-gonic/gin"
)

func GetApiRouters(orderHandler *handlers.OrderHandler, apiToken string) http.Handler  {

	e := gin.New()
	e.Use(gin.Recovery())

	if apiToken != "" {
		e.Use(middleware.SetCredentials(apiToken))
	}

	setPingRouter(e)
	setOrderRouter(e, orderHandler)

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
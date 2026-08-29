package router

import (
	"rhino-api/api/api_const"
	"rhino-api/api/server/handlers"

	"github.com/gin-gonic/gin"
)

func setTaskInstrRouter(e*gin.Engine) {

	r := e.Group(api_const.RouteTaskInstr)
	{
		r.POST("", handlers.IssueInstruction)
		r.POST(api_const.SubRouteQuery, handlers.FindInstruction)
		r.GET(api_const.SubRouteStatisUpdate, handlers.StatisTaskInstrStock)
	}
}
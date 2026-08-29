package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func ApiRouter(positionServiceUrl string) http.Handler  {

	e := gin.New()
	e.Use(gin.Recovery())

	setTaskInstrRouter(e)
	setTradeInstrRouter(e)
	setAssetUnitRouter(e)
	setPositionRouter(e, positionServiceUrl)

	return e
}
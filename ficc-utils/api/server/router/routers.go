package router

import (
	"ficc-utils/api/api_const"
	"ficc-utils/common/server/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetApiRouters(dataQryServiceUrl, loginServiceUrl, loginClientID, loginClientKey, capitalServiceUrl, positionServiceUrl, tradeOrdersServiceUrl, hisTradeOrdersServiceUrl, execReportServiceUrl, pricingHubYtmConvUrl, pricingHubAuth, dataQryUrl, appId, appSecret string, queryCycle time.Duration, local *time.Location) http.Handler {

	e := gin.New()
	e.Use(gin.Recovery())

	credProvider := NewCredProvider(dataQryServiceUrl, queryCycle, local)
	setLoginRouter(e, loginServiceUrl, loginClientID, loginClientKey, credProvider)
	setPingRouter(e)
	setPricingRouter(e, pricingHubYtmConvUrl, pricingHubAuth, dataQryUrl)
	setCurrentOrdersExportRouter(e, tradeOrdersServiceUrl)
	setHistoryOrdersExportRouter(e, hisTradeOrdersServiceUrl)
	setInternalCapitalRouter(e, capitalServiceUrl, positionServiceUrl, appId, appSecret)
	setInternalSecuritiesRouter(e, dataQryUrl)

	// 这里加调用定价服务的handler在setCredential之前，即需要鉴权
	// 设置权限控制
	e.Use(middleware.SetCredentials(credProvider.CredFunc))

	// 资金查询、持仓查询在这之后增加，需要鉴权
	setCapitalRouter(e, capitalServiceUrl, positionServiceUrl, appId, appSecret)
	setPositionRouter(e, positionServiceUrl)
	setContractPositionRouter(e, dataQryUrl)
	setCurrentOrdersQueryRouter(e, tradeOrdersServiceUrl)
	setHistoryOrdersQueryRouter(e, hisTradeOrdersServiceUrl)
	//setCurrentExecReportQueryRouter(e, execReportServiceUrl)
	setSecuritiesRouter(e, dataQryUrl)

	return e
}

func setPingRouter(e *gin.Engine) {
	r := e.Group(api_const.RoutePing)
	{
		r.GET("", func(c *gin.Context) {
			c.String(200, "pong")
		})
	}
}


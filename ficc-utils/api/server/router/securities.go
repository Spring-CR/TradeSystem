package router

import (
	"ficc-utils/api/api_const"
	"ficc-utils/common/domain_error"
	"ficc-utils/common/server/middleware"
	"ficc-utils/common/utils/data_qry"
	"ficc-utils/common/utils/request"
	"log"

	"github.com/gin-gonic/gin"
)

func setSecuritiesRouter(e *gin.Engine, dataQryUrl string) {
	setSecuritiesRouterWithPath(e, dataQryUrl, api_const.RouteInterestRateBondSecurities)
}

func setInternalSecuritiesRouter(e *gin.Engine, dataQryUrl string) {
	setSecuritiesRouterWithPath(e, dataQryUrl, api_const.RouteInterestRateBondSecuritiesInternal)
}


func setSecuritiesRouterWithPath(e *gin.Engine, dataQryUrl string, routePath string) {

	r := e.Group(routePath)
	{
		r.GET("", func(c *gin.Context) {

			if routePath == api_const.RouteInterestRateBondSecurities {
				account, de := request.GetQueryAsInt(c, api_const.ParamAccount, false)
				if middleware.ProcessDomainError(de, c) {
					return
				}
				log.Printf("qrySymbolsT0, account: %d", account)
			}

			symbolsT0, err := data_qry.GetSymbolsT0(dataQryUrl)
			if err != nil {
				de := domain_error.Build(domain_error.CANNOT_QUERY_SECURITIES_INFO_ERR_CODE, err)
				middleware.ProcessDomainError(de, c)
				return
			}

			middleware.ResponseJson(c, symbolsT0)
		})
	}
}
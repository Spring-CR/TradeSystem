package router

import (
	"ficc-utils/api/api_const"
	"ficc-utils/api/options"
	"ficc-utils/common/domain_error"
	"ficc-utils/common/server/middleware"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

func setCurrentExecReportQueryRouter(e *gin.Engine, execReportServiceUrl string) {

	r := e.Group(api_const.RouteCurrentExecReport)
	{
		r.POST("", func(c *gin.Context) {
			handleOrdersQuery(c, execReportServiceUrl)
		})
	}
}

func handleExecReportQuery(c *gin.Context, execReportServiceUrl string) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		de := domain_error.Build(domain_error.CANNOT_EXPORT_ORDERS_INFO_ERR_CODE, err)
		middleware.ProcessDomainError(de, c)
		return
	}

	result, err := forwardQueryExecReportQueryService(execReportServiceUrl, body)
	if err != nil {
		de := domain_error.Build(domain_error.CANNOT_QUERY_ORDERS_INFO_ERR_CODE, err, "")
		middleware.ProcessDomainError(de, c)
		return
	}

	if result.Code != 0 {
		de := domain_error.Build(domain_error.CANNOT_QUERY_ORDERS_INFO_ERR_CODE, nil, result.Message)
		middleware.ProcessDomainError(de, c)
		return
	}

	middleware.ResponseJson(c, result.Data)
}

// 调用成交信息查询服务（execReportServiceUrl），并将结果返回
func forwardQueryExecReportQueryService(execReportServiceUrl string, payload []byte) (result *options.GenericQueryResult[*options.ExecReport], err error) {
	log.Printf("forwardQueryExecReportQueryService execReportServiceUrl:%s", execReportServiceUrl)
	return ForwardQueryTradeOrdersService[*options.ExecReport](execReportServiceUrl, payload)
}
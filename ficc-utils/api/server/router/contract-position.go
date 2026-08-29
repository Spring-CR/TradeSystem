package router

import (
	"ficc-utils/api/api_const"
	"ficc-utils/api/options"
	"ficc-utils/common/domain_error"
	"ficc-utils/common/server/middleware"
	"ficc-utils/common/utils/data_qry"
	"ficc-utils/common/utils/request"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func setContractPositionRouter(e *gin.Engine, dataQryUrl string) {

	r := e.Group(api_const.RouteContractPosition)
	{
		r.GET("", func(c *gin.Context) {
			account, de := request.GetQueryAsInt(c, api_const.ParamAccount, false)
			if middleware.ProcessDomainError(de, c) {
				return
			}
			symbol, de := request.GetQueryAsString(c, api_const.ParamSymbol, true)
			if middleware.ProcessDomainError(de, c) {
				return
			}

			log.Printf("QueryContractPosition account:%d, symbol:%s\n", account, symbol)

			result, err := invokeContractPositionService(dataQryUrl, account, symbol)
			if err != nil {
				de := domain_error.Build(domain_error.CANNOT_QUERY_CAPITAL_INFO_ERR_CODE, err, []int{account})
				middleware.ProcessDomainError(de, c)
				return
			}

			middleware.ResponseJson(c, result)
		})
	}
}

// 调用data_qry服务（dataQryUrl），并将结果返回, "tableName": "ContractPosition"
func invokeContractPositionService(dataQryUrl string, account int, symbol string) (result []*options.ContractPositionOut, err error){
	key := fmt.Sprintf("%d", account)
	if symbol != "" {
		key = fmt.Sprintf("%d-%s", account, symbol)
	}
	//data := []*options.ContractPosition{}
	data := options.GenericQueryResult[*options.ContractPosition]{}
	err = data_qry.QryTableData(dataQryUrl, "ContractPosition", key, &data)
	for _, d := range data.Data {
		result = append(result, d.ToContractPositionOut())
	}
	return
}
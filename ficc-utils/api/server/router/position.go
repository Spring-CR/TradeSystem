package router

import (
	"bytes"
	"encoding/json"
	"ficc-utils/api/api_const"
	"ficc-utils/api/options"
	"ficc-utils/common/domain_error"
	"ficc-utils/common/server/middleware"
	"ficc-utils/common/utils/request"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 测试环境：http://olts-dev.gf.com.cn/unitrade/test/titans/ficc/order_report/api/v1/report/position
func setPositionRouter(e *gin.Engine, positionServiceUrl string) {

	r := e.Group(api_const.RoutePosition)
	{
		r.GET("", func(c *gin.Context) {
			account, de := request.GetQueryAsInt(c, api_const.ParamAccount, false)
			if middleware.ProcessDomainError(de, c) {
				return
			}
			planCode, de := request.GetQueryAsString(c, api_const.ParamPlanCode, true)
			if middleware.ProcessDomainError(de, c) {
				return
			}
			symbol, de := request.GetQueryAsString(c, api_const.ParamSymbol, true)
			if middleware.ProcessDomainError(de, c) {
				return
			}
			pageNum, de := request.GetQueryAsInt(c, api_const.ParamPageNum, false)
			if middleware.ProcessDomainError(de, c) {
				return
			}
			pageSize, de := request.GetQueryAsInt(c, api_const.ParamPageSize, false)
			if middleware.ProcessDomainError(de, c) {
				return
			}

			log.Printf("QueryPosition account:%d, planCode:%s, symbol:%s, pageIndex:%d, pageSize:%d\n", account, planCode, symbol, pageNum, pageSize)

			result, err := invokePositionService(positionServiceUrl, []int{account}, planCode, symbol, pageNum, pageSize)
			if err != nil {
				de := domain_error.Build(domain_error.CANNOT_QUERY_CAPITAL_INFO_ERR_CODE, err, []int{account})
				middleware.ProcessDomainError(de, c)
				return
			}

			middleware.ResponseJson(c, result)
		})
	}
}

// 基于交易对手id，调用titans的持仓服务（positionServiceUrl），并将结果返回
func invokePositionService(positionServiceUrl string, accounts []int, planCode, symbol string, pageNum, pageSize int) (result *options.PositionResult, err error){

	input := make(map[string]any)
	//input["select_fields"] = []string{"account","counterpartyID","counterparty","symbol","symbolName","currency","planCode","ultraContractCode","securityExchange","securityType","parValue","netPositionT0","netPositionT1","longAvailablePositionT0","longAvailablePositionT1","shortAvailablePositionT0","shortAvailablePositionT1","longCleanPriceCost","longDirtyPriceCost","longDirtyPriceWithFeeCost","shortCleanPriceCost","shortDirtyPriceCost","shortDirtyPriceWithFeeCost","longAvgCleanPrice","longAvgDirtyPrice","longAvgDirtyPriceWithFee","shortAvgCleanPrice","shortAvgDirtyPrice","ShortAvgDirtyPriceWithFee","maxLongMarginOccupancy","maxShortMarginOccupancy","maxMarginOccupancy"}
	fieldConditions := []map[string]any{}
	accountField := map[string]any{
		"field": "account",
		"field_type": 2, //整形
		"value_type": 2, //集合
		"value": accounts,
	}
	fieldConditions = append(fieldConditions, accountField)

	if planCode != "" {
		planCodeField := map[string]any{
			"field": "planCode",
			"field_type": 3, //字符串
			"value_type": 0, //单值
			"value": planCode,
		}
		fieldConditions = append(fieldConditions, planCodeField)
	}
	if symbol != "" {
		symbolField := map[string]any{
			"field": "symbol",
			"field_type": 3, //字符串
			"value_type": 0, //单值
			"value": symbol,
		}
		fieldConditions	= append(fieldConditions, symbolField)
	}

	// 添加非零持仓条件
	netPositionT0 := map[string]any{
		"field": "netPositionT0",
		"field_type": 2, //整形
		"value_type": 0, //单值
		"value": 0,
		"is_not": true,
	}
	netPositionT1 := map[string]any{
		"field": "netPositionT1",
		"field_type": 2, //整形
		"value_type": 0, //单值
		"value": 0,
		"is_not": true,
	}
	subFieldConditions := map[string]any{
		"sub_field_conditions": []map[string]any{netPositionT0, netPositionT1},
		"operator": "OR",
	}
	fieldConditions	= append(fieldConditions, subFieldConditions)

	input["field_conditions"] = fieldConditions
	input["limit"] = pageSize
	input["offset"] = (pageNum-1) * pageSize

	inputJson, err := json.Marshal(input)
	if err != nil {
		return
	}
	log.Printf("positionServiceUrl input:%s", inputJson)
	resp, err := http.Post(positionServiceUrl, "application/json", bytes.NewReader(inputJson))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	log.Printf("positionServiceUrl:%s", positionServiceUrl)

	body, err := io.ReadAll(resp.Body)
	log.Printf("invokePositionService body:%s", body)

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}
	return
}
package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"ficc-utils/api/api_const"
	"ficc-utils/api/options"
	"ficc-utils/common/domain_error"
	"ficc-utils/common/server/middleware"
	"ficc-utils/common/utils/request"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setCapitalRouter(e *gin.Engine, capitalServiceUrl, positionServiceUrl, appId, appSecret string) {
	setCapitalRouterWithPath(e, api_const.RouteCapital, capitalServiceUrl, positionServiceUrl, appId, appSecret)
}

func setInternalCapitalRouter(e *gin.Engine, capitalServiceUrl, positionServiceUrl, appId, appSecret string) {
	setCapitalRouterWithPath(e, api_const.RouteCapitalInternal, capitalServiceUrl, positionServiceUrl, appId, appSecret)
}
func setCapitalRouterWithPath(e *gin.Engine, routePath, capitalServiceUrl, positionServiceUrl, appId, appSecret string) {

	r := e.Group(routePath)
	{
		r.GET("", func(c *gin.Context) {

			account, de := request.GetQueryAsInt(c, api_const.ParamAccount, false)
			if middleware.ProcessDomainError(de, c) {
				return
			}

			// 基于交易对手id，调用titans的资金服务（capitalServiceUrl），并将结果返回
			result, err := invokeCapitalService(capitalServiceUrl, positionServiceUrl,account, appId, appSecret)
			if err != nil {
				de := domain_error.Build(domain_error.CANNOT_QUERY_CAPITAL_INFO_ERR_CODE, err, account)
				middleware.ProcessDomainError(de, c)
				return
			}

			if routePath == api_const.RouteCapitalInternal {
				middleware.ResponseJson(c, result)
			} else {
				middleware.ResponseJson(c, result.CapitalResult)
			}
		})
	}
}

// 基于交易对手id，调用titans的资金服务（capitalServiceUrl），并将结果返回
func invokeCapitalService(capitalServiceUrl, positionServiceUrl string, account int, appId, appSecret string) (result *options.InternalCapitalResult, err error) {
	// 构造请求URL
	requestUrl := fmt.Sprintf("%s/%d?keyCtptyId=%d", capitalServiceUrl, account, account)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("TITANS资金查询创建HTTP请求失败: %v", err)
	}

	// 添加认证头
	req.Header.Add("AppId", appId)
	req.Header.Add("AppSecret", appSecret)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("TITANS资金查询请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TITANS资金查询返回错误状态码: %d", resp.StatusCode)
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("TITANS资金查询读取响应失败: %v", err)
	}

	log.Printf("invokeCapitalService TITANS资金查询返回结果body: %s", body)
	// 解析响应JSON
	var serviceResp options.CapitalServiceResponse
	if err := json.Unmarshal(body, &serviceResp); err != nil {
		return nil, fmt.Errorf("TITANS资金查询解析响应JSON失败: %v", err)
	}

	// 筛选符合条件的记录：purpose为TRS或MIXTURE，且currency为CNY
	for _, data := range serviceResp.Data {
		if (data.Purpose == "TRS" || data.Purpose == "MIXTURE") && data.Currency == "CNY" {
			// 转换为目标结构体
			result = &options.InternalCapitalResult{
				CapitalResult: options.CapitalResult{
					Account:         data.KeyCtptyId,
					Currency:        data.Currency,
					TotalBalance:    data.TotalBalance,
					AvailableAmount: data.TotalBalance,
				},
				KeyCapAcctId: data.KeyCapAcctId,
				CapAcctCode:  data.CapAcctCode,
			}
			marginOccupancy, err := invokeMarginService(positionServiceUrl, []int{result.Account})
			if err != nil {
				return nil, fmt.Errorf("TITANS持仓保证金占用查询失败 (account: %d), %v", result.Account, err)
			}
			for _, m := range marginOccupancy {
				if m.Account == result.Account && m.MaxMarginOccupancy >1e-9 {
					result.AvailableAmount = math.Round((result.TotalBalance -  m.MaxMarginOccupancy)*100)/100
				}
			}
			return result, nil
		}
	}

	// 如果没有找到符合条件的记录，返回错误
	return nil, fmt.Errorf("TITANS资金查询未找到符合条件的资金账户 (account: %d)", account)
}

func invokeMarginService(positionServiceUrl string, accounts []int) (result []*options.MarginOccupancy, err error){
	input := make(map[string]any)
	input["select_fields"] = []string{"account","sum(maxLongMarginOccupancy) as maxLongMarginOccupancy, sum(maxShortMarginOccupancy) as maxShortMarginOccupancy"}
	fieldConditions := []map[string]any{}
	accountField := map[string]any{
		"field": "account",
		"field_type": 2, //整形
		"value_type": 2, //集合
		"value": accounts,
	}
	fieldConditions = append(fieldConditions, accountField)
	input["field_conditions"] = fieldConditions
	input["group_fields"] = []string{"account"}
	input["limit"] = 100000000
	input["offset"] = 0

	inputJson, err := json.Marshal(input)
	if err != nil {
		return
	}
	log.Printf("invokeMarginService positionServiceUrl(%s) input:%s", positionServiceUrl, inputJson)
	resp, err := http.Post(positionServiceUrl, "application/json", bytes.NewReader(inputJson))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	log.Printf("invokeMarginService body:%s\n", body)
	output := options.MarginResponse{}
	err = json.Unmarshal(body, &output)
	if err != nil {
		return
	}
	if output.Code != 0 {
		err = errors.New(output.Message)
		return
	}
	result = output.Data
	if len(result) > 0 {
		for _, data := range result {
			if data.MaxLongMarginOccupancy > data.MaxShortMarginOccupancy {
				data.MaxMarginOccupancy = data.MaxLongMarginOccupancy
			} else {
				data.MaxMarginOccupancy = data.MaxShortMarginOccupancy
			}
		}
	}
	return result, nil
}
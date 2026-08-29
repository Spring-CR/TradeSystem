package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"ficc-utils/api/api_const"
	"ficc-utils/api/options"
	"ficc-utils/common/domain_error"
	"ficc-utils/common/server/middleware"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

/*{
    "asOf": "2025-07-03",
    "inputType": "YTM",
    "ytm": 0.016389,
    "dirtyPrice": 0,
    "cleanPrice": 0,
    "startDate": "2025-05-25",
    "dayCount": "ACTACT",
    "parValue": 100,
    "couponType": "FIXED",
    "cashTimes": 1,
    "couponRate": 0.0167,
    "issuePrice": 100,
    "maturityDate": "2035-05-25",
    "bondType": "GOVERNMENT_BOND"
}*/

/*{
    "serviceId": "32d4a343-066b-42a7-8f3e-957e88896bb9",
    "errCode": {
        "code": 200,
        "chs": "成功",
        "eng": "success",
        "exceptionClass": null
    },
    "errMsg": null,
    "data": {
        "ytm": 0.016389,
        "dirtyPrice": 100.4591,
        "cleanPrice": 100.2806
    },
    "timestamp": 1760332964570
}*/

/*{
    "serviceId": "c1aab98b-daad-46af-ba56-79fdea261e6d",
    "errCode": {
        "code": 400,
        "chs": "失败",
        "eng": "failure",
        "exceptionClass": null
    },
    "errMsg": "JSON parse error: Cannot deserialize value of type `com.gf.udps.common.enumration.finance.BondCouponTypeEnum` from String \"FIXEDd\": not one of the values accepted for Enum class: [ZERO, FLOAT, FIXED]; nested exception is com.fasterxml.jackson.databind.exc.InvalidFormatException: Cannot deserialize value of type `com.gf.udps.common.enumration.finance.BondCouponTypeEnum` from String \"FIXEDd\": not one of the values accepted for Enum class: [ZERO, FLOAT, FIXED]\n at [Source: (PushbackInputStream); line: 10, column: 19] (through reference chain: com.gf.udps.pricinghub.interfaces.dto.BondYTMConvReq[\"couponType\"])",
    "data": null,
    "timestamp": 1760337148061
}*/

/*{
  "total": 1,
  "data": [
    {
      "AssetType": "SPT_BD",
      "BusinessType": "FICC_BOND",
      "ChnBondPoolIndicator": "Y",
      "CreatedTimestamp": "2025-09-29 00:00:00.123",
      "Currency": "CNY",
      "DelistDate": "2029-10-31",
      "IssueDate": "2009-10-28",
      "ListDate": "2009-11-10",
      "ParValue": 100,
      "SecurityExchange": "NIB",
      "SecurityFullName": "国家开发银行2009年第十九期金融债券",
      "SecurityID": 1548651,
      "SecurityName": "09国开19",
      "SecurityType": "BOND",
      "Symbol": "090219.IB",
      "Symbol2": "DC20091123170004419",
      "TradeDate": "20250929",
      "UpdatedTimestamp": "2025-09-29 00:00:00.223"
    }
  ],
  "displayLen": 1
}*/
/*{
  "total": 1,
  "data": [
    {
      "AssetType": "SPT_BD",
      "BondType": "POLICY_BANK_BOND",
      "BusinessType": "FICC_BOND",
      "CashTimes": 2,
      "ChnBondPoolIndicator": "Y",
      "CouponRate": 0.048,
      "CouponType": "FIXED",
      "CreatedTimestamp": "2025-09-29 00:00:00.123",
      "Currency": "CNY",
      "DayCount": "ACT/ACT(ISMA)",
      "DelistDate": "2029-10-31",
      "IssueDate": "2009-10-28",
      "IssuePrice": 100,
      "ListDate": "2009-11-10",
      "MaturityDate": "2029-11-04 00:00:00",
      "ParValue": 100,
      "SecurityExchange": "NIB",
      "SecurityFullName": "国家开发银行2009年第十九期金融债券",
      "SecurityID": 1548651,
      "SecurityName": "09国开19",
      "SecurityType": "BOND",
      "StartDate": "2009-11-04 00:00:00",
      "Symbol": "090219.IB",
      "Symbol2": "DC20091123170004419",
      "TradeDate": "20250929",
      "UpdatedTimestamp": "2025-09-29 00:00:00.223"
    }
  ],
  "displayLen": 1
}*/


type YtmConvInput struct {
	AsOf string `json:"asOf"`
	InputType string `json:"inputType"`
	Ytm float64 `json:"ytm"`
	DirtyPrice float64 `json:"dirtyPrice"`
	CleanPrice float64 `json:"cleanPrice"`
	StartDate string `json:"startDate"`
	DayCount string `json:"dayCount"`
	ParValue float64 `json:"parValue"`
	CouponType string `json:"couponType"`
	CashTimes int32 `json:"cashTimes"`
	CouponRate float64 `json:"couponRate"`
	IssuePrice float64 `json:"issuePrice"`
	MaturityDate string `json:"maturityDate"`
	BondType string `json:"bondType"`
}

type YtmOutput struct {
	ErrCode struct {
		Code int   `json:"code"`
		Chs string `json:"chs"`
		Eng string `json:"eng"`
		ExceptionClass string `json:"exceptionClass"`
	} `json:"errCode"`
	ErrMsg string `json:"errMsg"`
	Data *options.PricingResult `json:"data"`
}

type BondInfo struct {
	Symbol string `json:"Symbol"`
	BondType string `json:"BondType"`
	CachTimes int32 `json:"CashTimes"`
	CouponType string `json:"CouponType"`
	CouponRate float64 `json:"CouponRate"`
	StartDate string `json:"StartDate"`
	MaturityDate string `json:"MaturityDate"`
	Parvalue float64 `json:"Parvalue"`
	IssuePrice float64 `json:"IssuePrice"`
	DayCount string `json:"DayCount"`
	TradeDate string `json:"TradeDate"`
}

type BondInfoOutput struct {
	Code string `json:"code"`
	Msg string `json:"msg"`
	Stack string `json:"stack"`
	Total int32 `json:"total"`
	Data []BondInfo `json:"data"`
}
func setPricingRouter(e *gin.Engine, pricingHubYtmConvUrl, pricingHubAuth, dataQryUrl string) {

	r := e.Group(api_const.RoutePricing)
	{
		r.POST("", func(c *gin.Context) {
			var form options.PricingForm
			if err := c.ShouldBindJSON(&form); err != nil {
				de := domain_error.Build(domain_error.CANNOT_QUERY_PRICING_INFO_ERR_CODE, err)
				middleware.ProcessDomainError(de, c)
				return
			}
			log.Printf("QueryPricing tradeDate:%s, ytm:%f, symbol:%s\n", form.TradeDate, form.Ytm, form.Symbol)

			result, err := invokePricingService(dataQryUrl, pricingHubYtmConvUrl, pricingHubAuth, form.TradeDate, form.Ytm, form.Symbol)
			if err != nil {
				de := domain_error.Build(domain_error.CANNOT_QUERY_PRICING_INFO_ERR_CODE, err)
				middleware.ProcessDomainError(de, c)
				return
			}

			middleware.ResponseJson(c, result)
		})
	}
}
func invokePricingService(dataQryUrl, pricingHubYtmConvUrl, pricingHubAuth string, tradeDate string, ytm float64, symbol string) (result *options.PricingResult, err error) {
	bondInfo, err := reqBondInfo(dataQryUrl, symbol)
	if err!= nil {
		return
	}
	ytmConvInput := YtmConvInput{
		AsOf: tradeDate,
		InputType: "YTM",
		Ytm: ytm,
		DirtyPrice: 0,
		CleanPrice: 0,
		StartDate: bondInfo.StartDate[0:10],			//起息日
		DayCount: dayCountConvert(bondInfo.DayCount),	//计息基准 /ACT365, BUSINESS244, ACTACT_ISDA, ACT365_NOLEAP, ACT360, BUSINESS252, ACTACT, ACT365F, THIRTY360, THIRTY365, BUSINESS242
		ParValue: bondInfo.Parvalue,					//票面价值
		CouponType: bondInfo.CouponType,				//息票类型
		CashTimes: bondInfo.CachTimes,					//息票付息次数
		CouponRate: bondInfo.CouponRate,				//息票利率
		IssuePrice: bondInfo.IssuePrice,				//发行价
		MaturityDate:bondInfo.MaturityDate[0:10],		//‌兑付日
		BondType: bondInfo.BondType,					//债券类型
	}
	return reqYtmPricingConv(pricingHubYtmConvUrl, pricingHubAuth, &ytmConvInput)
}

func reqYtmPricingConv(pricingHubYtmConvUrl, pricingHubAuth string, reqYtmConvInput *YtmConvInput) (ytmPriceInfo *options.PricingResult, err error) {
	cli := &http.Client{}
	reqBody, err := json.Marshal(reqYtmConvInput)
	req, err := http.NewRequest("POST", pricingHubYtmConvUrl, bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", pricingHubAuth)
	log.Printf("reqYtmPricingConv reqBody:%s\n", reqBody)
	resp, err := cli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	ytmOutput := YtmOutput{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &ytmOutput)
	if err != nil {
		return
	}
	if ytmOutput.ErrCode.Code != 200 {
		err = errors.New(ytmOutput.ErrMsg)
		return
	}
	ytmPriceInfo = ytmOutput.Data
	return
}

func reqBondInfo(dataQryUrl, symbol string) (bondInfo BondInfo, err error) {
	cli := &http.Client{}
	req, err := http.NewRequest("GET", dataQryUrl, nil)
	qry := req.URL.Query()
	qry.Add("collection", "Security")
	qry.Add("key", symbol)
	req.URL.RawQuery = qry.Encode()
	resp, err := cli.Do(req)
	if err != nil {
		err = fmt.Errorf("req bonfinfo %s error:%v", symbol, err)
		return
	}
	defer resp.Body.Close()

	bondInfoOutput := BondInfoOutput{}
	body, _ := io.ReadAll(resp.Body)
	log.Printf("reqBondInfo bond body: %s\n", body)
	err = json.Unmarshal(body, &bondInfoOutput)
	if err != nil {
		err = fmt.Errorf("req bonfinfo %s error:%v", symbol, err)
		return
	}
	if bondInfoOutput.Code != "" {
		err = fmt.Errorf("req bonfinfo %s error, errcode:%s, msg:%s", symbol, bondInfoOutput.Code, bondInfoOutput.Msg)
		return
	}
	if len(bondInfoOutput.Data) == 0 {
		err = fmt.Errorf("req bonfinfo %s empty", symbol)
		return
	}

	bondInfo = bondInfoOutput.Data[0]
	return
}

/*
Titans 枚举值：
	Actual/360 ACT/365 ACT/365NoLeap ACT/ACT(ISMA) Actual/Actual (ISMA) ACT/ACT
PricingHub 枚举值：
	//ACT365, BUSINESS244, ACTACT_ISDA, ACT365_NOLEAP, ACT360, BUSINESS252, ACTACT, ACT365F, THIRTY360, THIRTY365, BUSINESS242
	模型实际支持列表：ACT360 ACT365 ACT365F ACTACT ACTACT_ISDA THIRTY360 兜底默认值为 ACTACT
*/
func dayCountConvert(dayCount string) string {
	dayCount = strings.ToUpper(dayCount)
	dayCount = strings.ReplaceAll(dayCount, "ACTUAL", "ACT")
	dayCount = strings.ReplaceAll(dayCount, "/", "")
	switch dayCount {
	case "ACTACT", "ACT360", "ACT365", "ACT365F", "ACTACT_ISDA", "THIRTY360":
		return dayCount
	case "ACT365NOLEAP":
		return "ACT365F"
	case "ACTACT(ISMA)", "ACTACT (ISMA)":
		return "ACTACT_ISMA"
	}
	return "ACTACT"
}



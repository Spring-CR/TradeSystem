package ficc_fut

import (
	"fmt"
	"log"
	"math"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-core/schema"
	"strconv"
	"strings"

	json "github.com/bytedance/sonic"
)

func (a *FiccFurAPIAdapter) RefineAndValidate(tradeOrder *schema.TradeOrder, trade bool) *domain_error.Error {

	// 正确设置交易对手
	// 检查交易对手是否有业务资格
	counterpartyID, ok, err := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "account", enum.AttrValueType_INT)
	if err != nil {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_ERR_CODE, err)
	}
	if !ok {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COUNTERPARTY_CANNNOT_BE_EMPTY_ERR_CODE, nil)
	}
	counterpartyIDStr := strconv.Itoa(counterpartyID.(int))
	valList, ok, de := a.autoSyncRepo.Get("Counterparty", counterpartyIDStr)
	if de != nil {
		return de
	}
	if !ok || len(valList) == 0 {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COUNTERPARTY_NOT_FOUND_ERR_CODE, nil, counterpartyID)
	}
	// 设置交易对手短名
	ctpyShortName, _, _ := attrutil.GetAttrValue(valList[0], "Counterparty", enum.AttrValueType_STRING)
	tradeOrder.ExtendAttrMap["counterparty"] = ctpyShortName

	allowedBusinessTypes, _, _ := attrutil.GetAttrValue(valList[0], "AllowedBusinessTypes", enum.AttrValueType_STRING)
	if !strings.Contains(allowedBusinessTypes.(string), "TRS") {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.CUST_QUAL_NOT_CORRECT_ERR_CODE, nil)
	}

	//-----------------------------------------------------------------------------------------------------------------------------------------------

	// 正确设置业务方案
	// 业务方案编号校验
	planCode, ok, err := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "planCode", enum.AttrValueType_STRING)
	if err != nil {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_ERR_CODE, err)
	}
	if !ok || planCode == "" {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.BUSINESS_PLANCODE_CANNNOT_BE_EMPTY_ERR_CODE, nil)
	}
	valList, ok, de = a.autoSyncRepo.Get("BusinessPlan", fmt.Sprintf("%v-%v", counterpartyID, planCode))
	if de != nil {
		return de
	}
	if !ok || len(valList) == 0 {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.BUSINEDD_PLAN_NOT_FOUND_ERR_CODE, nil, planCode, counterpartyID)
	}
	businessPlanData := valList[len(valList)-1]
	// 设置大合约编号
	ultraContractCode, _, _ := attrutil.GetAttrValue(businessPlanData, "UltraContractCode", enum.AttrValueType_STRING)
	tradeOrder.ExtendAttrMap["ultraContractCode"] = ultraContractCode
	// 设置业务类型：C、N、S
	businessType, _, _ := attrutil.GetAttrValue(businessPlanData, "BusinessType", enum.AttrValueType_STRING)
	tradeOrder.ExtendAttrMap["businessType"] = businessType

	planID, ok, _ := attrutil.GetAttrValue(valList[len(valList)-1], "PlanID", enum.AttrValueType_INT)
	if !ok {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.MISS_FIELD_ERR_CODE, nil, "业务方案:"+planCode.(string), "PlanID")
	}
	log.Printf("planID: %d\n", planID)

	ordSource, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "ordSource", enum.AttrValueType_INT)
	// 检查交易权限
	if ordSource != "titans" {

		// 获取交易对手权限
		valList, ok, de = a.autoSyncRepo.Get("CounterpartyAuthority", strconv.Itoa(counterpartyID.(int)))
		if de != nil {
			return de
		}
		if !ok || len(valList) == 0 {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COUNTERPARTY_AUTH_NOT_CONFIG_ERR_CODE, nil, counterpartyID)
		}

		authData := valList[0]

		switch businessType {

		case "C":

			authIndicator, ok, err := attrutil.GetAttrValue(authData, "ChnFutIndicator", enum.AttrValueType_INT)
			if err != nil {
				return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_ERR_CODE, err)
			}
			if !ok {
				return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COUNTERPARTY_AUTH_NOT_CONFIG_ERR_CODE, nil, counterpartyID)
			}
			if authIndicator.(int) != 1 {
				return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COUNTERPARTY_AUTH_NOT_ENOUGH2_ERR_CODE, nil, "境内期货")
			}

		case "N":

			authIndicator, ok, err := attrutil.GetAttrValue(authData, "NorthFutIndicator", enum.AttrValueType_INT)
			if err != nil {
				return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_ERR_CODE, err)
			}
			if !ok {
				return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COUNTERPARTY_AUTH_NOT_CONFIG_ERR_CODE, nil, counterpartyID)
			}
			if authIndicator.(int) != 1 {
				return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COUNTERPARTY_AUTH_NOT_ENOUGH2_ERR_CODE, nil, "北向跨境期货")
			}

		case "S":

			authIndicator, ok, err := attrutil.GetAttrValue(authData, "SouthFutIndicator", enum.AttrValueType_INT)
			if err != nil {
				return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_ERR_CODE, err)
			}
			if !ok {
				return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COUNTERPARTY_AUTH_NOT_CONFIG_ERR_CODE, nil, counterpartyID)
			}
			if authIndicator.(int) != 1 {
				return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COUNTERPARTY_AUTH_NOT_ENOUGH2_ERR_CODE, nil, "南向跨境期货")
			}
		}
	}

	//-----------------------------------------------------------------------------------------------------------------------------------------------

	// 正确设置symbol
	valList, ok, de = a.autoSyncRepo.Get("Security", fmt.Sprintf("%v-%v", businessType.(string), tradeOrder.ExtendAttrMap["symbol2"]))
	if de != nil {
		return de
	}
	if !ok || len(valList) == 0 {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.SYMBOL_NOT_FOUND_ERR_CODE, nil, tradeOrder.ExtendAttrMap["symbol2"])
	}
	symbolData := valList[len(valList)-1]
	_symbol, _, _ := attrutil.GetAttrValue(symbolData, "Symbol", enum.AttrValueType_STRING)
	symbol := _symbol.(string)
	if len(symbol) > 0 {
		tradeOrder.Symbol = symbol
	}
	tradeOrder.ExtendAttrMap["symbol"] = symbol

	// 设置标的名称
	symbolName, _, _ := attrutil.GetAttrValue(symbolData, "SecurityName", enum.AttrValueType_STRING)
	tradeOrder.ExtendAttrMap["symbolName"] = symbolName

	// 设置securityID
	securityID, _, _ := attrutil.GetAttrValue(symbolData, "SecurityID", enum.AttrValueType_INT)
	tradeOrder.ExtendAttrMap["securityID"] = securityID

	// 设置交易所信息
	securityExchange, _, _ := attrutil.GetAttrValue(symbolData, "SecurityExchange", enum.AttrValueType_STRING)
	tradeOrder.ExtendAttrMap["securityExchange"] = securityExchange
	securityExchange2, _, _ := attrutil.GetAttrValue(symbolData, "SecurityExchange2", enum.AttrValueType_STRING)
	tradeOrder.ExtendAttrMap["securityExchange2"] = securityExchange2
	tradeOrder.SecurityExchange = securityExchange.(string)

	// 设置合约乘数
	contractMultiplier, _, _ := attrutil.GetAttrValue(symbolData, "ContractMultiplier", enum.AttrValueType_FLOAT)
	tradeOrder.ExtendAttrMap["contractMultiplier"] = contractMultiplier

	// 设置标的币种
	currency, _, _ := attrutil.GetAttrValue(symbolData, "Currency", enum.AttrValueType_STRING)
	tradeOrder.ExtendAttrMap["currency"] = currency
	tradeOrder.Currency = currency.(string)

	// 设置期货类型
	futType, _, _ := attrutil.GetAttrValue(symbolData, "FutType", enum.AttrValueType_STRING)
	tradeOrder.ExtendAttrMap["futType"] = futType

	// 设置标的类型
	securityType, _, _ := attrutil.GetAttrValue(symbolData, "SecurityType", enum.AttrValueType_STRING)
	tradeOrder.SecurityType = securityType.(string)

	// 设置期货标的品种代码
	productCode, _, _ := attrutil.GetAttrValue(symbolData, "ProductCode", enum.AttrValueType_STRING)
	tradeOrder.ExtendAttrMap["productCode"] = productCode

	// 设置汇率
	exchangeRateCNY, _, _ := attrutil.GetAttrValue(symbolData, "ExchangeRateCNY", enum.AttrValueType_FLOAT)
	tradeOrder.ExtendAttrMap["exchangeRateCNY"] = exchangeRateCNY

	// 设置汇率转换的人民币价格
	tradeOrder.ExtendAttrMap["priceCNY"] = tradeOrder.Price * exchangeRateCNY.(float64)

	//-----------------------------------------------------------------------------------------------------------------------------------------------

	// 正确设置资金账户
	valList, ok, de = a.autoSyncRepo.Get("CapitalAccount", strconv.Itoa(counterpartyID.(int))+"-TRS-CNY")
	if de != nil {
		return de
	}
	if !ok || len(valList) == 0 {
		valList, ok, de = a.autoSyncRepo.Get("CapitalAccount", strconv.Itoa(counterpartyID.(int))+"-MIXTURE-CNY")
		if de != nil {
			return de
		}
	}
	if !ok || len(valList) == 0 {
		// 资金账号不通过
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.CAP_ACCT_NOT_FOUND_ERR_CODE, nil)
	}
	capitalAccountData := valList[len(valList)-1]
	capitalAcctID, _, _ := attrutil.GetAttrValue(capitalAccountData, "CapAcctID", enum.AttrValueType_INT)
	tradeOrder.ExtendAttrMap["capitalAcctID"] = int64(capitalAcctID.(int))

	//-----------------------------------------------------------------------------------------------------------------------------------------------

	// 正确设置佣金
	valList, ok, de = a.autoSyncRepo.Get("Commission", fmt.Sprintf("%v-%v-%v", counterpartyID, planCode, productCode))
	if de != nil {
		return de
	}
	if !ok || len(valList) == 0 {
		// 资金账号不通过
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COMMISSION_NOT_FOUND_ERR_CODE, nil)
	}
	commissionData := valList[len(valList)-1]
	commissionType, _, _ := attrutil.GetAttrValue(commissionData, "CommissionType", enum.AttrValueType_STRING)
	tradeOrder.ExtendAttrMap["commissionType"] = commissionType
	var commissionValue interface{}
	switch tradeOrder.Side {
	case "1":
		commissionValue, _, _ = attrutil.GetAttrValue(commissionData, "BuyCommissionValue", enum.AttrValueType_FLOAT)
	case "2":
		commissionValue, _, _ = attrutil.GetAttrValue(commissionData, "SellCommissionValue", enum.AttrValueType_FLOAT)
	}
	tradeOrder.ExtendAttrMap["commissionValue"] = commissionValue

	//-----------------------------------------------------------------------------------------------------------------------------------------------

	// 正确初保比率
	valList, ok, de = a.autoSyncRepo.Get("MarginThreshold", fmt.Sprintf("%v-%v-%v", counterpartyID, planCode, productCode))
	if de != nil {
		return de
	}
	if !ok || len(valList) == 0 {
		log.Printf("cannot fond MarginThreshold by %s, used %s instend!", fmt.Sprintf("%v-%v-%v", counterpartyID, planCode, productCode), fmt.Sprintf("%v-%v-%v", counterpartyID, planCode, ""))
		valList, ok, de = a.autoSyncRepo.Get("MarginThreshold", fmt.Sprintf("%v-%v-%v", counterpartyID, planCode, ""))
	}
	if !ok || len(valList) == 0 {
		// 初保比率参数找不到
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.MARGIN_NOT_FOUND_ERR_CODE, nil)
	}
	marginThresholdData := valList[len(valList)-1]
	marginRatio, _, _ := attrutil.GetAttrValue(marginThresholdData, "MarginRatio", enum.AttrValueType_FLOAT)
	tradeOrder.ExtendAttrMap["marginRatio"] = marginRatio

	//-----------------------------------------------------------------------------------------------------------------------------------------------

	// 正确设置交易通道
	var channelCode string
	switch securityExchange {
	case "3", "4", "5", "6", "8", "R":
		channelCode = "olts-fut"
	default:
		channelCode = "stars-fut"
	}
	tradeOrder.ChannelCode = channelCode

	//-----------------------------------------------------------------------------------------------------------------------------------------------

	// 额度检查
	if tradeOrder.ApproveStatus != int(enum.ApproveStatus_Approved) {

		var limitLongQty, limitShortQty, limitNotional interface{}
		limitLongQty = math.MaxFloat64
		limitShortQty = math.MaxFloat64
		limitNotional = math.MaxFloat64

		valList, _, _ = a.autoSyncRepo.Get("RiskSecurity", productCode.(string))
		if len(valList) > 0 {
			limitLongQty, _, _ = attrutil.GetAttrValue(valList[0], "LongMaxTradeLotValue", enum.AttrValueType_FLOAT)
			limitShortQty, _, _ = attrutil.GetAttrValue(valList[0], "ShortMaxTradeLotValue", enum.AttrValueType_FLOAT)
		}

		valList, _, _ = a.autoSyncRepo.Get("RiskCounterparty", counterpartyIDStr)
		if len(valList) > 0 {
			limitNotional, _, _ = attrutil.GetAttrValue(valList[0], "MaxNotionalValue", enum.AttrValueType_FLOAT)
		}

		tradeOrder.ExtendAttrMap["limitLongQty"] = limitLongQty
		tradeOrder.ExtendAttrMap["limitShortQty"] = limitShortQty
		tradeOrder.ExtendAttrMap["limitNotional"] = limitNotional
	}

	//-----------------------------------------------------------------------------------------------------------------------------------------------

	// 设置含费priceCNY
	priceCNYWithFee := GetPriceWithFee(tradeOrder.ExtendAttrMap["priceCNY"].(float64), commissionType.(string), commissionValue.(float64), tradeOrder.Side)
	tradeOrder.ExtendAttrMap["priceCNYWithFee"] = priceCNYWithFee

	// 设置tradeDate
	channelDetails, ok := a.channelMap[channelCode]
	if ok {
		tradeDate := channelDetails.GetCurrentExchangeDate()
		tradeOrder.ExtendAttrMap["tradeDate"] = tradeDate

		log.Printf("get tradeDate:%v\n", tradeDate)
	} else {
		log.Printf("cannot get channelDetails for order:%s\n", tradeOrder.AppOrdID)
	}

	
	de = a.checkTradeTime(tradeOrder, trade)
	if de != nil {
		return de
	}

	extendAttr, _ := json.Marshal(tradeOrder.ExtendAttrMap)
	tradeOrder.ExtendAttr = string(extendAttr)

	return nil
}

func GetPriceWithFee(price float64, commissionType string, commissionValue float64, side string) float64 {

	lastPrice := price

	if side == "1" {
		switch commissionType {
		case "FEE_BY_RATE":
			price *= 1 + commissionValue
		case "FEE_BY_PER_SHARE":
			price += commissionValue
		}
	}

	if side == "2" {
		switch commissionType {
		case "FEE_BY_RATE":
			price *= 1 - commissionValue
		case "FEE_BY_PER_SHARE":
			price -= commissionValue
		}
	}

	if price <= 0 {
		price = lastPrice
	}

	return price
}

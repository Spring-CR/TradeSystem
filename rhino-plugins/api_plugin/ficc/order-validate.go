package ficc

import (
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	"rhino-plugins/api_plugin/util"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/manucorporat/try"
)

var (
	// 因为每天会重启，所以不需要清空逻辑
	dumplicatOrderCheckMap     = make(map[string]bool)
	dumplicatOrderCheckMapLock = &sync.Mutex{}
)

func (a *TitansFiccAPIAdapter) RefineAndValidate(tradeOrder *schema.TradeOrder, trade bool) *domain_error.Error {

	// 检查校验时间
	var txTime int64
	if trade {
		timeNow := timeutil.ConvertTimeToMilliseconds(time.Now())
		txTime = tradeOrder.TransactTime
		if timeNow > txTime {
			txTime = timeNow
		}
		if txTime < a.tradeTimeBegin || txTime > a.tradeTimeEnd {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.NOT_IN_TRADING_TIME_ERR_CODE, nil)
		}
	}

	// 检查id长度
	if len(tradeOrder.AppOrdID) > 128 {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.APP_ORD_ID_TOO_LONG_ERR_CODE, nil, 128)
	}

	// 检查交易对手是否有业务资格
	counterpartyID, ok, err := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "account", enum.AttrValueType_INT)
	if err != nil {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_ERR_CODE, err)
	}
	if !ok {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COUNTERPARTY_CANNNOT_BE_EMPTY_ERR_CODE, nil)
	}
	valList, ok, de := a.autoSyncRepo.Get("Counterparty", strconv.Itoa(counterpartyID.(int)))
	if de != nil {
		return de
	}
	if !ok || len(valList) == 0 {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COUNTERPARTY_NOT_FOUND_ERR_CODE, nil, counterpartyID)
	}
	// 设置investorID，在FIX的时候，默认没有设置的
	ctpyShortName, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "counterparty", enum.AttrValueType_STRING)
	if ctpyShortName == "" {
		ctpyShortName, _, _ = attrutil.GetAttrValue(valList[0], "Counterparty", enum.AttrValueType_STRING)
		tradeOrder.ExtendAttrMap["counterparty"] = ctpyShortName
	}
	if ctpyShortName != "" {
		// 根据交易对手短名，设置InvestorID
		if ctpyShortName == "广发全球资本FICC" {
			tradeOrder.ExtendAttrMap["investorID"] = a.twoInvestorID[0]
		} else {
			tradeOrder.ExtendAttrMap["investorID"] = a.twoInvestorID[1]
		}
	}
	allowedBusinessTypes, _, _ := attrutil.GetAttrValue(valList[0], "AllowedBusinessTypes", enum.AttrValueType_STRING)
	if !strings.Contains(allowedBusinessTypes.(string), "TRS") {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.CUST_QUAL_NOT_CORRECT_ERR_CODE, nil)
	}

	// 设置counterpartyID
	//tradeOrder.ExtendAttrMap["counterpartyID"] = counterpartyID
	valList, ok, de = a.autoSyncRepo.Get("CounterpartyAuthority", strconv.Itoa(counterpartyID.(int)))
	if de != nil {
		return de
	}
	if !ok || len(valList) == 0 {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COUNTERPARTY_AUTH_NOT_CONFIG_ERR_CODE, nil, counterpartyID)
	}

	ordSource, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "ordSource", enum.AttrValueType_INT)

	if ordSource != "titans" {
		xBondTrsIndicator, ok, err := attrutil.GetAttrValue(valList[len(valList)-1], "XbondTrsIndicator", enum.AttrValueType_INT)
		if err != nil {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_ERR_CODE, err)
		}
		if !ok {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COUNTERPARTY_AUTH_NOT_CONFIG_ERR_CODE, nil, counterpartyID)
		}
		if xBondTrsIndicator.(int) != 1 {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COUNTERPARTY_AUTH_NOT_ENOUGH_ERR_CODE, nil)
		}
	}

	// 业务方案编号校验
	planCode, ok, err := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "planCode", enum.AttrValueType_STRING)
	if err != nil {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_ERR_CODE, err)
	}
	if !ok || planCode == "" {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.BUSINESS_PLANCODE_CANNNOT_BE_EMPTY_ERR_CODE, nil)
	}
	valList, ok, de = a.autoSyncRepo.Get("BusinessPlan", fmt.Sprintf("%v-%v", strconv.Itoa(counterpartyID.(int)), planCode))
	if de != nil {
		return de
	}
	if !ok || len(valList) == 0 {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.BUSINEDD_PLAN_NOT_FOUND_ERR_CODE, nil, planCode, counterpartyID)
	}
	ultraContractCode, ok, _ := attrutil.GetAttrValue(valList[len(valList)-1], "UltraContractCode", enum.AttrValueType_STRING)
	if ok {
		// 设置大合约编号
		tradeOrder.ExtendAttrMap["ultraContractCode"] = ultraContractCode
	}
	planID, ok, _ := attrutil.GetAttrValue(valList[len(valList)-1], "PlanID", enum.AttrValueType_INT)
	if !ok {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.MISS_FIELD_ERR_CODE, nil, "业务方案:"+planCode.(string), "PlanID")
	}

	// 标的代码校验（开仓检查，平仓不检查）
	symbol, ok, err := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "symbol", enum.AttrValueType_STRING)
	if err != nil {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_ERR_CODE, err)
	}
	if !ok || symbol == "" {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.SYMBOL_CANNNOT_BE_EMPTY_ERR_CODE, nil)
	}
	valList, ok, de = a.autoSyncRepo.Get("Security", symbol.(string))
	if de != nil {
		return de
	}
	if !ok || len(valList) == 0 {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.SYMBOL_NOT_FOUND_ERR_CODE, nil, symbol)
	}
	_, ok, _ = attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "openClose", enum.AttrValueType_STRING)
	if ok {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_VALUE_NOT_CORRECT_ERR_CODE, nil, "本系统最新API已不支持【开平仓标志】，您仅需正确设置订单的【交易方向】，详情请联系业务人员或查阅最新版API文档")
	}
	symbolData := valList[len(valList)-1]
	//if openClose == "O" {
	if tradeOrder.Side == "1" {
		chnBondPoolIndicator, ok, err := attrutil.GetAttrValue(symbolData, "ChnBondPoolIndicator", enum.AttrValueType_STRING)
		if err != nil {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_ERR_CODE, err)
		}
		if !ok {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.MISS_FIELD_ERR_CODE, nil, "标的", "ChnBondPoolIndicator")
		}
		if chnBondPoolIndicator != "Y" {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.SYMBOL_NOT_IN_WHITE_POOL_ERR_CODE, nil, symbol)
		}
	} else {
		overSoldIndicator, _, _ := attrutil.GetAttrValue(symbolData, "OverSoldIndicator", enum.AttrValueType_STRING)
		if overSoldIndicator == "Y" {
			tradeOrder.ExtendAttrMap["allowOverSold"] = true
		} else {
			tradeOrder.ExtendAttrMap["allowOverSold"] = false
		}
	}

	// 券面总额校验
	quantity, ok, err := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "quantity", enum.AttrValueType_FLOAT)
	if err != nil {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_ERR_CODE, err)
	}
	if !ok || quantity.(float64) <= 0 {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.NUM_VALUE_CANNNOT_BE_EMPTY_ERR_CODE, nil, "券面总额")
	}

	// 意向到期收益率校验
	if ordSource.(string) != "FIX" {

		// http下单，检查重复的订单号
		if tradeOrder.ID <= 0 && trade { // 只有id大于0，才可能插入数据库
			dumplicatOrderCheckMapLock.Lock()
			if dumplicatOrderCheckMap[tradeOrder.AppOrdID] {
				// 重复订单
				dumplicatOrderCheckMapLock.Unlock()
				return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.DUPLICATE_ORDER_ERR_CODE, nil, tradeOrder.AppOrdID)
			} else {
				dumplicatOrderCheckMap[tradeOrder.AppOrdID] = true
				dumplicatOrderCheckMapLock.Unlock()
			}
		}

		ytm, ok, err := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "ytm", enum.AttrValueType_FLOAT)
		if err != nil {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_ERR_CODE, err)
		}
		if !ok || ytm.(float64) <= 0 {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.NUM_VALUE_CANNNOT_BE_EMPTY_ERR_CODE, nil, "意向到期收益率%")
		}
	} else { // 检查广发通电话、交易对手是否有关联
		phoneNum, _, err := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "phoneNum", enum.AttrValueType_STRING)
		if err != nil {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_ERR_CODE, err)
		}
		valList, ok, de = a.autoSyncRepo.Get("CounterpartyPhoneNum", strconv.Itoa(counterpartyID.(int))+"-"+phoneNum.(string))
		if de != nil {
			return de
		}
		if !ok || len(valList) == 0 {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.PRODUCT_ACCT_NOT_FOUND_ERR_CODE, nil, counterpartyID, phoneNum)
		}
	}

	// 意向净价校验
	price, ok, err := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "price", enum.AttrValueType_FLOAT)
	if err != nil {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_ERR_CODE, err)
	}
	if !ok || price.(float64) <= 0 {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.NUM_VALUE_CANNNOT_BE_EMPTY_ERR_CODE, nil, "意向净价")
	}

	// 意向全价校验
	/*var dirtyPrice interface{} = 0.0
	if ordSource.(string) != "FIX" {
		dirtyPrice, ok, err = attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "dirtyPrice", enum.AttrValueType_FLOAT)
		if err != nil {
			return domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		}
		if !ok || dirtyPrice.(float64) <= 0 {
			return domain_error.Build(domain_error.NUM_VALUE_CANNNOT_BE_EMPTY_ERR_CODE, nil, "意向全价")
		}
	}
	if dirtyPrice == 0.0 {
		dirtyPrice = price.(float64) * posificc.DirtyPriceRation
	}*/

	// 设置标的名称债券面值和币种
	parValue, ok, err := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "parValue", enum.AttrValueType_INT)
	if (!ok || err != nil || parValue.(int) <= 0) && symbolData != nil {
		parValue, _, _ = attrutil.GetAttrValue(symbolData, "ParValue", enum.AttrValueType_INT)
		if parValue.(int) > 0 {
			tradeOrder.ExtendAttrMap["parValue"] = parValue
		}
	}
	if parValue == 0 {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.PAR_VALUE_NOT_FOUD_ERR_CODE, nil, symbol)
	}

	val, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "settlType", enum.AttrValueType_STRING)
	settlType := val.(string)
	if settlType != "T+0" && settlType != "T+1" {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.SETTLE_SPEED_NOT_CORRECT_ERR_CODE, nil)
	}
	if trade && settlType == "T+0" && tradeOrder.Side == "2" && txTime > a.t0SellTimeEnd {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.NOT_IN_T0_TRADING_TIME_ERR_CODE, nil, a.t0SellEndTimeStr)
	}

	dirtyPrice, ok, err := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "dirtyPrice", enum.AttrValueType_FLOAT)
	if err != nil {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_ERR_CODE, err)
	}
	if !ok || dirtyPrice.(float64) <= 0 {
		// 如果意向全价不存在，需要调用计算器把全价计算出来
		var ytm float64
		try.This(func() {
			dirtyPrice, ytm, err = computeDirtyPrice(a.applicationCfg, settlType, tradeOrder, parValue.(int), price.(float64), symbolData)
		}).Catch(func(err try.E) {
			log.Printf("error occur while run computeDirtyPrice! error:%v\n", err)
			domain_error.Build(domain_error.GENERIC_ERR_CODE, fmt.Errorf("error occur while run computeDirtyPrice! error:%v", err))
		})
		if err != nil {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_ERR_CODE, err)
		}
		tradeOrder.ExtendAttrMap["dirtyPrice"] = dirtyPrice
		tradeOrder.ExtendAttrMap["ytm"] = ytm
	}

	currency, ok, de := util.GetStringValueInField(tradeOrder.ExtendAttrMap, "currency")
	if (!ok || de != nil || currency == "") && symbolData != nil {
		currency, _, _ = util.GetStringValueInField(symbolData, "Currency")
		if currency != "" {
			tradeOrder.ExtendAttrMap["currency"] = currency
		}
	}
	symbolName, ok, de := util.GetStringValueInField(tradeOrder.ExtendAttrMap, "symbolName")
	if (!ok || de != nil || symbolName == "") && symbolData != nil {
		symbolName, _, _ = util.GetStringValueInField(symbolData, "SecurityName")
		if symbolName != "" {
			tradeOrder.ExtendAttrMap["symbolName"] = symbolName
		}
	}
	var securityType string
	var bondType string
	if symbolData != nil {
		securityType, _, _ = util.GetStringValueInField(symbolData, "SecurityType")
		if securityType != "" {
			tradeOrder.ExtendAttrMap["securityType"] = securityType
		}
		securityID, _, _ := attrutil.GetAttrValue(symbolData, "SecurityID", enum.AttrValueType_INT)
		if securityID.(int) > 0 {
			tradeOrder.ExtendAttrMap["securityID"] = int64(securityID.(int))
		} else {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.SECURITY_ID_NOT_FOUD_ERR_CODE, nil, symbol)
		}
		// 债券类型
		bondType, _, _ = util.GetStringValueInField(symbolData, "BondType")
	}

	// 资金账户校验
	capitalAcctID, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "capitalAcctID", enum.AttrValueType_INT)
	if capitalAcctID.(int) <= 0 {
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
		capitalAcctID, _, _ = attrutil.GetAttrValue(capitalAccountData, "CapAcctID", enum.AttrValueType_INT)
		if capitalAcctID.(int) <= 0 {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.NUM_VALUE_CANNNOT_BE_EMPTY_ERR_CODE, nil, "资金账户ID")
		}
	} else {
		valList, ok, de = a.autoSyncRepo.Get("CapitalAccount", strconv.Itoa(capitalAcctID.(int))+"-"+strconv.Itoa(counterpartyID.(int))+"-TRS-CNY")
		if de != nil {
			return de
		}
		if !ok || len(valList) == 0 {
			valList, ok, de = a.autoSyncRepo.Get("CapitalAccount", strconv.Itoa(capitalAcctID.(int))+"-"+strconv.Itoa(counterpartyID.(int))+"-MIXTURE-CNY")
			if de != nil {
				return de
			}
		}
		if !ok || len(valList) == 0 {
			// 资金账号不通过
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.CAP_ACCT_NOT_FOUND_ERR_CODE, nil)
		}
	}
	tradeOrder.ExtendAttrMap["capitalAcctID"] = int64(capitalAcctID.(int))

	val, _, _ = attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "clOrdID", enum.AttrValueType_STRING)
	if val.(string) == "" {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.CLORDID_CANNOT_BE_EMPTY_ERR_CODE, nil)
	}

	val, _, _ = attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "side", enum.AttrValueType_STRING)
	if val.(string) != "1" && val.(string) != "2" {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_VALUE_NOT_CORRECT_ERR_CODE, nil, "交易方向需设置为 1（买入）、2（卖出）")
	}

	val, _, _ = attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "ordType", enum.AttrValueType_STRING)
	if val.(string) != "2" {
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.ONLY_SUPPORT_PRICE_LIMIT_ORDER_ERR_CODE, nil)
	}

	if ordSource.(string) == "titans" {
		if tradeOrder.HandlInst != "1" && tradeOrder.HandlInst != "2" && tradeOrder.HandlInst != "3" && tradeOrder.HandlInst != "4" {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_VALUE_NOT_CORRECT_ERR_CODE, nil, "交易效率需要设置为 1（快速交易）、2（策略交易）、3（普通交易）、4（补录交易）")
		}
	} else {
		if tradeOrder.HandlInst != "1" && tradeOrder.HandlInst != "2" && tradeOrder.HandlInst != "3" {
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.GENERIC_VALUE_NOT_CORRECT_ERR_CODE, nil, "交易效率需要设置为 1（快速交易）、2（策略交易）、3（普通交易）")
		}
	}

	tradeOrder.ExtendAttrMap["limitCheckResult"] = 0

	//if openClose == "O" {
	// 支持空头之后，初始和预估保证金都不用设置了
	// if tradeOrder.Side == "1" {
	// 	initMarginRatio, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "initMarginRatio", enum.AttrValueType_FLOAT)
	// 	initMarginAmount := quantity.(float64) / float64(parValue.(int)) * dirtyPrice.(float64) * initMarginRatio.(float64) / 100.0
	// 	tradeOrder.ExtendAttrMap["initMarginAmount"] = initMarginAmount
	// 	tradeOrder.ExtendAttrMap["estiFrozenAmount"] = initMarginAmount
	// }

	//log.Printf("tradeOrder.ApproveStatus=%v\n", tradeOrder.ApproveStatus)
	if tradeOrder.ApproveStatus != int(enum.ApproveStatus_Approved) {
		// 仅对开仓做限额检查
		// 限额检查：TRS_DOMESTIC_BOND_SINGLE_NOTIONAL
		valList, _, de = a.autoSyncRepo.Get("RiskLimit", "TRS_DOMESTIC_BOND_SINGLE_NOTIONAL")
		if de != nil {
			return de
		}
		// log.Printf("RiskLimit:%v\n", valList)
		if len(valList) > 0 {
			metricData := valList[len(valList)-1]
			enable, _, _ := attrutil.GetAttrValue(metricData, "EnabledIndicator", enum.AttrValueType_INT)
			// log.Printf("EnabledIndicator: %v\n", enable)
			if enable.(int) > 0 {
				threshold, _, _ := attrutil.GetAttrValue(metricData, "Threshold", enum.AttrValueType_INT)
				// log.Printf("threshold: %v, quantity:%v\n", threshold, quantity)
				if threshold.(int) > 0 {
					// 订单名额不能大于10亿元
					if int(quantity.(float64)) > threshold.(int) {
						tradeOrder.ExtendAttrMap["limitCheckResult"] = 2
						de = domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.QUOTA_LIMIT_EXCEEDED_ERR_CODE, nil, fmt.Sprintf("订单券面总额大于%v亿元", threshold.(int)/100000000))
						return de
					} else {
						tradeOrder.ExtendAttrMap["limitCheckResult"] = 1
					}
				}
			}
		}
	}

	// 20260331, 设置佣金和利差
	//var handlInstName string
	// switch tradeOrder.HandlInst {
	// case "1":
	// 	handlInstName = "FAST"
	// case "3":
	// 	handlInstName = "NORMAL"
	// }
	// handlInstName := tradeOrder.HandlInst
	// valList, ok, de = a.autoSyncRepo.Get("CommissionParam", fmt.Sprintf("%v-%v-%v-%v", planID, securityType, bondType, handlInstName))
	// if de != nil {
	// 	return de
	// }
	// if !ok || len(valList) == 0 {
	// 	log.Printf("fail to get commission %v-%v-%v-%v", planID, securityType, bondType, handlInstName)
	// 	return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COMMISSION_NOT_FOUND_ERR_CODE, nil)
	// }
	// commissionData := valList[len(valList)-1]

	handlInstKey := tradeOrder.HandlInst
	switch tradeOrder.HandlInst {
	case "2":
		tradeOrder.AlgName = "FullMarket"
	case "3":
		tradeOrder.HandlInst = "2"
		tradeOrder.AlgName = "InsideBrokers"
	case "4":
		tradeOrder.HandlInst = "3"
		handlInstKey = "3"
		tradeOrder.AlgName = ""
	}
	
	commissionRate, ok := a.getComissionRate(planID, securityType, bondType, handlInstKey, tradeOrder.Side)
	var commissionRate1 float64
	var commissionRate1Ok bool
	if tradeOrder.AlgName != "FullMarket" && !ok {
		log.Printf("fail to get commission %v-%v-%v-%v", planID, securityType, bondType, handlInstKey)
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COMMISSION_NOT_FOUND_ERR_CODE, nil)
	} else if !ok {
		commissionRate, ok = a.getComissionRate(planID, securityType, bondType, "3", tradeOrder.Side)
		if !ok {
			log.Printf("fail to get commission %v-%v-%v-%v", planID, securityType, bondType, "3")
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COMMISSION_NOT_FOUND_ERR_CODE, nil)
		}
		commissionRate1, commissionRate1Ok = a.getComissionRate(planID, securityType, bondType, "1", tradeOrder.Side)
		if !commissionRate1Ok {
			log.Printf("fail to get commission %v-%v-%v-%v", planID, securityType, bondType, "1")
			return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.COMMISSION_NOT_FOUND_ERR_CODE, nil)
		}
	}

	if tradeOrder.Side == "1" { // 买单
		tradeOrder.ExtendAttrMap["commissionRate"] = commissionRate
		tradeOrder.ExtendAttrMap["dirtyPriceWithFee"] = dirtyPrice.(float64) * (1 + commissionRate)
	} else { // 卖单
		tradeOrder.ExtendAttrMap["commissionRate"] = commissionRate
		tradeOrder.ExtendAttrMap["dirtyPriceWithFee"] = dirtyPrice.(float64) * (1 - commissionRate)
	}

	if commissionRate1Ok {
		tradeOrder.ExtendAttrMap["commissionRate1"] = commissionRate1
		tradeOrder.ExtendAttrMap["commissionRate3"] = commissionRate
	}

	valList, ok, de = a.autoSyncRepo.Get("InterestParam", fmt.Sprintf("%v-%v-%v-%v", planID, securityType, bondType, "LONG"))
	if de != nil {
		return de
	}
	if !ok || len(valList) == 0 {
		log.Printf("fail to get interest %v-%v-%v", planID, securityType, bondType)
		return domain_error.BuildWithDetails(domain_error.ERROR, tradeOrder, domain_error.INTEREST_NOT_FOUND_ERR_CODE, nil)
	}
	interestData := valList[len(valList)-1]
	scale := 1.0
	commissionType, _, _ := attrutil.GetAttrValue(interestData, "CommissionType", enum.AttrValueType_STRING)
	if commissionType == "FEE_BY_RATE" {
		scale = 100.0
	}
	val, _, _ = attrutil.GetAttrValue(interestData, "Spread", enum.AttrValueType_FLOAT)
	tradeOrder.ExtendAttrMap["spread"] = val.(float64) / scale

	// 增加三个成交信息
	tradeOrder.ExtendAttrMap["avgPrice"] = tradeOrder.Price
	tradeOrder.ExtendAttrMap["avgDirtyPrice"] = tradeOrder.ExtendAttrMap["dirtyPrice"]
	tradeOrder.ExtendAttrMap["avgDirtyPriceWithFee"] = tradeOrder.ExtendAttrMap["dirtyPriceWithFee"]

	extendAttr, _ := json.Marshal(tradeOrder.ExtendAttrMap)
	tradeOrder.ExtendAttr = string(extendAttr)
	return nil
}

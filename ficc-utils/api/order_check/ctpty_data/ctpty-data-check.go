package ctpty_data

import (
	"ficc-utils/common/utils/data_qry"
	"ficc-utils/common/utils/wechat"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
)

var (
	ctptyNames data_qry.CachedData[map[int]string]
	authCtptyIds data_qry.CachedData[[]int]
	capitalAccounts data_qry.CachedData[map[int][]CapitalAccount]
	businessPlans data_qry.CachedData[map[int][]BusinessPlan]
	commissionParams data_qry.CachedData[map[int][]CommissionParam]
	interestParams  data_qry.CachedData[map[int][]InterestParam]

	handleSpeeds = []string{"1", "2", "3"}
	bondTypes  = []string{"GOVERNMENT_BOND", "POLICY_BANK_BOND", "LOCAL_TREASURY_BOND"}
)

type Counterparty struct {
	CounterpartyID int
	Counterparty string
}

type CtptyAuthority struct {
	CounterpartyID int
	XbondTrsIndicator int
}

type CapitalAccount struct {
	CapAcctCode string
	CapAcctID int
	CapAcctName string
	Counterparty string
	CounterpartyID int
	Currency string
	Purpose string
}

type BusinessPlan struct {
	BusinessType string
	CounterpartyID int
	PlanCode string
	PlanId int
	TradeDate string
	UltraContractCode string
	UltraContractID int
}

type CommissionParam struct {
	AlgSource string
	AlgType string
	BondType string
	BuyCommissionRate float64
	BuyConvertCommissionRate float64
	BuyCoverStampDutyRate float64
	BuyStampDutyRate float64
	SellCommissionRate float64
	SellShortCommissionRate float64
	SellShortStampDutyRate float64
	SellStampDutyRate float64
	CommissionRateSource string
	CommissionType string
	FutureVariety int
	HandlInst string
	HkCommissionRate float64
	ID int64
	Notional string
	PlanId int
	SecurityExchange string
	SecurityType string
	TradeDate string
}

type InterestParam struct {
	BondType string
	BookingType string
	FixedRate float64
	ID int
	IncludeLendInterestFlag string
	InterestDirection int
	InterestRatio float64
	InterestType string
	LongShort string
	LookbackPeriod int
	PlanId int
	RateDate string
	Remark string
	SecurityExchange string
	SecurityID int64
	SecurityType string
	Spread float64
	TradeDate string
	KeyInterestID int64
}

func getPlanIds() []int {
	var planIds []int
	for _, plans := range businessPlans.Data {
		for _, plan := range plans {
			planIds = append(planIds, plan.PlanId)
		}
	}
	return planIds
}

func getPlanCode(planId int) string {
	for _, plans := range businessPlans.Data {
		for _, plan := range plans {
			if plan.PlanId == planId {
				return plan.PlanCode
			}
		}
	}
	return strconv.Itoa(planId)
}

func getCtptyName(ctptyId int) string {
	name := ctptyNames.Data[ctptyId]
	name = strings.TrimSpace(name)
	if name == "" {
		return strconv.Itoa(ctptyId)
	}
	return name
}

func getHandlInstName(handlInst string) string {
	switch handlInst {
	case "1": return "快速交易"
	case "2": return "策略交易"
	case "3": return "普通交易"
	default: return handlInst
	}
}

func getBondTypeName(bondType string) string {
	switch bondType {
	case "GOVERNMENT_BOND": return "国债"
	case "POLICY_BANK_BOND": return "政策性银行债"
	case "LOCAL_TREASURY_BOND": return "地方政府债"
	case "IBNCD": return "同业存单"
	default: return bondType
	}
}

func syncAllCachedData(dataQryUrl string) error {
	err := syncCounterparties(dataQryUrl)
	if err != nil {
		return fmt.Errorf("syncCounterparties error: %v", err)
	}

	err = syncCtptyAuthority(dataQryUrl)
	if err != nil {
		return fmt.Errorf("syncCtptyAuthority error: %v", err)
	}
	slices.Sort(authCtptyIds.Data)
	fmt.Printf("syncAllCachedData syncCtptyAuthority:%+v\n", authCtptyIds)

	err = syncCapitalAccount(dataQryUrl, authCtptyIds.Data)
	if err != nil {
		return err
	}

	err = syncBusinessPlans(dataQryUrl, authCtptyIds.Data)
	if err != nil {
		return err
	}

	planIds := getPlanIds()
	err = syncCommissionParam(dataQryUrl, planIds)
	if err != nil {
		return err
	}

	err = syncInterestParam(dataQryUrl, planIds)
	if err != nil {
		return err
	}
	return nil
}

func qryCounterparties(dataQryUrl string) (map[int]string, error) {
	var qryData map[string][]Counterparty
	var data = map[int]string{}
	err := data_qry.QryTableData(dataQryUrl, "Counterparty", "ALL_MAP_DATA", &qryData)
	if err != nil {
		return nil, err
	}
	for _, v := range qryData {
		for _, v1 := range v {
			if v1.Counterparty != "" {
				data[v1.CounterpartyID] = v1.Counterparty
			} else {
				data[v1.CounterpartyID] = strconv.Itoa(v1.CounterpartyID)
			}
		}
	}
	return data, nil
}

func syncCounterparties(dataQryUrl string) error {
	return data_qry.SyncCachedData("counterparties", &ctptyNames, func() (map[int]string, error) {
		return qryCounterparties(dataQryUrl)
	})
}

func qryCtptyAuthority(dataQryUrl string) ([]int, error) {
	var qryData map[string][]CtptyAuthority
	var data []int
	err := data_qry.QryTableData(dataQryUrl, "CounterpartyAuthority", "ALL_MAP_DATA", &qryData)
	if err != nil {
		return nil, err
	}
	for _, v := range qryData {
		for _, v1 := range v {
			if v1.XbondTrsIndicator == 1 {
				data = append(data, v1.CounterpartyID)
			}
		}
	}
	return data, nil
}

func syncCtptyAuthority(dataQryUrl string) error {
	return data_qry.SyncCachedData("counterpartyAuthority", &authCtptyIds, func() ([]int, error) {
		return qryCtptyAuthority(dataQryUrl)
	})
}

func qryCapitalAccount(dataQryUrl string, ctptyIds []int) (map[int][]CapitalAccount, error) {
	var qryData map[string][]CapitalAccount
	var data = map[int][]CapitalAccount{}
	err := data_qry.QryTableData(dataQryUrl, "CapitalAccount", "ALL_MAP_DATA", &qryData)
	if err != nil {
		return nil, err
	}
	for _, v := range qryData {
		for _, v1 := range v {
			if slices.Contains(ctptyIds, v1.CounterpartyID) {
				data[v1.CounterpartyID] = append(data[v1.CounterpartyID], v1)
			}
		}
	}
	return data, nil
}

func syncCapitalAccount(dataQryUrl string, ctptyIds []int) error {
	return data_qry.SyncCachedData("CapitalAccount", &capitalAccounts, func() (map[int][]CapitalAccount, error) {
		return qryCapitalAccount(dataQryUrl, ctptyIds)
	})
}

func qryBusinessPlans(dataQryUrl string, ctptyIds []int) (map[int][]BusinessPlan, error) {
	var qryData map[string][]BusinessPlan
	var data = map[int][]BusinessPlan{}
	err := data_qry.QryTableData(dataQryUrl, "BusinessPlan", "ALL_MAP_DATA", &qryData)
	if err != nil {
		return nil, err
	}
	for k, v := range qryData {
		id, err := strconv.Atoi(k)
		if err == nil && slices.Contains(ctptyIds, id) {
			for _, v1 := range v {
				data[v1.CounterpartyID] = append(data[v1.CounterpartyID], v1)
			}
		}

	}
	return data, nil
}

func syncBusinessPlans(dataQryUrl string, ctptyIds []int) error {
	return data_qry.SyncCachedData("business_plans", &businessPlans, func() (map[int][]BusinessPlan, error) {
		return qryBusinessPlans(dataQryUrl, ctptyIds)
	})
}

func qryCommissionParam(dataQryUrl string, planIds []int) (map[int][]CommissionParam, error) {
	var qryData map[string][]CommissionParam
	var data = map[int][]CommissionParam{}
	err := data_qry.QryTableData(dataQryUrl, "CommissionParam", "ALL_MAP_DATA", &qryData)
	if err != nil {
		return nil, err
	}
	for _, v := range qryData {
		for _, v1 := range v {
			if slices.Contains(planIds, v1.PlanId) {
				data[v1.PlanId] = append(data[v1.PlanId], v1)
			}
		}
	}
	return data, nil
}

func syncCommissionParam(dataQryUrl string, planIds []int) error {
	return data_qry.SyncCachedData("commission_param", &commissionParams, func() (map[int][]CommissionParam, error) {
		return qryCommissionParam(dataQryUrl, planIds)
	})
}

func qryInterestParam(dataQryUrl string, planIds []int) (map[int][]InterestParam, error) {
	var qryData map[string][]InterestParam
	var data = map[int][]InterestParam{}
	err := data_qry.QryTableData(dataQryUrl, "InterestParam", "ALL_MAP_DATA", &qryData)
	if err != nil {
		return nil, err
	}
	for _, v := range qryData {
		for _, v1 := range v {
			if slices.Contains(planIds, v1.PlanId) {
				data[v1.PlanId] = append(data[v1.PlanId], v1)
			}
		}
	}
	return data, nil
}

func syncInterestParam(dataQryUrl string, planIds []int) error {
	return data_qry.SyncCachedData("interest_param", &interestParams, func() (map[int][]InterestParam, error) {
		return qryInterestParam(dataQryUrl, planIds)
	})
}

func checkTrsCtptyData(dataQryUrl string)  (passed bool, checkInfo string, err error)  {
	passed = true
	err = syncAllCachedData(dataQryUrl)
	if err != nil {
		return false, "", err
	}

	// 资金账户检查
	accCheckInfo := ""
	for _, ctptyId := range authCtptyIds.Data {
		if _, ok := capitalAccounts.Data[ctptyId]; !ok {
			passed = false
			accCheckInfo += fmt.Sprintf("交易对手：%s\n", getCtptyName(ctptyId))
		}
	}
	if accCheckInfo != "" {
		checkInfo += "【以下交易对手缺少资金账户】\n"
		checkInfo += accCheckInfo
	}

	// 业务方案检查
	planCheckInfo := ""
	for _, ctptyId := range authCtptyIds.Data {
		if _, ok := businessPlans.Data[ctptyId]; !ok {
			passed = false
			planCheckInfo += fmt.Sprintf("交易对手：%s\n", getCtptyName(ctptyId))
		}
	}
	if planCheckInfo != "" {
		if len(checkInfo) > 0 {
			checkInfo += "\n"
		}
		checkInfo += "【以下交易对手缺少业务方案】\n"
		checkInfo += planCheckInfo
	}

	// 佣金费率参数检查
	commissionCheckInfo1 := ""
	commissionCheckInfo2 := ""
	planIds := getPlanIds()
	for _, planId := range planIds {
		params, _ := commissionParams.Data[planId]
		//if !ok {
		//	passed = false
		//	commissionCheckInfo1 += fmt.Sprintf("业务方案：%s\n", getPlanCode(planId))
		//	continue
		//}

		paramMap := map[string]struct{}{}
		for _, param := range params {
			paramMap[param.BondType + "|" + param.HandlInst] = struct{}{}
		}


		for _, bondType := range bondTypes {
			missingHandleSpeed := []string{}
			for _, handleSpeed := range handleSpeeds {
				if _, ok := paramMap[bondType + "|" + handleSpeed]; !ok {
					missingHandleSpeed = append(missingHandleSpeed, getHandlInstName(handleSpeed))
				}
			}
			if len(missingHandleSpeed) > 0 {
				commissionCheckInfo2 += fmt.Sprintf("业务方案：%s 【%s】⊗【%s】\n", getPlanCode(planId), getBondTypeName(bondType), strings.Join(missingHandleSpeed, ","))
			}
		}
	}
	if len(commissionCheckInfo1) > 0 || len(commissionCheckInfo2) > 0 {
		if len(checkInfo) > 0 {
			checkInfo += "\n"
		}
		checkInfo += "【以下业务方案缺少佣金费率参数】\n"
		checkInfo += commissionCheckInfo1 + commissionCheckInfo2
	}


	// 利差参数检查
	interestCheckInfo := ""
	for _, planId := range planIds {
		paramMap := map[string]struct{}{}
		missingBondTypes := []string{}
		params, _ := interestParams.Data[planId]
		for _, param := range params {
			if param.LongShort == "LONG" {
				paramMap[param.BondType] = struct{}{}
			}
		}
		for _, bondType := range bondTypes {
			if _, ok := paramMap[bondType]; !ok {
				missingBondTypes = append(missingBondTypes, getBondTypeName(bondType))
			}
		}
		if len(missingBondTypes) > 0 {
			interestCheckInfo += fmt.Sprintf("业务方案：%s 【%s】\n", getPlanCode(planId), strings.Join(missingBondTypes, ","))
		}
	}
	if interestCheckInfo != "" {
		if len(checkInfo) > 0 {
			checkInfo += "\n"
		}
		checkInfo += "【以下业务方案缺少利差参数】\n"
		checkInfo += interestCheckInfo
	}

	return passed, checkInfo, nil
}

func Task_TrsCtptyData(positionServiceUrl, dataQryUrl, WebhookUrl string) error {
	passed, checkInfo, err := checkTrsCtptyData(dataQryUrl)
	if err != nil {
		fmt.Printf("checkTrsCtptyData error:%v\n", err)
		return err
	}

	if !passed {
		msg := checkInfo
		fmt.Printf("Run Task Task_TrsCtptyData not passed:\n%s", msg)
		err = wechat.SendToWeChat(WebhookUrl, msg)
		if err != nil {
			fmt.Printf("Run Task Task_TrsCtptyData SendToWeChat error: %v", err)
			return err
		}
		return nil
	}

	log.Println("Run Task Task_TrsCtptyData check passed!")
	return nil
}

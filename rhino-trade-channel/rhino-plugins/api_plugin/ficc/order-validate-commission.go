package ficc

import (
	"fmt"
	"rhino-plugins/api_plugin/util"
)

func (a *TitansFiccAPIAdapter) getComissionRate(planID interface{}, securityType, bondType, handlInst, side string) (commissionRate float64, ok bool) {
	valList, ok2, de := a.autoSyncRepo.Get("CommissionParam", fmt.Sprintf("%v-%v-%v-%v", planID, securityType, bondType, handlInst))
	if de != nil {
		return 
	}
	if !ok2 || len(valList) == 0 {
		return
	}
	commissionData := valList[len(valList)-1]
	if side == "1" { // 买单
		commissionRate, ok, _ = util.GetFloatValueInField(commissionData, "BuyCommissionRate")
	} else { // 卖单
		commissionRate, ok, _ = util.GetFloatValueInField(commissionData, "SellCommissionRate")
	}
	return
}

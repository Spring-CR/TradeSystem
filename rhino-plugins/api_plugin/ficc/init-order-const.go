package ficc

import (
	"errors"
	"fmt"
	"rhino-common/domain_error"
	"rhino-core/schema"
	"strconv"
	"strings"
)

func (a *TitansFiccAPIAdapter) initOrderConst(configMap map[string]*schema.ApplicationCfgItem) {
	configItem, ok := configMap["ExecUser"]
	if !ok || configItem.ConfigItemValue == "" {
		domain_error.ProcessSevereError(true, 5, nil, errors.New("ExecUser not config"), "ExecUser not config")
	}
	a.execUser = configItem.ConfigItemValue

	configItem, ok = configMap["TraderID"]
	if !ok || configItem.ConfigItemValue == "" {
		domain_error.ProcessSevereError(true, 5, nil, errors.New("TraderID not config"), "TraderID not config")
	}
	a.traderID = configItem.ConfigItemValue

	configItem, ok = configMap["InvestorID"]
	if !ok || configItem.ConfigItemValue == "" {
		domain_error.ProcessSevereError(true, 5, nil, errors.New("InvestorID not config"), "InvestorID not config")
	}
	strs := strings.Split(configItem.ConfigItemValue, ",")
	if len(strs) != 2 {
		domain_error.ProcessSevereError(true, 5, nil, fmt.Errorf("InvestorID is not correct:%s", configItem.ConfigItemValue), "InvestorID is not correct")
	}
	a.twoInvestorID[0] = strs[0]
	a.twoInvestorID[1] = strs[1]

	configItem, ok = configMap["InitMarginRatio"]
	if !ok || configItem.ConfigItemValue == "" {
		domain_error.ProcessSevereError(true, 5, nil, errors.New("InitMarginRatio not config"), "InitMarginRatio not config")
	}
	val, err := strconv.ParseFloat(configItem.ConfigItemValue, 64)
	if err != nil || val <= 0 {
		domain_error.ProcessSevereError(true, 5, nil, fmt.Errorf("InitMarginRatio is not correct:%s", configItem.ConfigItemValue), "InitMarginRatio is not correct")
	}
	a.initMarginRatio = val
}

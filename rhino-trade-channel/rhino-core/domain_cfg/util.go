package domain_cfg

import (
	"rhino-common/domain_error"
	"rhino-core/schema"
	"strings"
)

func checkCfgItems(cfgItemsToCheck []*schema.TradeChannelCfgItem, checkers[]*CfgItemChecker) (itemValues map[string]string, de *domain_error.Error) {
	itemValues = map[string]string{}
	mCfgItemsToCheck := map[string]*schema.TradeChannelCfgItem{}
	for _, cfgItem := range cfgItemsToCheck {
		mCfgItemsToCheck[cfgItem.ConfigItemName] = cfgItem
	}
	// 执行每个checker的检查逻辑
	for _, checker := range checkers {
		// 对于需要程序化设置的配置项，不进行检查
		if checker.SetByProgram > 0 {
			continue
		}
		cfgItemToCheck, ok := mCfgItemsToCheck[checker.ConfigItemName]
		value := ""
		if ok { // 如果入参有同名的配置项
			if checker.Required > 0 && checker.ConfigItemDefaultValue == "" && cfgItemToCheck.ConfigItemValue == "" {
				de = domain_error.Build(domain_error.REQUIRED_CONFIG_MISSING_ERR_CODE, nil, checker.ConfigItemName, checker.Description)
				return
			}
			if cfgItemToCheck.ConfigItemValue == "" {
				value = checker.ConfigItemDefaultValue
			} else {
				value = cfgItemToCheck.ConfigItemValue
			}
		} else {
			if checker.Required > 0 && checker.ConfigItemDefaultValue == "" {
				de = domain_error.Build(domain_error.REQUIRED_CONFIG_MISSING_ERR_CODE, nil, checker.ConfigItemName, checker.Description)
				return
			}
			value = checker.ConfigItemDefaultValue
		}
		itemValues[checker.ConfigItemName] = value
	}

	return
}


type HostAndPort struct {
	Host string
	Port string
}
func parseHostAndPortPaires(address string) (hostAndPortPaires[]*HostAndPort) {
	strs := strings.Split(address, ",")
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}
		hostAndPortStrs := strings.Split(str, ":")
		hostAndPort := &HostAndPort{Host:  strings.TrimSpace(hostAndPortStrs[0])}
		if hostAndPort.Host == "" {
			continue
		}
		if len(hostAndPortStrs) > 1 {
			hostAndPort.Port = strings.TrimSpace(hostAndPortStrs[1])
		}
		if hostAndPort.Port == "" {
			hostAndPort.Port  = "80"
		}
		hostAndPortPaires = append(hostAndPortPaires, hostAndPort)
	}
	return
}
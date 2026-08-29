package domain_cfg

import "rhino-core/schema"

func (c *ApplicationCfg) ConfigOrderFromExtendAttrItems(orderProperties map[string]interface{}, order *schema.TradeOrder) {
	if len(orderProperties) == 0 || len(c.extendAttrItems) == 0 {
		return
	}
}

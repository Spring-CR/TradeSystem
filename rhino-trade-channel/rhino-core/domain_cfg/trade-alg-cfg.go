package domain_cfg

import "rhino-core/schema"

type TradeAlgorithmCfg struct {
	tradeAlgorithm          *schema.TradeAlgorithm
	tradeAlgorithmAttrItems []*schema.TradeAlgorithm
}

func NewTradeAlgorithmCfg() {

}

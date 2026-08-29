package order_capital

import (
	"log"
	"rhino-common/domain_error"
	"rhino-core/adapter_registry"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
)

type CapitalCalculator struct {
	applicationCfg *domain_cfg.ApplicationCfg
	capitalAdapter CapitalAdapter
}

func NewCapitalCalculator(applicationCfg *domain_cfg.ApplicationCfg) (cc *CapitalCalculator, de *domain_error.Error) {
	var capitalAdapter CapitalAdapter
	adapterPath := applicationCfg.GetOrdCapitalAdapterPath()
	if adapterPath != "" {
		// 从注册表获取适配器的构造函数（目前，对于apiAdpater，是无参数函数，有参函数需要根据特殊情况来处理了）
		_capitalAdapter, de, err := adapter_registry.CallAdapterFunction(adapterPath, applicationCfg)
		if err != nil {
			domain_error.ProcessSevereError(true, 5, nil, err, "fail to construct CapitalAdapter from "+adapterPath)
		}
		if de != domain_error.NilDomainError {
			domain_error.ProcessSevereError(true, 5, nil, err, "fail to construct CapitalAdapter from "+adapterPath)
		}
		// 获取apiAdapter
		capitalAdapter = _capitalAdapter.(CapitalAdapter)
		log.Printf("finish get CapitalAdapter:%s\n", adapterPath)
	}

	cc = &CapitalCalculator{
		applicationCfg: applicationCfg,
		capitalAdapter: capitalAdapter,
	}

	return
}

func (cc *CapitalCalculator) AcquireOrderCapital(force bool, order *schema.TradeOrder) (freezedAmount float64, turnToReviewOnErr bool, de *domain_error.Error) {
	if cc.capitalAdapter.WouldOrderCapitalIncrease(order) {
		return
	}
	log.Printf("AcquireOrderCapital, AppOrdID:%s\n", order.AppOrdID)
	return cc.capitalAdapter.FreezeOrderCapitalOnTradeBegin(force, order)
}

func (cc *CapitalCalculator) RollbackFreezedOrderCapitalForTradeError(order *schema.TradeOrder) {
	if cc.capitalAdapter.WouldOrderCapitalIncrease(order) {
		return
	}
	cc.capitalAdapter.RollbackFreezedOrderCapitalForTradeError(order)
}
package domain_trade_instr

import (
	"rhino-common/context"
	"rhino-common/domain_error"
	"rhino-instr/store"
)

func FindTradeDetails() {
	//store.Find
}

func FindTradeDeskOrderIdsByParentKey(parentKey string) (ids []string, de *domain_error.Error) {
	var err error
	ids, err = store.FindTradeDeskOrderIdsByParentKey(context.DB, parentKey)
	if err != nil {
		de = domain_error.Build(domain_error.CANNOT_GET_TRADE_DESK_ORDER_ID_ERR_CODE, err, parentKey)
		return
	}
	return
}

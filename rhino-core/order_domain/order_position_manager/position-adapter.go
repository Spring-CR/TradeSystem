package order_position_manager

import (
	"rhino-common/domain_error"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/types"
)

type PositionAdapter interface {
	HasSufficientQuota(positionUnits []*PositionUnit, tradeOrder *schema.TradeOrder, metadata interface{}) (sufficient bool, de *domain_error.Error)
	CalculateOrderFreezeQuotaInPositionUnit(positionUnit *PositionUnit, tradeOrder *schema.TradeOrder, metadata interface{}) (freezeQuota float64, freezeQuotaBuf float64, ok bool)
	AfterFreezeQuotaInPositionUnit(positionUnit *PositionUnit, freezeQty float64, tradeOrder *schema.TradeOrder, metadata interface{}, quotaLocker *QuotaLocker, lastPositionUnit bool)
	CalculateIncreaseQuotaInPositionUnit(positionUnit *PositionUnit, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, metadata interface{}) (increaseQuota float64, ok bool)
	AfterIncreaseQuotaInPositionUnit(positionUnit *PositionUnit, increaseQty float64, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, metadata interface{}, quotaIncrementer *QuotaIncrementer, lastPositionUnit bool)
	CalculateUnfreezeQuotaInPositionUnit(positionUnit *PositionUnit, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, metadata interface{}, quotaLocker *QuotaLocker) (unfreezeQuota float64, ok bool)
	AfterUnfreezeQuotaInPositionUnit(positionUnit *PositionUnit, unfreezeQuota float64, tradeOrder *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, metadata interface{}, quotaLocker *QuotaLocker, lastPositionUnit bool, rollback bool)
	AfterUpdateQuota(tradeResp *types.TradeActionRespReturn, metadata interface{})
	GetPositionCalculatorKey(tradeOrder *schema.TradeOrder) (key string)
	GetPositionCalculatorConstructParams(tradeOrder *schema.TradeOrder) (param *PositionCalculatorConstructParam)
	LoadInitPositionRecords(applicationConfig *domain_cfg.ApplicationCfg) (params []*PositionCalculatorConstructParam)
	GeneralizePositionRecord(metadata interface{}, insert bool) map[string]interface{}
}

type PositionChangeEvent struct {
	InsertOrUpdate int // 0 - insert, 1 - update
	PositionData   map[string]interface{}
}

type PositionCalculatorConstructParam struct {
	Key           string
	Metadata      interface{}
	PositionUnits []*PositionUnit
}

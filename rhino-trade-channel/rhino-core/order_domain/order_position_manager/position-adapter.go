package order_position_manager

import (
	"rhino-common/domain_error"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/types"
	"rhino-data/datamap"
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
	GetQuotaNotEnoughHandler() func(tradeOrder *schema.TradeOrder, de *domain_error.Error)
	PreparePositionAdjustmentParams(tradeOrder *schema.TradeOrder) (mockTradeOrder *schema.TradeOrder, mockTradeActionResp *schema.TradeActionResp, de *domain_error.Error)
	GetPositionCalculatorKeysForPurgingTask(tradeOrdersToKeep []*schema.TradeOrder, tradeActionLatestRespsToKeep []*schema.TradeActionLatestResp, tradeActionRespsToKeep []*schema.TradeActionResp, tradeOrdersToArchive []*schema.TradeOrder, tradeActionLatestRespsToArchive []*schema.TradeActionLatestResp, tradeActionRespsToArchive []*schema.TradeActionResp, mockTradeOrders map[string][]*schema.TradeOrder, purgingLog *schema.DataPurgingLog) (positionCalculatorKeysToPurge []string)
	ReloadPositionRecordsForPurgingTask(applicationConfig *domain_cfg.ApplicationCfg, purgingLog *schema.DataPurgingLog) (params []*PositionCalculatorConstructParam)
	ProcessPositionBaseDataSyncEvent(event *datamap.DataChangeEvent)
	ParsePositionRecordFromReposRecord(extendAttrMap map[string]interface{}) (key string, positionRecord interface{})
	GetPositionCalculatorConstructParamForPositionRecord(positionRecord interface{}) (param *PositionCalculatorConstructParam)
	PreparePositionAdjustmentParamsByPositionBaseDiff(positionRecordBase, positionRecordCurr interface{}) (mockTradeOrder *schema.TradeOrder, mockTradeActionResp *schema.TradeActionResp, de *domain_error.Error)
	PrepareForRecover(positionRecord interface{})
	UpdatePositionBaseDynamically()(dynamicallyUpdate bool)
}

type PositionChangeEvent struct {
	TableName      string
	InsertOrUpdate int // 0 - insert, 1 - update, 2 - delete
	PositionData   map[string]interface{}
}

type PositionCalculatorConstructParam struct {
	Key           string
	Metadata      interface{}
	PositionUnits []*PositionUnit
}

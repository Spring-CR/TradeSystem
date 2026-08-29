package order_position

import (
	"rhino-common/utils/quota"
	"rhino-core/schema"
	"rhino-core/types"
	"rhino-data/datamap"
	"sync"
)

type PositionAdapter interface {
	// 返回多个position key，其中第一个是需要进行平仓校验
	GetPositionKeyList(order *schema.TradeOrder) (keyList []string, forceList []bool)
	WouldOrderPositionIncrease(order *schema.TradeOrder) bool
	WouldOrderPositionDecrease(order *schema.TradeOrder) bool
	IsOrderFinished(tradeActionResp *schema.TradeActionResp) bool
	IsTradeActionRespInBeginStatus(tradeActionResp *schema.TradeActionResp) bool
	GetQtyNotEnoughErrMsgPrefix(order *schema.TradeOrder) string
	GetQuotaMetadata(positionBaseRecord map[string]interface{}) map[string]interface{}
	ComputeReleaseQuota(quotaAcquire *quota.QuotaAcquire[*schema.TradeOrder]) float64
	GetOrderQuota(order *schema.TradeOrder) float64
	GetBaseQuantity(order *schema.TradeOrder, positionKey string) (baseQuota float64, positionBaseRecord map[string]interface{})
	CouldIncreaseQuota(order *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, positionKey string) bool
	UpdatePositionAfterAcquireOrderQuota(order *schema.TradeOrder, positionKey string, qc *quota.QuotaControl[*schema.TradeOrder, *schema.TradeActionResp], runInTradeEngine bool, newCreate bool, keyIndex int) (events []*PositionChangeEvent)
	AfterUpdatePositionByTradeResp(order *schema.TradeOrder, tradeResp *types.TradeActionRespReturn, positionKey string, qc *quota.QuotaControl[*schema.TradeOrder, *schema.TradeActionResp], isOrderPositionDecrease bool, runInTradeEngine bool, newCreate bool, keyIndex int) (events []*PositionChangeEvent)
	ProcessDataSyncEvent(event *datamap.DataChangeEvent, mapLock *sync.RWMutex, positionMap map[string]*quota.QuotaControl[*schema.TradeOrder, *schema.TradeActionResp]) (events []*PositionChangeEvent)
}

type PositionChangeEvent struct {
	InsertOrUpdate int // 0 - insert, 1 - update
	PositionData   map[string]interface{}
	InPositionMap  bool
}

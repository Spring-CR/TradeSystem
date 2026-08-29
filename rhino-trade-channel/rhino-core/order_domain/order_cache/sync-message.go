package order_cache

import (
	"rhino-core/schema"
	"rhino-core/types"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type OrderCacheSyncMessageType int

const (
	OrderCacheSyncMessageType_AddRootTradeOrder            = 1
	OrderCacheSyncMessageType_AddTradeActionForDirectOrder = 2
	OrderCacheSyncMessageType_UpdateByTradeActionResp      = 3
	OrderCacheSyncMessageType_UpdateTradeOrderDraft        = 4
	OrderCacheSyncMessageType_DeleteTradeOrderDraft        = 5
	OrderCacheSyncMessageType_Reset                        = 6
	OrderCacheSyncMessageType_SyncOrder                    = 7
	OrderCacheSyncMessageType_SyncTradeActionLatestResp    = 8
	OrderCacheSyncMessageType_UpdateTradeOrderAttributes   = 9
	OrderCacheSyncMessageType_AdjustPosition               = 10
)

type OrderCacheSyncMessage struct {
	MessageType                     OrderCacheSyncMessageType
	TradeOrder                      *schema.TradeOrder
	TradeActionLatestResp           *schema.TradeActionLatestResp
	TradeActionResp                 *schema.TradeActionResp
	UpdateAttributes                map[string]interface{}
	AppOrdID                        string
	AppOrdIDMap                     map[string]bool
	ClOrdIDMap                      map[string]bool
	ComposeOrdIDMap                 map[string]bool
	ActionKeyMap                    map[string]bool
	TradeOrdersToKeep               []*schema.TradeOrder
	TradeActionLatestRespsToKeep    []*schema.TradeActionLatestResp
	TradeActionRespsToKeep          []*schema.TradeActionResp
	TradeOrdersToArchive            []*schema.TradeOrder
	TradeActionLatestRespsToArchive []*schema.TradeActionLatestResp
	TradeActionRespsToArchive       []*schema.TradeActionResp
	PurgingLog                      *schema.DataPurgingLog
}

func (m *OrderCacheSyncMessage) GetMsgTime() int64 {
	switch m.MessageType {
	case OrderCacheSyncMessageType_AddRootTradeOrder:
		return m.TradeOrder.DBInsertTime
	case OrderCacheSyncMessageType_UpdateByTradeActionResp:
		return m.TradeActionResp.MsgTime
	case OrderCacheSyncMessageType_AddTradeActionForDirectOrder:
		return m.TradeActionLatestResp.ActionMsgTime
	}
	return 0
}

// // 按消息时间从小到大排序
// func sortOrderCacheSyncMessages(orderCacheSyncMessages []*OrderCacheSyncMessage) {
// 	if len(orderCacheSyncMessages) == 0 {
// 		return
// 	}
// 	sort.Slice(orderCacheSyncMessages, func(i, j int) bool {
// 		t1 := orderCacheSyncMessages[i].GetMsgTime()
// 		t2 := orderCacheSyncMessages[j].GetMsgTime()
// 		if t1 != t2 {
// 			return t1 < t2
// 		}
// 		return int(orderCacheSyncMessages[i].MessageType) < int(orderCacheSyncMessages[j].MessageType)
// 	})
// }

func newOrderCacheSyncMessageForAddRootTradeOrder(newOrder *types.TraceableTradeOrder) (jsData []byte, err error) {
	lock := newOrder.GetLock()
	lock.RLock()
	defer lock.RUnlock()

	msg := &OrderCacheSyncMessage{
		MessageType: OrderCacheSyncMessageType_AddRootTradeOrder,
		TradeOrder:  newOrder.GetBasicInfo(),
	}
	jsData, err = json.Marshal(msg)

	return
}

func newOrderCacheSyncMessageForUpdateByTradeActionResp(tradeActionResp *schema.TradeActionResp) (jsData []byte, err error) {

	msg := &OrderCacheSyncMessage{
		MessageType:     OrderCacheSyncMessageType_UpdateByTradeActionResp,
		TradeActionResp: tradeActionResp,
	}
	jsData, err = json.Marshal(msg)

	return
}

func newOrderCacheSyncMessageForAddTradeActionForDirectOrder(tradeAction *types.TraceableTradeActionResp) (jsData []byte, err error) {
	lock := tradeAction.GetLock()
	lock.RLock()
	defer lock.RUnlock()

	msg := &OrderCacheSyncMessage{
		MessageType:           OrderCacheSyncMessageType_AddTradeActionForDirectOrder,
		TradeActionLatestResp: tradeAction.GetTradeActionLatestResp(),
	}
	jsData, err = json.Marshal(msg)

	return
}

func newOrderCacheSyncMessageForUpdateTradeOrderDraft(tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder) (jsData []byte, err error) {

	msg := &OrderCacheSyncMessage{
		MessageType:           OrderCacheSyncMessageType_UpdateTradeOrderDraft,
		TradeOrder:            order,
		TradeActionLatestResp: tradeActionLatestResp,
	}
	jsData, err = json.Marshal(msg)

	return
}

func newOrderCacheSyncMessageForDeleteTradeOrderDraft(tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder) (jsData []byte, err error) {

	msg := &OrderCacheSyncMessage{
		MessageType:           OrderCacheSyncMessageType_DeleteTradeOrderDraft,
		TradeOrder:            order,
		TradeActionLatestResp: tradeActionLatestResp,
	}
	jsData, err = json.Marshal(msg)

	return
}

func newOrderCacheSyncMessageForUpdateTradeOrderAttributes(appOrdID string, updateAttrs map[string]interface{}) (jsData []byte, err error) {

	msg := &OrderCacheSyncMessage{
		MessageType:      OrderCacheSyncMessageType_UpdateTradeOrderAttributes,
		AppOrdID:         appOrdID,
		UpdateAttributes: updateAttrs,
	}
	jsData, err = json.Marshal(msg)

	return
}

func newResetMessage(appOrdIDMap, clOrdIDMap, composeOrdIDMap, actionKeyMap map[string]bool, tradeOrdersToKeep []*schema.TradeOrder, tradeActionLatestRespsToKeep []*schema.TradeActionLatestResp, tradeActionRespsToKeep []*schema.TradeActionResp,
	tadeOrdersToArchive []*schema.TradeOrder, tradeActionLatestRespsToArchive []*schema.TradeActionLatestResp, tradeActionRespsToArchive []*schema.TradeActionResp, purgingLog *schema.DataPurgingLog) (jsData []byte, err error) {
	msg := &OrderCacheSyncMessage{
		MessageType:                     OrderCacheSyncMessageType_Reset,
		AppOrdIDMap:                     appOrdIDMap,
		ClOrdIDMap:                      clOrdIDMap,
		ComposeOrdIDMap:                 composeOrdIDMap,
		ActionKeyMap:                    actionKeyMap,
		TradeOrdersToKeep:               tradeOrdersToKeep,
		TradeActionLatestRespsToKeep:    tradeActionLatestRespsToKeep,
		TradeActionRespsToKeep:          tradeActionRespsToKeep,
		TradeOrdersToArchive:            tadeOrdersToArchive,
		TradeActionLatestRespsToArchive: tradeActionLatestRespsToArchive,
		TradeActionRespsToArchive:       tradeActionRespsToArchive,
		PurgingLog:                      purgingLog,
	}
	jsData, err = json.Marshal(msg)
	return
}

func newOrderCacheSyncMessageForSyncOrder(tradeOrder *schema.TradeOrder) (jsData []byte, err error) {
	msg := &OrderCacheSyncMessage{
		MessageType: OrderCacheSyncMessageType_SyncOrder,
		TradeOrder:  tradeOrder,
	}
	jsData, err = json.Marshal(msg)
	return
}

func newOrderCacheSyncMessageForSyncTradeActionLatestResp(tradeActionLatestResp *schema.TradeActionLatestResp) (jsData []byte, err error) {
	msg := &OrderCacheSyncMessage{
		MessageType:           OrderCacheSyncMessageType_SyncTradeActionLatestResp,
		TradeActionLatestResp: tradeActionLatestResp,
	}
	jsData, err = json.Marshal(msg)
	return
}

func newOrderCacheSyncMessageForAdjustPosition(mockTradeOrder *schema.TradeOrder, mockTradeActionResp *schema.TradeActionResp) (jsData []byte, err error) {
	msg := &OrderCacheSyncMessage{
		MessageType:     OrderCacheSyncMessageType_AdjustPosition,
		TradeOrder:      mockTradeOrder,
		TradeActionResp: mockTradeActionResp,
	}
	jsData, err = json.Marshal(msg)
	return
}

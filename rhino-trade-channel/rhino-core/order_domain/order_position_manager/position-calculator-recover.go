package order_position_manager

import (
	"rhino-core/schema"
	"rhino-core/types"
)

type PositionCalculatorEventType int

const (
	PositionCalculatorEvent_FreezeQuota               = 0
	PositionCalculatorEvent_RollbackFreezeQuota       = 1
	PositionCalculatorEvent_UpdatePositionByTradeResp = 2
)

type PositionCalculatorEvent struct {
	eventType  PositionCalculatorEventType
	tradeOrder *schema.TradeOrder
	tradeResp  *types.TradeActionRespReturn
}

func (c *PositionCalculator) SetSuccessor(successor *PositionCalculator) {
	c.successor = successor
}

func (c *PositionCalculator) AddFreezeQuotaEvent(tradeOrder *schema.TradeOrder) {
	evt := &PositionCalculatorEvent{eventType: PositionCalculatorEvent_FreezeQuota, tradeOrder: tradeOrder}
	c.calEventList = append(c.calEventList, evt)
}

func (c *PositionCalculator) AddRollbackFreezeQuotaEvent(tradeOrder *schema.TradeOrder) {
	evt := &PositionCalculatorEvent{eventType: PositionCalculatorEvent_RollbackFreezeQuota, tradeOrder: tradeOrder}
	c.calEventList = append(c.calEventList, evt)
}

func (c *PositionCalculator) AddUpdatePositionByTradeRespEvent(tradeResp *types.TradeActionRespReturn) {
	evt := &PositionCalculatorEvent{eventType: PositionCalculatorEvent_UpdatePositionByTradeResp, tradeResp: tradeResp}
	c.calEventList = append(c.calEventList, evt)
}

func (successor *PositionCalculator) RecoverFromParent(parent *PositionCalculator) {
	
	for _, evt := range parent.calEventList {
		switch evt.eventType {
		case PositionCalculatorEvent_FreezeQuota:
			successor.FreezeQuota(true, evt.tradeOrder)
		case PositionCalculatorEvent_RollbackFreezeQuota:
			successor.RollbackFreezeQuota(evt.tradeOrder)
		case PositionCalculatorEvent_UpdatePositionByTradeResp:
			successor.UpdatePositionByTradeResp(evt.tradeResp)
		}
	}
}

package order_position_manager

import (
	"log"
	"rhino-common/domain_error"
	"rhino-common/utils/logger"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	"rhino-core/types"
	"sync"
	"time"
)

type PositionCalculator struct {
	key             string          // 通常按交易对手(产品)+标的维度来组织
	lock            *sync.RWMutex   // 锁
	metadata        interface{}     // 元数据。由插件自定义
	positionUnits   []*PositionUnit // 持仓单元列表
	positionAdapter PositionAdapter // 持仓适配器
	orderLog        *logger.OrderLog
	persisFunc      func(evt *PositionChangeEvent)
	calEventList    []*PositionCalculatorEvent
	successor       *PositionCalculator // 发生replay时会被设置
}

func NewPositionCalculator(key string, tradeOrder *schema.TradeOrder, positionAdapter PositionAdapter, orderLog *logger.OrderLog, persisFunc func(evt *PositionChangeEvent)) *PositionCalculator {
	param := positionAdapter.GetPositionCalculatorConstructParams(tradeOrder)
	inst := &PositionCalculator{
		key:             key,
		lock:            &sync.RWMutex{},
		metadata:        param.Metadata,
		positionUnits:   param.PositionUnits,
		positionAdapter: positionAdapter,
		orderLog:        orderLog,
		persisFunc:      persisFunc,
	}
	return inst
}

func NewPositionCalculatorByConstructParam(param *PositionCalculatorConstructParam, positionAdapter PositionAdapter, orderLog *logger.OrderLog, persisFunc func(evt *PositionChangeEvent)) *PositionCalculator {
	inst := &PositionCalculator{
		key:             param.Key,
		lock:            &sync.RWMutex{},
		metadata:        param.Metadata,
		positionUnits:   param.PositionUnits,
		positionAdapter: positionAdapter,
		orderLog:        orderLog,
		persisFunc:      persisFunc,
	}
	return inst
}

func (c *PositionCalculator) HasSufficientQuota(tradeOrder *schema.TradeOrder) (sufficient bool, de *domain_error.Error) {

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.successor != nil {
		return c.HasSufficientQuota(tradeOrder)
	}

	sufficient, de = c.positionAdapter.HasSufficientQuota(c.positionUnits, tradeOrder, c.metadata)

	return
}

func (c *PositionCalculator) FreezeQuota(force bool, tradeOrder *schema.TradeOrder) (sufficient bool, de *domain_error.Error) {

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.successor != nil {
		return c.FreezeQuota(force, tradeOrder)
	}

	if !force {
		sufficient, de = c.positionAdapter.HasSufficientQuota(c.positionUnits, tradeOrder, c.metadata)
		if de != nil {
			return
		}
	}
	
	if tradeOrder.QuotaValidateTime == 0 {
		tradeOrder.QuotaValidateTime = timeutil.ConvertTimeToMicroseconds(time.Now())
	}

	n := len(c.positionUnits)
	for i, positionUnit := range c.positionUnits {
		positionUnit.freezeQuota(c.positionAdapter, c.metadata, tradeOrder, i == n-1)
	}

	if c.persisFunc != nil {
		c.persisFunc(&PositionChangeEvent{InsertOrUpdate: 0, PositionData: c.positionAdapter.GeneralizePositionRecord(c.metadata, true)})
	}

	c.AddFreezeQuotaEvent(tradeOrder)

	return
}

func (c *PositionCalculator) RollbackFreezeQuota(tradeOrder *schema.TradeOrder) {

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.successor != nil {
		c.successor.RollbackFreezeQuota(tradeOrder)
		return
	}

	n := len(c.positionUnits)
	for i, positionUnit := range c.positionUnits {
		positionUnit.rollbackFreezeQuota(c.positionAdapter, c.metadata, tradeOrder, i == n-1)
	}

	if c.persisFunc != nil {
		c.persisFunc(&PositionChangeEvent{InsertOrUpdate: 1, PositionData: c.positionAdapter.GeneralizePositionRecord(c.metadata, false)})
	}

	c.AddRollbackFreezeQuotaEvent(tradeOrder)
}

func (c *PositionCalculator) UpdatePositionByTradeResp(tradeResp *types.TradeActionRespReturn) {
	
	log.Println("start to UpdatePositionByTradeResp")

	c.lock.Lock()
	defer c.lock.Unlock()

	log.Println("enter UpdatePositionByTradeResp")

	if c.successor != nil {
		log.Println("successor of UpdatePositionByTradeResp")
		c.successor.UpdatePositionByTradeResp(tradeResp)
		return
	}

	c.increaseQuotaByFillReport(tradeResp)
	c.unfreezeQuotaByTradeResp(tradeResp)

	log.Println("before AfterUpdateQuota")

	c.positionAdapter.AfterUpdateQuota(tradeResp, c.metadata)

	log.Println("after AfterUpdateQuota")

	if c.persisFunc != nil {
		c.persisFunc(&PositionChangeEvent{InsertOrUpdate: 1, PositionData: c.positionAdapter.GeneralizePositionRecord(c.metadata, false)})
	}

	log.Println("before UpdatePositionByTradeResp")

	c.AddUpdatePositionByTradeRespEvent(tradeResp)

	log.Println("after UpdatePositionByTradeResp")
}

func (c *PositionCalculator) increaseQuotaByFillReport(tradeResp *types.TradeActionRespReturn) {
	tradeOrder := tradeResp.GetTradeOrder()
	tradeActionResp := tradeResp.CurrentTradeActionResp
	n := len(c.positionUnits)
	for i, positionUnit := range c.positionUnits {
		positionUnit.increaseQuota(c.positionAdapter, c.metadata, tradeOrder, tradeActionResp, i == n-1)
	}
}

func (c *PositionCalculator) unfreezeQuotaByTradeResp(tradeResp *types.TradeActionRespReturn) {
	tradeOrder := tradeResp.GetTradeOrder()
	tradeActionResp := tradeResp.CurrentTradeActionResp
	n := len(c.positionUnits)
	for i, positionUnit := range c.positionUnits {
		positionUnit.unfreezeQuota(c.positionAdapter, c.metadata, tradeOrder, tradeActionResp, i == n-1)
	}
}

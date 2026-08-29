// 同步逻辑，作为master角色
package order_cache

import (
	"darkpool-common/bean"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"rhino-common/domain_error"
	"rhino-common/utils/kafka"
	"rhino-core/schema"
	"rhino-core/types"

	"github.com/IBM/sarama"
)

func (c *OrderCache) configForMaster() {

	msgLogPath := filepath.Join(c.applicationCfg.GetWorkingDir(), "order_cace")
	os.MkdirAll(msgLogPath, 0755)
	msgLogPath = filepath.Join(msgLogPath, "order_status_msg.log")

	system, business := c.applicationCfg.GetSystemAndBusinessCodes()
	sendingTopic := fmt.Sprintf("%s-%s-order-status", system, business)

	brokers := c.applicationCfg.GetKafkaBrokers()

	producer, err := kafka.NewTenaciousProducerWithBuffer(msgLogPath, true, 4096*2, sendingTopic, brokers, 1024)
	if err != nil {
		domain_error.ProcessSevereError(true, 3, nil, err, "fail to init kafka producer client in OrderCache")
	}

	consumer, err := kafka.NewTenaciousConsumerWithBuffer(fmt.Sprintf("%s-%s-order-status-ack", system, business), brokers, 4096, sarama.OffsetNewest, c.onOrderStatusAckMessageReceived)
	if err != nil {
		domain_error.ProcessSevereError(true, 3, nil, err, "fail to init kafka consumer client in OrderCache")
	}

	c.producer = producer
	c.consumer = consumer

	log.Println("start strictRecover")
	// 数据恢复
	c.strictRecover()

	log.Println("finish configForMaster")
}

func (c *OrderCache) onOrderStatusAckMessageReceived(data []byte) {

}

// Todo 还需要补充关于order draft update、order draft delete的相关handler
func (c *OrderCache) afterAddRootTradeOrder(order *types.TraceableTradeOrder) {

	if !c.master {
		c.initPositionForNewOrder(order.GetBasicInfo())
	}

	if !c.master {
		if c._afterAddRootTradeOrder != nil {
			c._afterAddRootTradeOrder(order)
		}
		return
	}

	jsData, _ := newOrderCacheSyncMessageForAddRootTradeOrder(order)
	log.Printf("======>afterAddRootTradeOrder send data:%s\n", jsData)
	c.producer.SendMessage(jsData)
}

func (c *OrderCache) afterUpdateByTradeActionResp(tradeActionResp *schema.TradeActionResp, traceableTradeActionResp *types.TraceableTradeActionResp) {
	if !c.master {
		if c._afterUpdateByTradeActionResp != nil {
			c._afterUpdateByTradeActionResp(tradeActionResp, traceableTradeActionResp)
		}
		return
	}

	jsData, _ := newOrderCacheSyncMessageForUpdateByTradeActionResp(tradeActionResp)
	c.producer.SendMessage(jsData)
}

func (c *OrderCache) afterAddTradeActionForDirectOrder(tradeAction *types.TraceableTradeActionResp) {
	if !c.master {
		if c._afterAddTradeActionForDirectOrder != nil {
			c._afterAddTradeActionForDirectOrder(tradeAction)
		}
		return
	}

	jsData, _ := newOrderCacheSyncMessageForAddTradeActionForDirectOrder(tradeAction)
	c.producer.SendMessage(jsData)
}

func (c *OrderCache) afterUpdateTradeOrderDraft(tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder) {
	if !c.master {
		if c._afterUpdateTradeOrderDraft != nil {
			c._afterUpdateTradeOrderDraft(tradeActionLatestResp, order)
		}
		return
	}

	jsData, _ := newOrderCacheSyncMessageForUpdateTradeOrderDraft(tradeActionLatestResp, order)
	c.producer.SendMessage(jsData)
}

func (c *OrderCache) afterDeleteTradeOrderDraft(tradeActionLatestResp *schema.TradeActionLatestResp, order *schema.TradeOrder) {
	if !c.master {
		if c._afterDeleteTradeOrderDraft != nil {
			c._afterDeleteTradeOrderDraft(tradeActionLatestResp, order)
		}
		return
	}

	jsData, _ := newOrderCacheSyncMessageForDeleteTradeOrderDraft(tradeActionLatestResp, order)
	c.producer.SendMessage(jsData)
}

func (c *OrderCache) afterUpdateTradeOrderAttributes(appOrdID string, updateAttrs map[string]interface{}) {
	if !c.master {
		if c._afterUpdateTradeOrderAttributes != nil {
			c._afterUpdateTradeOrderAttributes(appOrdID, updateAttrs)
		}
		return
	}

	jsData, _ := newOrderCacheSyncMessageForUpdateTradeOrderAttributes(appOrdID, updateAttrs)
	c.producer.SendMessage(jsData)
}

func (c *OrderCache) SyncOrder(tradeOrder *schema.TradeOrder) {
	if c.IsMaster() {
		jsData, _ := newOrderCacheSyncMessageForSyncOrder(tradeOrder)
		log.Printf("SyncOrder, send message:%s\n", jsData)
		c.producer.SendMessage(jsData)
	} else {
		order, ok := c.GetOrderByAppOrdID(tradeOrder.AppOrdID)
		if !ok {
			domain_error.ProcessSevereError(false, 0, nil, fmt.Errorf("fail to SyncOrder as appOrdID %s not found in cache", tradeOrder.AppOrdID), "error in SyncOrder")
			return
		}

		originalOrder := order.GetBasicInfo()
		
		if tradeOrder.OrdStatusUpdateTime < originalOrder.OrdStatusUpdateTime {
			return
		}

		err := bean.Copy(tradeOrder).To(originalOrder)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "error in SyncOrder")
		}

		if c._afterSyncOrder != nil {
			c._afterSyncOrder(originalOrder)
		}
	}
}

func (c *OrderCache) SyncTradeActionLatestResp(tradeActionLatestResp *schema.TradeActionLatestResp) {
	if c.IsMaster() {
		jsData, _ := newOrderCacheSyncMessageForSyncTradeActionLatestResp(tradeActionLatestResp)
		c.producer.SendMessage(jsData)
	} else {
		c.tradeActionRespMapLock.Lock()
		defer c.tradeActionRespMapLock.Unlock()

		tradeActionResp, ok := c.tradeActionRespMap[tradeActionLatestResp.ClOrdID]
		if !ok {
			domain_error.ProcessSevereError(false, 0, nil, fmt.Errorf("fail to SyncTradeActionLatestResp as clOrdID %s not found in cache", tradeActionLatestResp.ClOrdID), "error in SyncTradeActionLatestResp")
			return
		}

		originalTradeActionLatestResp := tradeActionResp.GetTradeActionLatestResp()

		if tradeActionLatestResp.MsgTime < originalTradeActionLatestResp.MsgTime {
			return
		}

		err := bean.Copy(tradeActionLatestResp).To(originalTradeActionLatestResp)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "error in SyncTradeActionLatestResp")
		}

		if c._afterSyncTradeActionLatestResp != nil {
			c._afterSyncTradeActionLatestResp(originalTradeActionLatestResp)
		}
	}
}

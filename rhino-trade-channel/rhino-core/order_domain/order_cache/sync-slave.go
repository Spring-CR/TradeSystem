// 同步逻辑，作为slave角色
package order_cache

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"rhino-common/domain_error"
	"rhino-common/utils/kafka"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

func (c *OrderCache) getHistoricalOrderCacheSyncMessages() (newestOffest, oldestOffset int64, syncMsgs []*OrderCacheSyncMessage, err error) {

	system, business := c.applicationCfg.GetSystemAndBusinessCodes()
	brokers := c.applicationCfg.GetKafkaBrokers()
	receivingTopic := fmt.Sprintf("%s-%s-order-status", system, business)

	beginTime := time.Now()
	log.Printf("Start to GetHistoricMessages, brokers:%v,  receivingTopic:%s\n", brokers, receivingTopic)
	// 获取历史数据
	newestOffest, oldestOffset, hisDatas, err1 := kafka.GetHistoricMessages(brokers, receivingTopic)
	if err != nil {
		domain_error.ProcessSevereError(true, 3, nil, err, "fail to get historic data for OrderCache")
		err = err1
		return
	}

	log.Printf("kafka msg size:%d, cost:%v\n", len(hisDatas), time.Since(beginTime))

	sep := []byte("\n")
	for _, hisData := range hisDatas {
		if len(hisData) == 0 {
			continue
		}
		// Todo
		// 对于不用buffer的情况
		// syncMsg:=&OrderCacheSyncMessage{}
		// err = json.Unmarshal(hisData, syncMsg)
		// if err != nil {
		// 	domain_error.ProcessSevereError(true, 3, nil, err, fmt.Sprintf("fail to unmarshal OrderCacheSyncMessage for OrderCache, data=%s", hisData))
		// }
		// syncMsgs = append(syncMsgs, syncMsg)

		var data []byte
		for len(hisData) > 0 {

			end := bytes.Index(hisData, sep)
			if end < 0 {
				data = hisData
			} else {
				data = hisData[:end]
			}

			if len(data) > 0 {

				log.Printf("===> sync raw data:%s\n", data)

				syncMsg := &OrderCacheSyncMessage{}
				err = json.Unmarshal(data, syncMsg)
				if err != nil {
					domain_error.ProcessSevereError(true, 3, nil, err, fmt.Sprintf("fail to unmarshal OrderCacheSyncMessage for OrderCache, data=%s", data))
				}
				syncMsgs = append(syncMsgs, syncMsg)
			}

			if end < 0 {
				break
			} else {
				hisData = hisData[end+1:]
			}
		}
		//log.Printf("======> len(syncMsgs)=%d\n", len(syncMsgs))
	}

	log.Printf("Parse data time cost:%v\n", time.Since(beginTime))
	// 如果检测到有reset的syncMsg，需要截断
	n := len(syncMsgs)
	if n == 0 {
		return
	}

	// 因为20260803之后涉及多轮归档了，不能只判断到OrderCacheSyncMessageType_Reset就直接截断
	// i := n - 1
	// for ; i >= 0; i-- {
	// 	if syncMsgs[i].MessageType == OrderCacheSyncMessageType_Reset {
	// 		break
	// 	}
	// }
	// syncMsgs = syncMsgs[i+1:]
	appOrdIDMap := make(map[string]bool)
	clOrdIDMap := make(map[string]bool)
	composeOrdIDMap := make(map[string]bool)
	actionKeyMap := make(map[string]bool)
	for _, syncMsg := range syncMsgs {
		if syncMsg.MessageType != OrderCacheSyncMessageType_Reset {
			continue
		}
		for k, v := range syncMsg.AppOrdIDMap {
			appOrdIDMap[k] = v
		}
		for k, v := range syncMsg.ClOrdIDMap {
			clOrdIDMap[k] = v
		}
		for k, v := range syncMsg.ComposeOrdIDMap {
			composeOrdIDMap[k] = v
		}
		for k, v := range syncMsg.ActionKeyMap {
			actionKeyMap[k] = v
		}
	}

	var newSyncMsgs []*OrderCacheSyncMessage
	for _, syncMsg := range syncMsgs {
		switch syncMsg.MessageType {
		case OrderCacheSyncMessageType_AddRootTradeOrder:
			tradeOrder := syncMsg.TradeOrder
			if tradeOrder == nil {
				continue
			}
			if appOrdIDMap[tradeOrder.AppOrdID] || clOrdIDMap[tradeOrder.ClOrdID] {
				continue
			}

		case OrderCacheSyncMessageType_AddTradeActionForDirectOrder:
			tradeActionLatestResp := syncMsg.TradeActionLatestResp
			if tradeActionLatestResp == nil {
				continue
			}
			if actionKeyMap[tradeActionLatestResp.ActionKey] {
				continue
			}

		case OrderCacheSyncMessageType_UpdateByTradeActionResp: //其实最重要是需要这个消息类型
			tradeActionResp :=syncMsg.TradeActionResp
			if tradeActionResp == nil {
				continue
			}
			if clOrdIDMap[tradeActionResp.ClOrdID] || clOrdIDMap[tradeActionResp.OrigClOrdID] {
				continue
			}

		case OrderCacheSyncMessageType_UpdateTradeOrderDraft:
		case OrderCacheSyncMessageType_DeleteTradeOrderDraft:
		case OrderCacheSyncMessageType_Reset:
			continue
		case OrderCacheSyncMessageType_SyncOrder:
		case OrderCacheSyncMessageType_SyncTradeActionLatestResp:
		case OrderCacheSyncMessageType_UpdateTradeOrderAttributes:
			continue
		}

		newSyncMsgs = append(newSyncMsgs, syncMsg)
	}

	syncMsgs = newSyncMsgs

	return
}

func (c *OrderCache) configForSlave() {

	msgLogPath := filepath.Join(c.applicationCfg.GetWorkingDir(), "order_cace")
	os.MkdirAll(msgLogPath, 0755)
	msgLogPath = filepath.Join(msgLogPath, "order_status_ack_msg.log")

	system, business := c.applicationCfg.GetSystemAndBusinessCodes()
	sendingTopic := fmt.Sprintf("%s-%s-order-status-ack", system, business)
	receivingTopic := fmt.Sprintf("%s-%s-order-status", system, business)

	brokers := c.applicationCfg.GetKafkaBrokers()

	producer, err := kafka.NewTenaciousProducerWithBuffer(msgLogPath, true, 256, sendingTopic, brokers, 1024)
	if err != nil {
		domain_error.ProcessSevereError(true, 3, nil, err, "fail to init kafka producer client in OrderCache")
	}

	c.producer = producer

	// 防止重复收到tradeActionLatestResp消息
	c.tradeActionLatestRespKeyMap = make(map[string]bool)
	c.tradeActionLatestRespKeyMapLock = &sync.RWMutex{}
	c.rootOrderKeyMap = make(map[string]bool)

	//newestOffest, oldestOffset, syncMsgs, _ := c.getHistoricalOrderCacheSyncMessages()
	/*newestOffest, oldestOffset, syncMsgs := c.strictRecover()

	beginTime := time.Now()
	// 从历史数据中恢复
	if len(syncMsgs) > 0 {
		c.recoverFromSyncMessages(syncMsgs)
	}

	log.Printf("recoverFromSyncMessages time cost:%v\n", time.Since(beginTime))*/
	newestOffest, oldestOffset := c.recoverForSlave()

	// 初始化数据
	log.Printf("receivingTopic:%s\n", receivingTopic)
	offset := sarama.OffsetOldest
	if newestOffest > oldestOffset {
		// 重复传一个消息
		offset = newestOffest - 1
	}
	consumer, err := kafka.NewTenaciousConsumerWithBuffer(receivingTopic, c.applicationCfg.GetKafkaBrokers(), 1000, offset, c.onOrderStatusMessageReceived)
	if err != nil {
		domain_error.ProcessSevereError(true, 3, nil, err, "fail to init kafka consumer client in OrderCache")
	}
	c.consumer = consumer
	c.consumer.Start()
}

// 对于master，因为在其完全启动前，数据库不会再插入新记录，因此，strictRecover对于master而言是充分的。
// 对于slave，因为不能确保在加载数据库记录执行恢复时，是否还有新的记录持续插入，因此仅仅使用strictRecover是不充分的。
func (c *OrderCache) recoverForSlave() (newestOffest, oldestOffset int64) {
	newestOffest, oldestOffset, syncMsgs := c.strictRecover()
	// 对于slave，strictRecover并不完全充分
	beginTime := time.Now()
	// 从历史数据中恢复
	if len(syncMsgs) > 0 {
		c.recoverFromSyncMessages(syncMsgs)
	}
	log.Printf("recoverFromSyncMessages time cost:%v\n", time.Since(beginTime))
	return
}

func (c *OrderCache) onOrderStatusMessageReceived(data []byte) {

	if len(data) == 0 {
		return
	}
	syncMsg := &OrderCacheSyncMessage{}
	err := json.Unmarshal(data, syncMsg)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to unmarshal OrderCacheSyncMessage for OrderCache, data=%s", data))
		return
	}

	c.updateByOrderCacheSyncMessage(syncMsg)
}

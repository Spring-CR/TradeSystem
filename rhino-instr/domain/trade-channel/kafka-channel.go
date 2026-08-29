package trade_channel

import (
	"encoding/json"
	"log"
	"math"
	"rhino-common/context"

	"fmt"
	"rhino-common/domain_error"
	"rhino-instr/schema"
	"rhino-instr/store"
	"time"

	"github.com/Shopify/sarama"
)

var(
	DefaultTradeChannel *KafkaTradeChannel
)

type KafkaTradeChannel struct {
	tradeReqTopic  string
	tradeRespTopic string
	brokers        []string
	producer       sarama.SyncProducer
	consumer       sarama.PartitionConsumer
}

// NewKafkaTradeChannel creates a new KafkaTradeChannel instance
func NewKafkaTradeChannel(tradeReqTopic, tradeRespTopic string, brokers []string) (*KafkaTradeChannel, *domain_error.Error) {
	inst := &KafkaTradeChannel{tradeReqTopic: tradeReqTopic, tradeRespTopic: tradeRespTopic, brokers: brokers}
	de := inst.initProducer()
	if de != nil {
		return nil, de
	}
	log.Println("kafka initProducer...")
	de = inst.initConsumer()
	log.Println("kafka initConsumer...")
	return inst, de
}

func (c *KafkaTradeChannel) initProducer() *domain_error.Error {
	// 配置 Kafka 生产者
	config := sarama.NewConfig()

	basicKafkaConfig(config)

	// 同步需要配置
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second

	// 创建同步生产者
	producer, err := sarama.NewSyncProducer(c.brokers, config)
	if err != nil {
		return domain_error.Build(domain_error.CANNOT_CREATE_PRODUCER_ERR_CODE, fmt.Errorf("failed to start Kafka producer: %v", err))
	}
	c.producer = producer
	return nil
}

func (c *KafkaTradeChannel) PublishTradeInstr(tradeInstr *schema.TradeInstr) (int64, *domain_error.Error) {

	jsonData, _ := json.MarshalIndent(tradeInstr, "", "  ")

	/// 构建要发送的 Kafka 消息
	kafkaMessage := &sarama.ProducerMessage{
		Topic: c.tradeReqTopic,
		Key:   sarama.StringEncoder("key"),
		Value: sarama.ByteEncoder(jsonData),
	}

	log.Printf("交易台下单报文：%s\n", jsonData)

	// 发送消息
	_, offset, err := c.producer.SendMessage(kafkaMessage)
	if err != nil {
		defer func() {
			c.producer.Close()
			// re-init
			for i := 0; i < 5; i++ {
				de := c.initProducer()
				if de == nil {
					return
				}
				domain_error.ReportIfErrorHappen(de)
				time.Sleep(5 * time.Second)
			}
		}()
		return 0, domain_error.Build(domain_error.CANNOT_PUBLISH_TRADE_MSG_ERR_CODE, err)
	}

	return offset, nil
}

func (c *KafkaTradeChannel) initConsumer() *domain_error.Error {

	var maxOffset int64
	var err error
	maxOffset, err = store.GetMaxKafkaOffset(context.DB)
	for err != nil {
		err = fmt.Errorf("fail to GetMaxKafkaOffset, error:%+v ", err)
		log.Println(err.Error())
		de := domain_error.Build(domain_error.DATABASE_OPERATION_ERR_CODE, err)
		domain_error.ReportIfErrorHappen(de)
		time.Sleep(5 * time.Second)
		maxOffset, err = store.GetMaxKafkaOffset(context.DB)
	}

	// Configuration for the Sarama client
	config := sarama.NewConfig()

	basicKafkaConfig(config)

	config.Consumer.Return.Errors = true // 启用错误通道返回
	config.Consumer.Retry.Backoff = 2 * time.Second
	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // 因consumer.ConsumePartition已经设置了偏移量，这里可以省略
	config.Consumer.Offsets.AutoCommit.Enable = false // 禁用自动提交
	//config.Consumer.Offsets.CommitInterval = 0  // 禁用自动提交

	client, err := sarama.NewClient(c.brokers, config)
	if err != nil {
		return domain_error.Build(domain_error.KAFKA_ERR_CODE, err)
	}

	// 获取当前的start和end偏移量
	// 获取起始偏移量
	partition := int32(0)
    startOffset, err := client.GetOffset(c.tradeRespTopic, partition, sarama.OffsetOldest)
    if err != nil {
        return domain_error.Build(domain_error.KAFKA_ERR_CODE, err)
    }
    // 获取结束偏移量
    endOffset, err := client.GetOffset(c.tradeRespTopic, partition, sarama.OffsetNewest)
    if err != nil {
        return domain_error.Build(domain_error.KAFKA_ERR_CODE, err)
    }

	log.Printf("maxOffset:%d, startOffset:%d, endOffset:%d\n", maxOffset, startOffset, endOffset)
	if maxOffset + 1 > endOffset {
		maxOffset = endOffset - 1
		log.Printf("reset1 maxOffset=%d\n", maxOffset)
	} else if maxOffset + 1 < startOffset {
		maxOffset = startOffset - 1
		log.Printf("reset2 maxOffset=%d\n", maxOffset)
	}

	// 创建消费者组
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return domain_error.Build(domain_error.CANNOT_CREATE_CONSUMER_ERR_CODE, err)
	}
	log.Printf("c.tradeRespTopic:%s, maxOffset:%d\n", c.tradeRespTopic, maxOffset)

	if maxOffset == 0 {
        partitionConsumer, err := consumer.ConsumePartition(c.tradeRespTopic, 0, sarama.OffsetOldest)
        if err != nil {
            return domain_error.Build(domain_error.CANNOT_CREATE_CONSUMER_ERR_CODE, err)
        }
        log.Println("ConsumePartition from the oldest offset...")
        c.consumer = partitionConsumer
    } else {
        partitionConsumer, err := consumer.ConsumePartition(c.tradeRespTopic, 0, maxOffset+1)
        if err != nil {
            return domain_error.Build(domain_error.CANNOT_CREATE_CONSUMER_ERR_CODE, err)
        }
        log.Println("ConsumePartition from the specified offset...")
        c.consumer = partitionConsumer
    }

	return nil
}

func basicKafkaConfig(config *sarama.Config) {
	config.Metadata.Retry.Max = math.MaxInt
	config.Metadata.Retry.Backoff = 2 * time.Second
	config.Net.DialTimeout = 30 * time.Second
	config.Net.ReadTimeout = 30 * time.Second
	config.Net.WriteTimeout = 30 * time.Second
	config.Net.MaxOpenRequests = 16
}

func (c *KafkaTradeChannel) KeepListening(onReceived func(m *sarama.ConsumerMessage) bool) {
	log.Println("KeepListening...")
	go func() {
		for {
			select {
			case msg := <-c.consumer.Messages():
				log.Printf("Consumed message: %s, offset: %d\n", string(msg.Value), msg.Offset)
				// 出错时重试10次
				processed := false
				for i := 0; i < context.RetryIntervalTimes; i++ {
					ok := onReceived(msg)
					if ok {
						log.Printf("finish process message, offset:%d\n", msg.Offset)
						processed = true
						break
					}
					time.Sleep( time.Duration(context.RetryIntervalSeconds) * time.Second)
				}
				if !processed {
                    log.Printf("Failed to process message after %d retries, offset: %d\n", context.RetryIntervalTimes, msg.Offset)
                }
			case err := <-c.consumer.Errors():
				log.Printf("Consumer error: %s", err)
			}
		}
	}()
}

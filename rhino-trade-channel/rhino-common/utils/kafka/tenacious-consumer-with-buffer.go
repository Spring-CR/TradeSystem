package kafka

import (
	"bytes"
	"log"
	"rhino-common/domain_error"
	"time"

	"github.com/IBM/sarama"
)

type TenaciousConsumerWithBuffer struct {
	consumeTopic   string
	brokers        []string
	consumeBuffer  int
	consumeOffset  int64
	consumer       sarama.PartitionConsumer
	onDataReceived func(data []byte)
}

func NewTenaciousConsumerWithBuffer(consumeTopic string, brokers []string, consumeBuffer int, consumeOffset int64, onDataReceived func(data []byte)) (*TenaciousConsumerWithBuffer, error) {
	inst := &TenaciousConsumerWithBuffer{
		consumeTopic:   consumeTopic,
		brokers:        brokers,
		consumeBuffer:  consumeBuffer,
		consumeOffset:  consumeOffset,
		onDataReceived: onDataReceived,
	}
	err := inst.resetConsumer()
	if err != nil {
		return nil, err
	}
	return inst, err
}

func (c *TenaciousConsumerWithBuffer) resetConsumer() error {
/*
	// 配置 Kafka 生产者
	config := sarama.NewConfig()

	// 设置最大重试次数，确保生产者在 Kafka broker 不可用时会重试，但不会无限重试
	config.Metadata.Retry.Max = 2 // 设置最大重试次数为 3 次，可以根据实际情况调整

	// 设置元数据重试的间隔时间，减少重试频率，以免产生过多的请求
	config.Metadata.Retry.Backoff = 200 * time.Millisecond // 500 秒的重试间隔

	// 设置网络连接超时，快速检测 broker 的不可用情况
	config.Net.DialTimeout = 2 * time.Second // 连接超时为 2 秒

	// 设置读取超时，防止在读取响应时等待过长时间
	config.Net.ReadTimeout = 2 * time.Second // 读取超时为 5 秒

	// 设置写入超时，防止在写入数据时等待过长时间
	config.Net.WriteTimeout = 2 * time.Second // 写入超时为 5 秒

	// 设置每个请求的最大并发请求数，避免并发请求过多导致性能瓶颈
	config.Net.MaxOpenRequests = 16 // 最多允许 5 个并发请求

	//	config.Consumer.Return.Errors = true // 启用错误通道返回
	config.Consumer.Retry.Backoff = 200 * time.Millisecond

	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // 取最老的消息
	config.Consumer.Offsets.AutoCommit.Enable = false     // 禁用自动提交, 如果要取最新的消息，则不能禁用自动提交
	config.Consumer.Offsets.CommitInterval = 0            // 禁用自动提交
*/
	config := getConsumerConfig()
	client, err := sarama.NewClient(c.brokers, config)
	if err != nil {
		return err
	}

	// 创建消费者组
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return err
	}

	partitionConsumer, err := consumer.ConsumePartition(c.consumeTopic, 0, c.consumeOffset)
	if err != nil {
		return err
	}
	log.Printf("ConsumePartition from offset:%v\n", c.consumeOffset)
	c.consumer = partitionConsumer

	return nil
}

func (c *TenaciousConsumerWithBuffer) Start() {
	ch := make(chan []byte, c.consumeBuffer)
	go func() {
		for {
			select {
			case msg := <-c.consumer.Messages():
				if msg != nil {
					ch <- msg.Value
				} else {
					log.Printf("======> kafka message is nil, topic=%s\n", c.consumeTopic)
					time.Sleep(2 * time.Second)
					err1 := c.resetConsumer()
					errCount:=0
					for err1 != nil {
						errCount++
						if errCount > 120 || err1 == sarama.ErrOffsetOutOfRange {
							c.consumeOffset = sarama.OffsetOldest
						}
						log.Printf("fail to resetConsumer, err=%v\n", err1)
						time.Sleep(2 * time.Second)
						err1 = c.resetConsumer()
					}
				}
			case err := <-c.consumer.Errors():
				log.Printf("Consumer error: %s", err)
				domain_error.ProcessSevereError(false, 0, nil, err, "error while comsuming kafka message")
				time.Sleep(2 * time.Second)
				err1 := c.resetConsumer()
				errCount:=0
				for err1 != nil {
					errCount++
					if errCount > 120 || err1 == sarama.ErrOffsetOutOfRange {
						c.consumeOffset = sarama.OffsetOldest
					}
					log.Printf("fail to resetConsumer, err=%v\n", err1)
					time.Sleep(2 * time.Second)
					err1 = c.resetConsumer()
				}
			}
		}
	}()
	go func() {
		splitor := []byte("\n")
		for {
			data := <-ch
			lines := bytes.Split(data, splitor)
			for _, line := range lines {
				c.onDataReceived(line)
			}
		}
	}()
}

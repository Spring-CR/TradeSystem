package kafka

import (
	"log"
	"rhino-common/domain_error"
	"time"

	"github.com/IBM/sarama"
)

type SimpleConsumer struct {
	consumeTopic   string
	brokers        []string
	consumeOffset  int64
	consumer       sarama.PartitionConsumer
	onDataReceived func(data []byte, msgSeq int64, msgTime time.Time)
}

func NewSimpleConsumer(consumeTopic string, brokers []string, consumeOffset int64, onDataReceived func(data []byte, msgSeq int64, msgTime time.Time)) (*SimpleConsumer, error) {
	inst := &SimpleConsumer{
		consumeTopic:   consumeTopic,
		brokers:        brokers,
		consumeOffset:  consumeOffset,
		onDataReceived: onDataReceived,
	}
	err := inst.resetConsumer()
	if err != nil {
		return nil, err
	}
	return inst, err
}

func (c *SimpleConsumer) resetConsumer() error {

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

	c.consumer = partitionConsumer

	return nil
}

func (c *SimpleConsumer) Start() {
	
	go func() {
		for {
			select {
			case msg, ok := <-c.consumer.Messages():
				if !ok {
					domain_error.ProcessSevereError(false, 0, nil, nil, "error while comsuming kafka message")
					time.Sleep(time.Second)
					continue
				}

				c.onDataReceived(msg.Value, msg.Offset, msg.Timestamp)
				
			case err := <-c.consumer.Errors():
				log.Printf("Consumer error: %s", err)
				domain_error.ProcessSevereError(false, 0, nil, err, "error while comsuming kafka message")
				time.Sleep(2 * time.Second)
			}
		}
	}()
}

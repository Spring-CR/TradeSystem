package kafka

import (
	"github.com/IBM/sarama"
)

type SimpleProducer struct {
	producer           sarama.SyncProducer
	sendingTopic       string
	brokers            []string
}

func NewSimpleProducer(sendingTopic string, brokers []string) (*SimpleProducer, error) {
	inst := &SimpleProducer{
		sendingTopic:       sendingTopic,
		brokers:            brokers,
	}
	err := inst.resetProducer()
	if err != nil {
		return nil, err
	}

	return inst, err
}

func (k *SimpleProducer) resetProducer() error {
	config := getProducerConfig()
	var err error
	// 创建同步生产者
	k.producer, err = sarama.NewSyncProducer(k.brokers, config)
	if err != nil {
		return err
	}
	return err
}

func (k *SimpleProducer) SendMessage(data []byte) error{
	kafkaMessage := &sarama.ProducerMessage{
		Topic: k.sendingTopic,
		Key:   sarama.StringEncoder("key"),
		Value: sarama.ByteEncoder(data),
	}
	_, _, err := k.producer.SendMessage(kafkaMessage)
	return err
}

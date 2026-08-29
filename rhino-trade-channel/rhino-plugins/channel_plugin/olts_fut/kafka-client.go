package olts_fut

import (
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/utils/kafka"
	"time"

	json "github.com/bytedance/sonic"
)


type kafkaClient struct {
	brokers      []string
	productTopic string
	consumeTopic string
	producer     *kafka.SimpleProducer
	consumer     *kafka.SimpleConsumer
	producerDFD  *kafka.SimpleProducer
}

func newKafkaClient(brokers []string, productTopic string, consumeTopic string, consumeOffset int64, onDataReceived func(data []byte, msgSeq int64, msgTime time.Time)) *kafkaClient {
	producer, err := kafka.NewSimpleProducer(productTopic, brokers)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, fmt.Sprintf("fail to NewSimpleProducer,  topic:%s, brokers:%v\n", productTopic, brokers))
	}
	consumer, err := kafka.NewSimpleConsumer(consumeTopic, brokers, consumeOffset, onDataReceived)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, fmt.Sprintf("fail to NewSimpleConsumer,  topic:%s, consumeOffset:%d, brokers:%v\n", consumeTopic, consumeOffset, brokers))
	}
	producerDFD, err := kafka.NewSimpleProducer(consumeTopic, brokers)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, fmt.Sprintf("fail to NewSimpleProducer,  topic:%s, brokers:%v\n", productTopic, brokers))
	}
	inst := &kafkaClient{
		brokers:      brokers,
		productTopic: productTopic,
		consumeTopic: consumeTopic,
		producer:     producer,
		consumer:     consumer,
		producerDFD:  producerDFD,
	}
	consumer.Start()
	return inst
}

func (c *kafkaClient) send(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to send trading request for JSON marshal error:%v", err))
		return err
	}

	log.Printf("===> Send trading request:%s\n", data)
	err = c.producer.SendMessage(data)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to send trading request error:%v, request:%s", err, data))
		return err
	}

	return err
}

func (c *kafkaClient) sendDFD(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to send DFD for JSON marshal error:%v", err))
		return err
	}

	log.Printf("===> Send DFD response:%s\n", data)
	err = c.producerDFD.SendMessage(data)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to send DFD error:%v, request:%s", err, data))
		return err
	}

	return err
}

func (c *kafkaClient) reset() {
	kafka.PurgeMessageForTopics(c.brokers, []string{c.productTopic, c.consumeTopic})
}

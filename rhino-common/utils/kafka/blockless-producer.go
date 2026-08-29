package kafka

import (
	"log"
	"rhino-common/domain_error"
	"time"

	"github.com/IBM/sarama"
)

// 无阻塞kafka生产者
// 用同步producer来实现无阻塞效果，因为用于消息恢复，对kafka的错误需要立即感知
// 当消息缓存满、或者kafka连接异常时，将把最新消息缓存到临时变量，并尝试通过ticker持续写入
type BlocklessProducer struct {
	producer sarama.SyncProducer
	brokers  []string
	input    chan *sarama.ProducerMessage
	retryMsg *sarama.ProducerMessage
}

func NewBlocklessProducer(brokers []string, inputBuf int) (*BlocklessProducer, error) {
	inst := &BlocklessProducer {
		brokers: brokers,
		input: make(chan *sarama.ProducerMessage, inputBuf),
	}
	err := inst.resetProducer()
	if err != nil {
		return nil, err
	}
	inst.doSendMessage()
	return inst, err
}

func (k *BlocklessProducer) resetProducer() error {

	// 配置 Kafka 生产者
	config := sarama.NewConfig()

	// 设置最大重试次数，确保生产者在 Kafka broker 不可用时会重试，但不会无限重试
	config.Metadata.Retry.Max = 2  // 设置最大重试次数为 3 次，可以根据实际情况调整

	// 设置元数据重试的间隔时间，减少重试频率，以免产生过多的请求
	config.Metadata.Retry.Backoff = 200 * time.Millisecond  // 1 秒的重试间隔

	// 设置网络连接超时，快速检测 broker 的不可用情况
	config.Net.DialTimeout = 2 * time.Second  // 连接超时为 2 秒

	// 设置读取超时，防止在读取响应时等待过长时间
	config.Net.ReadTimeout = 2 * time.Second  // 读取超时为 5 秒

	// 设置写入超时，防止在写入数据时等待过长时间
	config.Net.WriteTimeout = 2 * time.Second  // 写入超时为 5 秒

	// 设置每个请求的最大并发请求数，避免并发请求过多导致性能瓶颈
	config.Net.MaxOpenRequests = 16  // 最多允许 5 个并发请求

	// 设置消息发送的重试策略
	config.Producer.Retry.Max = 2  // 设置最大重试次数为 3 次（重试次数可根据业务需求调整）
	config.Producer.Retry.Backoff = 200 * time.Millisecond  // 设置重试的间隔时间为 500 毫秒

	// 设置生产者对消息的响应类型。这里使用 WaitForAll 表示等待所有副本的确认
	//config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.RequiredAcks = sarama.NoResponse

	// 设置生产者的消息批处理配置，控制每次发送消息的最大数量和最大时间间隔
	//config.Producer.Flush.Messages = 1000  // 每次发送最多 1000 条消息
	//config.Producer.Flush.Frequency = 100 * time.Millisecond  // 每 100 毫秒刷新一次

	// 同步需要配置
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 2 * time.Second

	var err error
	// 创建同步生产者
	k.producer, err = sarama.NewSyncProducer(k.brokers, config)
	if err != nil {
		return err
	}

	log.Println("finish resetProducer")
	return err
}

func (k *BlocklessProducer) SendMessage(topic string, data []byte) {
	/// 构建要发送的 Kafka 消息
	kafkaMessage := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder("key"),
		Value: sarama.ByteEncoder(data),
	}
	k.sendMessage(kafkaMessage)
}

func (k *BlocklessProducer) sendMessage(kafkaMessage *sarama.ProducerMessage) {
	select {
	case k.input <- kafkaMessage:
	default:
		//log.Println(" >>>>>> send too fast!")
		k.retryMsg = kafkaMessage
	}
}

func (k *BlocklessProducer) doSendMessage()  {
	t := time.NewTicker(time.Second)
	go func() {
		for{
			select {
			case msg := <-k.input:
				_, _, err := k.producer.SendMessage(msg)
				domain_error.ProcessSevereError(false, 0, nil, err, "fail to send kafka message")
				if err == nil {
					k.retryMsg = nil
					//log.Printf("send msg success!")
				} else {
					k.retryMsg = msg
					//log.Printf("send msg failed!")
				}
			case <-t.C:
				if k.retryMsg!=nil {
					k.sendMessage(k.retryMsg)
				}
			}
		}
	}()
}
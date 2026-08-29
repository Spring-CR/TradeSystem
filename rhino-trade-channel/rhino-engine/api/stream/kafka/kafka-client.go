package kafka

import (
	"errors"
	"log"
	"math"
	"os"
	"path/filepath"
	"rhino-common/domain_error"
	"rhino-common/utils/kafka"
	"rhino-core/domain_cfg"
	"rhino-core/store/app_store"
	"rhino-engine/api/stream"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

// 注意：在运维上需要求 Kafka Topic只能清空数据，不能删除重建
// 如果topic删除重建，要重启trade-engine

type KafkaClient struct {
	lock           *sync.Mutex
	cfg            *domain_cfg.ApplicationCfg
	noIngress      bool // 还没有api请求入内。 加入了noIngress之后，可以确保在开始期间，即时kafka topic被强制重建，请求信息可以继续过来；否则，会被minOffset拦截
	minOffset      int64
	brokers        []string
	reqTopic       string
	respTopic      string
	producer       *kafka.TenaciousProducerWithBuffer
	consumer       sarama.PartitionConsumer
	msgChan        <-chan *stream.IngressMessage
}

// 1、首次重启，从tradeOrder表取最大的MsgSeq作为minOffset
// 2、每次得到一次消息，需要把消息的offset设置为minOffset
// 3、当发现消息的offset小于等于记录的minOffset时，要忽略该消息
func NewKafkaClient(cfg *domain_cfg.ApplicationCfg) *KafkaClient {
	brokers := cfg.GetKafkaBrokers()
	systemCode, businessCode := cfg.GetSystemAndBusinessCodes()
	inst := &KafkaClient{
		lock:      &sync.Mutex{},
		cfg:       cfg,
		noIngress: true,
		brokers:   brokers,
		reqTopic:  systemCode + "-" + businessCode + "-req",
		respTopic: systemCode + "-" + businessCode + "-resp",
		minOffset: -1,
	}

	// 初始化，第一次需要从db中得到最大的消息offset
	msgSeq, err := app_store.GetMaxMsgSeqOfTradeOrder(cfg.GetAppDB(), systemCode, businessCode)
	domain_error.ProcessSevereError(true, 2, nil, err, "fail to get GetMaxMsgSeqOfTradeOrder")

	log.Printf("the max msgSeq for create kafka client: %d\n", msgSeq)

	// 查看topic的最新和最久的offset
	newestOffset, oldestOffset, err := kafka.GetNewestAndOldestOffset(cfg.GetKafkaBrokers(), inst.reqTopic)
	domain_error.ProcessSevereError(true, 2, nil, err, "fail to GetNewestAndOldestOffset for "+inst.reqTopic)
	if msgSeq >= int(newestOffset) {
		log.Printf("检测到数据库保存的最大msgSeq=%d 大于等于Topic(%s)的newestOffset=%d\n", msgSeq, inst.reqTopic, newestOffset)
		msgSeq = int(oldestOffset - 1)
		log.Printf("调整之后的msgSeq=%d\n", msgSeq)
		// 因系统能保证下单是幂等的，因此即时重复消费数据也不会影响业务正确性
	}

	inst.minOffset = int64(msgSeq)
	err = inst.resetKafka()
	log.Printf("finish resetKafka, inst.minOffset: %d", inst.minOffset)
	domain_error.ProcessSevereError(true, 2, nil, err, "fail to resetKafka")

	inst.initProducer()

	return inst
}

func (k *KafkaClient) resetKafka() error {
	// err := k.resetProducer()
	// if err != nil {
	// 	return err
	// }
	err := k.resetConsumer()
	if err != nil {
		return err
	}
	return err
}

func basicKafkaConfig(config *sarama.Config) {
	config.Metadata.Retry.Max = math.MaxInt // Todo，有争议，后面看看要不要设置为2
	config.Metadata.Retry.Backoff = 200 * time.Millisecond
	config.Net.DialTimeout = 2 * time.Second
	config.Net.ReadTimeout = 2 * time.Second
	config.Net.WriteTimeout = 2 * time.Second
	config.Net.MaxOpenRequests = 1 // 不需要并发
}

func (k *KafkaClient) resetConsumer() error {

	var err error

	// Configuration for the Sarama client
	config := sarama.NewConfig()

	basicKafkaConfig(config)

	config.Consumer.Return.Errors = true // 启用错误通道返回
	config.Consumer.Retry.Backoff = 2 * time.Second
	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // 从最老的消息开始重传，然后通过minOffset过滤，这样操作最安全，不会引起bug
	config.Consumer.Offsets.AutoCommit.Enable = false     // 禁用自动提交
	//config.Consumer.Offsets.CommitInterval = 0  // 禁用自动提交

	client, err := sarama.NewClient(k.brokers, config)
	if err != nil {
		return err
	}

	// 创建消费者组
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return err
	}

	partitionConsumer, err := consumer.ConsumePartition(k.reqTopic, 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}
	log.Println("ConsumePartition from the oldest offset...")
	k.consumer = partitionConsumer

	log.Println("finish resetConsumer")

	return nil
}

func (k *KafkaClient) PrepareIngressMessageChannel(buffer int) (msgChan <-chan *stream.IngressMessage) {

	k.lock.Lock()
	defer k.lock.Unlock()

	if k.msgChan != nil {
		return k.msgChan
	}

	ch := make(chan *stream.IngressMessage, buffer)
	k.msgChan = ch
	msgChan = ch

	go func() {
		for {
			select {
			case msg := <-k.consumer.Messages():

				if msg == nil {
					domain_error.ProcessSevereError(false, 0, nil, errors.New("kafka message is nil"), "error in PrepareIngressMessageChannel")
					time.Sleep(2 * time.Second)
					continue
				}

				// 从最老的消息开始重传，然后通过minOffset过滤，这样操作最安全，不会引起bug
				if k.noIngress && msg.Offset <= k.minOffset {
					//log.Printf("======>received old message with offset: %d\n", msg.Offset)
					continue
				}

				//log.Printf("=====> received message with offset: %d at: %d msg:%s\n", msg.Offset, timeutil.ConvertTimeToMilliseconds(time.Now()), msg.Value)

				ingressMsg := &stream.IngressMessage{
					Data:    msg.Value,
					MsgTime: msg.Timestamp,
					MsgSeq:  msg.Offset,
				}
				ch <- ingressMsg

				if k.noIngress {
					k.noIngress = false
				}

			case err := <-k.consumer.Errors():
				log.Printf("Consumer error: %s", err)
				domain_error.ProcessSevereError(false, 0, nil, err, "error while comsuming kafka message")
				time.Sleep(3 * time.Second)
			}
		}
	}()

	return
}

func (k *KafkaClient) SendMessage(data []byte) (msgSeq int64, err error) {

	k.producer.SendMessage(data)

	return
}

func (k *KafkaClient) initProducer() {

	messageLogFilePath := filepath.Join(k.cfg.GetWorkingDir(), "stream_api")
	os.MkdirAll(messageLogFilePath, 0755)
	messageLogFilePath = filepath.Join(messageLogFilePath, "trade_resp.log")

	var err error
	k.producer, err = kafka.NewTenaciousProducerWithBuffer(messageLogFilePath, true, 0, k.respTopic, k.cfg.GetKafkaBrokers(), 0)
	if err != nil {
		domain_error.ProcessSevereError(true, 3, nil, err, "fail to initProducer")
	}
}

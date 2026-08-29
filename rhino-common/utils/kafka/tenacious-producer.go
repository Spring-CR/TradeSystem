package kafka

import (
	"bytes"
	"log"
	"os"
	"rhino-common/domain_error"
	"time"

	"github.com/IBM/sarama"
	"github.com/hpcloud/tail"
)

// 健壮的kafka生产者
// 消息先写入文件系统，再通过inotify事件读取消息并发送到kafka，kafka侧则采用最严格的数据确认（ACK）
// 当kafka broker出现异常，producer会一直重试直到broker恢复
// 限制：当前消息体仅能支持json，并且需要加入换行符
type TenaciousProducer struct {
	producer           sarama.SyncProducer
	messageLogFilePath string
	writeBufferBytes   int
	sendingTopic       string
	brokers            []string
	input              chan []byte
	restSignal         chan bool
	//flushSignal        chan bool
	bytesBuffer        *bytes.Buffer
	messageLogFile     *os.File
	tail               *tail.Tail
}

func NewTenaciousProducer(messageLogFilePath string, resetMessageLogOnStart bool, writeBufferBytes int, sendingTopic string, brokers []string, inputBuf int) (*TenaciousProducer, error) {
	inst := &TenaciousProducer{
		messageLogFilePath: messageLogFilePath,
		writeBufferBytes:   writeBufferBytes,
		sendingTopic:       sendingTopic,
		brokers:            brokers,
		input:              make(chan []byte, inputBuf),
		bytesBuffer:        bytes.NewBuffer(nil),
		restSignal:         make(chan bool),
		//flushSignal:        make(chan bool),
	}
	err := inst.resetProducer()
	if err != nil {
		return nil, err
	}

	var messageLogFile *os.File
	if resetMessageLogOnStart {
		messageLogFile, err = os.OpenFile(messageLogFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return nil, err
		}
	} else {
		messageLogFile, err = os.OpenFile(messageLogFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
	}
	inst.messageLogFile = messageLogFile

	tail, err := tail.TailFile(messageLogFilePath, tail.Config{
		Follow: true, // 跟踪文件变化
		ReOpen: true, // 文件被轮换时重新打开
	})
	if err != nil {
		return nil, err
	}
	inst.tail = tail

	inst.doSendMessage()
	return inst, err
}

func (k *TenaciousProducer) resetProducer() error {
/*
	// 配置 Kafka 生产者
	config := sarama.NewConfig()

	// 设置最大重试次数，确保生产者在 Kafka broker 不可用时会重试，但不会无限重试
	config.Metadata.Retry.Max = 2 // 设置最大重试次数为 3 次，可以根据实际情况调整

	// 设置元数据重试的间隔时间，减少重试频率，以免产生过多的请求
	config.Metadata.Retry.Backoff = 200 * time.Millisecond // 1 秒的重试间隔

	// 设置网络连接超时，快速检测 broker 的不可用情况
	config.Net.DialTimeout = 2 * time.Second // 连接超时为 2 秒

	// 设置读取超时，防止在读取响应时等待过长时间
	config.Net.ReadTimeout = 2 * time.Second // 读取超时为 5 秒

	// 设置写入超时，防止在写入数据时等待过长时间
	config.Net.WriteTimeout = 2 * time.Second // 写入超时为 5 秒

	// 设置每个请求的最大并发请求数，避免并发请求过多导致性能瓶颈
	//config.Net.MaxOpenRequests = 16 / 16 // 最多允许 5 个并发请求
	config.Net.MaxOpenRequests = 1 // 根本就不需要并发

	// 设置消息发送的重试策略
	config.Producer.Retry.Max = 2                          // 设置最大重试次数为 3 次（重试次数可根据业务需求调整）
	config.Producer.Retry.Backoff = 200 * time.Millisecond // 设置重试的间隔时间为 500 毫秒

	// 设置生产者对消息的响应类型。这里使用 WaitForAll 表示等待所有副本的确认
	//config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.RequiredAcks = sarama.WaitForLocal
	//config.Producer.RequiredAcks = sarama.NoResponse

	// 设置生产者的消息批处理配置，控制每次发送消息的最大数量和最大时间间隔
	// 采用默认的批处理模式，在出错的时候，无法进行连续的重试
	//config.Producer.Flush.Messages = 1000  // 每次发送最多 1000 条消息
	//config.Producer.Flush.Frequency = 100 * time.Millisecond  // 每 100 毫秒刷新一次

	// 同步需要配置
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 2 * time.Second
*/
	config := getProducerConfig()
	var err error
	// 创建同步生产者
	k.producer, err = sarama.NewSyncProducer(k.brokers, config)
	if err != nil {
		return err
	}

	log.Println("finish resetProducer")
	return err
}

func (k *TenaciousProducer) SendMessage(data []byte) {
	k.input <- data
}

func (k *TenaciousProducer) reset() {
	//k.flushSignal <- true
	time.Sleep(5*time.Second)
	k.flush()
}

func (k *TenaciousProducer) doSendMessage() {

	// go func(){
	// 	for {
	// 		n := len(k.input)
	// 		log.Printf("Input chan len=%d\n", n)
	// 		time.Sleep(5*time.Second)
	// 	}
	// }()

	go func() {
		t := time.NewTicker(time.Second)
		for {
			select {
			case msg := <-k.input:
				if len(msg) == 0 {
					continue
				}
				k.bytesBuffer.Write(msg)
				if msg[len(msg)-1] != '\n' {
					k.bytesBuffer.WriteByte('\n')
				}
				if k.bytesBuffer.Len() >= k.writeBufferBytes {
					//k.flushSignal <- true
					k.flush()
				}
			case <-t.C:
				k.flush()
			case <-k.restSignal:
				k.reset()
			}
		}
	}()

	// go func() {
	// 	for {
	// 		select {
	// 		case <-k.flushSignal:
	// 			k.flush()
	// 		case <-t.C:
	// 			k.flush()
	// 		}
	// 	}
	// }()

	go func() {
		for {
			select {
			case line := <-k.tail.Lines:
				k.kafkaMessageSend(line.Text)
			}
		}
	}()
}

func (k *TenaciousProducer) kafkaMessageSend(data string) {
	for {
		kafkaMessage := &sarama.ProducerMessage{
			Topic: k.sendingTopic,
			Key:   sarama.StringEncoder("key"),
			Value: sarama.StringEncoder(data),
		}
		_, _, err := k.producer.SendMessage(kafkaMessage)
		domain_error.ProcessSevereError(false, 0, nil, err, "TenaciousProducer :: fail to send kafka message")
		if err == nil {
			return
		} else {
			time.Sleep(time.Second)
		}
	}
}

func (k *TenaciousProducer) flush() {
	if k.bytesBuffer.Len() == 0 {
		return
	}
	for {
		data := k.bytesBuffer.Bytes()
		n, err := k.messageLogFile.Write(data)

		if n < len(data) {
			log.Printf("write Len:%d < data-Len:%d\n", n, len(data))
		}

		domain_error.ProcessSevereError(false, 0, nil, err, "TenaciousProducer :: fail to flush message")
		if err == nil {
			k.bytesBuffer.Reset()
			return
		} else {
			time.Sleep(time.Second)
		}
	}
}

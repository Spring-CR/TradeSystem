package kafka

import (
	"fmt"
	"log"
	"rhino-common/domain_error"
	"time"

	"github.com/IBM/sarama"
)

func GetNewestMessage(brokers []string, topic string) (msg []byte, offset int64, exist bool, err error) {

	// 配置 Kafka 生产者
	config := sarama.NewConfig()
	// 连接到 Kafka 集群
	var client sarama.Client
	client, err = sarama.NewClient(brokers, config)
	if err != nil {
		return
	}
	defer func() {
		if client != nil {
			client.Close()
		}
	}()

	// 创建消费者
	var consumer sarama.Consumer
	consumer, err = sarama.NewConsumerFromClient(client)
	if err != nil {
		return
	}
	defer func() {
		if consumer != nil {
			consumer.Close()
		}
	}()

	// 获取该主题的唯一分区
	partition := int32(0) // 在rhino项目中，kafka都是只有一个分区
	// 获取该分区的最新偏移量
	offset, err = client.GetOffset(topic, partition, sarama.OffsetNewest)
	if err != nil {
		return
	}

	// 获取该分区的最早偏移量
	var oldestOffset int64
	oldestOffset, err = client.GetOffset(topic, partition, sarama.OffsetOldest)
	if err != nil {
		return
	}

	log.Printf("OffsetNewest:%d, OffsetOldest:%d\n", offset, oldestOffset)

	if offset == oldestOffset {
		exist = false
		log.Println("no data")
		return
	}

	offset = offset - 1

	// 如果偏移量小于等于0，说明主题可能为空或已被清空
	if offset <= 0 {
		// 主题为空或消息已删除，选择从最早的偏移量开始
		offset = sarama.OffsetOldest
	}

	// 获取最新的消息偏移量前一条的消息（即最新一条消息）
	var consumerPartition sarama.PartitionConsumer
	consumerPartition, err = consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		return
	}
	defer func() {
		consumerPartition.Close()
	}()

	// 读取消息
	t := time.NewTimer(3 * time.Second)
	select {
	case kafkaMsg := <-consumerPartition.Messages():
		// 成功获取最新消息
		fmt.Printf("Received latest kafka message: %s\n", string(kafkaMsg.Value))
		msg = kafkaMsg.Value
		exist = true
	case e := <-consumerPartition.Errors():
		// 如果出错了，打印错误信息
		fmt.Printf("Error occurred while consuming the message, error:%v\n", e.Err)
		err = e.Err
		return
	case <-t.C:
		exist = false
		return
	}

	return
}

func GetMessage(brokers []string, topic string, offset int64) (msg []byte, msgTime time.Time, exist bool, err error) {

	// 配置 Kafka 消费者
	config := sarama.NewConfig()
	// 连接到 Kafka 集群
	var client sarama.Client
	client, err = sarama.NewClient(brokers, config)
	if err != nil {
		return
	}
	defer func() {
		if client != nil {
			client.Close()
		}
	}()

	// 创建消费者
	var consumer sarama.Consumer
	consumer, err = sarama.NewConsumerFromClient(client)
	if err != nil {
		return
	}
	defer func() {
		if consumer != nil {
			consumer.Close()
		}
	}()

	// 获取该主题的唯一分区
	partition := int32(0) // 在rhino项目中，kafka都是只有一个分区

	// 获取该分区的最新偏移量
	newestOffset, err := client.GetOffset(topic, partition, sarama.OffsetNewest)
	if err != nil {
		return
	}

	// 获取该分区的最早偏移量
	oldestOffset, err := client.GetOffset(topic, partition, sarama.OffsetOldest)
	if err != nil {
		return
	}

	// 以下3种情况将无数据：1、newestOffset和oldestOffset相等；2、检索偏移量小于最老偏移量；3、检索偏移量大于等于最大偏移量（kafka的newestOffset其实是下一条消息的偏移偏移量，因此是不能相等的）
	if newestOffset == oldestOffset || offset < oldestOffset || offset >= newestOffset {
		exist = false
		log.Println("no data")
		return
	}

	// 获取最新的消息偏移量前一条的消息（即最新一条消息）
	var consumerPartition sarama.PartitionConsumer
	consumerPartition, err = consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		return
	}
	defer func() {
		consumerPartition.Close()
	}()

	// 读取消息
	t := time.NewTimer(3 * time.Second)
	select {
	case kMsg := <-consumerPartition.Messages():
		// 成功获取最新消息
		//fmt.Printf("Received kafka message: %s\n", string(msg.Value))

		if kMsg == nil {
			domain_error.ProcessSevereError(false, 0, nil, fmt.Errorf("kafka msg is nil"), "error in GetMessage")
			exist = false
			return
		}

		msg = kMsg.Value
		msgTime = kMsg.Timestamp
		exist = true
	case e := <-consumerPartition.Errors():
		// 如果出错了，打印错误信息
		fmt.Printf("Error occurred while consuming the message, error:%v\n", e.Err)
		err = e.Err
		return
	case <-t.C:
		log.Println("no data")
		exist = false
		return
	}

	return
}

func GetNewestAndOldestOffset(brokers []string, topic string) (newestOffset int64, oldestOffset int64, err error) {

	// 配置 Kafka 生产者
	config := sarama.NewConfig()
	// 连接到 Kafka 集群
	var client sarama.Client
	client, err = sarama.NewClient(brokers, config)
	if err != nil {
		return
	}
	defer func() {
		if client != nil {
			client.Close()
		}
	}()

	// 创建消费者
	var consumer sarama.Consumer
	consumer, err = sarama.NewConsumerFromClient(client)
	if err != nil {
		return
	}
	defer func() {
		if consumer != nil {
			consumer.Close()
		}
	}()

	// 获取该主题的唯一分区
	partition := int32(0) // 在rhino项目中，kafka都是只有一个分区

	// 获取该分区的最新偏移量
	newestOffset, err = client.GetOffset(topic, partition, sarama.OffsetNewest)
	if err != nil {
		return
	}

	// 获取该分区的最早偏移量
	oldestOffset, err = client.GetOffset(topic, partition, sarama.OffsetOldest)
	if err != nil {
		return
	}

	return
}

func GetHistoricMessages(brokers []string, topic string) (newestOffset, oldestOffset int64, messages [][]byte, err error) {

	newestOffset, oldestOffset, err = GetNewestAndOldestOffset(brokers, topic)
	if err != nil {
		return
	}

	if newestOffset == oldestOffset {
		return
	}

	config := sarama.NewConfig()
	// 连接到 Kafka 集群
	var client sarama.Client
	client, err = sarama.NewClient(brokers, config)
	if err != nil {
		return
	}
	defer func() {
		if client != nil {
			client.Close()
		}
	}()

	// 创建消费者
	var consumer sarama.Consumer
	consumer, err = sarama.NewConsumerFromClient(client)
	if err != nil {
		return
	}
	defer func() {
		if consumer != nil {
			consumer.Close()
		}
	}()

	// 获取该主题的唯一分区
	partition := int32(0) // 在rhino项目中，kafka都是只有一个分区

	// 获取最新的消息偏移量前一条的消息（即最新一条消息）
	var consumerPartition sarama.PartitionConsumer
	consumerPartition, err = consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
	if err != nil {
		return
	}
	defer func() {
		consumerPartition.Close()
	}()

	// 读取消息

	for {

		select {
		case kMsg := <-consumerPartition.Messages():
			// 成功获取最新消息
			msg := kMsg.Value
			log.Printf("===> msg time:%v\n, ===> data:%s\n", kMsg.Timestamp.Format(time.DateTime), msg)
			messages = append(messages, msg)
			if kMsg.Offset == newestOffset-1 {
				return
			}
		case e := <-consumerPartition.Errors():
			// 如果出错了，打印错误信息
			fmt.Printf("Error occurred while consuming the message, error:%v\n", e.Err)
			err = e.Err
			return
		}
	}

	return
}

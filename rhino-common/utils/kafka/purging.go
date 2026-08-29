package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func PurgeMessageForTopics(brokers []string, topics []string) (err error) {

	beginTime := time.Now()
	maxWaitTime := 20 * time.Minute

	admin, err := createAdminClient(brokers)
	if err != nil {
		return fmt.Errorf("fail to create admin client while try to clear kafka message, error:%v", err)
	}
	defer admin.Close()

	// 配置参数对象
	newConfigEntry := &sarama.ConfigEntry{
		Name:  "retention.ms", // Kafka保留时间参数
		Value: "1",            // 1秒（单位：毫秒）
	}

	for _, topic := range topics {
		// 配置资源定义
		resource := sarama.ConfigResource{
			Type: sarama.TopicResource,
			Name: topic,
		}
		err := admin.AlterConfig(sarama.TopicResource, resource.Name, map[string]*string{newConfigEntry.Name: &newConfigEntry.Value}, false) //
		if err != nil {
			return fmt.Errorf("fail to set retention.ms=1000 for topic %s, error: %v", topic, err)
		}

		log.Printf("success set retention.ms=1000 for topic %s\n", topic)
	}

	for {

		if time.Since(beginTime) > maxWaitTime {
			return fmt.Errorf("time out for checking kafka purging result")
		}

		confirm := true
		for _, topic := range topics {
			log.Printf("===> for topic topic %s\n", topic)
			newestOffset, oldestOffset, err := GetNewestAndOldestOffset(brokers, topic)
			if err != nil {
				return fmt.Errorf("fail to get newest and oldest offset for topic %s, error:%v", topic, err)
			}

			if newestOffset == oldestOffset {
				log.Printf("===> topic:%s, newestOffset:%d, oldestOffset:%d\n", topic, newestOffset, oldestOffset)
				continue
			}

			_, oldestMsgTime, exist, err := GetMessage(brokers, topic, oldestOffset)
			if err != nil {
				return fmt.Errorf("fail to get oldest message for topic %s, error:%v", topic, err)
			}
			if !exist {
				// 没数据时，不进行后续时间比对
				log.Println("没数据时，不进行后续时间比对")
				continue
			}

			log.Printf("===> topic:%s, newestOffset:%d, oldestOffset:%d, oldestMsgTime:%v, beginTime:%v\n", topic, newestOffset, oldestOffset, oldestMsgTime, beginTime)

			if oldestMsgTime.Before(beginTime) {
				confirm = false
				break
			}
		}

		if confirm {
			break
		}

		time.Sleep(10 * time.Second)
	}

	newConfigEntry.Value = "172800000"
	for _, topic := range topics {
		// 配置资源定义
		resource := sarama.ConfigResource{
			Type: sarama.TopicResource,
			Name: topic,
		}
		err := admin.AlterConfig(sarama.TopicResource, resource.Name, map[string]*string{newConfigEntry.Name: &newConfigEntry.Value}, false) //
		if err != nil {
			return fmt.Errorf("fail to set retention.ms=172800000 for topic %s, error: %v", topic, err)
		}

		log.Printf("success set retention.ms=172800000 for topic %s\n", topic)
	}

	log.Printf("Finish purging kafka message for topics:%v\n", topics)

	return nil
}

func createAdminClient(brokers []string) (sarama.ClusterAdmin, error) {
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

	return sarama.NewClusterAdmin(brokers, config)
}

package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

func getProducerConfig() *sarama.Config{
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

	return config
}

func getConsumerConfig() *sarama.Config {
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

	return config
}
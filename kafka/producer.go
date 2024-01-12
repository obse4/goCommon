package kafka

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/obse4/goCommon/logger"
)

type KafkaProducerConfig struct {
	Name                 string    `yaml:"name"`                 // 自定义名称
	Brokers              []string  `yaml:"brokers"`              // broker集群地址 例如["127.0.0.1:9092", "127.0.0.2:9092"]
	Topic                string    `yaml:"topic"`                // topic 默认发送消息的topic，可在SendMessage中重新配置
	Compression          string    `yaml:"compression"`          // 消息压缩方式，可选的值包括"gzip"、"snappy"、"lz4"、"zstd"、"none" 默认"gzip"
	Timeout              int       `yaml:"timeout"`              // 发送消息的超时时间，单位为毫秒 默认30000ms
	BatchSize            int       `yaml:"batchSize"`            // 批量发送的消息数量，超过这个数量后就发送 默认0，即实时发送 与BatchTime连用生效
	BatchTime            int       `yaml:"batchTime"`            // 批量发送的时间间隔，超过这个时间就发送，单位为毫秒 默认0，即实时发送 与BatchSize连用生效
	WaitForAll           bool      `yaml:"waitForAll"`           // 是否等待服务器所有副本都保存成功后再返回 默认false
	MaxRetries           int       `yaml:"maxRetries"`           // 重试的最大次数 默认3
	RetryBackoff         int       `yaml:"retryBackoff"`         // 两次重试之间的时间间隔，单位为毫秒 默认100ms
	NewManualPartitioner bool      `yaml:"newManualPartitioner"` // 是否手动分配分区 默认false
	Producer             *Producer // 生产者指针
}

type Producer struct {
	name     string
	producer sarama.SyncProducer
	topic    string
}

type ProducerMessage struct {
	Topic     string
	Key       string
	Value     string
	Headers   map[string][]byte
	Metadata  interface{}
	Offset    int64
	Partition int32
}

func (p Producer) SendMessage(msg ProducerMessage) error {
	var headers []sarama.RecordHeader

	for k, v := range msg.Headers {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(k),
			Value: v,
		})
	}

	if msg.Topic != "" {
		p.topic = msg.Topic
	}

	message := &sarama.ProducerMessage{
		Topic:     p.topic,
		Key:       sarama.StringEncoder(msg.Key),
		Value:     sarama.StringEncoder(msg.Value),
		Headers:   headers,
		Metadata:  msg.Metadata,
		Offset:    msg.Offset,
		Partition: msg.Partition,
		Timestamp: time.Now(),
	}
	partition, offset, err := p.producer.SendMessage(message)

	if err != nil {
		logger.Error("Kafka Producer new message err %s\n", err.Error())

		return err
	}

	logger.Debug("Kafka Producer %s new message success partition: %d, offset: %d", p.name, partition, offset)

	return nil
}

func NewKafkaProducer(config *KafkaProducerConfig) *Producer {
	saramaConfig := sarama.NewConfig()

	// 配置压缩方式
	switch config.Compression {
	case "gzip":
		saramaConfig.Producer.Compression = sarama.CompressionGZIP
	case "snappy":
		saramaConfig.Producer.Compression = sarama.CompressionSnappy
	case "zstd":
		saramaConfig.Producer.Compression = sarama.CompressionZSTD
	case "none":
		saramaConfig.Producer.Compression = sarama.CompressionNone
	default:
		saramaConfig.Producer.Compression = sarama.CompressionGZIP
	}

	// 配置发送消息超市时间
	if config.Timeout != 0 {
		saramaConfig.Net.DialTimeout = time.Duration(config.Timeout) * time.Millisecond
	}

	// 触发flush的最大上限，即使不到间隔时间也会发送
	if config.BatchSize != 0 {
		saramaConfig.Producer.Flush.Messages = config.BatchSize
	}

	// 触发flush的频率，到达时间即刻发送
	if config.BatchTime != 0 {
		saramaConfig.Producer.Flush.Frequency = time.Duration(config.BatchTime) * time.Millisecond
	}

	// 是否等待服务器所有副本都保存成功后再返回
	if config.WaitForAll {
		saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	} else {
		saramaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	}

	if config.MaxRetries != 0 {
		saramaConfig.Producer.Retry.Max = config.MaxRetries
	}

	if config.RetryBackoff != 0 {
		saramaConfig.Producer.Retry.Backoff = time.Duration(config.RetryBackoff) * time.Millisecond
	}

	// 配置返回成功发送的消息
	saramaConfig.Producer.Return.Successes = true

	if config.NewManualPartitioner {
		saramaConfig.Producer.Partitioner = sarama.NewManualPartitioner
	}

	producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)

	if err != nil {
		logger.Fatal("Kafka new producer %s err %s", config.Name, err.Error())

	}

	config.Producer = &Producer{
		producer: producer,
		topic:    config.Topic,
		name:     config.Name,
	}

	logger.Info("Kafka new producer %s success", config.Name)

	return config.Producer
}

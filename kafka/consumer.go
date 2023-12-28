package kafka

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/obse4/goCommon/logger"
)

type KafkaConsumerConfig struct {
	Name              string    // 自定义名称
	Brokers           []string  // broker集群地址 例如["127.0.0.1:9092", "127.0.0.2:9092"]
	Topics            []string  // topics
	AutoOffsetReset   string    // 开始消费的位置，可能的值包括'earliest'、'latest' 默认‘latest’
	GroupId           string    // 消费组id
	MaxWaitTime       int       // 从Kafka获取记录的最大等待时间（毫秒）默认250ms
	SessionTimeout    int       // 消费者组会话的超时时间（毫秒）默认10000ms
	HeartbeatInterval int       // 心跳间隔时间（毫秒）默认3000ms
	Consumer          *Consumer // 消费者指针
}

type Consumer struct {
	name     string
	topics   []string
	consumer sarama.ConsumerGroup
	handler  ConsumerGroupHandler
}

type ConsumerGroupHandler interface {
	sarama.ConsumerGroupHandler
}

// 新建消费者
func NewKafkaConsumer(config *KafkaConsumerConfig) {
	saramaConfig := sarama.NewConfig()

	saramaConfig.Consumer.Return.Errors = true

	if config.AutoOffsetReset == "earliest" {
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	if config.MaxWaitTime != 0 {
		saramaConfig.Consumer.MaxWaitTime = time.Duration(config.MaxWaitTime) * time.Millisecond
	}

	if config.SessionTimeout != 0 {
		saramaConfig.Consumer.Group.Session.Timeout = time.Duration(config.SessionTimeout) * time.Millisecond
	}

	if config.HeartbeatInterval != 0 {
		saramaConfig.Consumer.Group.Heartbeat.Interval = time.Duration(config.HeartbeatInterval) * time.Millisecond
	}

	consumer, err := sarama.NewConsumerGroup(config.Brokers, config.GroupId, saramaConfig)
	if err != nil {
		logger.Fatal("Kafka new consumer %s err %s", config.Name, err.Error())
		return
	}

	config.Consumer = &Consumer{
		name:     config.Name,
		topics:   config.Topics,
		consumer: consumer,
	}
	logger.Info("Kafka new consumer %s success", config.Name)
}

// 注册处理方法
// type ConsumeFunction func(msg *ConsumerMessage, sess ConsumerGroupSession, claim ConsumerGroupClaim) error
// msg 单条消息指针
// mark 处理成功是否标记此条消息已被消费
func (c *Consumer) RegisterHandle(f ConsumeFunction, mark bool) {
	var handle = consumerGroupHandler{
		name:        c.name,
		consumeFunc: f,
		mark:        mark,
	}
	c.handler = handle

	logger.Info("Kafka consumer %s register handle success", c.name)
}

// 停止消费
func (c *Consumer) Close() error {
	if err := c.consumer.Close(); err != nil {
		logger.Error("Kafka close consumer %s err %s", c.name, err.Error())
		return err
	}
	logger.Info("Kafka colse consumer %s", c.name)
	return nil
}

// 开始消费
func (c *Consumer) Consume(ctx context.Context) error {
	if c.handler == nil {
		logger.Error("Kafka consumer %s has no register handle", c.name)
		return fmt.Errorf("kafka consumer %s has no register handle", c.name)
	}
	// 使用消费者组从主题中消费消息
	go func() {
		logger.Info("Kafka consumer %s start consume", c.name)
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	Loop:
		for {
			select {
			case <-ctx.Done():
				logger.Warn("Kafka consumer %s context cancelled", c.name)
				break Loop
			case <-sigterm:
				logger.Warn("Kafka consumer %s service interrupt", c.name)
				break Loop
			case <-signals:
				logger.Warn("Kafka consumer %s service interrupt", c.name)
			default:
				// `Consume` 应该在一个无限循环中被调用，当服务器端的重新平衡发生时，
				// 消费者会话将需要被重新创建以获取新的声明
				if err := c.consumer.Consume(ctx, c.topics, c.handler); err != nil {
					if err != sarama.ErrClosedConsumerGroup {
						logger.Error("Kafka consumer %s consume error:%v", c.name, err)
					}
					break Loop
				}
			}
		}
	}()

	// 等待上下文取消。这可能是从另一个函数发出的信号，表示这个消费者应该停止。
	<-ctx.Done()

	return c.consumer.Close()
}

type consumerGroupHandler struct {
	name        string
	consumeFunc ConsumeFunction
	mark        bool
}

type ConsumeFunction func(msg *ConsumerMessage, sess ConsumerGroupSession, claim ConsumerGroupClaim) error

type ConsumerGroupSession interface {
	sarama.ConsumerGroupSession
}

type ConsumerGroupClaim interface {
	sarama.ConsumerGroupClaim
}

type ConsumerMessage struct {
	sarama.ConsumerMessage
}

// Setup 在新会话开始之前，ConsumeClaim 之前运行
func (h consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

// Cleanup 在会话结束时运行，所有的 ConsumeClaim goroutines 都已经退出
func (h consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim 必须启动一个消费者循环，处理 ConsumerGroupClaim 的 Messages()
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		logger.Debug("Kafka consumer %s receive message topic %s val %s", h.name, msg.Topic, string(msg.Value))
		// 在这里处理你的消息
		// 标记消息已处理

		err := h.consumeFunc(&ConsumerMessage{ConsumerMessage: *msg}, sess, claim)

		if err != nil {
			logger.Error("Kafka consumer %s receive message topic %s val %s handle err %v", h.name, msg.Topic, string(msg.Value), err)
			return err
		}

		if h.mark {
			sess.MarkMessage(msg, "")
		}
	}
	return nil
}

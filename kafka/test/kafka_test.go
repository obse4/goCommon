package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/obse4/goCommon/kafka"
)

func TestNewProducer(t *testing.T) {
	var kafkaProducer = kafka.KafkaProducerConfig{
		Name:    "kafka_producer_test",
		Brokers: []string{"127.0.0.1:9092"},
		Topic:   "test",
	}

	kafka.NewKafkaProducer(&kafkaProducer)

	for i := 0; i < 50; i++ {
		err := kafkaProducer.Producer.SendMessage(kafka.ProducerMessage{
			Topic: "test",
			Value: fmt.Sprintf("obse4^-^ %d", i),
		})
		if err != nil {
			fmt.Printf("new message err %s\n", err.Error())
			t.Fail()
		}
	}

}

func TestNewConsumer(t *testing.T) {
	var kafkaConsumer = kafka.KafkaConsumerConfig{
		Name:            "kafka_consumer_test",
		Brokers:         []string{"127.0.0.1:9092"},
		Topics:          []string{"test"},
		GroupId:         "consumer_test",
		AutoOffsetReset: "earliest",
		BlockingPool:    4,
	}
	kafka.NewKafkaConsumer(&kafkaConsumer)

	go func() {
		kafkaConsumer.Consumer.RegisterHandle(func(msg *kafka.ConsumerMessage, sess kafka.ConsumerGroupSession, claim kafka.ConsumerGroupClaim) error {
			time.Sleep(time.Second)
			fmt.Printf("topic %s val %s, timestamp %s\n", msg.Topic, msg.Value, msg.Timestamp.Format("2006-01-02 15:04:05"))
			return nil
		}, true)

		kafkaConsumer.Consumer.Consume(context.TODO())
	}()

	// 在另一个goroutine中执行关闭操作
	go func() {
		time.AfterFunc(5*time.Second, func() {
			fmt.Println("close")
			kafkaConsumer.Consumer.Close()
		})
	}()

	time.Sleep(time.Second * 10)
}

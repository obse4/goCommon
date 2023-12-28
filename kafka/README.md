# goCommon kafka

## 使用
- producer

```
import (
    "github.com/obse4/goCommon/kafka"
    "fmt"
)

func main() {
    var kafkaProducer = kafka.KafkaProducerConfig{
		Name:    "kafka_producer_test",
		Brokers: []string{"127.0.0.1:9092"},
		Topic:   "test",
	}

    kafka.NewKafkaProducer(&kafkaProducer)

    err := kafkaProducer.Producer.SendMessage(kafka.ProducerMessage{
		Topic: "test",
		Value: "obse4^-^",
	})

	if err != nil {
		fmt.Printf("new message err %s\n", err.Error())
	}
}
```
- consumer

```
import (
    "github.com/obse4/goCommon/kafka"
    "fmt"
)

func main() {
    var kafkaConsumer = kafka.KafkaConsumerConfig{
		Name:            "kafka_consumer_test",
		Brokers:         []string{"127.0.0.1:9092"},
		Topics:          []string{"test"},
		GroupId:         "consumer_test",
		AutoOffsetReset: "earliest",
	}
	kafka.NewKafkaConsumer(&kafkaConsumer)

    // func (*kafka.Consumer).RegisterHandle(f kafka.ConsumeFunction, mark bool)
    // f func(msg *kafka.ConsumerMessage, sess kafka.ConsumerGroupSession, claim kafka.ConsumerGroupClaim) error
    // msg 消息指针，sess、claim不常用
    // mark 标记已消费消息
    go func() {
		kafkaConsumer.Consumer.RegisterHandle(func(msg *kafka.ConsumerMessage, sess kafka.ConsumerGroupSession, claim kafka.ConsumerGroupClaim) error {
			fmt.Printf("topic %s val %s, timestamp %s\n", msg.Topic, msg.Value, msg.Timestamp.Format("2006-01-02 15:04:05"))
			return nil
		}, true)

		kafkaConsumer.Consumer.Consume(context.TODO())
	}()

    //  测试Consumer.Close()
    // 在另一个goroutine中执行关闭操作
	go func() {
		time.AfterFunc(5*time.Second, func() {
			fmt.Println("close")
			kafkaConsumer.Consumer.Close()
		})
	}()

    // 10s后执行结束
    time.Sleep(time.Second * 10)
}
```
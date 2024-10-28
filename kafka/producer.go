package Kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"strconv"
)

type Producer struct {
	Client *kafka.Writer
}

type SetProducerCfg func(writer *kafka.Writer)

func SetProducerTopic(topic string) SetProducerCfg {
	return func(writer *kafka.Writer) {
		writer.Topic = topic
	}
}

func SetProducerAddr(addr string) SetProducerCfg {
	return func(writer *kafka.Writer) {
		writer.Addr = kafka.TCP(addr)
	}
}

// NewProducer 新建生产者
func NewProducer(cfg ...SetProducerCfg) *Producer {
	writer := kafka.Writer{
		Addr:     kafka.TCP("23.95.15.178:9092"),
		Topic:    "tp1",
		Balancer: &kafka.LeastBytes{},
	}
	for _, opt := range cfg {
		opt(&writer)
	}
	producer := &Producer{
		Client: &writer,
	}
	return producer
}

func (producer *Producer) WriteData(key *[]byte, value *[]byte) error {
	err := producer.Client.WriteMessages(context.Background(), kafka.Message{
		Key:   *key,
		Value: *value,
	})
	if err != nil {
		log.Println(err)
	}
	return err
}

// CreateTopicWithRetention 创建话题并设置过期时间
func (producer *Producer) CreateTopicWithRetention(topic string, retention string, addr string) {
	conn, err := kafka.Dial("tcp", addr)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	configEntries := make([]kafka.ConfigEntry, 1)
	configEntries[0] = kafka.ConfigEntry{
		ConfigName:  "retention.ms",
		ConfigValue: retention,
	}

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
			ConfigEntries:     configEntries,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}

func (producer *Producer) GetTopic() {
	conn, err := kafka.Dial("tcp", "23.95.15.178:9092")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 获取所有 topic
	topics, err := conn.ReadPartitions("t")

	// 打印所有 topic 的名称
	for _, topic := range topics {
		fmt.Println("Topic:", topic.Topic)
	}
}

package Kafka

import (
	"github.com/segmentio/kafka-go"
)

type SetConsumerGroupConfig func(config *kafka.ReaderConfig)

// SetConsumerTopic 设置消费者的topic
func SetConsumerTopic(str string) SetConsumerGroupConfig {
	return func(config *kafka.ReaderConfig) {
		config.Topic = str
	}
}

// SetConsumerBrokers 设置消费者broker
func SetConsumerBrokers(broker []string) SetConsumerGroupConfig {
	return func(config *kafka.ReaderConfig) {
		config.Brokers = broker
	}
}

// SetConsumerGroupID 设置消费者groupID
func SetConsumerGroupID(str string) SetConsumerGroupConfig {
	return func(config *kafka.ReaderConfig) {
		config.GroupID = str
	}
}

// SetConsumerMaxBytes 设置消费者maxBytes
func SetConsumerMaxBytes(max int) SetConsumerGroupConfig {
	return func(config *kafka.ReaderConfig) {
		config.MaxBytes = max
	}
}

// NewConsumer 新建消费者
func NewConsumer(cfg ...SetConsumerGroupConfig) *kafka.Reader {
	// make a new reader that consumes from topic-A
	var config = kafka.ReaderConfig{
		Brokers:  []string{"23.95.15.178:9092"},
		GroupID:  "0",
		Topic:    "tp1",
		MaxBytes: 10e6, // 10MB
	}
	for _, c := range cfg {
		c(&config)
	}

	r := kafka.NewReader(config)
	return r
}

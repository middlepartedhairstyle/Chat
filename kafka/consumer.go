package Kafka

import (
	"github.com/segmentio/kafka-go"
)

type SetConsumerGroupConfig func(config *kafka.ReaderConfig)

// SetTopic 设置消费者的topic
func SetTopic(str string) SetConsumerGroupConfig {
	return func(config *kafka.ReaderConfig) {
		config.Topic = str
	}
}

// SetBrokers 设置消费者broker
func SetBrokers(broker []string) SetConsumerGroupConfig {
	return func(config *kafka.ReaderConfig) {
		config.Brokers = broker
	}
}

// SetGroupID 设置消费者groupID
func SetGroupID(str string) SetConsumerGroupConfig {
	return func(config *kafka.ReaderConfig) {
		config.GroupID = str
	}
}

// SetMaxBytes 设置消费者maxBytes
func SetMaxBytes(max int) SetConsumerGroupConfig {
	return func(config *kafka.ReaderConfig) {
		config.MaxBytes = max
	}
}

// NewConsumer 新建消费者
func NewConsumer(cfg ...SetConsumerGroupConfig) *kafka.Reader {
	// make a new reader that consumes from topic-A
	var config = kafka.ReaderConfig{
		Brokers:  []string{"23.95.15.178:9092"},
		GroupID:  "f",
		Topic:    "tp1",
		MaxBytes: 10e6, // 10MB
	}
	for _, c := range cfg {
		c(&config)
	}

	r := kafka.NewReader(config)
	return r
}

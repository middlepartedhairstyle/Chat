package Kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
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

package Kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

func Producer(key string, value string) error {
	// make a writer that produces to topic-A, using the least-bytes distribution
	w := &kafka.Writer{
		Addr:     kafka.TCP("23.95.15.178:9092"),
		Topic:    "tp1",
		Balancer: &kafka.LeastBytes{},
	}

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(value),
		},
	)
	fmt.Println(key)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err = w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
	return nil
}

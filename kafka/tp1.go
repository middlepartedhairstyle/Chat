package Kafka

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func Write(msg []byte) bool {
	err := KafkaTp1.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return true
	}
	_, err = KafkaTp1.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err = KafkaTp1.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
	return true
}

func Read() {
	KafkaTp1.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := KafkaTp1.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	if err := KafkaTp1.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}

package Kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

var KafkaTp1 *kafka.Conn

const (
	Topic1 = "tp1"
)

const (
	Partition0 int = iota
	Partition1
	Partition2
)

func Init() {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "23.95.15.178:9092", Topic1, Partition0)
	if err != nil {
		fmt.Println("init kafka err:", err)
	} else {
		KafkaTp1 = conn
		fmt.Println("init kafka success")
	}
}

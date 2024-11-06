package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	Kafka "github.com/middlepartedhairstyle/HiWe/kafka"
	"strconv"
)

const (
	UserChatMessageType = iota + 1
	UserMessageType
)

type Info struct {
	Types uint8           `json:"type"`
	Data  json.RawMessage `json:"data"`
}

type DisposeInfo interface {
	Marshal() ([]byte, error)
}

func NewInfo() *Info {
	return new(Info)
}

func (info *Info) CheckType() Information {
	switch info.Types {
	case UserChatMessageType:
		msg := NewUserChatMessage()
		err := json.Unmarshal(info.Data, msg)
		if err != nil {
			return nil
		}
		return msg
	case UserMessageType:
		return nil
	default:
		return nil
	}
}

// WriteKafka 将消息写入kafka
func (info *Info) WriteKafka(disposeInfo DisposeInfo, args ...interface{}) error {
	var topic string
	var userID uint
	if len(args) >= 2 {
		topic = args[0].(string)
		userID = args[1].(uint)
		info.Data, _ = disposeInfo.Marshal()
		producer := Kafka.NewProducer(Kafka.SetProducerTopic(topic))
		if !producer.GetTopic(topic) {
			producer.CreateTopicWithRetention(topic)
		}
		fmt.Println(topic)
		key := []byte(UserMessageBaseGroup + strconv.Itoa(int(userID))) //u5,u1
		message, _ := info.Marshal()
		err := producer.WriteData(&key, &message)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	}
	return errors.New("marshal error")
}

func (info *Info) ReadKafka(userID uint, ws *WebSocketClient) {
	topic := fmt.Sprintf("%s%s%v", UserMessageBaseTopic, "tp", strconv.Itoa(int(userID/maxUser+1)))
	key := UserMessageBaseGroup + strconv.Itoa(int(userID))

	//检查topic是否存在
	producer := Kafka.NewProducer(Kafka.SetProducerTopic(topic))
	if !producer.GetTopic(topic) {
		producer.CreateTopicWithRetention(topic)
	}

	consumer := Kafka.NewConsumer(Kafka.SetConsumerTopic(topic), Kafka.SetConsumerGroupID(UserMessageBaseGroup+strconv.Itoa(int(userID))))
	for {
		message, err := consumer.ReadMessage(context.Background())
		if err != nil {
			continue
		}
		if string(message.Key) == key {
			ws.messageList <- message.Value
		} else {
			if err = consumer.CommitMessages(context.Background(), message); err != nil {
				fmt.Printf("提交偏移量失败: %v\n", err)
			}
		}
	}
}

func (info *Info) Marshal() ([]byte, error) {
	return json.Marshal(info)
}

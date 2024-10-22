package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	Kafka "github.com/middlepartedhairstyle/HiWe/kafka"
	"github.com/segmentio/kafka-go"
	"net/http"
	"strconv"
	"time"
)

type WebSocketClient struct {
	Conn        *websocket.Conn
	Context     *gin.Context
	FromId      uint `json:"from_id"`
	messageList chan []byte
}

// NewWebSocketClient 创建一个新的webSocket连接
func NewWebSocketClient(c *gin.Context, check bool, userId uint) (*WebSocketClient, error) {
	var ws = &WebSocketClient{}
	var err error
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return check
		},
	}
	ws.Conn, err = upGrader.Upgrade(c.Writer, c.Request, nil)
	ws.FromId = userId
	ws.Context = c
	ws.messageList = make(chan []byte, 50)
	return ws, err
}

// Close 关闭连接
func (ws *WebSocketClient) Close() error {
	return ws.Conn.Close()
}

// ReadMessage 读取消息
func (ws *WebSocketClient) ReadMessage() ([]byte, Information, error) {
	msg := NewInfo()
	_, data, err := ws.Conn.ReadMessage()
	if err != nil {
		return nil, nil, err
	}
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return nil, nil, err
	}
	return data, msg.CheckType(), nil
}

// WriteMessage 写入消息
func (ws *WebSocketClient) WriteMessage(messageType int, message []byte) error {
	tm := time.Now().Format(time.DateTime)
	data := fmt.Sprintf("[ws][%s]%s", tm, message)
	err := ws.Conn.WriteMessage(messageType, []byte(data))
	if err != nil {
		return err
	}
	return nil
}

// SendMessage 发送消息
func (ws *WebSocketClient) SendMessage(fromId uint) {
	var producers = make(map[string]*Kafka.Producer) //[topic]*Producer
	for {
		msg, ucm, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			err = ws.Close()
			if err != nil {
				return
			}
			break
		}
		//查看消息类型格式是否正确
		if ucm == nil {
			fmt.Println(err)
			continue
		}
		ucm.MessageDispose(producers, fromId, msg)
	}

}

// GetMessage 获取消息
func (ws *WebSocketClient) GetMessage(id uint) {
	// 查找用户消息
	go func(id uint) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovering from panic in GetMessage:", r)
			}
		}()

		var message []byte // 消息
		var tpId = id/maxUserNum + 1
		topic := fmt.Sprintf("%s%s%d", ChatWithFriend, "tp", tpId) // 例如 ftp1, ftp2
		consumer := Kafka.NewConsumer(Kafka.SetConsumerTopic(topic), Kafka.SetConsumerGroupID(strconv.Itoa(int(id))))
		defer func(consumer *kafka.Reader) {
			err := consumer.Close()
			if err != nil {

			}
		}(consumer) // 确保消费者关闭

		for {
			m, err := consumer.ReadMessage(context.Background())
			if err != nil {
				fmt.Println("Error reading message:", err)
				continue
			}
			if string(m.Key) == strconv.Itoa(int(id)) {
				message = m.Value
				ws.messageList <- message
			} else {
				if err = consumer.CommitMessages(context.Background(), m); err != nil {
					fmt.Printf("提交偏移量失败: %v\n", err)
				}
			}
		}
	}(id)

	// 消息写入 websocket
	for {
		select {
		case msg := <-ws.messageList:
			err := ws.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Println("Error writing message to websocket:", err)
				continue
			}
		}
	}
}

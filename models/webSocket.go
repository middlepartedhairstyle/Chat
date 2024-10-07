package models

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	Kafka "github.com/middlepartedhairstyle/HiWe/kafka"
	"github.com/middlepartedhairstyle/HiWe/mySQL"
	"net/http"
	"strconv"
	"time"
)

type WebSocketClient struct {
	Conn        *websocket.Conn
	Context     *gin.Context
	messageList chan UserMessage
	UserMessage
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
	ws.FromID = userId
	ws.Context = c
	ws.messageList = make(chan UserMessage, 50)
	return ws, err
}

// Close 关闭连接
func (ws *WebSocketClient) Close() error {
	return ws.Conn.Close()
}

// ReadMessage 读取消息
func (ws *WebSocketClient) ReadMessage() (UserMessage, error) {
	var msg UserMessage
	_, data, err := ws.Conn.ReadMessage()
	if err != nil {
		return msg, err
	}
	msg, err = ws.Unmarshal(data)
	if err != nil {
		return msg, err
	}
	return msg, nil
}

// WriteMessage 写入消息
func (ws *WebSocketClient) WriteMessage(messageType int, message UserMessage) error {
	tm := time.Now().Format(time.DateTime)
	msg, err := message.Marshal(message)
	if err != nil {
		return err
	}
	data := fmt.Sprintf("[ws][%s]%s", tm, msg)
	err = ws.Conn.WriteMessage(messageType, []byte(data))
	if err != nil {
		return err
	}
	return nil
}

// SendMessage 发送消息
func (ws *WebSocketClient) SendMessage() {
	for {
		msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			continue
		}

		bytes, err := ws.Marshal(msg)
		if err != nil {
			return
		}
		err = Kafka.Producer(ws.MessageChannel(msg.Media, msg.ToID), string(bytes))
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

// GetMessage 获取消息
func (ws *WebSocketClient) GetMessage(id uint) {

	go func(id uint) {
		r := Kafka.NewConsumer(Kafka.SetGroupID(ChatWithFriend + strconv.FormatUint(uint64(id), 10)))
		//查找好友
		var friendIDs []string
		for _, item := range mySQL.FindAllFriendId(id) {
			friendIDs = append(friendIDs, ChatWithFriend+strconv.Itoa(int(item)))
		}
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				break
			}
			var message string
			for index, key := range friendIDs {
				if string(m.Key) == key {
					message = string(m.Value)
					if err != nil {
						fmt.Println(err)
						break
					}
					var data UserMessage
					data, err = data.Unmarshal([]byte(message))
					if err != nil {
						return
					}
					ws.messageList <- data
					break // 找到匹配后可以退出循环
				}
				if index == len(friendIDs)-1 {
					if err = r.CommitMessages(context.Background(), m); err != nil {
						fmt.Printf("提交偏移量失败: %v\n", err)
					}
				}
			}
		}
	}(id)
	//消息写入websocket
	for {
		select {
		case msg := <-ws.messageList:
			err := ws.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				break
			}
		}
	}
}

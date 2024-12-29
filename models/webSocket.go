package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	Kafka "github.com/middlepartedhairstyle/HiWe/kafka"
	"net/http"
	"sync"
	"time"
)

type WebSocketClient struct {
	Conn               *websocket.Conn
	Context            *gin.Context
	Ctx                context.Context
	Mu                 *sync.Mutex
	Cancel             context.CancelFunc
	PingTicker         *time.Ticker
	ConnState          bool //连接状态
	FromId             uint `json:"from_id"`
	messageList        chan []byte
	GroupChangeMessage chan uint //用户群聊发生变化，如用户加入新群聊
}

var GroupChangeMessage = make(map[uint]*chan uint)

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
	ws.PingTicker = time.NewTicker(30 * time.Second) //设置ping时间间隔为30秒
	ws.ConnState = true
	ws.Ctx, ws.Cancel = context.WithCancel(context.Background())
	ws.Mu = new(sync.Mutex)
	ws.messageList = make(chan []byte, 50)
	ws.GroupChangeMessage = make(chan uint, 3)
	GroupChangeMessage[ws.FromId] = &ws.GroupChangeMessage
	return ws, err
}

// Close 关闭连接
func (ws *WebSocketClient) Close() error {
	ws.PingTicker.Stop()
	err := ws.Conn.Close()
	ws.ConnState = false
	ws.Cancel()
	if err != nil {
		return err
	}
	return err
}

func (ws *WebSocketClient) Ping() bool {
	return ws.ConnState
}

func (ws *WebSocketClient) SetConnState(state bool) {
	ws.ConnState = state
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

// SendMessage 发送消息(待完善，连接断开后关闭所有producers连接)
func (ws *WebSocketClient) SendMessage(fromId uint) {
	var producers = make(map[string]*Kafka.Producer) //[topic]*Producer
	var infoVerify = make(map[string]uint)           //例如[f1]2,[g]2
	for {
		msg, ucm, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			err = ws.Close()
			if err != nil {
				return
			}
			fmt.Println("close ")
			break
		}
		//查看消息类型格式是否正确
		if ucm == nil {
			fmt.Println(err)
			continue
		}
		ucm.MessageDispose(producers, fromId, msg, infoVerify)
	}

}

// GetMessage 获取消息
func (ws *WebSocketClient) GetMessage(id uint) {
	defer func() {
		_ = ws.Close()
		fmt.Println("close get message")
	}()

	//从kafka中读取新消息
	userChatMessage := NewUserChatMessage()
	info := NewInfo()
	go userChatMessage.GetFriendMessage(id, ws)
	go userChatMessage.GetGroupMessage(id, ws)
	go info.ReadKafka(id, ws)

	// 消息写入 websocket
	for ws.Ping() {
		select {
		case msg := <-ws.messageList:
			err := ws.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Println("Error writing message to websocket:", err)
				continue
			}
		case <-ws.PingTicker.C:
			fmt.Println("test1")
			err := ws.Conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				fmt.Println("Error writing ping message to websocket:", err)
				ws.SetConnState(false)
				return
			}
		}
	}
}

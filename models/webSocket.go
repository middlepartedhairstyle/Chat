package models

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type WebSocketClient struct {
	Conn        *websocket.Conn
	Context     *gin.Context
	messageList chan UserMessage
	UserMessage
}

// NewWebSocketClient 创建一个新的webSocket连接
func NewWebSocketClient(c *gin.Context, check bool, userId uint64) (ws *WebSocketClient, err error) {
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

// SendMessage 发送消息
func (ws *WebSocketClient) SendMessage() {

}

// GetMessage 获取消息
func (ws *WebSocketClient) GetMessage() {

}

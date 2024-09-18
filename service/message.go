package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/middlepartedhairstyle/HiWe/redis"
	"time"
)

const (
	PublishKey = "websocket" //订阅消息者管道，后期更改为用户具体管道
)

// SendMsg 发送消息
func SendMsg(ws *websocket.Conn, c *gin.Context) {
	for {
		messageType, data, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		if messageType == websocket.TextMessage {
			message := string(data)
			err = redis.Publish(c, PublishKey, message)
			fmt.Println("发送消息:", message)
			if err != nil {
				fmt.Println(err)
				break
			}

		}

	}
}

func GetMsg(ws *websocket.Conn, c *gin.Context) {
	for {
		m, err := redis.Subscribe(c, PublishKey)
		if err != nil {
			fmt.Println(err)
			break
		}
		tm := time.Now().Format("2006-01-02 15:04:05")
		m = fmt.Sprintf("[ws][%s]%s", tm, m)
		err = ws.WriteMessage(websocket.TextMessage, []byte(m))
		fmt.Println("订阅消息:", m)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

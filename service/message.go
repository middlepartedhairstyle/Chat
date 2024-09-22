package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/redis"
	"gopkg.in/fatih/set.v0"
	"net/http"
	"strconv"
	"time"
)

const (
	ChatWithFriend = "chatWithFriend" //订阅消息者管道，后期更改为用户具体管道
)

// SendMsg 发送消息
func SendMsg(ws *websocket.Conn, c *gin.Context, channel string) {
	for {
		messageType, data, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		if messageType == websocket.TextMessage {
			message := string(data)
			err = redis.Publish(c, ChatWithFriend+channel, message)
			fmt.Println("发送消息:", message)
			if err != nil {
				fmt.Println(err)
				break
			}

		}

	}

	defer func(ws *websocket.Conn) {
		err := ws.Close()
		fmt.Println("close")
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
}

// GetMsg 获取消息
func GetMsg(ws *websocket.Conn, c *gin.Context, channel string) {
	for {
		m, err := redis.Subscribe(c, ChatWithFriend+channel)
		if err != nil {
			fmt.Println(err)
			break
		}
		tm := time.Now().Format(time.RFC3339Nano)
		m = fmt.Sprintf("[ws][%s]%s", tm, m)
		err = ws.WriteMessage(websocket.TextMessage, []byte(m))
		fmt.Println("订阅消息:", m)
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	defer func(ws *websocket.Conn) {
		err := ws.Close()
		fmt.Println("close")
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
}

func Chat(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	//context := query.Get("context")
	fromID, _ := strconv.Atoi(query.Get("from_id")) //信息发送者id
	//toID, _ := strconv.Atoi(query.Get("to_id"))               //信息发向者id
	//messageType, _ := strconv.Atoi(query.Get("message_type")) //消息类型
	//media, _ := strconv.Atoi(query.Get("media"))              //发送消息类型，例如私聊，群聊等
	token := query.Get("token")

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return models.CheckToken(uint64(fromID), token)
		}}).Upgrade(writer, request, nil)

	if err != nil {
		return
	}

	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	note := &models.Note{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	models.RwLocker.Lock()
	models.ClientMap[fromID] = note
	models.RwLocker.Unlock()

	go SendProc(note)

	go GetProc(note)
}

func SendProc(note *models.Note) {
	for {
		select {
		case data := <-note.DataQueue:
			err := note.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
func GetProc(note *models.Note) {
	for {
		messageType, data, err := note.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		if messageType == websocket.TextMessage {
			message := string(data)
			//err = redis.Publish(c, PublishKey, message)
			fmt.Println("发送消息:", message)
			if err != nil {
				fmt.Println(err)
				break
			}

		}
	}
}

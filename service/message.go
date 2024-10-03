package service

//
//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/gorilla/websocket"
//	"github.com/middlepartedhairstyle/HiWe/models"
//	"github.com/middlepartedhairstyle/HiWe/mySQL"
//	"github.com/middlepartedhairstyle/HiWe/redis"
//	"strconv"
//	"time"
//)
//
//const (
//	ChatWithFriend = "chatWithFriend" //订阅消息者管道，后期更改为用户具体管道
//)
//
//// SendMsg 发送消息
//func SendMsg(ws *websocket.Conn, c *gin.Context) {
//	var data models.UserMessage
//	for {
//		messageType, d, err := ws.ReadMessage()
//		if err != nil {
//			fmt.Println(err)
//			break
//		}
//		if messageType == websocket.TextMessage {
//			data.ToJson(d)
//			err = redis.Publish(c, ChatWithFriend+strconv.FormatUint(data.ToID, 10), string(d))
//			fmt.Println("发送消息:", string(d))
//			if err != nil {
//				fmt.Println(err)
//				break
//			}
//
//		}
//	}
//
//	defer func(ws *websocket.Conn) {
//		err := ws.Close()
//		fmt.Println("close")
//		if err != nil {
//			fmt.Println(err)
//		}
//	}(ws)
//}
//
//// GetMsg 获取消息
//func GetMsg(ws *websocket.Conn, c *gin.Context, id uint64) {
//	var message chan models.UserMessage
//	message = make(chan models.UserMessage, 50)
//
//	//查找好友
//	for _, item := range mySQL.SelectAllFriend(uint(id)) {
//		fmt.Println(item)
//		go func(id uint64) {
//			for {
//				m, err := redis.Subscribe(c, ChatWithFriend+strconv.FormatUint(id, 10))
//				if err != nil {
//					fmt.Println(err)
//					break
//				}
//				var data models.UserMessage
//				data.ToJson([]byte(m))
//
//				message <- data
//			}
//		}(uint64(item))
//	}
//
//	for {
//		msg := <-message
//		tm := time.Now().Format(time.RFC3339Nano)
//		msgs := fmt.Sprintf("[ws][%s]%s", tm, string(msg.FromJson()))
//		err := ws.WriteMessage(websocket.TextMessage, []byte(msgs))
//		fmt.Println("订阅消息:", msgs)
//		if err != nil {
//			fmt.Println(err)
//			break
//		}
//	}
//
//	defer func(ws *websocket.Conn) {
//		err := ws.Close()
//		fmt.Println("close")
//		if err != nil {
//			fmt.Println(err)
//		}
//	}(ws)
//}

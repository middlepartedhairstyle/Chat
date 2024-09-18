package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/middlepartedhairstyle/HiWe/service"
	"net/http"
)

const (
	CSRF = "CSRF-TOKEN" //防止跨站点伪造请求，后期更改为用户交流密钥
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		csrf := r.Header.Get("csrf-token")
		return csrf == CSRF
	},
}

func SendMsgController(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	service.SendMsg(ws, c)
}

func GetMsgController(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)

	service.GetMsg(ws, c)
}

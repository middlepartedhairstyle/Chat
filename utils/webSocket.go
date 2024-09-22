package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func UpGraderFriend(c *gin.Context, check bool) (ws *websocket.Conn, err error) {
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			fmt.Println(check)
			return check
		},
	}
	ws, err = upGrader.Upgrade(c.Writer, c.Request, nil)
	return ws, err
}

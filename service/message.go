package service

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/utils"
)

// Chat 聊天
func (h *HTTPServer) Chat(c *gin.Context) {
	userId, _ := utils.StringToUint(c.Query("user_id"))
	client, err := models.NewWebSocketClient(c, true, userId)
	if err != nil {
		return
	}
	go client.SendMessage(userId)
	go client.GetMessage(userId)
}

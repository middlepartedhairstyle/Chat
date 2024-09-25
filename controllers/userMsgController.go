package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/service"
	"github.com/middlepartedhairstyle/HiWe/utils"
	"strconv"
)

func ChatWithFriendController(c *gin.Context) {
	//获取信息
	var friend models.Friend
	friend.Id = utils.ToUint64(c.Query("id"))
	friend.UserID = utils.ToUint64(c.Query("user_id"))
	friend.FriendID = utils.ToUint64(c.Query("friend_id"))
	friend.UserToken = c.GetHeader("token")

	//判断是否为好友
	if friend.IsFriend() {

		//判断用户token,正确升级为webSocket
		ws, err1 := utils.UpGraderFriend(c, models.CheckToken(friend.UserID, friend.UserToken))
		if err1 != nil {
			fmt.Println(err1)
			return
		}

		//好友消息通道
		channel := strconv.FormatUint(friend.Id, 10)

		//开启聊天协程
		go service.SendMsg(ws, c, channel)
		go service.GetMsg(ws, c, channel)

	} else {
		return
	}
}

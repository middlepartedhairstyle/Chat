package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/utils"
)

func AddFriendController(c *gin.Context) {
	var friend models.Friend
	friend.UserID = utils.ToUint64(c.Query("user_id"))
	friend.FriendID = utils.ToUint64(c.Query("friend_id"))
	token := c.GetHeader("token")
	if models.CheckToken(friend.UserID, token) {
		b, msg := friend.AddFriend()

		if b {
			utils.Success(c, msg, gin.H{
				"user_id":   friend.FriendID,
				"friend_id": friend.UserID,
				"id":        friend.Id,
			})
		} else {
			utils.Fail(c, msg, gin.H{
				"user_id":   friend.FriendID,
				"friend_id": friend.UserID,
			})
		}
	} else {
		utils.Fail(c, "用户未登录", gin.H{
			"user_id":   friend.UserID,
			"friend_id": friend.FriendID,
		})
	}
}

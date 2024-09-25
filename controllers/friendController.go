package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/utils"
)

func AddFriendController(c *gin.Context) {
	var friendInfo models.FriendInfo
	friendInfo.UserId = utils.ToUint64(c.Query("user_id"))
	friendInfo.FriendId = utils.ToUint64(c.Query("friend_id"))
	token := c.GetHeader("token")
	if models.CheckToken(friendInfo.UserId, token) {
		b, msg := friendInfo.AddFriend()

		if b {
			utils.Success(c, msg, gin.H{
				"user_id":   friendInfo.FriendId,
				"friend_id": friendInfo.UserId,
				"id":        friendInfo.ID,
			})
		} else {
			utils.Fail(c, msg, gin.H{
				"user_id":   friendInfo.FriendId,
				"friend_id": friendInfo.UserId,
			})
		}
	} else {
		utils.Fail(c, "用户未登录", gin.H{
			"user_id":   friendInfo.UserId,
			"friend_id": friendInfo.FriendId,
		})
	}
}

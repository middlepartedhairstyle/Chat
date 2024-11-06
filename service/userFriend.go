package service

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/mySQL"
	"github.com/middlepartedhairstyle/HiWe/utils"
)

// GetFriendList 获取好友列表
func GetFriendList(c *gin.Context) {
	var user models.UserBaseInfo
	//获取用户数据
	user.Id, _ = utils.StringToUint(c.Query("id"))
	//查询用户好友列表
	var friendList []mySQL.Friends //存放列表数据
	friendList, _ = user.GetFriendList()
	utils.Success(c, "成功", friendList)
}

// GetRequestFriendList 获取好友请求添加列表(用于首次登录)
func GetRequestFriendList(c *gin.Context) {
	var user models.UserBaseInfo
	var friendList []mySQL.RequestFriend
	var b bool
	user.Id, _ = utils.StringToUint(c.Query("id"))
	friendList, b = user.GetRequestFriendList()
	if b {
		utils.Success(c, "成功", friendList)
	} else {
		utils.Fail(c, "失败", friendList)
	}
}

// RequestAddFriend 添加好友
func RequestAddFriend(c *gin.Context) {
	var user models.UserBaseInfo
	var fromId uint
	var toId uint
	fromId, _ = utils.StringToUint(c.Query("from_id"))
	toId, _ = utils.StringToUint(c.Query("to_id"))
	user.Id = fromId
	err := user.RequestAddFriend(fromId, toId)
	if err {
		utils.Success(c, "成功", gin.H{
			"from_id": fromId,
			"to_id":   toId,
		})
	} else {
		utils.Fail(c, "失败", gin.H{
			"from_id": fromId,
			"to_id":   toId,
		})
	}
}

// DisposeAddFriend 处理好友请求
func DisposeAddFriend(c *gin.Context) {
	var user models.UserBaseInfo
	var friend mySQL.Friends
	var requestId uint
	var state uint8
	friend.UserOneID, _ = utils.StringToUint(c.Query("from_id"))
	friend.UserTwoID, _ = utils.StringToUint(c.Query("to_id"))
	requestId, _ = utils.StringToUint(c.Query("request_id")) //请求好友id
	state, _ = utils.StringToUint8(c.Query("state"))
	user.Id = friend.UserTwoID

	b, s := user.DisposeAddFriend(friend, requestId, state)
	if b {
		utils.Success(c, "成功", gin.H{
			"state":   s,
			"from_id": friend.UserOneID,
			"to_id":   friend.UserTwoID,
		})
	} else {
		utils.Fail(c, "失败", gin.H{
			"state": s,
			"id":    friend.UserTwoID,
		})
	}
}

// DeleteFriend 删除好友
func DeleteFriend(c *gin.Context) {}

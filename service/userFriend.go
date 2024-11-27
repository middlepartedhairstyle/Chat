package service

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/mySQL/tables"
	"github.com/middlepartedhairstyle/HiWe/utils"
)

const (
	AlreadyFriend        = 21101 //已经是好友
	DisposeAddFriendFail = 21102 //处理好友请求失败
)

// GetFriendList 获取好友列表
func (h *HTTPServer) GetFriendList(c *gin.Context) {
	var user models.UserBaseInfo
	//获取用户数据
	user.Id, _ = utils.StringToUint(c.Query("id"))
	//查询用户好友列表
	var friendList []tables.Friends //存放列表数据
	friendList, _ = user.GetFriendList()
	utils.Success(c, SUCCESS, friendList)
}

// GetRequestFriendList 获取好友请求添加列表(用于首次登录)
func (h *HTTPServer) GetRequestFriendList(c *gin.Context) {
	var user models.UserBaseInfo
	var friendList []tables.RequestFriend
	var b bool
	user.Id, _ = utils.StringToUint(c.Query("id"))
	friendList, b = user.GetRequestFriendList()
	if b {
		utils.Success(c, SUCCESS, gin.H{
			"friend_list": friendList,
		})
	} else {
		utils.Fail(c, ServerError, gin.H{
			"friend_list": friendList,
		})
	}
}

// RequestAddFriend 请求添加好友
func (h *HTTPServer) RequestAddFriend(c *gin.Context) {
	var user models.UserBaseInfo
	var fromId uint
	var toId uint
	fromId, _ = utils.StringToUint(c.Query("from_id"))
	toId, _ = utils.StringToUint(c.Query("to_id"))
	user.Id = fromId
	err := user.RequestAddFriend(fromId, toId)
	if err {
		utils.Success(c, SUCCESS, gin.H{
			"from_id": fromId,
			"to_id":   toId,
		})
	} else {
		utils.Fail(c, AlreadyFriend, gin.H{
			"from_id": fromId,
			"to_id":   toId,
		})
	}
}

// DisposeAddFriend 处理好友请求
func (h *HTTPServer) DisposeAddFriend(c *gin.Context) {
	var user models.UserBaseInfo
	var friend tables.Friends
	var requestId uint
	var state uint8
	friend.UserOneID, _ = utils.StringToUint(c.Query("from_id"))
	friend.UserTwoID, _ = utils.StringToUint(c.Query("to_id"))
	requestId, _ = utils.StringToUint(c.Query("request_id")) //请求好友id
	state, _ = utils.StringToUint8(c.Query("state"))
	user.Id = friend.UserTwoID

	b, s := user.DisposeAddFriend(friend, requestId, state)
	if b {
		utils.Success(c, SUCCESS, gin.H{
			"state":   s,
			"from_id": friend.UserOneID,
			"to_id":   friend.UserTwoID,
		})
	} else {
		utils.Fail(c, DisposeAddFriendFail, gin.H{
			"state":   s,
			"from_id": friend.UserOneID,
			"to_id":   friend.UserTwoID,
		})
	}
}

// DeleteFriend 删除好友
func DeleteFriend(c *gin.Context) {}

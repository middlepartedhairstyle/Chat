package service

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/internal/models"
	tables2 "github.com/middlepartedhairstyle/HiWe/internal/mySQL/tables"
	utils2 "github.com/middlepartedhairstyle/HiWe/internal/utils"
)

const (
	AlreadyFriend        = 21101 //已经是好友
	DisposeAddFriendFail = 21102 //处理好友请求失败
	ChangeFriendNoteFail = 21103 //更改好友备注错误
)

// GetFriendList 获取好友列表
func (h *HTTPServer) GetFriendList(c *gin.Context) {
	var user models.UserBaseInfo
	//获取用户数据
	user.Id, _ = utils2.StringToUint(c.Query("id"))
	//查询用户好友列表
	var friendList []tables2.Friends //存放列表数据
	friendList, _ = user.GetFriendList()
	utils2.Success(c, SUCCESS, friendList)
}

// GetRequestFriendList 获取好友请求添加列表(用于首次登录)
func (h *HTTPServer) GetRequestFriendList(c *gin.Context) {
	var user models.UserBaseInfo
	var friendList []tables2.RequestFriend
	var b bool
	user.Id, _ = utils2.StringToUint(c.Query("id"))
	friendList, b = user.GetRequestFriendList()
	if b {
		utils2.Success(c, SUCCESS, gin.H{
			"friend_list": friendList,
		})
	} else {
		utils2.Fail(c, ServerError, gin.H{
			"friend_list": friendList,
		})
	}
}

// RequestAddFriend 请求添加好友
func (h *HTTPServer) RequestAddFriend(c *gin.Context) {
	var user models.UserBaseInfo
	var fromId uint
	var toId uint
	fromId, _ = utils2.StringToUint(c.Query("from_id"))
	toId, _ = utils2.StringToUint(c.Query("to_id"))
	user.Id = fromId
	err := user.RequestAddFriend(fromId, toId)
	if err {
		utils2.Success(c, SUCCESS, gin.H{
			"from_id": fromId,
			"to_id":   toId,
		})
	} else {
		utils2.Fail(c, AlreadyFriend, gin.H{
			"from_id": fromId,
			"to_id":   toId,
		})
	}
}

// DisposeAddFriend 处理好友请求
func (h *HTTPServer) DisposeAddFriend(c *gin.Context) {
	var user models.UserBaseInfo
	var friend tables2.Friends
	var requestId uint
	var state uint8
	friend.UserOneID, _ = utils2.StringToUint(c.Query("from_id"))
	friend.UserTwoID, _ = utils2.StringToUint(c.Query("to_id"))
	requestId, _ = utils2.StringToUint(c.Query("request_id")) //请求好友id
	state, _ = utils2.StringToUint8(c.Query("state"))
	user.Id = friend.UserTwoID

	b, s := user.DisposeAddFriend(friend, requestId, state)
	if b {
		utils2.Success(c, SUCCESS, gin.H{
			"state":   s,
			"from_id": friend.UserOneID,
			"to_id":   friend.UserTwoID,
		})
	} else {
		utils2.Fail(c, DisposeAddFriendFail, gin.H{
			"state":   s,
			"from_id": friend.UserOneID,
			"to_id":   friend.UserTwoID,
		})
	}
}

func (h *HTTPServer) ChangeFriendNote(c *gin.Context) {
	var user models.UserBaseInfo
	friendID, _ := utils2.StringToUint(c.Query("friend_id"))
	note := c.Query("note")
	user.Id, _ = utils2.StringToUint(c.GetHeader("id"))
	if user.ChangeFriendNote(friendID, note) {
		utils2.Success(c, SUCCESS, gin.H{
			"friend_id": friendID,
			"note":      note,
		})
	} else {
		utils2.Fail(c, ChangeFriendNoteFail, gin.H{
			"err": "error",
		})
	}

}

// DeleteFriend 删除好友
func DeleteFriend(c *gin.Context) {}

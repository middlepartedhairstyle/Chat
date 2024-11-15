package service

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/utils"
)

const (
	CreateGroupFail        = 21201 //群创建失败
	AddGroupFail           = 21202 //添加群聊失败
	DisposeAddGroupFail    = 21203 //处理添加群聊失败
	GetCreateGroupListFail = 21204 //获取用户创建的全部群聊失败
	GetAllGroupListFail    = 21205 //获取用户加入的所有群聊失败
)

// CreateGroup 创建群聊
func CreateGroup(c *gin.Context) {
	var user models.UserBaseInfo
	id := c.GetHeader("id")
	groupName, _ := c.GetPostForm("group_name")
	user.Id, _ = utils.StringToUint(id)
	group, b := user.CreateGroup(groupName)
	if b {
		utils.Success(c, SUCCESS, gin.H{
			"create_group": group,
		})
	} else {
		utils.Fail(c, CreateGroupFail, "create group err")
	}
}

// GetCreateGroupList 获取用户创建的全部群聊
func GetCreateGroupList(c *gin.Context) {
	var user models.UserBaseInfo
	id := c.GetHeader("id")
	user.Id, _ = utils.StringToUint(id)
	groups, b := user.FindAllCreateGroup()
	if b {
		utils.Success(c, SUCCESS, gin.H{
			"group_list": groups,
		})
	} else {
		utils.Fail(c, GetCreateGroupListFail, "get create group list err")
	}
}

// GetAllGroupList 获取用户加入的全部群聊
func GetAllGroupList(c *gin.Context) {
	var user models.UserBaseInfo
	id := c.GetHeader("id")
	user.Id, _ = utils.StringToUint(id)
	groups, b := user.FindAllGroup()
	if b {
		utils.Success(c, SUCCESS, gin.H{
			"group_list": groups,
		})
	} else {
		utils.Fail(c, GetAllGroupListFail, "get group list err")
	}

}

// FindGroup 寻找群
func FindGroup(c *gin.Context) {
	var user models.UserBaseInfo
	info := c.Query("group_info")
	group := user.FindGroup(info)
	utils.Success(c, SUCCESS, gin.H{
		"group": group,
	})
}

// AddGroup 添加群聊
func AddGroup(c *gin.Context) {
	var user models.UserBaseInfo
	groupIDString := c.Query("group_id")
	user.Id, _ = utils.StringToUint(c.GetHeader("id"))
	groupID, _ := utils.StringToUint(groupIDString)
	result, t := user.AddGroup(groupID)
	//(t对应的含义)添加该群是否需要群主同意,0不需要,1需要······
	if result != nil {
		utils.Success(c, SUCCESS, gin.H{
			"type":   t,
			"result": result,
		})
	} else {
		utils.Fail(c, AddGroupFail, "add group err")
	}

}

// GetRequestGroupList 获取用户加群请求列表
func GetRequestGroupList(c *gin.Context) {}

// DisposeAddGroup 群主处理添加群的消息
func DisposeAddGroup(c *gin.Context) {
	var user models.UserBaseInfo
	groupIDString := c.Query("request_id")
	stateString := c.Query("state")
	groupID, _ := utils.StringToUint(groupIDString)
	state, _ := utils.StringToUint8(stateString)

	user.Id, _ = utils.StringToUint(c.GetHeader("id"))

	result, b := user.DisposeAddGroup(groupID, state)
	if b {
		utils.Success(c, SUCCESS, result)
	} else {
		utils.Fail(c, DisposeAddGroupFail, result)
	}
}

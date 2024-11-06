package service

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/utils"
)

// CreateGroup 创建群聊
func CreateGroup(c *gin.Context) {
	var user models.UserBaseInfo
	id := c.GetHeader("id")
	groupName, _ := c.GetPostForm("group_name")
	user.Id, _ = utils.StringToUint(id)
	group, b := user.CreateGroup(groupName)
	if b {
		utils.Success(c, "成功", group)
	} else {
		utils.Fail(c, "错误", group)
	}
}

// GetCreateGroupList 获取用户创建的全部群聊
func GetCreateGroupList(c *gin.Context) {
	var user models.UserBaseInfo
	id := c.GetHeader("id")
	user.Id, _ = utils.StringToUint(id)
	group := user.FindAllCreateGroup()
	utils.Success(c, "成功", group)
}

// GetAllGroupList 获取用户加入的全部群聊
func GetAllGroupList(c *gin.Context) {
	var user models.UserBaseInfo
	id := c.GetHeader("id")
	user.Id, _ = utils.StringToUint(id)
	groupUser := user.FindAllGroup()
	utils.Success(c, "成功", groupUser)
}

// FindGroup 寻找群
func FindGroup(c *gin.Context) {
	var user models.UserBaseInfo
	info := c.Query("group_info")
	groupUser := user.FindGroup(info)
	utils.Success(c, "成功", groupUser)
}

// AddGroup 添加群聊
func AddGroup(c *gin.Context) {
	var user models.UserBaseInfo
	groupIDString := c.Query("group_id")
	user.Id, _ = utils.StringToUint(c.GetHeader("id"))
	groupID, _ := utils.StringToUint(groupIDString)
	result := user.AddGroup(groupID)
	if result != nil {
		utils.Success(c, "成功", result)
	} else {
		utils.Fail(c, "失败", -1)
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
		utils.Success(c, "成功", result)
	} else {
		utils.Fail(c, "失败", result)
	}
}

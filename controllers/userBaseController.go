package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/service"
	"github.com/middlepartedhairstyle/HiWe/utils"
)

const (
	FAIL = "信息获取错误"
)

// RegisterController 用户注册控制器
func RegisterController(c *gin.Context) {
	var user models.UserBaseInfo
	//获取注册用户信息
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Fail(c, FAIL, err.Error())
	} else {
		//注册用户
		str, b := service.Register(&user)
		if b {
			//返回成功信息
			utils.Success(c, str, gin.H{
				"id":         user.ID,
				"created_at": user.CreatedAt,
				"email":      user.Email,
				"name":       user.Username,
				"token":      user.Token,
			})
		} else {
			//返回失败信息和用户邮箱
			utils.Fail(c, str, gin.H{
				"email": user.Email,
			})
		}
	}
}

// LoginController 用户登录控制器
func LoginController(c *gin.Context) {
	var user models.UserBaseInfo
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Fail(c, FAIL, err.Error())
	} else {
		str, b := service.Login(&user)
		if b {
			utils.Success(c, str, gin.H{
				"id":         user.ID,
				"created_at": user.CreatedAt,
				"email":      user.Email,
				"name":       user.Username,
				"token":      user.Token,
			})
		} else {
			utils.Fail(c, str, user.Email)
		}
	}
}

// SendCodeController 验证码发送控制器
func SendCodeController(c *gin.Context) {
	var user models.UserBaseInfo
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Fail(c, "fail", user)
	} else {
		utils.Success(c, "success", user)
	}
}

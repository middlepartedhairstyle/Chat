package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/service"
	"github.com/middlepartedhairstyle/HiWe/utils"
)

// RegisterController 用户注册控制器
func RegisterController(c *gin.Context) {
	var user models.UserBaseInfo
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Fail(c, "fail", user)
	} else {
		str, b := service.Register(&user)
		if b {
			utils.Success(c, str, user)
		} else {
			utils.Fail(c, str, user)
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

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/utils"
)

const (
	TokenNotUseful = 40001 //token无效
	ServerError    = 40004 //服务器错误
)

// CheckToken 核对用户token
func CheckToken(c *gin.Context) {
	var user models.UserBaseInfo
	user.Token = c.GetHeader("token")
	user.Id, _ = utils.StringToUint(c.GetHeader("id"))
	if user.CheckToken() {
		c.Next()
	} else {
		utils.Fail(c, TokenNotUseful, gin.H{
			"id": user.Id,
		})
		c.Abort()
	}
}

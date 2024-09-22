package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/controllers"
)

func Routers(router *gin.Engine) {
	router.POST("/register", controllers.RegisterController)
	router.POST("/login", controllers.LoginController)
	router.POST("/sendCode", controllers.SendCodeController)
	router.POST("/verifyCode", controllers.VerifyCodeController)

	//msg
	router.GET("/ChatWithFriend", controllers.ChatWithFriendController)

}

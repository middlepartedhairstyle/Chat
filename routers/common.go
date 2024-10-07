package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/middleware"
	"github.com/middlepartedhairstyle/HiWe/service"
)

func Routers(router *gin.Engine) {
	//用户基础功能
	router.POST("/register", middleware.RateLimiter("register"), service.Register)
	router.POST("/emailLogin", middleware.RateLimiter("emailLogin"), middleware.LimitLogin("emailLogin"), service.PassWordLogin)
	router.POST("/codeLogin", middleware.RateLimiter("codeLogin"), middleware.LimitLogin("codeLogin"), service.CodeLogin)
	router.POST("/sendCode", middleware.RateLimiter("sendCode"), service.SendCode)
	router.POST("/verifyCode", middleware.RateLimiter("verifyCode"), service.VerifyCode)
	//
	router.GET("getFriendList", middleware.RateLimiter("getFriendList"), middleware.CheckToken, service.GetFriendList)
	router.GET("getRequestAddFriendList", middleware.RateLimiter("getRequestAddFriendList"), middleware.CheckToken, service.GetRequestFriendList)
	router.POST("requestAddFriend", middleware.RateLimiter("requestAddFriend"), middleware.CheckToken, service.RequestAddFriend)
	router.POST("disposeAddFriend", middleware.RateLimiter("requestRemoveFriend"), middleware.CheckToken, service.DisposeAddFriend)

	//msg
	router.GET("/ChatWithFriend", middleware.CheckToken, service.Chat)

}

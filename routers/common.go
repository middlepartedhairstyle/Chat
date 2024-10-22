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
	//用户好友功能
	router.GET("getFriendList", middleware.RateLimiter("getFriendList"), middleware.CheckToken, service.GetFriendList)
	router.GET("getRequestAddFriendList", middleware.RateLimiter("getRequestAddFriendList"), middleware.CheckToken, service.GetRequestFriendList)
	router.POST("requestAddFriend", middleware.RateLimiter("requestAddFriend"), middleware.CheckToken, service.RequestAddFriend)
	router.POST("disposeAddFriend", middleware.RateLimiter("requestRemoveFriend"), middleware.CheckToken, service.DisposeAddFriend)
	//用户群功能
	router.POST("createGroup", middleware.RateLimiter("createGroup"), middleware.CheckToken, service.CreateGroup)
	router.POST("addGroup", middleware.RateLimiter("addGroup"), middleware.CheckToken, service.AddGroup)
	router.GET("findAllCreateGroup", middleware.RateLimiter("findAllCreateGroup"), middleware.CheckToken, service.GetCreateGroupList)
	router.GET("findAllGroup", middleware.RateLimiter("findAllGroup"), middleware.CheckToken, service.GetAllGroupList)
	router.GET("findGroup", middleware.RateLimiter("findGroup"), middleware.CheckToken, service.FindGroup)
	//msg
	router.GET("/ChatWithFriend", middleware.CheckToken, service.Chat)

}

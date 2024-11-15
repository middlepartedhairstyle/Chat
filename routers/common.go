package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/middleware"
	"github.com/middlepartedhairstyle/HiWe/service"
)

func Routers(router *gin.Engine) {
	//用户基础功能
	router.POST("/register", middleware.RateLimiter("rg"), service.Register)
	router.POST("/emailLogin", middleware.RateLimiter("el"), middleware.LimitLogin("el"), service.PassWordLogin)
	router.POST("/codeLogin", middleware.RateLimiter("cl"), middleware.LimitLogin("cl"), service.CodeLogin)
	router.POST("/sendCode", middleware.RateLimiter("sc"), service.SendCode)
	router.POST("/verifyCode", middleware.RateLimiter("vc"), service.VerifyCode)
	//用户好友功能
	router.GET("getFriendList", middleware.RateLimiter("gfl"), middleware.CheckToken, service.GetFriendList)
	router.GET("getRequestAddFriendList", middleware.RateLimiter("grafl"), middleware.CheckToken, service.GetRequestFriendList)
	router.POST("requestAddFriend", middleware.RateLimiter("raf"), middleware.CheckToken, service.RequestAddFriend)
	router.POST("disposeAddFriend", middleware.RateLimiter("rrf"), middleware.CheckToken, service.DisposeAddFriend)
	//用户群功能
	router.POST("createGroup", middleware.RateLimiter("cg"), middleware.CheckToken, service.CreateGroup)
	router.POST("addGroup", middleware.RateLimiter("ag"), middleware.CheckToken, service.AddGroup)
	router.POST("disposeAddGroup", middleware.RateLimiter("dag"), middleware.CheckToken, service.DisposeAddGroup)
	router.GET("findAllCreateGroup", middleware.RateLimiter("facg"), middleware.CheckToken, service.GetCreateGroupList)
	router.GET("findAllGroup", middleware.RateLimiter("fag"), middleware.CheckToken, service.GetAllGroupList)
	router.GET("findGroup", middleware.RateLimiter("fg"), middleware.CheckToken, service.FindGroup)
	//msg
	router.GET("/ChatWithFriend", middleware.CheckToken, service.Chat)

}

package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/middleware"
	"github.com/middlepartedhairstyle/HiWe/redis"
	"github.com/middlepartedhairstyle/HiWe/service"
)

const (
	register                = "rg"
	emailLogin              = "el"
	codeLogin               = "cl"
	sendCode                = "sc"
	verifyCode              = "vc"
	getFriendList           = "gfl"
	getRequestAddFriendList = "grafl"
	requestAddFriend        = "raf"
	disposeAddFriend        = "daf"
	createGroup             = "cg"
	addGroup                = "ag"
	disposeAddGroup         = "dag"
	findAllCreateGroup      = "facg"
	findAllGroup            = "fag"
	findGroup               = "fg"
	Chat                    = "c"
	ChangeUserProfilePhoto  = "cupp"
	GetUserProfilePhoto     = "gupp"
	ChangeUserName          = "cun"
	ChangeUserPassword      = "cup"
	ChangeUserEmail         = "cue"
)

func Routers(router *gin.Engine) {

	httpServer := service.NewHTTPServer(redis.Rdb)

	//用户基础功能
	router.POST("/register", middleware.RateLimiter(register), httpServer.Register)
	router.POST("/emailLogin", middleware.RateLimiter(emailLogin), middleware.LimitLogin(emailLogin), httpServer.PassWordLogin)
	router.POST("/codeLogin", middleware.RateLimiter(codeLogin), middleware.LimitLogin(codeLogin), httpServer.CodeLogin)
	router.POST("/sendCode", middleware.RateLimiter(sendCode), httpServer.SendCode)
	router.POST("/verifyCode", middleware.RateLimiter(verifyCode), httpServer.VerifyCode)
	router.POST("/changeUserName", middleware.RateLimiter(ChangeUserName), middleware.CheckToken, httpServer.ChangeUserName)
	router.POST("/changeUserPassword", middleware.RateLimiter(ChangeUserPassword), middleware.CheckToken, httpServer.ChangeUserPassword)
	router.POST("/changeUserEmail", middleware.RateLimiter(ChangeUserEmail), middleware.CheckToken, httpServer.ChangeUserEmail)
	//用户详细信息
	router.POST("/changeUserProfilePhoto", middleware.RateLimiter(ChangeUserProfilePhoto), middleware.LimitSizeMiddleware(1<<20), httpServer.ChangeUserProfilePhoto)
	router.GET("getUserProfilePhoto", middleware.RateLimiter(GetUserProfilePhoto), httpServer.GetUserProfilePhoto)
	//用户好友功能
	router.GET("getFriendList", middleware.RateLimiter(getFriendList), middleware.CheckToken, httpServer.GetFriendList)
	router.GET("getRequestAddFriendList", middleware.RateLimiter(getRequestAddFriendList), middleware.CheckToken, httpServer.GetRequestFriendList)
	router.POST("requestAddFriend", middleware.RateLimiter(requestAddFriend), middleware.CheckToken, httpServer.RequestAddFriend)
	router.POST("disposeAddFriend", middleware.RateLimiter(disposeAddFriend), middleware.CheckToken, httpServer.DisposeAddFriend)
	//用户群功能
	router.POST("createGroup", middleware.RateLimiter(createGroup), middleware.CheckToken, httpServer.CreateGroup)
	router.POST("addGroup", middleware.RateLimiter(addGroup), middleware.CheckToken, httpServer.AddGroup)
	router.POST("disposeAddGroup", middleware.RateLimiter(disposeAddGroup), middleware.CheckToken, httpServer.DisposeAddGroup)
	router.GET("findAllCreateGroup", middleware.RateLimiter(findAllCreateGroup), middleware.CheckToken, httpServer.GetCreateGroupList)
	router.GET("findAllGroup", middleware.RateLimiter(findAllGroup), middleware.CheckToken, httpServer.GetAllGroupList)
	router.GET("findGroup", middleware.RateLimiter(findGroup), middleware.CheckToken, httpServer.FindGroup)
	//msg
	router.GET("/chat", middleware.RateLimiter(Chat), middleware.CheckToken, httpServer.Chat)

}

package routers

import (
	"github.com/gin-gonic/gin"
	middleware2 "github.com/middlepartedhairstyle/HiWe/internal/middleware"
	"github.com/middlepartedhairstyle/HiWe/internal/redis"
	"github.com/middlepartedhairstyle/HiWe/internal/service"
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
	DeleteUser              = "du"
	ChangeFriendNote        = "cfn"
)

func Routers(router *gin.Engine) {

	httpServer := service.NewHTTPServer(redis.Rdb)

	//用户基础功能
	router.POST("/register", middleware2.RateLimiter(register), httpServer.Register)
	router.POST("/emailLogin", middleware2.RateLimiter(emailLogin), middleware2.LimitLogin(emailLogin), httpServer.PassWordLogin)
	router.POST("/codeLogin", middleware2.RateLimiter(codeLogin), middleware2.LimitLogin(codeLogin), httpServer.CodeLogin)
	router.POST("/sendCode", middleware2.RateLimiter(sendCode), httpServer.SendCode)
	router.POST("/verifyCode", middleware2.RateLimiter(verifyCode), httpServer.VerifyCode)
	router.POST("/changeUserName", middleware2.RateLimiter(ChangeUserName), middleware2.CheckToken, httpServer.ChangeUserName)
	router.POST("/changeUserPassword", middleware2.RateLimiter(ChangeUserPassword), middleware2.CheckToken, httpServer.ChangeUserPassword)
	router.POST("/changeUserEmail", middleware2.RateLimiter(ChangeUserEmail), middleware2.CheckToken, httpServer.ChangeUserEmail)
	router.POST("/deleteUser", middleware2.RateLimiter(DeleteUser), middleware2.CheckToken, httpServer.DeleteUser)
	//用户详细信息
	router.POST("/changeUserProfilePhoto", middleware2.RateLimiter(ChangeUserProfilePhoto), middleware2.LimitSizeMiddleware(1<<20), httpServer.ChangeUserProfilePhoto)
	router.GET("getUserProfilePhoto", middleware2.RateLimiter(GetUserProfilePhoto), httpServer.GetUserProfilePhoto)
	//用户好友功能
	router.GET("getFriendList", middleware2.RateLimiter(getFriendList), middleware2.CheckToken, httpServer.GetFriendList)
	router.GET("getRequestAddFriendList", middleware2.RateLimiter(getRequestAddFriendList), middleware2.CheckToken, httpServer.GetRequestFriendList)
	router.POST("requestAddFriend", middleware2.RateLimiter(requestAddFriend), middleware2.CheckToken, httpServer.RequestAddFriend)
	router.POST("disposeAddFriend", middleware2.RateLimiter(disposeAddFriend), middleware2.CheckToken, httpServer.DisposeAddFriend)
	router.POST("changeFriendNote", middleware2.RateLimiter(ChangeFriendNote), middleware2.CheckToken, httpServer.ChangeFriendNote)
	//用户群功能
	router.POST("createGroup", middleware2.RateLimiter(createGroup), middleware2.CheckToken, httpServer.CreateGroup)
	router.POST("addGroup", middleware2.RateLimiter(addGroup), middleware2.CheckToken, httpServer.AddGroup)
	router.POST("disposeAddGroup", middleware2.RateLimiter(disposeAddGroup), middleware2.CheckToken, httpServer.DisposeAddGroup)
	router.GET("findAllCreateGroup", middleware2.RateLimiter(findAllCreateGroup), middleware2.CheckToken, httpServer.GetCreateGroupList)
	router.GET("findAllGroup", middleware2.RateLimiter(findAllGroup), middleware2.CheckToken, httpServer.GetAllGroupList)
	router.GET("findGroup", middleware2.RateLimiter(findGroup), middleware2.CheckToken, httpServer.FindGroup)
	//msg
	router.GET("/chat", middleware2.RateLimiter(Chat), middleware2.CheckToken, httpServer.Chat)

}

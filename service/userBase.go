package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/mySQL"
	"github.com/middlepartedhairstyle/HiWe/redis"
	"github.com/middlepartedhairstyle/HiWe/utils"
	"time"
)

const (
	// EXISTS 用户创建
	EXISTS = "用户已存在"
	CREATE = "创建成功"

	// SUCCEED 登录部分
	SUCCEED = "登录成功"
	FAIL    = "密码错误"
	LOSS    = "未找到用户"
	//CodeFail 验证码部分
	CodeFail    = "验证码发送失败"
	CodeSuccess = "验证码发送成功"
	// VerifyFail 验证码验证部分
	VerifyFail    = "验证失败"
	VerifySuccess = "验证成功"
)

// Register 用户创建
func Register(c *gin.Context) {
	var user models.UserVerify
	err := c.ShouldBindJSON(&user)
	if err != nil {
		return
	}
	//校验验证码是否正确,防止绕过邮箱检查
	code, _ := redis.Rdb.Get(context.Background(), user.Email).Result()
	if user.Code == code {
		//判断用户是否已经注册
		if b, _ := user.EmailIsUser(); b {
			utils.Fail(c, EXISTS, gin.H{
				"username": user.Username,
				"email":    user.Email,
			})
		} else {
			//创建用户
			if user.CreateUser() {
				user.UserInfo() //获取注册信息
				utils.Success(c, CREATE, gin.H{
					"username":   user.Username,
					"email":      user.Email,
					"id":         user.Id,
					"token":      user.Token,
					"created_at": user.CreatedAt,
				})
			} else {
				utils.Fail(c, "创建失败", gin.H{
					"username": user.Username,
					"email":    user.Email,
				})
			}
		}

	} else {
		utils.Fail(c, "验证码错误", gin.H{
			"username": user.Username,
			"email":    user.Email,
		})
	}
	//删除验证码
	defer func() {
		redis.Rdb.Del(context.Background(), user.Email)
	}()
}

// PassWordLogin 用户密码登录
func PassWordLogin(c *gin.Context) {
	var user models.UserBaseInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		return
	}
	//确认是否为用户
	if b, _ := user.EmailIsUser(); b {
		if user.CheckPassword() {
			//登录次数限制清零
			redis.Rdb.Del(context.Background(), "emailLogin"+user.Email)
			//更新用户信息
			user.UpdateToken()
			//获取用户信息
			user.UserInfo()
			utils.Success(c, SUCCEED, gin.H{
				"username":   user.Username,
				"email":      user.Email,
				"token":      user.Token,
				"created_at": user.CreatedAt,
				"id":         user.Id,
			})
		} else {
			utils.Fail(c, FAIL, gin.H{
				"email": user.Email,
			})
		}
	} else {
		utils.Fail(c, LOSS, gin.H{
			"email": user.Email,
		})
	}
}

// CodeLogin 验证码登录
func CodeLogin(c *gin.Context) {
	var user models.UserVerify
	err := c.ShouldBindJSON(&user)
	if err != nil {
		return
	}
	//确认是否为用户
	if b, _ := user.EmailIsUser(); b {
		if user.CheckCode(user.Email) {
			//登录次数限制清零
			redis.Rdb.Del(context.Background(), "emailLogin"+user.Email)
			//更新用户信息
			user.UpdateToken()
			//获取用户信息
			user.UserInfo()
			utils.Success(c, SUCCEED, gin.H{
				"username":   user.Username,
				"email":      user.Email,
				"token":      user.Token,
				"created_at": user.CreatedAt,
				"id":         user.Id,
			})
		} else {
			utils.Fail(c, "验证码错误", gin.H{
				"email": user.Email,
			})
		}
	} else {
		utils.Fail(c, LOSS, gin.H{
			"email": user.Email,
		})
	}
}

// SendCode 发送验证码
func SendCode(c *gin.Context) {
	var user models.UserVerify
	err := c.ShouldBindJSON(&user)
	if err != nil {
		return
	}
	user.Code = utils.RandString()
	err = redis.Rdb.Set(context.Background(), user.Email, user.Code, time.Minute*5).Err()

	//邮箱发送验证码

	if err != nil {
		utils.Fail(c, CodeFail, gin.H{
			"email": user.Email,
		})
		return
	}
	utils.Success(c, CodeSuccess, gin.H{
		"email": user.Email, //测试代码（~~~）
		"code":  user.Code,
	})
}

// VerifyCode 验证验证码
func VerifyCode(c *gin.Context) {
	var user models.UserVerify
	err := c.ShouldBindJSON(&user)
	if err != nil {
		return
	}
	user.VerifyCode, _ = redis.Rdb.Get(context.Background(), user.Email).Result()
	fmt.Println(user.VerifyCode) //调试
	if user.VerifyCode == user.Code {
		redis.Rdb.Set(context.Background(), user.Email, user.VerifyCode, time.Minute*20)
		utils.Success(c, VerifySuccess, gin.H{
			"verify": true,
		})
	} else {
		utils.Fail(c, VerifyFail, gin.H{
			"verify": false,
		})
	}
}

// GetFriendList 获取好友列表
func GetFriendList(c *gin.Context) {
	var user models.UserBaseInfo
	//获取用户数据
	user.Id, _ = utils.StringToUint(c.Query("id"))
	//查询用户好友列表
	var friendList []mySQL.Friends //存放列表数据
	friendList, _ = user.GetFriendList()
	utils.Success(c, "成功", friendList)
}

// GetRequestFriendList 获取好友请求添加列表(用于首次登录)
func GetRequestFriendList(c *gin.Context) {
	var user models.UserBaseInfo
	var friendList []mySQL.RequestFriend
	var b bool
	user.Id, _ = utils.StringToUint(c.Query("id"))
	friendList, b = user.GetRequestFriendList()
	if b {
		utils.Success(c, "成功", friendList)
	} else {
		utils.Fail(c, "失败", friendList)
	}
}

// RequestAddFriend 添加好友
func RequestAddFriend(c *gin.Context) {
	var user models.UserBaseInfo
	var fromId uint
	var toId uint
	fromId, _ = utils.StringToUint(c.Query("from_id"))
	toId, _ = utils.StringToUint(c.Query("to_id"))
	user.Id = fromId
	err := user.RequestAddFriend(fromId, toId)
	if err {
		utils.Success(c, "成功", gin.H{
			"from_id": fromId,
			"to_id":   toId,
		})
	} else {
		utils.Fail(c, "失败", gin.H{
			"from_id": fromId,
			"to_id":   toId,
		})
	}
}

// DisposeAddFriend 处理好友请求
func DisposeAddFriend(c *gin.Context) {
	var user models.UserBaseInfo
	var friend mySQL.Friends
	var requestId uint
	var state uint8
	friend.UserOneID, _ = utils.StringToUint(c.Query("from_id"))
	friend.UserTwoID, _ = utils.StringToUint(c.Query("to_id"))
	requestId, _ = utils.StringToUint(c.Query("request_id")) //请求好友id
	state, _ = utils.StringToUint8(c.Query("state"))
	user.Id = friend.UserTwoID

	b, s := user.DisposeAddFriend(friend, requestId, state)
	if b {
		utils.Success(c, "成功", gin.H{
			"state":   s,
			"from_id": friend.UserOneID,
			"to_id":   friend.UserTwoID,
		})
	} else {
		utils.Fail(c, "失败", gin.H{
			"state": s,
			"id":    friend.UserTwoID,
		})
	}
}

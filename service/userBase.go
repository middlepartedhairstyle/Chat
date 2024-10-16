package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/models"
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

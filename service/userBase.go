package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/utils"
	"time"
)

const (
	SUCCESS        = 0     //返回消息为成功类型
	EXISTS         = 10001 //"用户已存在"
	CreateUserFail = 10002 //用户创建失败
	WrongPassword  = 10003 //"密码错误"
	NotFoundUser   = 10004 //"未找到用户"
	CodeSendFail   = 10005 //"验证码发送失败"
	VerifyFail     = 10006 //"验证失败"
	VerifyCodeFail = 10007 //验证码错误
	ServerError    = 40004 //服务器错误
	QueryNameError = 40005 //传参名字错误
)

type HTTPServer struct {
	redisDB *redis.Client
}

func NewHTTPServer(redisDB *redis.Client) *HTTPServer {
	return &HTTPServer{redisDB: redisDB}
}

// Register 用户创建
func (h *HTTPServer) Register(c *gin.Context) {
	var user models.UserVerify
	err := c.ShouldBindJSON(&user)
	if err != nil {
		return
	}
	//校验验证码是否正确,防止绕过邮箱检查
	code, _ := h.redisDB.Get(context.Background(), user.Email).Result()
	if user.Code == code {
		//判断用户是否已经注册
		if b, _ := user.EmailIsUser(); b {
			utils.Fail(c, EXISTS, gin.H{
				"err_msg": "用户已存在",
			})
		} else {
			//创建用户
			if user.CreateUser() {
				user.UserInfo() //获取注册信息
				utils.Success(c, SUCCESS, gin.H{
					"username":   user.Username,
					"email":      user.Email,
					"id":         user.Id,
					"token":      user.Token,
					"created_at": user.CreatedAt,
				})
			} else {
				utils.Fail(c, CreateUserFail, gin.H{
					"err_msg": "用户创建失败",
				})
			}
		}

	} else {
		utils.Fail(c, VerifyCodeFail, gin.H{
			"err_msg": "验证码错误",
		})
	}
	//删除验证码
	defer func() {
		h.redisDB.Del(context.Background(), user.Email)
	}()
}

// PassWordLogin 用户密码登录
func (h *HTTPServer) PassWordLogin(c *gin.Context) {
	var user models.UserBaseInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		return
	}
	//确认是否为用户
	if b, _ := user.EmailIsUser(); b {
		if user.CheckPassword() {
			//登录次数限制清零
			h.redisDB.Del(context.Background(), "el"+user.Email)
			//更新用户信息
			user.UpdateToken()
			//获取用户信息
			user.UserInfo()
			utils.Success(c, SUCCESS, gin.H{
				"username":   user.Username,
				"email":      user.Email,
				"token":      user.Token,
				"created_at": user.CreatedAt,
				"id":         user.Id,
			})
		} else {
			utils.Fail(c, WrongPassword, gin.H{
				"err_msg": "密码错误",
			})
		}
	} else {
		utils.Fail(c, NotFoundUser, gin.H{
			"err_msg": "用户不存在",
		})
	}
}

// CodeLogin 验证码登录
func (h *HTTPServer) CodeLogin(c *gin.Context) {
	var user models.UserVerify
	err := c.ShouldBindJSON(&user)
	if err != nil {
		return
	}
	//确认是否为用户
	if b, _ := user.EmailIsUser(); b {
		if user.CheckCode(user.Email) {
			//登录次数限制清零
			h.redisDB.Del(context.Background(), "el"+user.Email)
			//更新用户信息
			user.UpdateToken()
			//获取用户信息
			user.UserInfo()
			utils.Success(c, SUCCESS, gin.H{
				"username":   user.Username,
				"email":      user.Email,
				"token":      user.Token,
				"created_at": user.CreatedAt,
				"id":         user.Id,
			})
		} else {
			utils.Fail(c, VerifyCodeFail, gin.H{
				"err_msg": "验证码错误",
			})
		}
	} else {
		utils.Fail(c, NotFoundUser, gin.H{
			"err_msg": "用户不存在",
		})
	}
}

// SendCode 发送验证码
func (h *HTTPServer) SendCode(c *gin.Context) {
	var user models.UserVerify
	err := c.ShouldBindJSON(&user)
	if err != nil {
		return
	}
	user.Code = utils.RandString()
	err = h.redisDB.Set(context.Background(), user.Email, user.Code, time.Minute*5).Err()

	//邮箱发送验证码

	if err != nil {
		utils.Fail(c, CodeSendFail, gin.H{
			"err_msg": "验证码发送错误",
		})
		return
	}
	utils.Success(c, SUCCESS, gin.H{
		"send_code": true, //测试代码（~~~）
	})
}

// VerifyCode 验证验证码
func (h *HTTPServer) VerifyCode(c *gin.Context) {
	var user models.UserVerify
	err := c.ShouldBindJSON(&user)
	if err != nil {
		return
	}
	user.VerifyCode, _ = h.redisDB.Get(context.Background(), user.Email).Result()
	fmt.Println(user.VerifyCode) //调试
	if user.VerifyCode == user.Code {
		h.redisDB.Set(context.Background(), user.Email, user.VerifyCode, time.Minute*20)
		utils.Success(c, SUCCESS, gin.H{
			"verify": true,
		})
	} else {
		utils.Fail(c, VerifyFail, gin.H{
			"err_msg": "验证错误",
		})
	}
}

// ChangeUserName 更改用户名
func (h *HTTPServer) ChangeUserName(c *gin.Context) {
	var user models.UserBaseInfo
	ids := c.GetHeader("id")
	username, err := c.GetQuery("user_name")
	if !err {
		utils.Fail(c, QueryNameError, gin.H{
			"err_msg": "query name error",
		})
		return
	}
	user.Id, _ = utils.StringToUint(ids)
	user.Username = username
	if user.ChangeUserName() {
		utils.Success(c, SUCCESS, gin.H{
			"username": user.Username, //待改进
		})
	} else {
		utils.Fail(c, ServerError, gin.H{
			"err_msg": "data dispose err",
		})
	}

}

// ChangeUserPassword 更改用户密码
func (h *HTTPServer) ChangeUserPassword(c *gin.Context) {
	var user models.UserBaseInfo
	ids := c.GetHeader("id")
	err := c.ShouldBindJSON(&user)
	if err != nil {
		utils.Fail(c, QueryNameError, gin.H{
			"err_msg": "query name error",
		})
		return
	}
	user.Id, _ = utils.StringToUint(ids)
	if user.ChangePassword() {
		utils.Success(c, SUCCESS, gin.H{
			"token": user.Token,
		})
	} else {
		utils.Fail(c, ServerError, gin.H{
			"err_msg": "data dispose err",
		})
	}

}

// ChangeUserEmail 更改用户邮箱，需要传入新邮箱，新邮箱验证码
func (h *HTTPServer) ChangeUserEmail(c *gin.Context) {
	var code models.UserVerify
	ids := c.GetHeader("id")
	err := c.ShouldBindJSON(&code)
	if err != nil {
		utils.Fail(c, QueryNameError, gin.H{
			"err_msg": "query name error1",
		})
		return
	}
	code.Id, _ = utils.StringToUint(ids)
	if code.CheckCode(code.Email) {
		if code.ChangeEmail() {
			utils.Success(c, SUCCESS, gin.H{
				"email": code.Email,
				"token": code.Token,
			})
		} else {
			utils.Fail(c, QueryNameError, gin.H{
				"err_msg": "data dispose err2",
			})
		}
	} else {
		utils.Fail(c, VerifyCodeFail, gin.H{
			"err_msg": "data dispose err3",
		})
	}
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {}

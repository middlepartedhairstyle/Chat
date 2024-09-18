package service

import (
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/mySQL"
	"github.com/middlepartedhairstyle/HiWe/utils"
)

const (
	// EXISTS 用户创建
	EXISTS = "用户已存在"
	CREATE = "创建成功"

	// SUCCEED 登录部分
	SUCCEED = "登录成功"
	FAIL    = "登录失败"
	LOSS    = "未找到用户"
)

// Register 用户注册
func Register(user *models.UserBaseInfo) (string, bool) {
	//将用户注册密码加密存入数据库中，并创建token
	if user.PasswordMD5() && user.UpdateToken() {
		if user.CreatUser() {
			return CREATE, true
		} else {
			return EXISTS, false
		}
	} else {
		return EXISTS, false
	}
}

// Login 用户登录
func Login(user *models.UserBaseInfo) (string, bool) {
	if user.CheckPassword() {
		if user.UpdateToken() {
			if err := mySQL.DB.Table(mySQL.USERBASETABLE).Where("email=?", user.Email).Update("token", user.Token); err.Error == nil {
				return SUCCEED, true
			} else {
				return LOSS, false
			}
		}
		return FAIL, false
	}
	return LOSS, false
}

// SendCode 发送验证码
func SendCode(user *models.UserCaptcha) bool {
	if user.Email != "" {
		err, _ := utils.EmailSendCode(user.Email, user.Code)
		if !err {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

// VerifyCode 校验验证码
func VerifyCode(user *models.UserCaptcha) (string, bool) {
	str, err := user.VerifyCode()
	if err {
		return str, true
	} else {
		return str, false
	}
}

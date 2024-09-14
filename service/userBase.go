package service

import (
	"fmt"
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/utils"
)

const (
	EXISTS = "用户已存在"
	CREATE = "创建成功"
)

// Register 用户注册
func Register(user *models.UserBaseInfo) (string, bool) {
	if user.PasswordMD5() && user.UpdateToken() {
		fmt.Println(user.Password)
		if user.CreatUser() {
			return CREATE, true
		} else {
			return EXISTS, false
		}
	} else {
		return EXISTS, false
	}
}

// SendCode 发送验证码
func SendCode(user *models.UserBaseInfo) bool {
	if user.Email != "" {
		err, _ := utils.EmailSendCode(user.Email, utils.RandString())
		if !err {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

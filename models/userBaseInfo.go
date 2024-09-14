package models

import (
	"fmt"
	"github.com/middlepartedhairstyle/HiWe/mySQL"
	"github.com/middlepartedhairstyle/HiWe/utils"
	"gorm.io/gorm"
)

// UserBaseInfo 用户基础信息
type UserBaseInfo struct {
	gorm.Model
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Md5Num   string `json:"md5Num"`
	Token    string `json:"token"`
}

// Captcha 用户验证码
type Captcha struct {
	gorm.Model
	Code  uint32 `json:"code"`
	Email string `json:"email"`
}

// CreatUser 创建用户，将用户信息储存在数据库中
func (user *UserBaseInfo) CreatUser() bool {
	row := mySQL.DB.Table(mySQL.USERBASETABLE).Where("email = ?", user.Email).Select("1").Row()
	if row.Scan(nil) != nil {
		mySQL.DB.Table(mySQL.USERBASETABLE).Create(&user)
		return true
	} else {
		return false
	}
}

// PasswordMD5 生成MD5加密密码
func (user *UserBaseInfo) PasswordMD5() bool {
	if user.Password != "" {
		user.Md5Num = utils.RandString()
		(*user).Password = utils.MakePassword((*user).Password, (*user).Md5Num)
		return true
	} else {
		return false
	}
}

// UpdateToken 更新token
func (user *UserBaseInfo) UpdateToken() bool {
	if !utils.IsEmptyStruct(*user) {
		fmt.Println(user.Md5Num)
		user.Token = utils.MakeToken(user.Username, user.Password, user.Md5Num)
		return true
	} else {
		return false
	}
}

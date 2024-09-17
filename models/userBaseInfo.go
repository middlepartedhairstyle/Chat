package models

import (
	"database/sql"
	"errors"
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
	var exists int
	row := mySQL.DB.Table(mySQL.USERBASETABLE).Where("email = ?", user.Email).Select("1").Row()
	if err := row.Scan(&exists); err != nil && !errors.Is(err, sql.ErrNoRows) {
		// 处理扫描错误
		return false
	}

	if exists == 0 { // 如果用户不存在
		if err := mySQL.DB.Table(mySQL.USERBASETABLE).Create(&user).Error; err != nil {
			// 处理创建用户时的错误
			return false
		}
		return true
	} else {
		return false // 用户已存在
	}
}

// PasswordMD5 生成MD5加密密码
func (user *UserBaseInfo) PasswordMD5() bool {
	if user.Password != "" {
		(*user).Md5Num = utils.RandString()
		(*user).Password = utils.MakePassword((*user).Password, (*user).Md5Num)
		return true
	} else {
		return false
	}
}

// UpdateToken 更新token
func (user *UserBaseInfo) UpdateToken() bool {
	if !utils.IsEmptyStruct(*user) {
		user.Token = utils.MakeToken(user.Email, user.Password, utils.RandString())
		return true
	} else {
		return false
	}
}

// CheckPassword 查询用户邮箱密码
func (user *UserBaseInfo) CheckPassword() bool {
	var result struct {
		Password string
		Md5Num   string
	}
	err := mySQL.DB.Table(mySQL.USERBASETABLE).Where("email = ?", user.Email).Select("password,md5_num").Scan(&result)

	if err.Error != nil {
		return false
	}

	if utils.CheckPassword(user.Password, result.Md5Num, result.Password) {
		return true
	} else {
		return false
	}
}

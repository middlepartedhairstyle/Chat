package models

import (
	"context"
	"database/sql"
	"errors"
	"github.com/middlepartedhairstyle/HiWe/mySQL"
	"github.com/middlepartedhairstyle/HiWe/redis"
	"github.com/middlepartedhairstyle/HiWe/utils"
	"gorm.io/gorm"
	"time"
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
// 后期增加将通过了邮箱验证的邮箱存入其中，用户注册时从里面寻找，确定邮箱已近验证
type Captcha struct {
	gorm.Model
	Code  uint32 `json:"code"`
	Email string `json:"email"`
	Pass  bool   `json:"pass"`
}

// UserCaptcha 用户验证码
type UserCaptcha struct {
	Email string `json:"email"`
	Code  string `json:"code"`
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
		(*user).Password = utils.MakePasswordMd5((*user).Password, (*user).Md5Num)
		return true
	} else {
		return false
	}
}

// UpdateToken 更新token
func (user *UserBaseInfo) UpdateToken() bool {
	if !utils.IsEmptyStruct(*user) {
		user.Token = utils.MakeTokenMd5(user.Email, user.Password, utils.RandString())
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

	if utils.CheckPasswordMd5(user.Password, result.Md5Num, result.Password) {
		return true
	} else {
		return false
	}
}

// MakeVerifyCode 产生验证码
func (captcha *UserCaptcha) MakeVerifyCode() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute) // 创建带有超时的上下文
	defer cancel()                                                          // 确保在函数结束时取消上下文
	captcha.Code = utils.RandString()
	err := redis.Rdb.Set(ctx, captcha.Email, captcha.Code, time.Minute).Err()
	if err != nil {
		return false
	}
	return true
}

// VerifyCode 校验验证码
func (captcha *UserCaptcha) VerifyCode() (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute) // 创建带有超时的上下文
	defer cancel()                                                          // 确保在函数结束时取消上下文
	code, err := redis.Rdb.Get(ctx, captcha.Email).Result()
	if err != nil {
		return "验证码过期", false
	} else {
		if captcha.Code == code {
			return "验证成功", true
		} else {
			return "验证码错误", false
		}
	}
}

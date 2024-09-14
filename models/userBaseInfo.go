package models

import "gorm.io/gorm"

// UserBaseInfoTable 用户基础信息(数据库)
type UserBaseInfoTable struct {
	gorm.Model
	Email    string `gorm:"type:varchar(127);unique_index"`
	Username string `gorm:"type:varchar(20)"`
	Password string `gorm:"type:varchar(32)"`
	Md5Num   string `gorm:"type:varchar(5)"`
	Token    string `gorm:"type:varchar(32)"`
}

// CaptchaTable 用户验证码(数据库)
type CaptchaTable struct {
	gorm.Model
	Code         uint32 `gorm:"type:int(11)"`
	Email        string
	UserBaseInfo UserBaseInfoTable `gorm:"type:varchar(127);foreignKey:Email;references:Email"`
}

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

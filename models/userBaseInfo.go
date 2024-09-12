package models

import "gorm.io/gorm"

// UserBaseInfo 用户基础信息
type UserBaseInfo struct {
	gorm.Model
	Email    string `json:"email" gorm:"type:varchar(127);unique_index"`
	Username string `json:"username" gorm:"type:varchar(20)"`
	Password string `json:"password" gorm:"type:varchar(32)"`
	Md5Num   string `json:"md5Num" gorm:"type:varchar(5)"`
	Token    string `json:"token" gorm:"type:varchar(32)"`
}

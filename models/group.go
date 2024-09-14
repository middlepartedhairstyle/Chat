package models

import "gorm.io/gorm"

// GroupNumTable 用户群(数据库)
type GroupNumTable struct {
	gorm.Model
	UserBaseInfo  UserBaseInfoTable `gorm:"foreignKey:GroupLeaderID;references:ID"`
	GroupLeaderID uint
	GroupName     string `gorm:"type:varchar(25)"`
}

// GroupUserTable 用户群的用户(数据库)
type GroupUserTable struct {
	gorm.Model
	GroupNum     GroupNumTable     `gorm:"foreignKey:GroupID;references:ID"`
	UserBaseInfo UserBaseInfoTable `gorm:"foreignKey:UserID;references:ID"`
	GroupID      uint
	UserID       uint
	Level        uint8 `gorm:"type:tinyint(1)"`
}

// GroupNum 用户群
type GroupNum struct {
	gorm.Model
	GroupLeaderID uint   `json:"group_leader_id"`
	GroupName     string `json:"group_name"`
}

// GroupUser 用户群的用户
type GroupUser struct {
	gorm.Model
	GroupID uint  `json:"group_id"`
	UserID  uint  `json:"user_id"`
	Level   uint8 `json:"level"`
}

package models

import "gorm.io/gorm"

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

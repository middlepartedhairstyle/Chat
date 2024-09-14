package models

import "gorm.io/gorm"

// 消息类型
const (
	MessageTypeText uint8 = iota + 1 //1
	MessageTypeImage
	MessageTypeVoice
)

// UserMessageTable 用户消息(数据库)
type UserMessageTable struct {
	gorm.Model
	UserBaseInfoOne UserBaseInfoTable `gorm:"foreignKey:FromID;references:ID"`
	UserBaseInfo    UserBaseInfoTable `gorm:"foreignKey:ToID;references:ID"`
	FromID          uint
	ToID            uint
	MessageType     uint8  `gorm:"type:tinyint(1)"`
	Media           uint8  `gorm:"type:tinyint(1)"`
	Content         string `gorm:"type:text"`
}

// GroupMessageTable 群消息(数据库)
type GroupMessageTable struct {
	gorm.Model
	UserBaseInfo UserBaseInfoTable `gorm:"foreignKey:FromID;references:ID"`
	GroupNum     GroupNumTable     `gorm:"foreignKey:ToGroupID;references:ID"`
	FromID       uint
	ToGroupID    uint
	MessageType  uint8  `gorm:"type:tinyint(1)"`
	Media        uint8  `gorm:"type:tinyint(1)"`
	Content      string `gorm:"type:text"`
}

// UserMessage 用户消息
type UserMessage struct {
	gorm.Model
	FromID      uint   `json:"from_id"`
	ToID        uint   `json:"to_id"`
	MessageType uint8  `json:"message_type"`
	Media       uint8  `json:"media"`
	Content     string `json:"content"`
}

// GroupMessage 群消息
type GroupMessage struct {
	gorm.Model
	FromID      uint   `json:"from_id"`
	ToGroupID   uint   `json:"to_group_id"`
	MessageType uint8  `json:"message_type"`
	Media       uint8  `json:"media"`
	Content     string `json:"content"`
}

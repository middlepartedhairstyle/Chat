package models

import "gorm.io/gorm"

// 消息类型
const (
	MessageTypeText uint8 = iota + 1 //1
	MessageTypeImage
	MessageTypeVoice
)

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

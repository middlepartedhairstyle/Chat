package mySQL

import "gorm.io/gorm"

// FriendMessage 用户消息
type FriendMessage struct {
	gorm.Model
	FromID      uint   //发送者id
	FriendID    uint   //好友组队id
	MessageType uint8  `gorm:"type:tinyint(1)"` //消息类型
	Message     string `gorm:"type:text"`       //消息主体
}

// NewFriendMessage 创建好友消息类型
func NewFriendMessage(fromId uint, friendId uint, messageType uint8, message *string) *FriendMessage {
	return &FriendMessage{
		FromID:      fromId,
		FriendID:    friendId,
		MessageType: messageType,
		Message:     *message,
	}
}

// CreateFriendMessage 将好友消息放入数据库中
func (friendMessage *FriendMessage) CreateFriendMessage() bool {
	var count int64
	err := DB.Table(FriendT).Where("id=?", friendMessage.FriendID).Count(&count).Error
	if err != nil {
		return false
	}
	if count > 0 {
		err = DB.Table(FriendMessageT).Create(friendMessage).Error
		if err != nil {
			return false
		}
		return true
	}
	return false
}

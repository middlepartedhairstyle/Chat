package tables

import (
	"github.com/middlepartedhairstyle/HiWe/mySQL"
	"gorm.io/gorm"
)

// GroupMessage 群用户消息
type GroupMessage struct {
	gorm.Model
	FromID      uint   //发送者id
	GroupID     uint   //好友组队id
	MessageType uint8  `gorm:"type:tinyint(1)"` //消息类型
	Message     string `gorm:"type:text"`       //消息主体
}

func NewGroupMessage(fromId uint, groupId uint, messageType uint8, message *string) *GroupMessage {
	return &GroupMessage{
		FromID:      fromId,
		GroupID:     groupId,
		MessageType: messageType,
		Message:     *message,
	}
}

// CreateGroupMessage 将群消息放入数据库中
func (groupMessage *GroupMessage) CreateGroupMessage() bool {
	var count int64
	err := mySQL.DB.Table(mySQL.GroupNumT).Where("id=?", groupMessage.GroupID).Count(&count).Error
	if err != nil {
		return false
	}
	if count > 0 {
		err = mySQL.DB.Table(mySQL.GroupMessageT).Create(groupMessage).Error
		if err != nil {
			return false
		}
		return true
	}
	return false
}

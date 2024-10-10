package models

import (
	"fmt"
	"github.com/middlepartedhairstyle/HiWe/mySQL"
)

// 消息类型
const (
	MessageTypeText  uint8 = iota + 1 //文本
	MessageTypeImage                  //图片
	MessageTypeVoice                  //音频
)

const (
	MediaFriend uint8 = iota + 1 //好友消息
	MediaGroup                   //群消息
)

const (
	ChatWithFriend = "f" //订阅消息者管道，后期更改为用户具体管道
	ChatWithGroup  = "g"
)

const maxUserNum = 100
const maxGroupNum = 50

// UserChatMessage 用户消息
type UserChatMessage struct {
	FromID      uint   `json:"from_id"`      //用户ID
	ToID        uint   `json:"to_id"`        //好友ID,群ID,
	MessageType uint8  `json:"message_type"` //消息类型,如图片,文字等
	Media       uint8  `json:"media"`        //消息种类,如群消息和好友消息
	Message     string `json:"message"`      //消息主体
}

type Information interface {
	SetTopic(topic uint) string
	MessageDispose() bool
	SetConsumerID() uint
	GetInformation(opts ...string) map[string]interface{}
}

// GetInformation 获取消息信息
func (userMessage *UserChatMessage) GetInformation(opts ...string) map[string]interface{} {
	var info map[string]interface{} = make(map[string]interface{})
	for _, opt := range opts {
		switch opt {
		case "fromID":
			info["fromID"] = userMessage.FromID
		case "toID":
			info["toID"] = userMessage.ToID
		case "messageType":
			info["messageType"] = userMessage.MessageType
		case "media":
			info["media"] = userMessage.Media
		case "message":
			info["message"] = userMessage.Message
		}
	}
	return info
}

func NewUserChatMessage() *UserChatMessage {
	return &UserChatMessage{}
}

// SetTopic 设置消息话题发送点或获取点
func (userMessage *UserChatMessage) SetTopic(topic uint) string {
	switch userMessage.Media {
	case MediaFriend:
		var tpId = topic/maxUserNum + 1
		return fmt.Sprintf("%s%s%d", ChatWithFriend, "tp", tpId) //例如ftp1,ftp2
	case MediaGroup:
		var tpId = topic/maxGroupNum + 1
		return fmt.Sprintf("%s%s%d", ChatWithGroup, "tp", tpId) //例如gtp1,gtp2
	default:
		return ""
	}
}

// SetConsumerID 设置消费者id
func (userMessage *UserChatMessage) SetConsumerID() uint {
	switch userMessage.Media {
	//消息种类为好友消息
	case MediaFriend:
		var friend mySQL.Friends
		friend.ID = userMessage.ToID
		b := friend.FindTwoUserID()
		if b {
			if friend.UserOneID == userMessage.FromID {
				return friend.UserTwoID
			} else {
				return friend.UserOneID
			}
		} else {
			return 0
		}
	//消息种类为群消息
	case MediaGroup:
		return userMessage.FromID
	default:
		return 0
	}
}

// MessageDispose 消息处理
func (userMessage *UserChatMessage) MessageDispose() bool {

	switch userMessage.MessageType {
	//消息类型为文本
	case MessageTypeText:
		friendMessage := mySQL.NewFriendMessage(userMessage.ToID, userMessage.FromID, userMessage.MessageType, &(userMessage.Message))
		return friendMessage.CreateFriendMessage()
	//消息类型为图片
	case MessageTypeImage:
		return true
	//消息类型为语音
	case MessageTypeVoice:
		return true
	default:
		return false
	}
}

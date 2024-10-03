package models

import (
	"encoding/json"
	"strconv"
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

// UserMessage 用户消息
type UserMessage struct {
	FromID      uint64 `json:"from_id"`      //用户ID
	ToID        uint64 `json:"to_id"`        //好友,群ID
	MessageType uint8  `json:"message_type"` //消息类型，如图片，文字等
	Media       uint8  `json:"media"`        //消息种类，如群消息和好友消息
	Message     string `json:"message"`      //消息主体
}

func (userMessage *UserMessage) ToJson(data []byte) {
	err := json.Unmarshal(data, userMessage)
	if err != nil {
		return
	}
}

func (userMessage *UserMessage) FromJson() []byte {
	data, _ := json.Marshal(*userMessage)
	return data
}

func ToString(data []byte) string {
	return string(data)
}

// MessageChannel 消息频道
func (userMessage *UserMessage) MessageChannel() string {
	switch userMessage.Media {
	case MediaFriend:
		return ChatWithFriend + strconv.FormatUint(userMessage.FromID, 10)
	case MediaGroup:
		return ChatWithGroup + strconv.FormatUint(userMessage.FromID, 10)
	default:
		return ""
	}
}

// MessageDispose 消息处理
func (userMessage *UserMessage) MessageDispose() {
	switch userMessage.MessageType {
	case MessageTypeText:

	default:
		panic("unhandled default case")

	}
}

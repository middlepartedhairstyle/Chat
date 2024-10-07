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
	FromID      uint   `json:"from_id"`      //用户ID
	ToID        uint   `json:"to_id"`        //好友ID,群ID,
	MessageType uint8  `json:"message_type"` //消息类型,如图片,文字等
	Media       uint8  `json:"media"`        //消息种类,如群消息和好友消息
	Message     string `json:"message"`      //消息主体
}

// Unmarshal 将[]byte类型转换为UserMessage类型
func (userMessage *UserMessage) Unmarshal(data []byte) (UserMessage, error) {
	var msg UserMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return msg, err
	}
	return msg, nil
}

// Marshal 将UserMessage类型转换为[]byte类型
func (userMessage *UserMessage) Marshal(msg UserMessage) ([]byte, error) {
	data, err := json.Marshal(msg)
	return data, err
}

// ByteToString []byte转换为string类型
func ByteToString(data []byte) string {
	return string(data)
}

// MessageChannel 消息频道
func (userMessage *UserMessage) MessageChannel(media uint8, toId uint) string {
	switch media {
	//消息种类为好友消息
	case MediaFriend:
		return ChatWithFriend + strconv.FormatUint(uint64(toId), 10)
	//消息种类为群消息
	case MediaGroup:
		return ChatWithGroup + strconv.FormatUint(uint64(toId), 10)
	default:
		return ""
	}
}

// MessageDispose 消息处理
func (userMessage *UserMessage) MessageDispose() bool {
	switch userMessage.MessageType {
	//消息类型为文本
	case MessageTypeText:
		return true
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

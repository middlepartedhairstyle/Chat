package models

import (
	"encoding/json"
	"fmt"
)

const maxUser = 1000
const UserMessageBaseTopic = "u"
const UserMessageBaseGroup = "u"

// UserMessageBase 用户消息总体结构
type UserMessageBase struct {
	UserMessageTypes int8        `json:"user_message_types"` //1,用户好友请求确认类型;2,用户群请求类型等
	BaseMessage      interface{} `json:"base_message"`       //具体消息
}

type UserMessageBaseOpts func(*UserMessageBase)

func SetUserMessageTypes(userMessageTypes int8) UserMessageBaseOpts {
	return func(userMessageBase *UserMessageBase) {
		userMessageBase.UserMessageTypes = userMessageTypes
	}
}

func SetBaseMessage(baseMessage interface{}) UserMessageBaseOpts {
	return func(userMessageBase *UserMessageBase) {
		userMessageBase.BaseMessage = baseMessage
	}
}

func NewUserMessageBase(opts ...UserMessageBaseOpts) *UserMessageBase {
	userMessageBase := &UserMessageBase{
		UserMessageTypes: 0,
		BaseMessage:      nil,
	}
	for _, opt := range opts {
		opt(userMessageBase)
	}
	return userMessageBase
}

func (userMessageBase *UserMessageBase) Marshal() ([]byte, error) {
	return json.Marshal(userMessageBase)
}

func (userMessageBase *UserMessageBase) SetTopic(userId uint) string {
	var topic string
	topic = fmt.Sprintf("%s%s%v", UserMessageBaseTopic, "tp", userId/maxUser+1)
	return topic
}

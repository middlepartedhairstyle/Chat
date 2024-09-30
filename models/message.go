package models

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"sync"
)

// 消息类型
const (
	MessageTypeText  uint8 = iota + 1 //文本
	MessageTypeImage                  //图片
	MessageTypeVoice                  //音频
)

// UserMessage 用户消息
type UserMessage struct {
	FriendID    uint64 `json:"friend_id"`
	FromID      uint64 `json:"from_id"`
	ToID        uint64 `json:"to_id"`
	MessageType uint8  `json:"message_type"` //消息类型，如图片，文字等
	Message     string `json:"message"`
}

// GroupMessage 群消息
type GroupMessage struct {
	FromID      uint  `json:"from_id"`
	ToGroupID   uint  `json:"to_group_id"`
	MessageType uint8 `json:"message_type"` //消息类型，如图片，文字等
}

type Note struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

var ClientMap map[int]*Note = make(map[int]*Note)

var RwLocker sync.RWMutex

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

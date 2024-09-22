package models

import (
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
	FromID      uint   `json:"from_id"`
	ToID        uint   `json:"to_id"`
	MessageType uint8  `json:"message_type"` //消息类型，如图片，文字等
	Media       uint8  `json:"media"`        //发送消息类型，例如私聊，群聊等
	Content     string `json:"content"`
}

// GroupMessage 群消息
type GroupMessage struct {
	FromID      uint   `json:"from_id"`
	ToGroupID   uint   `json:"to_group_id"`
	MessageType uint8  `json:"message_type"` //消息类型，如图片，文字等
	Media       uint8  `json:"media"`        //发送消息类型，例如私聊，群聊等
	Content     string `json:"content"`
}

type Note struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

var ClientMap map[int]*Note = make(map[int]*Note)

var RwLocker sync.RWMutex

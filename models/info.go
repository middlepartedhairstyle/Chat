package models

import "encoding/json"

const (
	UserChatMessageType = iota + 1
)

type Info struct {
	Types uint8           `json:"type"`
	Data  json.RawMessage `json:"data"`
}

func NewInfo() *Info {
	return new(Info)
}

func (info *Info) CheckType() Information {
	switch info.Types {
	case UserChatMessageType:
		msg := NewUserChatMessage()
		err := json.Unmarshal(info.Data, msg)
		if err != nil {
			return nil
		}
		return msg
	default:
		return nil
	}
}

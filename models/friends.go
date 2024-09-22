package models

import (
	"fmt"
	"github.com/middlepartedhairstyle/HiWe/mySQL"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// Friend 好友
type Friend struct {
	Id        uint64 `json:"id"`
	UserID    uint64 `json:"user_id"`
	FriendID  uint64 `json:"friend_id"`
	UserToken string `json:"user_token"`
}

// FriendInfo 好友具体信息
type FriendInfo struct {
	Id       uint64         `json:"id"`
	CreateAt time.Time      `json:"create_at"`
	UpdateAt time.Time      `json:"update_at"`
	DeleteAt gorm.DeletedAt `json:"delete_at"`
	UserId   uint64         `json:"user_id"`
	FriendId uint64         `json:"friend_id"`
}

// CheckToken 检查token
func CheckToken(userId uint64, checkToken string) bool {
	var token string
	err := mySQL.DB.Table(mySQL.USERBASETABLE).Where("id=?", userId).Select("token").Scan(&token)
	if err.Error != nil {
		return false
	}
	if token == checkToken {
		return true
	}
	return false
}

// IsFriend 判断是否为好友
func (friend *Friend) IsFriend() bool {
	var info struct {
		User   uint64
		Friend uint64
	}
	err := mySQL.DB.Table(mySQL.USERFRIENDSTABLE).Where("id=?", friend.Id).Select("user,friend").Scan(&info)
	if err.Error != nil {
		fmt.Println(err.Error)
		return false
	}

	//确认是否为好友
	if friend.UserID == info.User && friend.FriendID == info.Friend {
		return true
	}
	return false
}

// Contrast 通道
func (friend *Friend) Contrast() string {
	if friend.UserID > friend.FriendID {
		return strconv.FormatUint(friend.UserID, 10) + strconv.FormatUint(friend.FriendID, 10)
	} else {
		return strconv.FormatUint(friend.FriendID, 10) + strconv.FormatUint(friend.UserID, 10)
	}
}

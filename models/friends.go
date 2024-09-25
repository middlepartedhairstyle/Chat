package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/middlepartedhairstyle/HiWe/mySQL"
	"gorm.io/gorm"
	"strconv"
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
	gorm.Model
	UserId   uint64 `json:"user_id"`
	FriendId uint64 `json:"friend_id"`
}

// IsFriend 判断是否为好友
func (friend *Friend) IsFriend() bool {
	var info struct {
		User   uint64
		Friend uint64
	}
	err := mySQL.DB.Table(mySQL.USERFRIENDSTABLE).Where("id=?", friend.Id).Select("user_id,friend_id").Scan(&info)
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

// AddFriend 使用用户id添加好友
func (friendInfo *FriendInfo) AddFriend() (bool, string) {

	//查看是否为好友
	var exists int
	//交换好友id顺序，以较大的id为第一个id放入数据库中
	if friendInfo.FriendId < friendInfo.UserId {
		friendInfo.FriendId, friendInfo.UserId = friendInfo.UserId, friendInfo.FriendId
	}

	row := mySQL.DB.Table(mySQL.USERFRIENDSTABLE).Where("user_id=? AND friend_id=?", friendInfo.UserId, friendInfo.FriendId).Select("1").Row()
	err := row.Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, "服务器错误"
	}

	//不为好友，创建好友
	if exists == 0 {
		err = mySQL.DB.Table(mySQL.USERFRIENDSTABLE).Create(friendInfo).Error
		if err != nil {
			return false, "服务器错误"
		}
		return true, "好友添加成功"
	}
	return false, "已为好友"
}

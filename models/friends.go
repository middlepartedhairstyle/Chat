package models

import (
	"github.com/middlepartedhairstyle/HiWe/mySQL"
)

const (
	Success string = "好友添加成功"
	Failed1 string = "添加失败"
)

// Friend 好友
type Friend struct {
	Id        uint   `json:"id"`         //用户与好友的组队id
	UserID    uint64 `json:"user_id"`    //用户自己的di
	FriendID  uint64 `json:"friend_id"`  //好友的id
	UserToken string `json:"user_token"` //用户的id
}

// IsFriend 判断是否为好友
func (friend *Friend) IsFriend() bool {
	var f mySQL.UserFriendsTable
	f.FriendID = friend.FriendID
	f.UserID = friend.UserID
	f.ID = friend.Id
	return f.IsFriend()
}

// AddFriend 使用用户id添加好友
func (friend *Friend) AddFriend() (bool, string) {
	var f mySQL.UserFriendsTable
	f.FriendID = friend.FriendID
	f.UserID = friend.UserID
	err := f.AddFriend()
	if err {
		return true, Success
	} else {
		return false, Failed1
	}
}

// ConfirmAddFriend 用于请求好友添加，后另一方确认是否成为好友
func (friend *Friend) ConfirmAddFriend() (bool, string) {
	return false, Failed1
}

// RequestAddFriend 用于请求成为好友,将请求存入数据库和redis
func (friend *Friend) RequestAddFriend() bool {
	var f mySQL.RequestAddFriend
	f.FromRequestID = friend.UserID
	f.ToRequestID = friend.FriendID
	//存入数据库
	err := f.RequestAddFriend()
	if !err {
		return false
	}

	//存入redis
	return true
}

// DeleteFriend 用于删除好友，不需要好友确认
func (friend *Friend) DeleteFriend() (bool, string) {
	return false, Failed1
}

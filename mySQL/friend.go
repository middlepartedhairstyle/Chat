package mySQL

import (
	"database/sql"
	"errors"
)

// IsFriend 判断是否为好友
func (userFriend *UserFriendsTable) IsFriend() bool {
	var info struct {
		UserID   uint64
		FriendID uint64
	}

	err := DB.Table(USERFRIENDSTABLE).Where("id=?", userFriend.ID).Select("user_id,friend_id").Scan(&info)
	if err.Error != nil {
		return false
	}

	//确认是否为好友
	if (userFriend.UserID == info.UserID && userFriend.FriendID == info.FriendID) || (userFriend.UserID == info.FriendID && userFriend.FriendID == info.UserID) {
		return true
	}
	return false
}

// AddFriend 添加好友
func (userFriend *UserFriendsTable) AddFriend() bool {
	//查看是否为好友
	var exists int
	//交换好友id顺序，以较大的id为第一个id放入数据库中
	if userFriend.FriendID < userFriend.UserID {
		userFriend.FriendID, userFriend.UserID = userFriend.UserID, userFriend.FriendID
	}

	row := DB.Table(USERFRIENDSTABLE).Where("user_id=? AND friend_id=?", userFriend.UserID, userFriend.FriendID).Select("1").Row()
	err := row.Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false
	}

	//不为好友，创建好友
	if exists == 0 {
		err = DB.Table(USERFRIENDSTABLE).Create(userFriend).Error
		if err != nil {
			return false
		}
		return true
	}
	return false
}


package mySQL

import (
	"database/sql"
	"errors"
	"fmt"
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

func SelectAllFriend(id uint) []uint {
	rows, err := DB.Table(USERFRIENDSTABLE).Where("user_id=? OR friend_id=?", id, id).Select("id").Rows()
	if err != nil {
		return nil
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {

		}
	}(rows)

	var friends []uint // 用于存储朋友 ID

	// 循环读取每一行数据
	for rows.Next() {
		var friendID uint
		err = rows.Scan(&friendID) // 将查询到的 ID 扫描到 friendID 变量中
		if err != nil {
			fmt.Println("读取行时出错:", err)
			continue // 如果读取出错，继续处理下一行
		}
		friends = append(friends, friendID) // 将 ID 添加到切片中
	}

	// 检查循环是否出错
	if err = rows.Err(); err != nil {
		fmt.Println("遍历 rows 时出错:", err)
		return nil
	}

	return friends // 返回所有朋友的 ID
}

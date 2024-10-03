package mySQL

import "gorm.io/gorm"

type Friends struct {
	gorm.Model
	UserOneID    uint   `gorm:"type:int(11);not null;index:idx_friends,unique"` //用户一id
	UserTwoID    uint   `gorm:"type:int(11);not null;index:idx_friends,unique"` //用户二id
	Relationship string `gorm:"type:varchar(10)"`                               //两者关系
	NoteOne      string `gorm:"type:varchar(20)"`                               //UserOne个UserTwo的备注
	NoteTwo      string `gorm:"type:varchar(20)"`                               //UserTwo个UserOne的备注
}

// GetFriendList 获取好友列表
func (friend *Friends) GetFriendList(userId uint) ([]Friends, bool) {
	var friendList []Friends
	// 使用 GORM 的 Find 方法来查询
	if err := DB.Table(Friend).Where("user_one_id = ? OR user_two_id = ?", userId, userId).Find(&friendList).Error; err != nil {
		return nil, false
	}
	// 返回好友列表
	return friendList, true
}

// AddFriend 添加好友
func (friend *Friends) AddFriend() bool {
	err := DB.Table(Friend).Create(&friend).Error
	if err != nil {
		return false
	}
	return true
}

// IsFriend 判断是否为好友
func (friend *Friends) IsFriend() bool {
	var count int64
	err := DB.Table(Friend).Where("(user_one_id = ? AND user_two_id = ?) OR (user_one_id = ? AND user_two_id = ?)", friend.UserOneID, friend.UserTwoID, friend.UserTwoID, friend.UserOneID).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

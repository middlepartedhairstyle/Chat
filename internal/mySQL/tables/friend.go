package tables

import (
	mySQL2 "github.com/middlepartedhairstyle/HiWe/internal/mySQL"
	"gorm.io/gorm"
)

type Friends struct {
	gorm.Model
	UserOneID    uint   `gorm:"type:int(11);not null;index:idx_friends,unique"` //用户一id
	UserTwoID    uint   `gorm:"type:int(11);not null;index:idx_friends,unique"` //用户二id
	Relationship string `gorm:"type:varchar(10)"`                               //两者关系
	NoteOne      string `gorm:"type:varchar(20)"`                               //UserOne给UserTwo的备注
	NoteTwo      string `gorm:"type:varchar(20)"`                               //UserTwo给UserOne的备注
}

type FriendOpt func(*Friends)

func SetFriendID(id uint) FriendOpt {
	return func(friends *Friends) {
		friends.ID = id
	}
}

func SetUserOneID(userOneID uint) FriendOpt {
	return func(friends *Friends) {
		friends.UserOneID = userOneID
	}
}
func SetUserTwoID(userTwoID uint) FriendOpt {
	return func(friends *Friends) {
		friends.UserTwoID = userTwoID
	}
}
func SetFriendRelationship(relationship string) FriendOpt {
	return func(friends *Friends) {
		friends.Relationship = relationship
	}
}
func SetNoteOne(noteOne string) FriendOpt {
	return func(friends *Friends) {
		friends.NoteOne = noteOne
	}
}
func SetNoteTwo(noteTwo string) FriendOpt {
	return func(friends *Friends) {
		friends.NoteTwo = noteTwo
	}
}

func NewFriend(opts ...FriendOpt) *Friends {
	friends := &Friends{
		UserOneID:    0,
		UserTwoID:    0,
		Relationship: "",
		NoteOne:      "",
		NoteTwo:      "",
	}
	for _, opt := range opts {
		opt(friends)
	}
	return friends
}

// GetFriendList 获取好友列表
func (friend *Friends) GetFriendList(userId uint) ([]Friends, bool) {
	var friendList []Friends
	// 使用 GORM 的 Find 方法来查询
	if err := mySQL2.DB.Table(mySQL2.FriendT).Where("user_one_id = ? OR user_two_id = ?", userId, userId).Find(&friendList).Error; err != nil {
		return nil, false
	}
	// 返回好友列表
	return friendList, true
}

// FindAllFriendId 获取所有好友id
func FindAllFriendId(userId uint) []uint {
	var friendList []uint
	if err := mySQL2.DB.Table(mySQL2.FriendT).Where("user_one_id = ? OR user_two_id = ?", userId, userId).Select("id").Find(&friendList).Error; err != nil {
		return nil
	}
	return friendList
}

// AddFriend 添加好友
func (friend *Friends) AddFriend() bool {
	err := mySQL2.DB.Table(mySQL2.FriendT).Create(&friend).Error
	if err != nil {
		return false
	}
	return true
}

// IsFriend 判断是否为好友
func (friend *Friends) IsFriend() bool {
	var count int64
	err := mySQL2.DB.Table(mySQL2.FriendT).Where("(user_one_id = ? AND user_two_id = ?) OR (user_one_id = ? AND user_two_id = ?)", friend.UserOneID, friend.UserTwoID, friend.UserTwoID, friend.UserOneID).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (friend *Friends) IsFriendUseFriendID() (uint, bool) {
	var count int64
	var toFriendID uint
	err := mySQL2.DB.Table(mySQL2.FriendT).Where("(id = ? AND user_two_id = ?) OR (user_one_id = ? AND id = ?)", friend.ID, friend.UserOneID, friend.UserOneID, friend.ID).Count(&count).Error
	if err != nil {
		return 0, false
	}
	if count > 0 {
		err = mySQL2.DB.Table(mySQL2.FriendT).Where("(id = ? AND user_two_id = ?) OR (user_one_id = ? AND id = ?)", friend.ID, friend.UserOneID, friend.UserOneID, friend.ID).Select("user_one_id").Scan(&toFriendID).Error
		if err != nil {
			return 0, false
		}
		if toFriendID != friend.UserOneID {
			return toFriendID, true
		} else {
			err = mySQL2.DB.Table(mySQL2.FriendT).Where("(id = ? AND user_two_id = ?) OR (user_one_id = ? AND id = ?)", friend.ID, friend.UserOneID, friend.UserOneID, friend.ID).Select("user_two_id").Scan(&toFriendID).Error
			if err != nil {
				return 0, false
			}
			return toFriendID, true
		}
	} else {
		return 0, false
	}
}

// FindTwoUserID 寻找两好友的id
func (friend *Friends) FindTwoUserID() bool {
	if err := mySQL2.DB.Table(mySQL2.FriendT).Where("id = ?", friend.ID).Select("user_one_id,user_two_id").Find(&friend).Error; err != nil {
		return false
	}
	return true
}

// ChangeNote 修改好友备注
func (friend *Friends) ChangeNote(userID uint, note string) bool {
	var count int64
	err := mySQL2.DB.Table(mySQL2.FriendT).Where("id = ?", friend.ID).Count(&count).Error
	if err != nil {
		return false
	}
	if count > 0 {
		err = mySQL2.DB.Table(mySQL2.FriendT).Where("id = ? and user_one_id = ?", friend.ID, userID).UpdateColumn("note_two", note).Error
		err = mySQL2.DB.Table(mySQL2.FriendT).Where("id = ? and user_two_id = ?", friend.ID, userID).UpdateColumn("note_one", note).Error
		if err != nil {
			return false
		}
		return true
	}
	return false
}

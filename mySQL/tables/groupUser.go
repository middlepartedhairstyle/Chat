package tables

import (
	"github.com/middlepartedhairstyle/HiWe/mySQL"
	"gorm.io/gorm"
	"time"
)

// GroupUser 用户群的用户(数据库)
type GroupUser struct {
	gorm.Model
	GroupID      uint   //群组id
	UserID       uint   //用户id
	Note         string `gorm:"type:varchar(20)"` //用户给群的备注
	Level        uint8  `gorm:"type:tinyint(1)"`  //用户在群中的等级
	Relationship uint8  `gorm:"type:tinyint(1)"`  //用户在群中的关系如,1群主,2管理员,3群员等
}

type GroupUserInfo struct {
	ID           uint
	GroupID      uint
	UserID       uint
	Level        uint8
	Relationship string
}

type GroupUserOpt func(*GroupUser)

func SetCreateAt(createAt time.Time) GroupUserOpt {
	return func(groupUser *GroupUser) {
		groupUser.CreatedAt = createAt
	}
}
func SetUpdateAt(updateAt time.Time) GroupUserOpt {
	return func(groupUser *GroupUser) {
		groupUser.UpdatedAt = updateAt
	}
}

func SetDeleteAt(deleteAt gorm.DeletedAt) GroupUserOpt {
	return func(groupUser *GroupUser) {
		groupUser.DeletedAt = deleteAt
	}
}

// SetGroupID 设置GroupID
func SetGroupID(groupID uint) GroupUserOpt {
	return func(groupUser *GroupUser) {
		groupUser.GroupID = groupID
	}
}

// SetUserID 设置UserID
func SetUserID(userID uint) GroupUserOpt {
	return func(groupUser *GroupUser) {
		groupUser.UserID = userID
	}
}

// SetNote 设置备注
func SetNote(note string) GroupUserOpt {
	return func(groupUser *GroupUser) {
		groupUser.Note = note
	}
}

// SetLevel 设置等级
func SetLevel(level uint8) GroupUserOpt {
	return func(groupUser *GroupUser) {
		groupUser.Level = level
	}
}

// SetRelationship 设置用户在群中的关系如,群主,管理员,群员等(安全要求中等)
func SetRelationship(relationship uint8) GroupUserOpt {
	return func(groupUser *GroupUser) {
		groupUser.Relationship = relationship
	}
}

// NewGroupUser 新建一个groupUser
func NewGroupUser(opts ...GroupUserOpt) *GroupUser {
	groupUser := &GroupUser{
		GroupID:      0,
		UserID:       0,
		Note:         "",
		Level:        0,
		Relationship: 3,
	}
	for _, opt := range opts {
		opt(groupUser)
	}
	return groupUser
}

// CreateGroupUser 创建群聊用户
func (groupUser *GroupUser) CreateGroupUser() bool {
	var count1 int64
	var count2 int64
	err := mySQL.DB.Table(mySQL.GroupUserT).Where("group_id=? and user_id=?", groupUser.GroupID, groupUser.UserID).Count(&count1).Error
	if err != nil {
		return false
	}
	err = mySQL.DB.Table(mySQL.GroupNumT).Where("id=?", groupUser.GroupID).Count(&count2).Error
	if err != nil {
		return false
	}
	if count1 <= 0 && count2 > 0 {
		err = mySQL.DB.Table(mySQL.GroupUserT).Create(groupUser).Error
		if err != nil {
			return false
		}
		return true
	} else {
		return false
	}

}

// FindAllGroup 寻找用户所有群聊包括创建的和加入的群聊
func (groupUser *GroupUser) FindAllGroup() ([]GroupUser, bool) {
	var groupUsers []GroupUser
	err := mySQL.DB.Table(mySQL.GroupUserT).Where("user_id = ?", groupUser.UserID).Find(&groupUsers).Error
	if err != nil {
		return []GroupUser{}, false
	}
	return groupUsers, true
}

// FindAllGroupID 寻找用户所有群聊包括创建的和加入的群聊
func (groupUser *GroupUser) FindAllGroupID() []uint {
	var groupIDs []uint
	err := mySQL.DB.Table(mySQL.GroupUserT).Where("user_id = ?", groupUser.UserID).Select("group_id").Scan(&groupIDs).Error
	if err != nil {
		return nil
	}
	return groupIDs
}

// FindGroupUserID 寻找用户所有群聊用户ID
func (groupUser *GroupUser) FindGroupUserID() []uint {
	var groupUserIDs []uint
	err := mySQL.DB.Table(mySQL.GroupUserT).Where("user_id = ?", groupUser.UserID).Select("id").Scan(&groupUserIDs).Error
	if err != nil {
		return nil
	}
	return groupUserIDs
}

// FindAllGroupUser 寻找用户加入群聊的所有成员
func (groupUser *GroupUser) FindAllGroupUser() []GroupUserInfo {

	var groupUsers []GroupUserInfo
	err := mySQL.DB.Table(mySQL.GroupUserT).Where("group_id = ?", groupUser.GroupID).Select("id,group_id,user_id,level,relationship").Find(&groupUsers).Error
	if err != nil {
		return nil
	}
	return groupUsers
}

// IsGroupUserGetID 判断是否为群用户并返回群用户id(groupUserID)
func (groupUser *GroupUser) IsGroupUserGetID() bool {
	var count int64
	err := mySQL.DB.Table(mySQL.GroupUserT).Where("group_id=? and user_id=?", groupUser.GroupID, groupUser.UserID).Count(&count).Error
	if err != nil {
		return false
	}
	if count > 0 {
		err = mySQL.DB.Table(mySQL.GroupUserT).Where("group_id=? and user_id=?", groupUser.GroupID, groupUser.UserID).Select("id").Scan(&groupUser.ID).Error
		if err != nil {
			return false
		}
		return true
	}
	return false
}

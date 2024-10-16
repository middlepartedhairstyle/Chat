package mySQL

import (
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
	Relationship string `gorm:"type:varchar(10)"` //用户在群中的关系如,群主,管理员,群员等
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
func SetRelationship(relationship string) GroupUserOpt {
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
		Relationship: "",
	}
	for _, opt := range opts {
		opt(groupUser)
	}
	return groupUser
}

// CreateGroupUser 创建群聊用户
func (groupUser *GroupUser) CreateGroupUser() bool {
	err := DB.Table(GroupUserT).Create(groupUser).Error
	if err != nil {
		return false
	}
	return true
}

// FindAllGroup 寻找用户所有群聊包括创建的和加入的群聊
func (groupUser *GroupUser) FindAllGroup() []GroupUser {
	var groupUsers []GroupUser
	err := DB.Table(GroupUserT).Where("user_id = ?", groupUser.UserID).Find(&groupUsers).Error
	if err != nil {
		return []GroupUser{}
	}
	return groupUsers
}

package mySQL

import (
	"gorm.io/gorm"
)

// GroupNum 用户群(数据库)
type GroupNum struct {
	gorm.Model
	GroupLeaderID uint   //群主id
	GroupName     string `gorm:"type:varchar(25)"` //群名
	Visible       bool   `gorm:"type:tinyint(1)"`  //该群是否可以被搜索
	Verify        uint8  `gorm:"type:tinyint(1)"`  //添加该群是否需要群主同意,0需要,1不需要······
}

const maxGroupNum = 100

type GroupNumOpt func(*GroupNum)

// SetGroupLeaderID 设置群主id
func SetGroupLeaderID(groupLeaderID uint) GroupNumOpt {
	return func(groupNum *GroupNum) {
		groupNum.GroupLeaderID = groupLeaderID
	}
}

// SetGroupName 设置群名
func SetGroupName(groupName string) GroupNumOpt {
	return func(groupNum *GroupNum) {
		groupNum.GroupName = groupName
	}
}

// SetVisible 设置该群是否可以被搜索
func SetVisible(visible bool) GroupNumOpt {
	return func(groupNum *GroupNum) {
		groupNum.Visible = visible
	}
}

// SetVerify 设置添加该群是否需要群主同意
func SetVerify(verify uint8) GroupNumOpt {
	return func(groupNum *GroupNum) {
		groupNum.Verify = verify
	}
}

// NewGroupNum 新建群结构体
func NewGroupNum(opts ...GroupNumOpt) *GroupNum {
	groupNum := &GroupNum{
		GroupLeaderID: 0,
		GroupName:     "",
		Visible:       false,
		Verify:        0,
	}
	for _, opt := range opts {
		opt(groupNum)
	}
	return groupNum
}

// CreateGroup 创建群聊
func (group *GroupNum) CreateGroup() bool {
	//检查是否超过最大群数限制
	var count int64
	err := DB.Table(GroupNumT).Where("group_leader_id = ?", group.GroupLeaderID).Count(&count).Error
	if count < maxGroupNum {
		//若小于最大群数则创建群
		err = DB.Table(GroupNumT).Create(group).Error
		if err != nil {
			return false
		}
		//创建成功将信息同时存入group_user_tables
		err = DB.Table(GroupNumT).Where("group_leader_id = ?", group.GroupLeaderID).First(&group).Error
		groupUser := NewGroupUser(SetCreateAt(group.CreatedAt), SetUpdateAt(group.UpdatedAt), SetGroupID(group.ID), SetUserID(group.GroupLeaderID), SetRelationship(group.GroupName))
		if groupUser.CreateGroupUser() {
			return true
		} else {
			return false
		}
	} else {
		//若大于最大群数则取消创建群
		return false
	}
}

// FindAllCreateGroup 寻找自己建立的全部群聊
func (group *GroupNum) FindAllCreateGroup() []GroupNum {
	var groups []GroupNum
	err := DB.Table(GroupNumT).Where("group_leader_id = ?", group.GroupLeaderID).Find(&groups).Error
	if err != nil {
		return nil
	}
	return groups
}

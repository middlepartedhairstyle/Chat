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
	Verify        uint8  `gorm:"type:tinyint(1)"`  //添加该群是否需要群主同意,0不需要,1需要······
}

const maxGroupNum = 100

type GroupNumOpt func(*GroupNum)

// SetGroupNumID 设置groupID
func SetGroupNumID(groupID uint) GroupNumOpt {
	return func(g *GroupNum) {
		g.ID = groupID
	}
}

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
		Visible:       true,
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
		groupUser := NewGroupUser(SetCreateAt(group.CreatedAt), SetUpdateAt(group.UpdatedAt), SetGroupID(group.ID), SetUserID(group.GroupLeaderID), SetNote(group.GroupName), SetRelationship(1))
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
func (group *GroupNum) FindAllCreateGroup() ([]GroupNum, bool) {
	var groups []GroupNum
	err := DB.Table(GroupNumT).Where("group_leader_id = ?", group.GroupLeaderID).Find(&groups).Error
	if err != nil {
		return nil, false
	}
	return groups, true
}

// UseGroupIDFind 使用群id找群
func (group *GroupNum) UseGroupIDFind() []GroupNum {
	var groups []GroupNum
	err := DB.Table(GroupNumT).Where("id = ? and visible = ?", group.ID, true).Find(&groups).Error
	if err != nil {
		return nil
	}
	return groups
}

// UseGroupNameFind	使用群名称找群
func (group *GroupNum) UseGroupNameFind() []GroupNum {
	var groups []GroupNum
	err := DB.Table(GroupNumT).Where("group_name = ? and visible = ?", group.GroupName, true).Find(&groups).Error
	if err != nil {
		return nil
	}
	return groups
}

// IsVerify 检查加入群聊是否需要验证
func (group *GroupNum) IsVerify() uint8 {
	var verify uint8
	err := DB.Table(GroupNumT).Where("id = ?", group.ID).Select("verify").Scan(&verify).Error
	if err != nil {
		return 127
	}
	return verify
}

// GetGroupLeaderID 获取群组id
func (group *GroupNum) GetGroupLeaderID() {
	err := DB.Table(GroupNumT).Where("id = ?", group.ID).Select("group_leader_id").Scan(&(group.GroupLeaderID)).Error
	if err != nil {
		return
	}
}

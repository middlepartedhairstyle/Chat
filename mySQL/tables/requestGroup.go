package tables

import (
	"fmt"
	"github.com/middlepartedhairstyle/HiWe/mySQL"
	"gorm.io/gorm"
)

// RequestAddGroup 存储请求添加群的请求数据
type RequestAddGroup struct {
	gorm.Model
	FromRequestID uint  `gorm:"type:int(11);not null;uniqueIndex:idx_group_request"` //userID
	ToRequestID   uint  `gorm:"type:int(11);not null;uniqueIndex:idx_group_request"` //接收请求的用户id
	AddGroupID    uint  `gorm:"type:int(11);not null;uniqueIndex:idx_group_request"` //GroupID
	State         uint8 `gorm:"type:tinyint(1)"`                                     //是否同意为好友,1为同意为好友,2为待定未确认,3为拒接成为好友
}

type RequestGroupOpt func(*RequestAddGroup)

func SetRequestAddGroupID(requestAddGroupID uint) RequestGroupOpt {
	return func(r *RequestAddGroup) {
		r.ID = requestAddGroupID
	}
}

func SetFromRequestID(fromRequestID uint) RequestGroupOpt {
	return func(r *RequestAddGroup) {
		r.FromRequestID = fromRequestID
	}
}
func SetToRequestID(toRequestID uint) RequestGroupOpt {
	return func(r *RequestAddGroup) {
		r.ToRequestID = toRequestID
	}
}

func SetAddGroupID(addGroupID uint) RequestGroupOpt {
	return func(r *RequestAddGroup) {
		r.AddGroupID = addGroupID
	}
}

func SetState(state uint8) RequestGroupOpt {
	return func(r *RequestAddGroup) {
		r.State = state
	}
}

// NewRequestAddGroup 新建一个群添加请求数据结构
func NewRequestAddGroup(opts ...RequestGroupOpt) *RequestAddGroup {
	requestGroup := &RequestAddGroup{
		FromRequestID: 0,
		AddGroupID:    0,
		State:         2,
	}
	for _, opt := range opts {
		opt(requestGroup)
	}

	return requestGroup
}

// GetRequestAddGroupList 获取加群请求列表
func (request *RequestAddGroup) GetRequestAddGroupList() []RequestAddGroup {
	var result []RequestAddGroup
	err := mySQL.DB.Table(mySQL.RequestAddGroupT).Where("to_request_id=?", request.ToRequestID).Find(&result).Error
	if err != nil {
		return nil
	}
	return result
}

// CreateRequestAddGroup 创建用户加群需求
func (request *RequestAddGroup) CreateRequestAddGroup() *RequestAddGroup {
	var count int64
	err := mySQL.DB.Table(mySQL.RequestAddGroupT).Where("from_request_id=? and to_request_id=? and add_group_id=?", request.FromRequestID, request.ToRequestID, request.AddGroupID).Count(&count).Error
	if err != nil {
		return nil
	}
	if count <= 0 {
		err = mySQL.DB.Table(mySQL.RequestAddGroupT).Create(request).Error
		if err != nil {
			return nil
		}
		return request
	}
	return nil
}

// ChangeState 改变请求状态
func (request *RequestAddGroup) ChangeState() bool {
	var state uint8
	var count int64
	err := mySQL.DB.Table(mySQL.RequestAddGroupT).Where("id=?", request.ID).Count(&count).Error
	if err != nil {
		return false
	}
	if count > 0 {
		err = mySQL.DB.Table(mySQL.RequestAddGroupT).Where("id=?", request.ID).Select("state").Scan(&state).Error
		if err != nil {
			return false
		}
		if state == 2 {
			err = mySQL.DB.Table(mySQL.RequestAddGroupT).Where("id=?", request.ID).Update("state", request.State).Error
			if err != nil {
				return false
			}
			err = mySQL.DB.Table(mySQL.RequestAddGroupT).Where("id=?", request.ID).Find(request).Error
			return true
		}
		fmt.Println(state)
		return false
	}
	return false
}

func (request *RequestAddGroup) ChickToUser() bool {
	var count int64
	err := mySQL.DB.Table(mySQL.RequestAddGroupT).Where("id=? and to_request_id=?", request.ID, request.ToRequestID).Count(&count).Error
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

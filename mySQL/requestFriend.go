package mySQL

import (
	"gorm.io/gorm"
)

const (
	Await  uint8 = iota + 1 //等待成为好友
	Agree                   //同意成为好友
	Refuse                  //拒绝成为好友
)

type RequestFriend struct {
	gorm.Model
	FromRequestID uint  `gorm:"type:int(11);not null;uniqueIndex:idx_friend_request"` //发送好友请求的用户id
	ToRequestID   uint  `gorm:"type:int(11);not null;uniqueIndex:idx_friend_request"` //接收好友请求的用户id
	State         uint8 `gorm:"type:tinyint(1)"`                                      //是否同意为好友,1为同意为好友,2为待定未确认,3为拒接成为好友
}

// InsertInto 插入添加好友请求
func (request *RequestFriend) InsertInto() bool {
	//默认为等待
	request.State = 1

	err := DB.Table(RequestAddFriendT).Create(request).Error
	if err != nil {
		return false
	}
	return true
}

// HaveRequest 使用from_id,to_id判断是否存在该好友请求
func (request *RequestFriend) HaveRequest() bool {
	var count int64
	err := DB.Table(RequestAddFriendT).Where("from_request_id=? AND to_request_id=?", request.FromRequestID, request.ToRequestID).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

// GetID 使用from_id,to_id获取好友请求id
func (request *RequestFriend) GetID() uint {
	var id uint
	err := DB.Table(RequestAddFriendT).Where("from_request_id=? AND to_request_id=?", request.FromRequestID, request.ToRequestID).Select("id").Find(&id).Error
	if err != nil {
		return 0
	}
	return id
}

// GetState 获取状态
func (request *RequestFriend) GetState() {
	err := DB.Table(RequestAddFriendT).Where("id", request.ID).Select("state").Find(&(*request).State).Error
	if err != nil {
		return
	}
}

// GetAllRequest 返回所有被请求信息
func (request *RequestFriend) GetAllRequest() ([]RequestFriend, bool) {
	var requestList []RequestFriend
	err := DB.Table(RequestAddFriendT).Where("to_request_id=? OR from_request_id=?", request.ToRequestID, request.ToRequestID).Find(&requestList).Error
	if err != nil {
		return nil, false
	} else {
		return requestList, true
	}
}

// SetState 设置状态
func (request *RequestFriend) SetState(state uint8) {
	err := DB.Table(RequestAddFriendT).Where("id", request.ID).Update("state", state).Error
	if err != nil {
		return
	}
}

// RemoveRequest 删除好友请求
func (request *RequestFriend) RemoveRequest() bool {
	if err := DB.Table(RequestAddFriendT).Where("id = ?", request.ID).Delete(&RequestAddFriendTable{}); err != nil {
		return false // 删除出错时返回 false
	}
	return true // 删除成功返回 true
}

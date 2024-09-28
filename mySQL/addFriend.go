package mySQL

import (
	"database/sql"
	"errors"
)

// RequestAddFriend 发起好友请求
func (request *RequestAddFriend) RequestAddFriend() bool{
	var exists int
	row := DB.Table(REQUESTADDFRIENDTABLE).Where("from_request_id=? AND to_request_id=?",request.FromRequestID,request.ToRequestID).Select("1").Row()
	err := row.Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		// 处理扫描错误
		return false
	}

	if exists == 0{
		request.State=false
		err=DB.Table(REQUESTADDFRIENDTABLE).Create(request).Error
		if err==nil {
			return true
		}
		return false
	}
	return false
}

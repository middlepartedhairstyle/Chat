package models

import (
	"context"
	"github.com/middlepartedhairstyle/HiWe/mySQL"
	"github.com/middlepartedhairstyle/HiWe/redis"
	"github.com/middlepartedhairstyle/HiWe/utils"
	"time"
)

// UserBaseInfo 用户基础信息
type UserBaseInfo struct {
	Id        uint      `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

// UserVerify 用户验证
type UserVerify struct {
	UserBaseInfo
	Code       string `json:"code"`
	VerifyCode string `json:"verify_code"`
}

// CheckCode 校验验证码
func (u *UserVerify) CheckCode(key string) bool {
	u.VerifyCode, _ = redis.Rdb.Get(context.Background(), key).Result()
	if u.VerifyCode == u.Code {
		return true
	} else {
		return false
	}
}

// EmailIsUser 查询该邮箱是否为用户
func (user *UserBaseInfo) EmailIsUser() (bool, error) {
	var u mySQL.UserBaseInfoTable
	u.Email = user.Email
	b, err := u.UseEmailSelect()
	//如果错误就返回错误
	return b, err
}

// CreateUser 创建用户
func (user *UserBaseInfo) CreateUser() bool {
	var u mySQL.UserBaseInfoTable
	u.Email = user.Email                                                     //用户邮箱
	u.Username = user.Username                                               //用户名
	u.Sale = utils.RandString()                                              //用户密码密钥
	u.Password = utils.MakePasswordSha256(user.Password, u.Sale)             //用户密码
	u.Token = utils.MakeToken(user.Email, user.Password, utils.RandString()) //用户token
	err := u.Create()
	//跟新redis token
	redis.UpdateToken(u.ID, u.Token)
	if err != nil {
		return false
	} else {
		return true
	}
}

// UserInfo 获取用户基本信息
func (user *UserBaseInfo) UserInfo() bool {
	var u mySQL.UserBaseInfoTable
	u.Email = user.Email
	b, err := u.UseEmailSelect()
	if err != nil || !b {
		return false
	} else {
		user.Token = u.Token
		user.Username = u.Username
		user.Email = u.Email
		user.CreatedAt = u.CreatedAt
		user.Id = u.ID
		return true
	}
}

// CheckPassword 验证用户密码正确性
func (user *UserBaseInfo) CheckPassword() bool {
	var u mySQL.UserBaseInfoTable
	u.Email = user.Email
	b, err := u.UseEmailSelect()
	if err != nil {
		return false
	}
	//验证密码准确性
	if b && utils.CheckPasswordSha256(user.Password, u.Sale, u.Password) {
		return true
	}
	return false
}

// UpdateToken 更新token(使用邮箱更新token)
func (user *UserBaseInfo) UpdateToken() bool {
	var u mySQL.UserBaseInfoTable
	u.Email = user.Email
	u.Token = utils.MakeToken(user.Email, user.Password, utils.RandString())
	//数据库token更新
	if u.UpdateToken() && u.FindId() {
		//redis token更新
		return redis.UpdateToken(u.ID, u.Token)
	}
	return false
}

// CheckToken 确认用户token(使用用户id确认)
func (user *UserBaseInfo) CheckToken() bool {
	var u mySQL.UserBaseInfoTable
	if redis.CheckToken(user.Id, user.Token) {
		return true
	} else {
		u.Token = user.Token
		u.ID = user.Id
		//return u.CheckToken()
		return false
	}
}

// GetFriendList 获取好友列表
func (user *UserBaseInfo) GetFriendList() ([]mySQL.Friends, bool) {
	var f mySQL.Friends
	return f.GetFriendList(user.Id)
}

// GetRequestFriendList 获取好友请求添加列表
func (user *UserBaseInfo) GetRequestFriendList() ([]mySQL.RequestFriend, bool) {
	var f mySQL.RequestFriend
	f.ToRequestID = user.Id
	return f.GetAllRequest()
}

// RequestAddFriend 发起添加好友请求(后期添加消息队列，实时了解)
func (user *UserBaseInfo) RequestAddFriend(fromId uint, toId uint) bool {
	//请求表
	var r mySQL.RequestFriend
	r.FromRequestID = fromId
	r.ToRequestID = toId
	//好友表
	var f mySQL.Friends
	f.UserOneID = fromId
	f.UserTwoID = toId
	//判断是否为存在该请求，是否为好友
	if !r.HaveRequest() && !f.IsFriend() {
		//不存在，存入数据库
		return r.InsertInto()
	}
	return false
}

// DisposeAddFriend 处理加好友请求,同意添加为好友,拒绝加好友(后期添加消息队列，实时了解)
func (user *UserBaseInfo) DisposeAddFriend(f mySQL.Friends, requestId uint, state uint8) (bool, uint8) {
	var r mySQL.RequestFriend
	//好友添加请求id
	r.ID = requestId
	r.FromRequestID = f.UserOneID
	r.ToRequestID = f.UserTwoID
	if r.ID == r.GetID() {
		//获取好友请求表内的状态
		r.GetState()
		if state == mySQL.Agree && r.State == mySQL.Await {
			//判断是否为好友
			if !f.IsFriend() {
				//添加好友
				if f.AddFriend() {
					r.SetState(state)

					//(后期添加消息队列，实时了解)

					r.GetState()
					return true, r.State
				}
			}
			//0表示已为好友
			return false, 0
		} else if state == mySQL.Refuse && r.State == mySQL.Await {
			r.SetState(state)

			//(后期添加消息队列，实时了解)

			r.GetState()
			return true, r.State
		} else {
			//4表示已经处理好友请求
			return false, 4
		}
	} else {
		//5表示请求数据与实际数据不符合
		return false, 5
	}
}

// DeleteFriend 删除好友(待完善)
func (user *UserBaseInfo) DeleteFriend() bool { return false }
func (user *UserBaseInfo) Delete()            {}
func (user *UserBaseInfo) Get()               {}

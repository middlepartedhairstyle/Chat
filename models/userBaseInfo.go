package models

import (
	"context"
	"github.com/middlepartedhairstyle/HiWe/mySQL/tables"
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
	var u tables.UserBaseInfo
	u.Email = user.Email
	b, err := u.UseEmailSelect()
	//如果错误就返回错误
	return b, err
}

// CreateUser 创建用户
func (user *UserBaseInfo) CreateUser() bool {
	var u tables.UserBaseInfo
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

// DeleteUser 用户注销且不能被找回
func (user *UserBaseInfo) DeleteUser() bool {
	var u tables.UserBaseInfo
	u.ID = user.Id
	return u.DeleteUser() && redis.DeleteToken(u.ID)
}

// UserInfo 获取用户基本信息
func (user *UserBaseInfo) UserInfo() bool {
	var u tables.UserBaseInfo
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
	var u tables.UserBaseInfo
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
	var u tables.UserBaseInfo
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
	var u tables.UserBaseInfo
	if redis.CheckToken(user.Id, user.Token) {
		return true
	} else {
		u.Token = user.Token
		u.ID = user.Id
		//return u.CheckToken()
		return false
	}
}

// ChangeEmail 更改用户邮箱,校验新邮箱是否被注册并且是否属于该用户
func (user *UserBaseInfo) ChangeEmail() bool {
	var u tables.UserBaseInfo
	u.Email = user.Email
	u.ID = user.Id
	u.UseUserIDSelectPassword()
	u.Token = utils.MakeToken(user.Email, user.Password, utils.RandString())

	if u.ChangeEmail() {
		redis.UpdateToken(u.ID, u.Token)
		user.Token = u.Token
		return true
	} else {
		return false
	}
}

// ChangeUserName 更改用户名
func (user *UserBaseInfo) ChangeUserName() bool {
	var u tables.UserBaseInfo
	u.ID = user.Id
	if u.ChangeBaseUserInfo("username", user.Username) {
		return true
	} else {
		return false
	}
}

// ChangePassword 更改用户密码,存储加密后密码和密钥同时更新token
func (user *UserBaseInfo) ChangePassword() bool {
	var u tables.UserBaseInfo
	u.ID = user.Id
	u.Sale = utils.RandString()
	u.Password = utils.MakePasswordSha256(user.Password, u.Sale)
	u.Token = utils.MakeToken(user.Email, user.Password, utils.RandString())
	if u.ChangePassword() {
		redis.UpdateToken(u.ID, u.Token)
		user.Token = u.Token
		return true
	} else {
		return false
	}
}

// GetFriendList 获取好友列表
func (user *UserBaseInfo) GetFriendList() ([]tables.Friends, bool) {
	var f tables.Friends
	return f.GetFriendList(user.Id)
}

// GetRequestFriendList 获取好友请求添加列表
func (user *UserBaseInfo) GetRequestFriendList() ([]tables.RequestFriend, bool) {
	var f tables.RequestFriend
	f.ToRequestID = user.Id
	return f.GetAllRequest()
}

// RequestAddFriend 发起添加好友请求(后期添加消息队列，实时了解)
func (user *UserBaseInfo) RequestAddFriend(fromId uint, toId uint) bool {
	//请求表
	var r tables.RequestFriend
	r.FromRequestID = fromId
	r.ToRequestID = toId
	//好友表
	var f tables.Friends
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
func (user *UserBaseInfo) DisposeAddFriend(f tables.Friends, requestId uint, state uint8) (bool, uint8) {
	var r tables.RequestFriend
	//好友添加请求id
	r.ID = requestId
	r.FromRequestID = f.UserOneID
	r.ToRequestID = f.UserTwoID
	if r.ID == r.GetID() {
		//获取好友请求表内的状态
		r.GetState()
		if state == tables.Agree && r.State == tables.Await {
			//判断是否为好友
			if !f.IsFriend() {
				//添加好友
				if f.AddFriend() {
					r.SetState(state)

					//(后期添加消息队列，实时了解)
					userMessageBase := NewUserMessageBase(SetUserMessageTypes(1), SetBaseMessage(map[string]uint8{"state": state}))
					info := NewInfo()
					err := info.WriteKafka(userMessageBase, userMessageBase.SetTopic(r.FromRequestID), r.FromRequestID)
					if err != nil {
						return true, r.State
					}

					r.GetState()
					return true, r.State
				}
			}
			//0表示已为好友
			return false, 0
		} else if state == tables.Refuse && r.State == tables.Await {
			r.SetState(state)

			//(后期添加消息队列，实时了解)
			userMessageBase := NewUserMessageBase(SetUserMessageTypes(1), SetBaseMessage(map[string]uint8{"state": state}))
			info := NewInfo()
			err := info.WriteKafka(userMessageBase, userMessageBase.SetTopic(r.FromRequestID), r.FromRequestID)
			if err != nil {
				return true, r.State
			}

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

// CreateGroup 新建群
func (user *UserBaseInfo) CreateGroup(groupName string) (*tables.GroupNum, bool) {
	groupNum := tables.NewGroupNum(tables.SetGroupLeaderID(user.Id), tables.SetGroupName(groupName))
	if groupNum.CreateGroup() {
		return groupNum, true
	}
	return nil, false
}

// FindAllCreateGroup 寻找用户创建的所有群聊
func (user *UserBaseInfo) FindAllCreateGroup() ([]tables.GroupNum, bool) {
	groupNum := tables.NewGroupNum(tables.SetGroupLeaderID(user.Id))
	return groupNum.FindAllCreateGroup()
}

// FindAllGroup 寻找用户加入的所有群聊
func (user *UserBaseInfo) FindAllGroup() ([]tables.GroupUser, bool) {
	groupUser := tables.NewGroupUser(tables.SetUserID(user.Id))
	return groupUser.FindAllGroup()
}

// FindGroup 使用群id或群名寻找群
func (user *UserBaseInfo) FindGroup(groupInfo string) []tables.GroupNum {
	info, err := utils.StringToUint(groupInfo)
	switch err {
	case nil:
		group := tables.NewGroupNum(tables.SetGroupNumID(info))
		return group.UseGroupIDFind()
	default:
		group := tables.NewGroupNum(tables.SetGroupName(groupInfo))
		return group.UseGroupNameFind()
	}
}

// AddGroup 添加群聊(待完善)
func (user *UserBaseInfo) AddGroup(groupID uint) (interface{}, uint8) {
	group := tables.NewGroupNum(tables.SetGroupNumID(groupID))
	switch group.IsVerify() {
	//不需要验证
	case 0:
		groupUser := tables.NewGroupUser(tables.SetUserID(user.Id), tables.SetGroupID(groupID))
		if groupUser.CreateGroupUser() {
			//添加消息队列，发送到对应的群topic中，key设置为gtp(groupID/max)
			//userMessageBase := NewUserMessageBase(SetUserMessageTypes(2), SetBaseMessage(*groupUser))
			//info := NewInfo()
			//err := info.WriteKafka(userMessageBase, userMessageBase.SetTopic(groupUser.UserID), groupUser.UserID)
			//if err != nil {
			//	return nil, 1
			//}
			msg := GroupChangeMessage[user.Id]
			*msg <- groupUser.GroupID
			return *groupUser, 0
		} else {
			return nil, 0
		}
	//需要验证
	case 1:
		group.GetGroupLeaderID()
		groupRequest := tables.NewRequestAddGroup(tables.SetFromRequestID(user.Id), tables.SetToRequestID(group.GroupLeaderID), tables.SetAddGroupID(group.ID))
		result := groupRequest.CreateRequestAddGroup()
		if result != nil {
			//在此处添加消息队列
			userMessageBase := NewUserMessageBase(SetUserMessageTypes(2), SetBaseMessage(result))
			info := NewInfo()
			err := info.WriteKafka(userMessageBase, userMessageBase.SetTopic(result.ToRequestID), result.ToRequestID)
			if err != nil {
				return nil, 1
			}
			return *result, 1
		} else {
			return nil, 1
		}
	default:
		return nil, 1
	}
}

// DisposeAddGroup 处理用户添加群聊，应使用冷热结合的方法，过久没有处理就存入数据库
func (user *UserBaseInfo) DisposeAddGroup(requestID uint, state uint8) (uint8, bool) {
	request := tables.NewRequestAddGroup(tables.SetRequestAddGroupID(requestID), tables.SetState(state), tables.SetToRequestID(user.Id))
	if request.ChickToUser() {
		switch request.State {
		case 1:
			if request.ChangeState() {
				//创建群用户，将被确认者加入群中
				groupUser := tables.NewGroupUser(tables.SetUserID(request.FromRequestID), tables.SetGroupID(request.AddGroupID))
				if groupUser.CreateGroupUser() {
					//kafka,将创建的groupUser数据发给请求者
					//将请求者加入该组的kafka，topic
					userMessageBase := NewUserMessageBase(SetUserMessageTypes(2), SetBaseMessage(groupUser))
					info := NewInfo()
					err := info.WriteKafka(userMessageBase, userMessageBase.SetTopic(groupUser.UserID), groupUser.UserID)
					if err != nil {
						//失败将状态重新归为初始化
						request.State = 2
						request.ChangeState()
						return state, false
					}
					return state, true
				} else {
					//失败将状态重新归为初始化
					request.State = 2
					request.ChangeState()
					return state, false
				}
			}
			return state, false
		case 3:
			if request.ChangeState() {
				//kafka，被确认者
				userMessageBase := NewUserMessageBase(SetUserMessageTypes(2), SetBaseMessage(map[string]uint{"拒绝": 3}))
				info := NewInfo()
				err := info.WriteKafka(userMessageBase, userMessageBase.SetTopic(request.FromRequestID), request.FromRequestID)
				if err != nil {
					//失败将状态重新归为初始化
					request.State = 2
					request.ChangeState()
					return state, false
				}
				return state, true
			}
			return state, false
		default:
			return 2, false
		}
	} else {
		return state, false
	}
}

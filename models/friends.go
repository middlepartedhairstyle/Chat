package models

const (
	Success string = "好友添加成功"
	Failed1 string = "添加失败"
)

// Friend 好友
type Friend struct {
	Id        uint   `json:"id"`         //用户与好友的组队id
	UserID    uint64 `json:"user_id"`    //用户自己的di
	FriendID  uint64 `json:"friend_id"`  //好友的id
	UserToken string `json:"user_token"` //用户的id
}

// GetFriendInfo 获取好友信息
func (friend *Friend) GetFriendInfo() {

}

// IsFriend 判断是否为好友
func (friend *Friend) IsFriend() bool {
	return friend.UserID == friend.UserID
}

// AddFriend 使用用户id添加好友
func (friend *Friend) AddFriend() (bool, string) {
	return friend.IsFriend(), Success
}

// ConfirmAddFriend 用于请求好友添加，后另一方确认是否成为好友
func (friend *Friend) ConfirmAddFriend() (bool, string) {
	return false, Failed1
}

// RequestAddFriend 用于请求成为好友,将请求存入数据库和redis
func (friend *Friend) RequestAddFriend() bool {
	return false
}

// DeleteFriend 用于删除好友，不需要好友确认
func (friend *Friend) DeleteFriend() (bool, string) {
	return false, Failed1
}

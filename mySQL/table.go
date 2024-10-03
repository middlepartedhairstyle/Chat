package mySQL

import "gorm.io/gorm"

const (
	UserBaseInfo     string = "user_base_info_tables"
	GroupNum         string = "group_num_tables"
	GroupMessage     string = "group_messages_tables"
	GroupUser        string = "group_user_tables"
	FriendMessage    string = "friend_message_tables"
	Friend           string = "friends_tables"
	RequestAddFriend string = "request_add_friend_tables"
	RequestAddGroup  string = "request_add_group_tables"
)

// UserBaseInfoTable 用户基础信息(数据库)
type UserBaseInfoTable struct {
	gorm.Model
	Email    string `gorm:"type:varchar(127);not null;unique:email"`
	Username string `gorm:"type:varchar(20)"` //用户名
	Password string `gorm:"type:varchar(64)"` //加密后的密码
	Sale     string `gorm:"type:varchar(6)"`  //加密密钥
	Token    string `gorm:"type:varchar(64)"`
}

// FriendsTable 用户好友表(数据库)
type FriendsTable struct {
	gorm.Model
	UserBaseInfoOne UserBaseInfoTable `gorm:"foreignKey:UserOneID;references:ID"`
	UserBaseInfo    UserBaseInfoTable `gorm:"foreignKey:UserTwoID;references:ID"`
	UserOneID       uint              `gorm:"type:int(11);not null;index:idx_friends,unique"` //用户一id
	UserTwoID       uint              `gorm:"type:int(11);not null;index:idx_friends,unique"` //用户二id
	Relationship    string            `gorm:"type:varchar(10)"`                               //两者关系
	NoteOne         string            `gorm:"type:varchar(20)"`                               //UserOne个UserTwo的备注
	NoteTwo         string            `gorm:"type:varchar(20)"`                               //UserTwo个UserOne的备注
}

// FriendMessageTable 用户消息(数据库)
type FriendMessageTable struct {
	gorm.Model
	UserBaseInfoOne UserBaseInfoTable `gorm:"foreignKey:FromID;references:ID"`
	UserBaseInfo    UserBaseInfoTable `gorm:"foreignKey:FriendID;references:ID"`
	FromID          uint              //发送者id
	FriendID        uint              //好友组队id
	MessageType     uint8             `gorm:"type:tinyint(1)"` //消息类型
	Message         string            `gorm:"type:text"`       //消息主体
}

// GroupMessageTable 群消息(数据库)
type GroupMessageTable struct {
	gorm.Model
	UserBaseInfo UserBaseInfoTable `gorm:"foreignKey:FromID;references:ID"`
	GroupNum     GroupNumTable     `gorm:"foreignKey:GroupID;references:ID"`
	FromID       uint              //发送信息用户id
	GroupID      uint              //群id
	MessageType  uint8             `gorm:"type:tinyint(1)"` //消息类型，文本图片
	Message      string            `gorm:"type:text"`       //消息主体
}

// GroupNumTable 用户群(数据库)
type GroupNumTable struct {
	gorm.Model
	UserBaseInfo  UserBaseInfoTable `gorm:"foreignKey:GroupLeaderID;references:ID"`
	GroupLeaderID uint              //群主id
	GroupName     string            `gorm:"type:varchar(25)"` //群名
	Visible       bool              `gorm:"type:tinyint(1)"`  //该群是否可以被搜索
}

// GroupUserTable 用户群的用户(数据库)
type GroupUserTable struct {
	gorm.Model
	GroupNum     GroupNumTable     `gorm:"foreignKey:GroupID;references:ID"`
	UserBaseInfo UserBaseInfoTable `gorm:"foreignKey:UserID;references:ID"`
	GroupID      uint              //群组id
	UserID       uint              //用户id
	Note         string            `gorm:"type:varchar(20)"` //用户给群的备注
	Level        uint8             `gorm:"type:tinyint(1)"`  //用户在群中的等级
	Relationship string            `gorm:"type:varchar(10)"` //用户在群中的关系如,群主,管理员,群员等
}

// RequestAddFriendTable 存储请求添加好友的请求数据(数据库)
type RequestAddFriendTable struct {
	gorm.Model
	UserBaseInfoOne UserBaseInfoTable `gorm:"foreignKey:FromRequestID;references:ID"`
	UserBaseInfo    UserBaseInfoTable `gorm:"foreignKey:ToRequestID;references:ID"`
	FromRequestID   uint              `gorm:"type:int(11);not null;uniqueIndex:idx_friend_request"` //发送好友请求的用户id
	ToRequestID     uint              `gorm:"type:int(11);not null;uniqueIndex:idx_friend_request"` //接收好友请求的用户id
	State           uint8             `gorm:"type:tinyint(1)"`                                      //是否同意为好友,1为同意为好友,2为待定未确认,3为拒接成为好友
}

// RequestAddGroupTable 存储请求添加群的请求数据(数据库)
type RequestAddGroupTable struct {
	gorm.Model
	GroupNum      GroupNumTable     `gorm:"foreignKey:FromRequestID;references:ID"`
	UserBaseInfo  UserBaseInfoTable `gorm:"foreignKey:AddGroupID;references:ID"`
	FromRequestID uint              `gorm:"type:int(11);not null;uniqueIndex:idx_group_request"`
	AddGroupID    uint              `gorm:"type:int(11);not null;uniqueIndex:idx_group_request"`
	State         uint8             `gorm:"type:tinyint(1)"` //是否同意为好友,1为同意为好友,2为待定未确认,3为拒接成为好友
}

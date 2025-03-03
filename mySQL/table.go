package mySQL

import "gorm.io/gorm"

const (
	UserBaseInfoT     string = "user_base_info_tables"
	GroupNumT         string = "group_num_tables"
	GroupMessageT     string = "group_message_tables"
	GroupUserT        string = "group_user_tables"
	FriendMessageT    string = "friend_message_tables"
	FriendT           string = "friends_tables"
	RequestAddFriendT string = "request_add_friend_tables"
	RequestAddGroupT  string = "request_add_group_tables"
	UserDetailedInfoT string = "user_detailed_info_tables"
)

// UserBaseInfoTable 用户基础信息(数据库)
type UserBaseInfoTable struct {
	gorm.Model
	Email                     string                  `gorm:"type:varchar(127);not null;unique:email"`
	Username                  string                  `gorm:"type:varchar(20)"` //用户名
	Password                  string                  `gorm:"type:varchar(64)"` //加密后的密码
	Sale                      string                  `gorm:"type:varchar(6)"`  //加密密钥
	Token                     string                  `gorm:"type:varchar(64)"`
	FriendsTablesOne          []FriendsTable          `gorm:"foreignKey:UserOneID;constraint:OnDelete:CASCADE;"`
	FriendsTablesTwo          []FriendsTable          `gorm:"foreignKey:UserTwoID;constraint:OnDelete:CASCADE;"`
	FriendMessageTablesOne    []FriendMessageTable    `gorm:"foreignKey:FromID;constraint:OnDelete:CASCADE;"`
	FriendMessageTablesTwo    []FriendMessageTable    `gorm:"foreignKey:FriendID;constraint:OnDelete:CASCADE;"`
	GroupMessageTables        []GroupMessageTable     `gorm:"foreignKey:FromID;constraint:OnDelete:CASCADE;"`
	GroupNumTables            []GroupNumTable         `gorm:"foreignKey:GroupLeaderID;constraint:OnDelete:CASCADE;"`
	GroupUserTables           []GroupUserTable        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	RequestAddFriendTablesOne []RequestAddFriendTable `gorm:"foreignKey:FromRequestID;constraint:OnDelete:CASCADE;"`
	RequestAddFriendTablesTwo []RequestAddFriendTable `gorm:"foreignKey:ToRequestID;constraint:OnDelete:CASCADE;"`
	RequestAddGroupTablesOne  []RequestAddGroupTable  `gorm:"foreignKey:FromRequestID;constraint:OnDelete:CASCADE;"`
	RequestAddGroupTablesTwo  []RequestAddGroupTable  `gorm:"foreignKey:ToRequestID;constraint:OnDelete:CASCADE;"`
	UserDetailedInfoTables    []UserDetailedInfoTable `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

// FriendsTable 用户好友表(数据库)
type FriendsTable struct {
	gorm.Model
	UserBaseInfoOne UserBaseInfoTable `gorm:"foreignKey:UserOneID;references:ID;constraint:OnDelete:CASCADE;"`
	UserBaseInfo    UserBaseInfoTable `gorm:"foreignKey:UserTwoID;references:ID;constraint:OnDelete:CASCADE;"`
	UserOneID       uint              `gorm:"type:int(11);not null;index:idx_friends,unique"` //用户一id
	UserTwoID       uint              `gorm:"type:int(11);not null;index:idx_friends,unique"` //用户二id
	Relationship    string            `gorm:"type:varchar(10)"`                               //两者关系
	NoteOne         string            `gorm:"type:varchar(20)"`                               //UserOne个UserTwo的备注
	NoteTwo         string            `gorm:"type:varchar(20)"`                               //UserTwo个UserOne的备注
}

// FriendMessageTable 用户消息(数据库)
type FriendMessageTable struct {
	gorm.Model
	UserBaseInfoOne UserBaseInfoTable `gorm:"foreignKey:FromID;references:ID;constraint:OnDelete:CASCADE;"`
	UserBaseInfo    UserBaseInfoTable `gorm:"foreignKey:FriendID;references:ID;constraint:OnDelete:CASCADE;"`
	FromID          uint              //发送者id
	FriendID        uint              //好友组队id
	MessageType     uint8             `gorm:"type:tinyint(1)"` //消息类型
	Message         string            `gorm:"type:text"`       //消息主体
}

// GroupMessageTable 群消息(数据库)
type GroupMessageTable struct {
	gorm.Model
	UserBaseInfo UserBaseInfoTable `gorm:"foreignKey:FromID;references:ID;constraint:OnDelete:CASCADE;"`
	GroupNum     GroupNumTable     `gorm:"foreignKey:GroupID;references:ID;constraint:OnDelete:CASCADE;"`
	FromID       uint              //发送信息用户id
	GroupID      uint              //群id
	MessageType  uint8             `gorm:"type:tinyint(1)"` //消息类型，文本图片
	Message      string            `gorm:"type:text"`       //消息主体
}

// GroupNumTable 用户群(数据库)
type GroupNumTable struct {
	gorm.Model
	UserBaseInfo  UserBaseInfoTable `gorm:"foreignKey:GroupLeaderID;references:ID;constraint:OnDelete:CASCADE;"`
	GroupLeaderID uint              //群主id
	GroupName     string            `gorm:"type:varchar(25)"` //群名
	Visible       bool              `gorm:"type:tinyint(1)"`  //该群是否可以被搜索
	Verify        uint8             `gorm:"type:tinyint(1)"`  //添加该群是否需要群主同意,0需要,1不需要······
}

// GroupUserTable 用户群的用户(数据库)
type GroupUserTable struct {
	gorm.Model
	GroupNum     GroupNumTable     `gorm:"foreignKey:GroupID;references:ID;constraint:OnDelete:CASCADE;"`
	UserBaseInfo UserBaseInfoTable `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`
	GroupID      uint              //群组id
	UserID       uint              //用户id
	Note         string            `gorm:"type:varchar(20)"` //用户给群的备注
	Level        uint8             `gorm:"type:tinyint(1)"`  //用户在群中的等级
	Relationship uint8             `gorm:"type:tinyint(1)"`  //用户在群中的关系如,1群主,2管理员,3群员等
}

// RequestAddFriendTable 存储请求添加好友的请求数据(数据库)
type RequestAddFriendTable struct {
	gorm.Model
	UserBaseInfoOne UserBaseInfoTable `gorm:"foreignKey:FromRequestID;references:ID;constraint:OnDelete:CASCADE;"`
	UserBaseInfo    UserBaseInfoTable `gorm:"foreignKey:ToRequestID;references:ID;constraint:OnDelete:CASCADE;"`
	FromRequestID   uint              `gorm:"type:int(11);not null;uniqueIndex:idx_friend_request"` //发送好友请求的用户id
	ToRequestID     uint              `gorm:"type:int(11);not null;uniqueIndex:idx_friend_request"` //接收好友请求的用户id
	State           uint8             `gorm:"type:tinyint(1)"`                                      //是否同意为好友,1为同意为好友,2为待定未确认,3为拒接成为好友
}

// RequestAddGroupTable 存储请求添加群的请求数据(数据库)
type RequestAddGroupTable struct {
	gorm.Model
	GroupNum        GroupNumTable     `gorm:"foreignKey:AddGroupID;references:ID;constraint:OnDelete:CASCADE;"`
	UserBaseInfoOne UserBaseInfoTable `gorm:"foreignKey:FromRequestID;references:ID;constraint:OnDelete:CASCADE;"`
	UserBaseInfo    UserBaseInfoTable `gorm:"foreignKey:ToRequestID;references:ID;constraint:OnDelete:CASCADE;"`
	FromRequestID   uint              `gorm:"type:int(11);not null;uniqueIndex:idx_group_request"` //发送请求的用户id
	ToRequestID     uint              `gorm:"type:int(11);not null;uniqueIndex:idx_group_request"` //接收请求的用户id
	AddGroupID      uint              `gorm:"type:int(11);not null;uniqueIndex:idx_group_request"`
	State           uint8             `gorm:"type:tinyint(1)"` //是否同意加群,1为同意,2为待定未确认,3为拒绝
}

// UserDetailedInfoTable 用户详细信息表
type UserDetailedInfoTable struct {
	gorm.Model
	UserBaseInfo UserBaseInfoTable `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`
	UserID       uint              `gorm:"type:int(11);not null;uniqueIndex:idx_user"`
	ProfilePhoto string            `gorm:"type:varchar(255)"` //用户头像
}

package mySQL

import "gorm.io/gorm"

const (
	GROUPNUMTABLE     string = "group_num_tables"
	GROUPMESSAGETABLE string = "group_messages_tables"
	GROUPUSERTABLE    string = "group_user_tables"
	USERBASETABLE     string = "user_base_info_tables"
	USERMESSAGETABLE  string = "user_messages_tables"
	USERFRIENDSTABLE  string = "user_friends_tables"
	REQUESTADDFRIENDTABLE  string = "request_add_friends"
)

const (
	Id            string = "id"
	CreatedAt            = "created_at"
	UpdatedAt            = "updated_at"
	DeletedAt            = "deleted_at"
	Email                = "email"
	Username             = "username"
	Password             = "password"
	Md5num               = "md5_num"
	Token                = "token"
	FromId               = "from_id"
	ToId                 = "to_id"
	MessageType          = "message_type"
	Media                = "media"
	Content              = "content"
	UserId               = "user_id"
	FriendId             = "friend_id"
	GroupId              = "group_id"
	Level                = "level"
	GroupName            = "group_name"
	GroupLeaderId        = "group_leader_id"
	ToGroupId            = "to_group_id"
)

// UserBaseInfoTable 用户基础信息(数据库)
type UserBaseInfoTable struct {
	gorm.Model
	Email    string `gorm:"type:varchar(127);not null;unique:email"`
	Username string `gorm:"type:varchar(20)"`
	Password string `gorm:"type:varchar(32)"`
	Md5Num   string `gorm:"type:varchar(6)"`
	Token    string `gorm:"type:varchar(32)"`
}

// UserFriendsTable 用户好友表(数据库)
type UserFriendsTable struct {
	gorm.Model
	UserBaseInfoOne UserBaseInfoTable `gorm:"foreignKey:UserID;references:ID"`
	UserBaseInfo    UserBaseInfoTable `gorm:"foreignKey:FriendID;references:ID"`
	UserID          uint64            `gorm:"type:int(11);not null;uniqueIndex:idx_user_friend"`
	FriendID        uint64            `gorm:"type:int(11);not null;uniqueIndex:idx_user_friend"`
}

// UserMessageTable 用户消息(数据库)
type UserMessageTable struct {
	gorm.Model
	UserBaseInfoOne UserBaseInfoTable `gorm:"foreignKey:FromID;references:ID"`
	UserBaseInfo    UserBaseInfoTable `gorm:"foreignKey:ToID;references:ID"`
	FromID          uint64
	ToID            uint64
	MessageType     uint8  `gorm:"type:tinyint(1)"`
	Media           uint8  `gorm:"type:tinyint(1)"`
	Content         string `gorm:"type:text"`
}

// GroupMessageTable 群消息(数据库)
type GroupMessageTable struct {
	gorm.Model
	UserBaseInfo UserBaseInfoTable `gorm:"foreignKey:FromID;references:ID"`
	GroupNum     GroupNumTable     `gorm:"foreignKey:ToGroupID;references:ID"`
	FromID       uint64
	ToGroupID    uint64
	MessageType  uint8  `gorm:"type:tinyint(1)"`
	Media        uint8  `gorm:"type:tinyint(1)"`
	Content      string `gorm:"type:text"`
}

// GroupNumTable 用户群(数据库)
type GroupNumTable struct {
	gorm.Model
	UserBaseInfo  UserBaseInfoTable `gorm:"foreignKey:GroupLeaderID;references:ID"`
	GroupLeaderID uint64
	GroupName     string `gorm:"type:varchar(25)"`
}

// GroupUserTable 用户群的用户(数据库)
type GroupUserTable struct {
	gorm.Model
	GroupNum     GroupNumTable     `gorm:"foreignKey:GroupID;references:ID"`
	UserBaseInfo UserBaseInfoTable `gorm:"foreignKey:UserID;references:ID"`
	GroupID      uint64
	UserID       uint64
	Level        uint8 `gorm:"type:tinyint(1)"`
}

// RequestAddFriend 存储请求添加好友的请求数据(数据库)
type RequestAddFriend struct {
	gorm.Model
	UserBaseInfoOne UserBaseInfoTable `gorm:"foreignKey:FromRequestID;references:ID"`
	UserBaseInfo    UserBaseInfoTable `gorm:"foreignKey:ToRequestID;references:ID"`
	FromRequestID   uint64            `gorm:"type:int(11);not null;uniqueIndex:idx_friend_request"`
	ToRequestID     uint64            `gorm:"type:int(11);not null;uniqueIndex:idx_friend_request"`
	State           bool              `gorm:"type:tinyint(1)"` //是否同意为好友
}

package mySQL

import "gorm.io/gorm"

const (
	CAPTCHATABLE      string = "captcha_tables"
	GROUPNUMTABLE     string = "group_num_tables"
	GROUPMESSAGETABLE string = "group_messages_tables"
	GROUPUSERTABLE    string = "group_user_tables"
	USERBASETABLE     string = "user_base_info_tables"
	USERMESSAGETABLE  string = "user_messages_tables"
)

// UserBaseInfoTable 用户基础信息(数据库)
type UserBaseInfoTable struct {
	gorm.Model
	Email    string `gorm:"type:varchar(127);unique_index"`
	Username string `gorm:"type:varchar(20)"`
	Password string `gorm:"type:varchar(32)"`
	Md5Num   string `gorm:"type:varchar(6)"`
	Token    string `gorm:"type:varchar(32)"`
}

// CaptchaTable 用户验证码(数据库)
type CaptchaTable struct {
	gorm.Model
	Code         uint32 `gorm:"type:int(11)"`
	Email        string
	UserBaseInfo UserBaseInfoTable `gorm:"type:varchar(127);foreignKey:Email;references:Email"`
}

// UserMessageTable 用户消息(数据库)
type UserMessageTable struct {
	gorm.Model
	UserBaseInfoOne UserBaseInfoTable `gorm:"foreignKey:FromID;references:ID"`
	UserBaseInfo    UserBaseInfoTable `gorm:"foreignKey:ToID;references:ID"`
	FromID          uint
	ToID            uint
	MessageType     uint8  `gorm:"type:tinyint(1)"`
	Media           uint8  `gorm:"type:tinyint(1)"`
	Content         string `gorm:"type:text"`
}

// GroupMessageTable 群消息(数据库)
type GroupMessageTable struct {
	gorm.Model
	UserBaseInfo UserBaseInfoTable `gorm:"foreignKey:FromID;references:ID"`
	GroupNum     GroupNumTable     `gorm:"foreignKey:ToGroupID;references:ID"`
	FromID       uint
	ToGroupID    uint
	MessageType  uint8  `gorm:"type:tinyint(1)"`
	Media        uint8  `gorm:"type:tinyint(1)"`
	Content      string `gorm:"type:text"`
}

// GroupNumTable 用户群(数据库)
type GroupNumTable struct {
	gorm.Model
	UserBaseInfo  UserBaseInfoTable `gorm:"foreignKey:GroupLeaderID;references:ID"`
	GroupLeaderID uint
	GroupName     string `gorm:"type:varchar(25)"`
}

// GroupUserTable 用户群的用户(数据库)
type GroupUserTable struct {
	gorm.Model
	GroupNum     GroupNumTable     `gorm:"foreignKey:GroupID;references:ID"`
	UserBaseInfo UserBaseInfoTable `gorm:"foreignKey:UserID;references:ID"`
	GroupID      uint
	UserID       uint
	Level        uint8 `gorm:"type:tinyint(1)"`
}

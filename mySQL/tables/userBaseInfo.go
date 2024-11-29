package tables

import (
	"github.com/middlepartedhairstyle/HiWe/mySQL"
	"gorm.io/gorm"
)

// UserBaseInfo 用户基础信息(数据库)
type UserBaseInfo struct {
	gorm.Model
	Email    string `gorm:"type:varchar(127);not null;unique:email"`
	Username string `gorm:"type:varchar(20)"` //用户名
	Password string `gorm:"type:varchar(64)"` //加密后的密码
	Sale     string `gorm:"type:varchar(6)"`  //加密密钥
	Token    string `gorm:"type:varchar(64)"`
}

// UseEmailSelect 使用邮箱查询用户数据
func (user *UserBaseInfo) UseEmailSelect() (bool, error) {
	var count int64
	err := mySQL.DB.Table(mySQL.UserBaseInfoT).Where("email=?", user.Email).Count(&count).Error
	//错误处理,不为空返回错误
	if err != nil {
		return false, err
	}
	//检查是否有该邮箱，没有就返回false
	if count <= 0 {
		return false, nil
	} else {
		//有该邮箱就将用户数据写入该结构体中
		err = mySQL.DB.Table(mySQL.UserBaseInfoT).Where("email=?", user.Email).Select("*").Scan(user).Error
		//如果错误返回err
		if err != nil {
			return true, err
		}
		return true, nil
	}
}

func (user *UserBaseInfo) UseUserIDSelectPassword() bool {
	if user.CheckUserID() {
		err := mySQL.DB.Table(mySQL.UserBaseInfoT).Where("id=?", user.ID).Select("password").Scan(&user.Password).Error
		if err != nil {
			return false
		}
		return true
	} else {
		return false
	}
}

// Create 创建新用户
func (user *UserBaseInfo) Create() error {
	err := mySQL.DB.Table(mySQL.UserBaseInfoT).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (user *UserBaseInfo) UpdateToken() bool {
	err := mySQL.DB.Table(mySQL.UserBaseInfoT).Where("email=?", user.Email).Update("token", user.Token).Error
	if err != nil {
		return false
	}
	return true
}

// CheckToken 确认用户token
func (user *UserBaseInfo) CheckToken() bool {
	var token string
	err := mySQL.DB.Table(mySQL.UserBaseInfoT).Where("id=?", user.ID).Select("token").Scan(&token).Error
	if err != nil {
		return false
	}
	if token == user.Token {
		return true
	}
	return false
}

// FindId 使用邮箱查找用户id
func (user *UserBaseInfo) FindId() bool {
	var id uint
	err := mySQL.DB.Table(mySQL.UserBaseInfoT).Where("email=?", user.Email).Select("id").Scan(&id).Error
	if err != nil {
		return false
	}
	if id != 0 {
		user.ID = id
		return true
	}
	return false
}

// UseIDFindEmail 使用用户id寻找用户邮箱
func (user *UserBaseInfo) UseIDFindEmail() bool {
	err := mySQL.DB.Table(mySQL.UserBaseInfoT).Where("id=?", user.ID).Select("email").Find(&user.Email).Error
	if err != nil {
		return false
	} else {
		return true
	}
}

// CheckUserID 使用用户id,核对是否存在该用户
func (user *UserBaseInfo) CheckUserID() bool {
	var count int64
	err := mySQL.DB.Table(mySQL.UserBaseInfoT).Where("id=?", user.ID).Count(&count).Error
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}

// ChangeBaseUserInfo 更改用户基本信息，仅支持单条信息改变(这里仅该改变用户名，用户邮箱，用户密码)
func (user *UserBaseInfo) ChangeBaseUserInfo(columnName string, value string) bool {
	if user.CheckUserID() {
		err := mySQL.DB.Table(mySQL.UserBaseInfoT).Where("id=?", user.ID).Update(columnName, value).Error
		if err != nil {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}

// ChangePassword 更改用户密码
func (user *UserBaseInfo) ChangePassword() bool {
	if user.CheckUserID() {
		err := mySQL.DB.Table(mySQL.UserBaseInfoT).Where("id=?", user.ID).Updates(map[string]interface{}{
			"password": user.Password,
			"sale":     user.Sale,
			"token":    user.Token,
		}).Error
		if err != nil {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}

func (user *UserBaseInfo) ChangeEmail() bool {
	var count int64
	err := mySQL.DB.Table(mySQL.UserBaseInfoT).Where("email=?", user.Email).Count(&count).Error
	//错误处理,不为空返回错误
	if err != nil {
		return false
	}
	//检查是否有该邮箱
	if count <= 0 && user.Email != "" && user.Token != "" {
		err = mySQL.DB.Table(mySQL.UserBaseInfoT).Where("id=?", user.ID).Updates(map[string]interface{}{
			"email": user.Email,
			"token": user.Token,
		}).Error
		if err != nil {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}

// DeleteUser 删除用户
func (user *UserBaseInfo) DeleteUser() bool {
	if user.CheckUserID() {
		err := mySQL.DB.Table(mySQL.UserBaseInfoT).Where("id=?", user.ID).Delete(user).Error
		if err != nil {
			return false
		}
		return true
	} else {
		return false
	}
}

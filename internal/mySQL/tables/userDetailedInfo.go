package tables

import (
	"fmt"
	mySQL2 "github.com/middlepartedhairstyle/HiWe/internal/mySQL"
	"github.com/middlepartedhairstyle/HiWe/internal/utils"
	"gorm.io/gorm"
	"time"
)

/*
	该文件是对数据库表格user_detailed_info_tables进行操作
*/

// UserDetailedInfo 用户详细信息表
type UserDetailedInfo struct {
	ID           uint           `json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
	UserID       uint           `json:"user_id"`
	ProfilePhoto string         `json:"profile_photo"`
}

// ChangeProfilePhoto 更改用户头像
func (UDI *UserDetailedInfo) ChangeProfilePhoto() bool {
	var count int64 = 0
	err := mySQL2.DB.Table(mySQL2.UserDetailedInfoT).Where("user_id=?", UDI.UserID).Count(&count).Error
	if err != nil {
		return false
	}
	if count <= 0 {
		UDI.CreatedAt = utils.GetTimeToUTC()
		UDI.UpdatedAt = utils.GetTimeToUTC()
		err = mySQL2.DB.Table(mySQL2.UserDetailedInfoT).Create(&UDI).Error
		if err != nil {
			return false
		}
		fmt.Println(UDI)
		return true
	} else {
		UDI.UpdatedAt = utils.GetTimeToUTC()
		err = mySQL2.DB.Table(mySQL2.UserDetailedInfoT).Where("user_id=?", UDI.UserID).Updates(map[string]interface{}{
			"profile_photo": UDI.ProfilePhoto,
			"updated_at":    UDI.UpdatedAt,
		}).Error
		if err != nil {
			return false
		}
		err = mySQL2.DB.Table(mySQL2.UserDetailedInfoT).Where("user_id=?", UDI.UserID).Find(&UDI).Error
		return true
	}
}

func (UDI *UserDetailedInfo) CheckProfilePhoto() bool {
	var count int64 = 0
	err := mySQL2.DB.Table(mySQL2.UserDetailedInfoT).Where("user_id=? AND profile_photo=?", UDI.UserID, UDI.ProfilePhoto).Count(&count).Error
	if err != nil {
		return false
	}
	if count <= 0 {
		return false
	}
	return true
}

package mySQL

// UseEmailSelect 使用邮箱查询用户数据
func (user *UserBaseInfoTable) UseEmailSelect() (bool, error) {
	var count int64
	err := DB.Table(UserBaseInfoT).Where("email=?", user.Email).Count(&count).Error
	//错误处理,不为空返回错误
	if err != nil {
		return false, err
	}
	//检查是否有该邮箱，没有就返回false
	if count <= 0 {
		return false, nil
	} else {
		//有该邮箱就将用户数据写入该结构体中
		err = DB.Table(UserBaseInfoT).Where("email=?", user.Email).Select("*").Scan(user).Error
		//如果错误返回err
		if err != nil {
			return true, err
		}
		return true, nil
	}
}

// Create 创建新用户
func (user *UserBaseInfoTable) Create() error {
	err := DB.Table(UserBaseInfoT).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (user *UserBaseInfoTable) UpdateToken() bool {
	err := DB.Table(UserBaseInfoT).Where("email=?", user.Email).Update("token", user.Token).Error
	if err != nil {
		return false
	}
	return true
}

// CheckToken 确认用户token
func (user *UserBaseInfoTable) CheckToken() bool {
	var token string
	err := DB.Table(UserBaseInfoT).Where("id=?", user.ID).Select("token").Scan(&token).Error
	if err != nil {
		return false
	}
	if token == user.Token {
		return true
	}
	return false
}

// FindId 使用邮箱查找用户id
func (user *UserBaseInfoTable) FindId() bool {
	var id uint
	err := DB.Table(UserBaseInfoT).Where("email=?", user.Email).Select("id").Scan(&id).Error
	if err != nil {
		return false
	}
	if id != 0 {
		user.ID = id
		return true
	}
	return false
}

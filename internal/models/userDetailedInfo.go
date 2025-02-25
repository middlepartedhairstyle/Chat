package models

import (
	"github.com/middlepartedhairstyle/HiWe/internal/mySQL/tables"
)

type UserDetailed struct {
	DetailedInfo *tables.UserDetailedInfo
}

func NewUserDetailedInfo() *UserDetailed {
	return &UserDetailed{
		DetailedInfo: &tables.UserDetailedInfo{},
	}
}

func (UDI UserDetailed) ChangeProfilePhoto(ext string, id string) (string, bool) {
	imageURL := "/images/profile_photo/" + id + ext
	UDI.DetailedInfo.ProfilePhoto = imageURL
	b := UDI.DetailedInfo.ChangeProfilePhoto()
	if b {
		return imageURL, true
	} else {
		return "", false
	}
}

func (UDI UserDetailed) GetProfilePhoto() bool {
	return UDI.DetailedInfo.CheckProfilePhoto()
}

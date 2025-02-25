package service

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/internal/models"
	utils2 "github.com/middlepartedhairstyle/HiWe/internal/utils"
	"path/filepath"
)

// ChangeUserProfilePhoto 用户更改头像
func (h *HTTPServer) ChangeUserProfilePhoto(c *gin.Context) {
	id := c.GetHeader("id")
	image, err := c.FormFile("profile_photo")
	if err != nil {
		utils2.Fail(c, ServerError, gin.H{
			"err_msg": "服务器错误",
		})
		return
	}
	ext := filepath.Ext(image.Filename)
	filePath := filepath.Join("./images/profile_photo", id+ext)
	err = c.SaveUploadedFile(image, filePath)
	if err != nil {
		utils2.Fail(c, ServerError, gin.H{
			"err_msg": "上传失败",
		})
		return
	}
	userID, _ := utils2.StringToUint(id)
	u := models.NewUserDetailedInfo()
	u.DetailedInfo.UserID = userID
	photoURL, b := u.ChangeProfilePhoto(ext, id)
	if b {
		utils2.Success(c, SUCCESS, gin.H{
			"msg":   "success",
			"image": photoURL,
		})
	} else {
		utils2.Fail(c, ServerError, gin.H{
			"err_msg": "上传失败",
		})
		return
	}

}

// GetUserProfilePhoto 用户获取头像
func (h *HTTPServer) GetUserProfilePhoto(c *gin.Context) {
	id := c.GetHeader("id")
	profilePhoto, _ := c.GetQuery("profile_photo")
	userID, _ := utils2.StringToUint(id)
	u := models.NewUserDetailedInfo()
	u.DetailedInfo.UserID = userID
	u.DetailedInfo.ProfilePhoto = profilePhoto
	if u.GetProfilePhoto() {
		c.File(profilePhoto)
	} else {
		utils2.Fail(c, ServerError, gin.H{
			"err_msg": "未找到该用户头像",
		})
	}
}

package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/controllers"
)

func Routers(router *gin.Engine) {
	router.POST("/register", controllers.RegisterController)
}

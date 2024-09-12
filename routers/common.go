package routers

import "github.com/gin-gonic/gin"

func Routers(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {})
}

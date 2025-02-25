package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/internal/utils"
)

// LimitSizeMiddleware 限制文件大小
func LimitSizeMiddleware(size int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > size {
			utils.Fail(c, ServerError, gin.H{
				"err_msg": "文件过大",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

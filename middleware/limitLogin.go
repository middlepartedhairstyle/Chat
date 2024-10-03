package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/redis"
	"net/http"
	"time"
)

const maxRequestOfLogin = 5
const expirationLogin = 1 * time.Minute // 请求计数的过期时间设置为1分钟

// LimitLogin 登录次数限制
func LimitLogin(router string) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.GetHeader("email")
		key := fmt.Sprintf("%s%s", router, email)
		count, _ := redis.Rdb.Incr(ctx, key).Result()
		// 设置过期时间，确保每分钟重置计数
		if count == 1 {
			redis.Rdb.Expire(ctx, key, expirationLogin)
		}

		if count > maxRequestOfLogin {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "大于最大次数，一个小时后再试"})
			c.Abort()
			return
		}
		c.Next()
	}
}

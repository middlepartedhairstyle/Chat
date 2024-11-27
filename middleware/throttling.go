package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/redis"
	"github.com/middlepartedhairstyle/HiWe/utils"
	"time"
)

const maxRequestsOfIP = 6          // IP每分钟最大请求次数
const maxRequestsOfUser = 5        //用户每分钟最大请求次数
const expiration = 1 * time.Minute // 请求计数的过期时间设置为1分钟
var ctx = context.Background()

// RateLimiter 请求频率限制
func RateLimiter(router string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIP := c.ClientIP()
		userToken := c.GetHeader("token")

		//没有token时
		if userToken == "" {
			key := fmt.Sprintf("%s%s", router, userIP)
			// 使用 Redis 的 INCR 命令增加请求计数
			count, err := redis.Rdb.Incr(ctx, key).Result()
			if err != nil {
				utils.Fail(c, ServerError, gin.H{"error": "内部服务器错误"})
				c.Abort()
				return
			}
			// 设置过期时间，确保每分钟重置计数
			if count == 1 {
				redis.Rdb.Expire(ctx, key, expiration)
			}

			if count > maxRequestsOfIP {
				time.Sleep(3 * time.Second) //对于异常用户进行延迟处理
				utils.Fail(c, RequestsTooFrequent, gin.H{"error": "请求过于频繁，请稍后再试"})
				c.Abort()
				return
			}
			c.Next() // 继续处理请求
		} else {
			key1 := fmt.Sprintf("%s%s", router, userIP)
			key2 := fmt.Sprintf("%s%s", router, userToken)

			count1, err1 := redis.Rdb.Incr(ctx, key1).Result()
			count2, err2 := redis.Rdb.Incr(ctx, key2).Result()
			if err1 != nil || err2 != nil {
				utils.Fail(c, ServerError, gin.H{"error": "内部服务器错误"})
				c.Abort()
				return
			}
			// 设置过期时间，确保每分钟重置计数
			if count1 == 1 {
				redis.Rdb.Expire(ctx, key1, expiration)
			}
			if count2 == 1 {
				redis.Rdb.Expire(ctx, key2, expiration)
			}

			if count1 > maxRequestsOfIP && count2 > maxRequestsOfUser {
				time.Sleep(3 * time.Second) //对于异常用户进行延迟处理
				utils.Fail(c, RequestsTooFrequent, gin.H{"error": "请求过于频繁，请稍后再试"})
				c.Abort()
				return
			}
			c.Next() // 继续处理请求
		}

	}

}

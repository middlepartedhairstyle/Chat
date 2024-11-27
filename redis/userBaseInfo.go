package redis

import (
	"context"
	"strconv"
	"time"
)

const maxTokenTime = 720 * time.Hour //设置token过期时间为30天

var ctx = context.Background()

// UpdateToken 更新token
func UpdateToken(id uint, token string) bool {
	err := Rdb.Set(ctx, "token"+strconv.Itoa(int(id)), token, maxTokenTime).Err()

	if err != nil {
		return false
	}
	return true
}

// CheckToken 校验token
func CheckToken(id uint, token string) bool {
	result, err := Rdb.Get(ctx, "token"+strconv.Itoa(int(id))).Result()
	if err != nil {
		return false
	}
	if result == token {
		return true
	}
	return false
}

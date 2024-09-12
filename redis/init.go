package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/middlepartedhairstyle/HiWe/utils"
)

var Rdb *redis.Client

func Init() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     utils.Cfg.Redis.Addr,
		Password: utils.Cfg.Redis.Password,
		DB:       utils.Cfg.Redis.Db,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	} else {
		Rdb = rdb
		fmt.Println("redis init success")
	}
}

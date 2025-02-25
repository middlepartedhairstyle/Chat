package main

import (
	"github.com/middlepartedhairstyle/HiWe/internal/mySQL"
	"github.com/middlepartedhairstyle/HiWe/internal/redis"
	"github.com/middlepartedhairstyle/HiWe/internal/routers"
	"github.com/middlepartedhairstyle/HiWe/internal/utils"
)

func init() {
	utils.ReadConfig("config/config.yaml")
	mySQL.Init() //mysql初始化
	redis.Init() //redis初始化
}

func main() {

	server := routers.NewServer(utils.Cfg.App.AppHost, utils.Cfg.App.AppPort)
	server.Run()

}

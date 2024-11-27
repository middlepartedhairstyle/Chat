package main

import (
	"github.com/middlepartedhairstyle/HiWe/mySQL"
	"github.com/middlepartedhairstyle/HiWe/redis"
	"github.com/middlepartedhairstyle/HiWe/routers"
	"github.com/middlepartedhairstyle/HiWe/utils"
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

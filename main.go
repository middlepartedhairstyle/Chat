package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/mySQL"
	"github.com/middlepartedhairstyle/HiWe/routers"
	"github.com/middlepartedhairstyle/HiWe/utils"
	"io"
	"os"
)

func init() {
	utils.ReadConfig("config/config.yaml")
	mySQL.Init()
	//redis.Init()//redis初始化
}

func main() {
	gin.SetMode(gin.DebugMode)
	logFile, _ := os.Create("./logs/gin.log")
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	router := gin.Default()
	routers.Routers(router)

	fmt.Println(router.Routes())

	err := router.Run(utils.Cfg.App.AppHost + ":" + utils.Cfg.App.AppPort)
	if err != nil {
		return
	}
}

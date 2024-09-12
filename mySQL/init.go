package mySQL

import (
	"errors"
	"fmt"
	"github.com/middlepartedhairstyle/HiWe/models"
	"github.com/middlepartedhairstyle/HiWe/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	//连接数据库MYSQL

	dsn := utils.Cfg.MySQL.User + ":" + utils.Cfg.MySQL.Password + "@tcp(" + utils.Cfg.MySQL.Host + ":" + utils.Cfg.MySQL.Port + ")/" + utils.Cfg.MySQL.DBName
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(errors.New("无法连接到数据库"))
	} else {
		DB = db
		fmt.Println("success connect to mysql")
	}

	//数据库表创建
	CreateTable(&(models.UserBaseInfo{})) //用户基础信息
}

func CreateTable(table interface{}) {
	err := DB.AutoMigrate(table)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("success create table")
	}
}

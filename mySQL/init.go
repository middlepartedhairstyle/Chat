package mySQL

import (
	"errors"
	"fmt"
	"github.com/middlepartedhairstyle/HiWe/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	//连接数据库MYSQL

	dsn := utils.Cfg.MySQL.User + ":" + utils.Cfg.MySQL.Password + "@tcp(" + utils.Cfg.MySQL.Host + ":" + utils.Cfg.MySQL.Port + ")/" + utils.Cfg.MySQL.DBName + "?parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(errors.New("无法连接到数据库"))
	} else {
		DB = db
		fmt.Println("success connect to mysql")
	}

	//数据库表创建
	CreateTable(&(UserBaseInfoTable{}))  //用户基础信息 //用户验证码
	CreateTable(&(GroupNumTable{}))      //用户群
	CreateTable(&(GroupUserTable{}))     //用户群用户
	CreateTable(&(GroupMessageTable{}))  //群消息
	CreateTable(&(FriendMessageTable{})) //用户消息
	CreateTable(&(FriendsTable{}))       //用户好友列表
	CreateTable(&(RequestAddFriendTable{}))
	CreateTable(&(RequestAddGroupTable{}))
	CreateTable(&(UserDetailedInfoTable{}))
}

func CreateTable(table interface{}) {
	err := DB.AutoMigrate(table)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("success create table")
	}
}

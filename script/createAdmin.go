package main

import (
	"gin-web/common"
	"gin-web/database"
	"gin-web/model"
	"github.com/joho/godotenv"
)

func main() {

	database.GetDb().Create(&model.Admin{
		Username: "peter",
		Password: common.HmacSha256("123"),
		Email:    "904801074",
	})

}

func init() {

	//加载配置文件
	err := godotenv.Load("./.env")
	if err != nil {
		panic("配置文件加载失败")
	}

}

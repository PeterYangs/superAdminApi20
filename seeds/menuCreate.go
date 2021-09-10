package main

import (
	"fmt"
	"gin-web/database"
	"gin-web/model"
	"github.com/joho/godotenv"
)

func main() {

	menu := []model.Menu{
		{Id: 1, Title: "管理员模块", Path: "", Sort: 100, Pid: 0},
		{Id: 3, Title: "管理员列表", Path: "/main/admin_list", Sort: 100, Pid: 1},
		{Id: 4, Title: "分类管理", Path: "", Sort: 100, Pid: 0},
		{Id: 5, Title: "分类列表", Path: "/main/category_list", Sort: 100, Pid: 4},
		{Id: 6, Title: "菜单管理", Path: "", Sort: 100, Pid: 0},
		{Id: 7, Title: "菜单列表", Path: "/main/menu_list", Sort: 100, Pid: 6},
		{Id: 8, Title: "角色列表", Path: "/main/role_list", Sort: 100, Pid: 1},
		{Id: 9, Title: "规则列表", Path: "/main/rule_list", Sort: 100, Pid: 1},
		{Id: 10, Title: "消息队列", Path: "", Sort: 100, Pid: 0},
		{Id: 11, Title: "即时队列", Path: "/main/queue_list", Sort: 100, Pid: 10},
		{Id: 12, Title: "延迟队列", Path: "/main/queue_delay_list", Sort: 100, Pid: 10},
		{Id: 13, Title: "文件管理", Path: "", Sort: 100, Pid: 0},
		{Id: 14, Title: "文件列表", Path: "/main/file_list", Sort: 100, Pid: 13},
		{Id: 15, Title: "日志管理", Path: "", Sort: 100, Pid: 0},
		{Id: 16, Title: "日志列表", Path: "/main/access_list", Sort: 100, Pid: 15},
	}

	re := database.GetDb().Create(&menu)

	if re.Error != nil {

		fmt.Println(re.Error)
	} else {

		fmt.Println("生成成功")
	}

}

func init() {

	//加载配置文件
	err := godotenv.Load("./.env")
	if err != nil {
		panic("配置文件加载失败")
	}

}

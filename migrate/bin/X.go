package main

import (
	"gin-web/migrate/migrations/migrate_2019_08_12_055619_create_admin_table"
	"gin-web/migrate/migrations/migrate_2019_08_12_666666_create_user_table"
	"gin-web/migrate/migrations/migrate_2020_08_12_666666_create_order_table"

	"github.com/joho/godotenv"
)

func init() {

	//加载配置文件
	err := godotenv.Load("./.env")
	if err != nil {
		panic("配置文件加载失败")
	}

}

func main() {

	//migrate_2019_08_12_055619_create_admin_table.Up()

	migrate_2019_08_12_055619_create_admin_table.Up()
	migrate_2019_08_12_666666_create_user_table.Up()
	migrate_2020_08_12_666666_create_order_table.Up()

}

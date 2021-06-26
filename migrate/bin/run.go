package main

import (
	"fmt"
	"gin-web/database"
	"gin-web/model"
	"github.com/PeterYangs/gcmd"
	"github.com/PeterYangs/tools"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
	"io/ioutil"
	"os"
	"time"
)

func init() {

	//加载配置文件
	err := godotenv.Load("./.env")
	if err != nil {
		panic("配置文件加载失败")
	}

}

func main() {

	//Down()
	//
	//return

	prompt := promptui.Select{
		Label: "选择类型",
		Items: []string{"创建数据库迁移", "执行迁移", "回滚迁移"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if result == "创建数据库迁移" {

		prompt := promptui.Select{
			Label: "选择类型",
			Items: []string{"创建表", "修改表"},
		}

		_, result, err = prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		if result == "创建表" {

			//CreateMigration()
			prompt := promptui.Prompt{
				Label: "输入表名",
				//Validate: validate,
			}

			result, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			CreateMigration(result, "create")

		}

		if result == "修改表" {

			//CreateMigration()
			prompt := promptui.Prompt{
				Label: "输入表名",
				//Validate: validate,
			}

			result, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			CreateMigration(result, "update")

		}

	}

	if result == "执行迁移" {

		prompt := promptui.Select{
			Label: "确定吗？",
			Items: []string{"是", "否"},
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		if result == "是" {

			Up()
		}

	}

	if result == "回滚迁移" {

		prompt := promptui.Select{
			Label: "确定吗？",
			Items: []string{"是", "否"},
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		if result == "是" {

			Down()
		}

	}

}

// CreateMigration 生成数据库迁移文件
func CreateMigration(table string, types string) {

	name := "migrate_" + tools.Date("Y_m_d_His", time.Now().Unix()) + "_" + types + "_" + table + "_table"

	os.Mkdir("./migrate/migrations/"+name, 755)

	if types == "create" {

		CreateTable(name, table)
	}

	if types == "update" {

		UpdateTable(name, table)
	}

}

func UpdateTable(pack string, table string) {

	f, err := os.OpenFile("./migrate/migrations/"+pack+"/migrate.go", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 644)

	if err != nil {

		panic(err)

	}

	defer f.Close()

	f.Write([]byte(`package ` + pack + `

import "gin-web/migrate"

func Up() {

	migrate.Table("` + table + `", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "` + pack + `"


	})

}

func Down() {


}



`))

}

func CreateTable(pack string, table string) {

	f, err := os.OpenFile("./migrate/migrations/"+pack+"/migrate.go", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 644)

	if err != nil {

		panic(err)

	}

	defer f.Close()

	f.Write([]byte(`package ` + pack + `

import "gin-web/migrate"

func Up() {

	migrate.Create("` + table + `", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "` + pack + `"

		createMigrate.BigIncrements("id")

		

	})

}

func Down() {

	migrate.DropIfExists("` + table + `")

}



`))

}

// Up 执行迁移
func Up() {

	fileInfo, _ := ioutil.ReadDir("./migrate/migrations")

	//for _, info := range fileInfo {

	f, _ := os.OpenFile("./migrate/bin/X.go", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 644)

	f.Write([]byte(`package main

import (
 ` + getPackageList(fileInfo) + `
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

   ` + getFuncList(fileInfo) + `

}
`))

	//}

	gcmd.Command("go run migrate/bin/X.go").Start()

}

func getPackageList(f []os.FileInfo) string {

	str := ""

	for _, info := range f {

		str += "\"gin-web/migrate/migrations/" + info.Name() + "\"\n"

	}

	return str

}

func getFuncList(f []os.FileInfo) string {

	str := ""

	for _, info := range f {

		str += info.Name() + ".Up()" + "\n"

	}

	return str

}

// Down 迁移回滚
func Down() {

	//database.GetDb().
	var migration model.Migrations

	re := database.GetDb().Order("id desc").First(&migration)

	if re.Error != nil {

		return
	}

	batch := migration.Batch

	migrations := make([]*model.Migrations, 0)

	database.GetDb().Model(&model.Migrations{}).Where("batch = ?", batch).Find(&migrations)

	f, _ := os.OpenFile("./migrate/bin/Y.go", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 644)

	f.Write([]byte(`package main

import (
 ` + getPackageListForDown(migrations) + `
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

   ` + getFuncListForDown(migrations) + `

}
`))

	gcmd.Command("go run migrate/bin/Y.go").Start()

	for _, m := range migrations {

		database.GetDb().Delete(m)

	}

}

func getPackageListForDown(m []*model.Migrations) string {

	str := ""

	for _, info := range m {

		str += "\"gin-web/migrate/migrations/" + info.Migration + "\"\n"

	}

	return str

}

func getFuncListForDown(m []*model.Migrations) string {

	str := ""

	for _, info := range m {

		str += info.Migration + ".Down()" + "\n"

	}

	return str

}

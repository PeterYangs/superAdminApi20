package main

import (
	"github.com/PeterYangs/gcmd"
	"io/ioutil"
	"os"
)

func main() {

	Up()

}

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

package main

import (
	"fmt"
	"gin-web/common"
	"gin-web/database"
	"gin-web/model"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {

	var admin model.Admin

	re := database.GetDb().Model(&model.Admin{}).Where("username=?", "root").First(&admin)

	password := createPassword(16)

	if re.Error == gorm.ErrRecordNotFound {

		admin.Username = "root"

		admin.Email = "root@superadmin.com"

		admin.Password = common.HmacSha256(password)

		admin.Status = 1

		database.GetDb().Create(&admin)

		common.UpdateOrCreateOne(database.GetDb(), &model.RoleDetail{}, map[string]interface{}{"admin_id": admin.Id}, &model.RoleDetail{AdminId: int(admin.Id), RoleId: 0})

	} else {

		admin.Password = common.HmacSha256(password)

		admin.Status = 1

		database.GetDb().Updates(&admin)

		common.UpdateOrCreateOne(database.GetDb(), &model.RoleDetail{}, map[string]interface{}{"admin_id": admin.Id}, &model.RoleDetail{AdminId: int(admin.Id), RoleId: 0})

	}

	fmt.Println("root密码为：", password)
}

func init() {

	//加载配置文件
	err := godotenv.Load("./.env")
	if err != nil {
		panic("配置文件加载失败")
	}

}

func createPassword(length int) string {

	password := ""

	c := []string{"q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "a", "s", "d", "f", "g", "h", "j", "k", "l", "z", "x", "c", "v", "b", "n", "m", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "!", "@", "#", "$", "%", "^", "&", "*"}

	for i := 0; i < length; i++ {

		password += c[common.MtRand(0, int64(len(c)-1))]
	}

	return password
}

package seeds

import (
	"fmt"
	"github.com/PeterYangs/superAdminCore/database"
	"github.com/PeterYangs/tools"
	"github.com/manifoldco/promptui"
	"gorm.io/gorm"
	"superadmin/common"
	"superadmin/model"
)

type Seeds struct {
}

func (s Seeds) GetName() string {

	return "数据填充"
}

func (s Seeds) ArtisanRun() {

	prompt := promptui.Select{
		Label: "选择类型",
		Items: []string{"重置管理员", "生成基础菜单"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch result {

	case "重置管理员":

		s.resetRoot()

	case "生成基础菜单":

		s.addMenu()

	}

}

func (s Seeds) resetRoot() {

	var admin model.Admin

	re := database.GetDb().Model(&model.Admin{}).Where("username=?", "root").First(&admin)

	password := s.createPassword(16)

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

func (s Seeds) addMenu() {

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

func (s Seeds) createPassword(length int) string {

	password := ""

	c := []string{"q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "a", "s", "d", "f", "g", "h", "j", "k", "l", "z", "x", "c", "v", "b", "n", "m", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "!", "@", "#", "$", "%", "^", "&", "*"}

	for i := 0; i < length; i++ {

		password += c[tools.MtRand(0, int64(len(c)-1))]
	}

	return password
}

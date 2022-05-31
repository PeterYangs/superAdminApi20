package menu

import (
	"github.com/PeterYangs/superAdminCore/v2/contextPlus"
	"github.com/PeterYangs/superAdminCore/v2/database"
	"github.com/PeterYangs/superAdminCore/v2/response"
	"gorm.io/gorm"
	"superadmin/common"
	"superadmin/model"
)

func GetFatherMenu(c *contextPlus.Context) *response.Response {

	menus := make([]*model.Menu, 0)

	database.GetDb().Model(&model.Menu{}).Where("pid = ?", 0).Find(&menus)

	return response.Resp().Api(1, "success", menus)

}

func Update(c *contextPlus.Context) *response.Response {

	type Form struct {
		Title string `json:"title" form:"title" binding:"required"`
		Pid   int    `json:"pid"  form:"pid"`
		Id    int    `json:"id" form:"id"`
		Path  string `json:"path" form:"path"`
		Sort  int    `json:"sort" form:"sort" binding:"required"`
		Rule  string `json:"rule" form:"rule"`
	}

	var form Form

	err := c.ShouldBindPlus(&form)

	if err != nil {

		return response.Resp().Api(2, err.Error(), "")
	}

	menu := model.Menu{
		Id:    uint(form.Id),
		Pid:   form.Pid,
		Title: form.Title,
		Path:  form.Path,
		Sort:  form.Sort,
		Rule:  form.Rule,
	}

	err = common.UpdateOrCreateOne(database.GetDb(), &model.Menu{}, map[string]interface{}{"id": form.Id}, &menu)

	if err != nil {

		return response.Resp().Api(2, err.Error(), "")

	}

	return response.Resp().Api(1, "success", "")

}

func List(c *contextPlus.Context) *response.Response {

	menus := make([]*model.Menu, 0)

	return response.Resp().Api(1, "success", GetMenu(0, &menus))

}

func GetMenu(pid int, m *[]*model.Menu) *[]*model.Menu {

	menus := make([]*model.Menu, 0)

	err := database.GetDb().Model(&model.Menu{}).Where("pid=?", pid).Order("sort asc").Find(&menus)

	if err.Error == gorm.ErrRecordNotFound {

		return nil
	}

	for _, menu := range menus {

		*m = append(*m, menu)

		GetMenu(int(menu.Id), m)

	}

	return m
}

func Detail(c *contextPlus.Context) *response.Response {

	type Form struct {
		Id int `json:"id" uri:"id"`
	}

	var form Form

	err := c.ShouldBindPlus(&form)

	if err != nil {

		return response.Resp().Api(2, err.Error(), "")
	}

	var r model.Menu

	database.GetDb().Where("id = ?", form.Id).First(&r)

	return response.Resp().Api(1, "success", r)

}

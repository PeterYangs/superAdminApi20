package menu

import (
	"gin-web/common"
	"gin-web/contextPlus"
	"gin-web/database"
	"gin-web/model"
	"gin-web/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	}

	var form Form

	err := c.ShouldBindPlus(&form)

	if err != nil {

		return response.Resp().Json(gin.H{"code": 2, "msg": err.Error()})

	}

	menu := model.Menu{
		Id:    uint(form.Id),
		Pid:   form.Pid,
		Title: form.Title,
		Path:  form.Path,
		Sort:  form.Sort,
	}

	err = common.UpdateOrCreateOne(database.GetDb(), &model.Menu{}, map[string]interface{}{"id": form.Id}, &menu)

	if err != nil {

		return response.Resp().Api(2, err.Error(), "")

	}

	return response.Resp().Api(1, "success", "")

}

func List(c *contextPlus.Context) *response.Response {

	menus := make([]*model.Menu, 0)

	//database.GetDb().Model(&model.Menu{}).Find(&menus)

	return response.Resp().Api(1, "success", getMenu(0, &menus))

}

func getMenu(pid int, m *[]*model.Menu) *[]*model.Menu {

	menus := make([]*model.Menu, 0)

	err := database.GetDb().Model(&model.Menu{}).Where("pid=?", pid).Order("sort asc").Find(&menus)

	if err.Error == gorm.ErrRecordNotFound {

		return nil
	}

	for _, menu := range menus {

		*m = append(*m, menu)

		getMenu(int(menu.Id), m)

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

		return response.Resp().Json(gin.H{"code": 2, "msg": err.Error()})

	}

	var r model.Menu

	database.GetDb().Where("id = ?", form.Id).First(&r)

	return response.Resp().Api(1, "success", r)

}

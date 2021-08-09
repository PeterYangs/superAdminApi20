package category

import (
	"gin-web/contextPlus"
	"gin-web/database"
	"gin-web/model"
	"gin-web/response"
	"gorm.io/gorm"
)

func List(c *contextPlus.Context) *response.Response {

	category := make([]*model.Category, 0)

	return response.Resp().Api(1, "success", getMenu(0, &category))

}

func getMenu(pid int, m *[]*model.Category) *[]*model.Category {

	menus := make([]*model.Category, 0)

	err := database.GetDb().Model(&model.Category{}).Where("pid=?", pid).Order("sort asc").Find(&menus)

	if err.Error == gorm.ErrRecordNotFound {

		return nil
	}

	for _, menu := range menus {

		*m = append(*m, menu)

		getMenu(int(menu.Id), m)

	}

	return m
}

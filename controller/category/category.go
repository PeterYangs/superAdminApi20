package category

import (
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/database"
	"github.com/PeterYangs/superAdminCore/response"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"superadmin/common"
	"superadmin/model"
)

func List(c *contextPlus.Context) *response.Response {

	category := make([]*model.Category, 0)

	return response.Resp().Api(1, "success", getMenu(0, &category))

}

//递归查询
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

func Update(c *contextPlus.Context) *response.Response {

	type Form struct {
		Id    int    `json:"id" form:"id"`
		Pid   int    `json:"pid" form:"pid"`
		Lv    int    `json:"lv" form:"lv"`
		Title string `json:"title" form:"title" binding:"required"`
		Img   string `json:"img" form:"img" `
		Sort  int    `json:"sort" form:"sort" binding:"required"`
		//Path  string `json:"path"`
	}

	var form Form

	err := c.ShouldBindPlus(&form)

	if err != nil {

		return response.Resp().Api(2, err.Error(), "")
	}

	category := model.Category{
		Id:    uint(form.Id),
		Pid:   form.Pid,
		Lv:    form.Lv,
		Title: form.Title,
		Img:   form.Img,
		Sort:  form.Sort,
	}

	if category.Id == 0 {

		category.Lv++
	}

	err = common.UpdateOrCreateOne(database.GetDb(), &model.Category{}, map[string]interface{}{"id": form.Id}, &category)

	if err != nil {

		return response.Resp().Api(2, err.Error(), "")
	}

	//fmt.Println(category.Base.CreatedAt)

	//if category.Id == 0 {

	if category.Pid == 0 {

		category.Path = cast.ToString(category.Id)

		//database.GetDb().Omit("CreatedAt").Save(&category)

		database.GetDb().Updates(&category)

	} else {

		var pCategory model.Category

		database.GetDb().Model(&model.Category{}).Where("id=?", category.Pid).First(&pCategory)

		category.Path = pCategory.Path + "," + cast.ToString(category.Id)

		database.GetDb().Save(&category)

	}

	//}

	return response.Resp().Api(1, "success", form)
}

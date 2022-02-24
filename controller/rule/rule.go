package rule

import (
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/database"
	"github.com/PeterYangs/superAdminCore/response"
	"github.com/spf13/cast"
	"superadmin/common"
	"superadmin/model"
)

func Update(c *contextPlus.Context) *response.Response {

	type Form struct {
		Title     string `json:"title" form:"title" binding:"required"`
		Rule      string `json:"rule" form:"rule" binding:"required"`
		GroupName string `json:"group_name"  form:"group_name" `
		Id        uint   `json:"id" form:"id"`
	}

	var form Form

	err := c.ShouldBind(&form)

	if err != nil {

		return response.Resp().Api(2, err.Error(), "")
	}

	r := model.Rule{
		GroupName: form.GroupName,
		Title:     form.Title,
		Rule:      form.Rule,
		Id:        form.Id,
	}

	if r.Id == 0 {

		database.GetDb().Create(&r)

	} else {

		database.GetDb().Model(&r).Updates(&r)

	}

	return response.Resp().Api(1, "测试消息", r)

}

func List(c *contextPlus.Context) *response.Response {

	rules := make([]*model.Rule, 0)

	tx := database.GetDb().Model(&model.Rule{})

	data := common.Paginate(tx, &rules, cast.ToInt(c.DefaultQuery("p", "1")), 10)

	return response.Resp().Api(1, "success", data)

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

	var r model.Rule

	database.GetDb().Where("id = ?", form.Id).First(&r)

	return response.Resp().Api(1, "success", r)
}

func Destroy(c *contextPlus.Context) *response.Response {

	type Form struct {
		Id int `json:"id" uri:"id"`
	}

	var form Form

	err := c.ShouldBindPlus(&form)

	if err != nil {

		return response.Resp().Api(2, err.Error(), "")
	}

	database.GetDb().Delete(&model.Rule{}, form.Id)

	return response.Resp().Api(1, "success", "")

}

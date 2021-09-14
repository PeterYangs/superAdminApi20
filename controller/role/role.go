package role

import (
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/database"
	"github.com/PeterYangs/superAdminCore/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"superadmin/common"
	"superadmin/model"
)

func GetRuleList(c *contextPlus.Context) *response.Response {

	rules := make([]*model.Rule, 0)

	database.GetDb().Find(&rules)

	group := make(map[string][]*model.Rule)

	for _, rule := range rules {

		group[rule.GroupName] = append(group[rule.GroupName], rule)

	}

	return response.Resp().Api(1, "success", group)
}

func Update(c *contextPlus.Context) *response.Response {

	type Form struct {
		Title string `json:"title" form:"title" binding:"required"`
		Rules []int  `json:"rules"  form:"rules" binding:"required"`
		Id    int    `json:"id" form:"id"`
	}

	var form Form

	err := c.ShouldBindPlus(&form)

	if err != nil {

		return response.Resp().Json(gin.H{"code": 2, "msg": err.Error()})

	}

	r := model.Role{

		Title: form.Title,
		Rules: form.Rules,
		Id:    uint(form.Id),
	}

	if r.Id == 0 {

		database.GetDb().Create(&r)

	} else {

		database.GetDb().Model(&r).Updates(&r)

	}

	return response.Resp().Api(1, "success", form)
}

func List(c *contextPlus.Context) *response.Response {

	roles := make([]*model.Role, 0)

	tx := database.GetDb().Model(&model.Role{})

	data := common.Paginate(tx, &roles, cast.ToInt(c.DefaultQuery("p", "1")), 10)

	return response.Resp().Api(1, "success", data)

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

	var r model.Role

	database.GetDb().Where("id = ?", form.Id).First(&r)

	return response.Resp().Json(gin.H{"data": r, "code": 1})

}

func Destroy(c *contextPlus.Context) *response.Response {

	type Form struct {
		Id int `json:"id" uri:"id"`
	}

	var form Form

	err := c.ShouldBindPlus(&form)

	if err != nil {

		return response.Resp().Json(gin.H{"code": 2, "msg": err.Error()})

	}

	database.GetDb().Delete(&model.Role{}, form.Id)

	return response.Resp().Api(1, "success", "")

}

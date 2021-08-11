package admin

import (
	"gin-web/common"
	"gin-web/contextPlus"
	"gin-web/database"
	"gin-web/model"
	"gin-web/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func GetRoleList(c *contextPlus.Context) *response.Response {

	roles := make([]*model.Role, 0)

	database.GetDb().Order("id desc").Find(&roles)

	return response.Resp().Api(1, "success", roles)
}

func List(c *contextPlus.Context) *response.Response {

	roles := make([]*model.Admin, 0)

	tx := database.GetDb().Model(&model.Admin{}).Order("id desc").Preload("RoleDetail.Role")

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

	var r model.Admin

	database.GetDb().Where("id = ?", form.Id).Preload("RoleDetail").First(&r)

	type res struct {
		model.Admin
		RoleId int `json:"role_id"`
	}

	rr := res{
		Admin:  r,
		RoleId: r.RoleDetail.RoleId,
	}

	return response.Resp().Json(gin.H{"data": rr, "code": 1})

}

func Info(c *contextPlus.Context) *response.Response {

	admin, _ := c.Session().Get("admin")

	return response.Resp().Api(1, "success", admin)
}

func SearchRule(c *contextPlus.Context) *response.Response {

	type Form struct {
		Keyword string `json:"keyword" form:"keyword" regex:"/admin.*"`
	}

	var form Form

	err := c.ShouldBindPlus(&form)

	if err != nil {

		return response.Resp().Api(1, "规则不匹配", []string{})

	}

	var rules = make([]*model.Rule, 0)

	database.GetDb().Model(model.Rule{}).Where("rule like ?", "%"+form.Keyword+"%").Limit(10).Find(&rules)

	return response.Resp().Api(1, "success", rules)

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

	database.GetDb().Delete(&model.Admin{}, form.Id)

	return response.Resp().Api(1, "success", "")

}

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

	tx := database.GetDb().Model(&model.Admin{})

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

		return response.Resp().Json(gin.H{"code": 2, "mgs": err.Error()})

	}

	var r model.Admin

	database.GetDb().Where("id = ?", form.Id).Preload("Role").First(&r)

	//var role model.RoleDetail
	//
	//database.GetDb().Where("admin_id=?",r.Id).First(&role)

	//r.RoleId= role.RoleId

	return response.Resp().Json(gin.H{"data": r})

}

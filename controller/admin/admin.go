package admin

import (
	"encoding/json"
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/database"
	"github.com/PeterYangs/superAdminCore/response"
	"github.com/PeterYangs/superAdminCore/route/allUrl"
	"superadmin/common"
	"superadmin/controller/menu"
	"superadmin/model"
	"superadmin/search"

	"github.com/PeterYangs/tools"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func GetRoleList(c *contextPlus.Context) *response.Response {

	roles := make([]*model.Role, 0)

	database.GetDb().Order("id desc").Find(&roles)

	return response.Resp().Api(1, "success", roles)
}

func List(c *contextPlus.Context) *response.Response {

	params := c.DefaultQuery("params", "")

	paramsMap := make(map[string]interface{})

	json.Unmarshal([]byte(params), &paramsMap)

	//fmt.Println(paramsMap)

	roles := make([]*model.Admin, 0)

	tx := database.GetDb().Model(&model.Admin{}).Order("id desc").Preload("RoleDetail.Role")

	if paramsMap["role_id"] != "" && len(paramsMap) > 0 {

		tx.Where("EXISTS( select * from role_detail where role_id = ? and  role_detail.admin_id = admin.id )", paramsMap["role_id"])
	}

	search.NewSearch(tx, paramsMap, []search.Field{{Key: "username", Condition: "like"}})

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

	all := allUrl.NewAllUrl()

	list := all.Search(form.Keyword)

	return response.Resp().Api(1, "success", list)

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

func GetMyMenu(c *contextPlus.Context) *response.Response {

	menus := make([]*model.Menu, 0)

	list := menu.GetMenu(0, &menus)

	admin, _ := c.Session().Get("admin")

	id := admin.(map[string]interface{})["id"].(float64)

	var r model.Admin

	database.GetDb().Where("id = ?", id).Preload("RoleDetail.Role").First(&r)

	//超级管理员显示所有菜单
	if r.RoleDetail.RoleId == 0 {

		return response.Resp().Api(1, "success", list)
	}

	rules := make([]*model.Rule, 0)

	var ids []string

	for _, rule := range r.RoleDetail.Role.Rules {

		ids = append(ids, cast.ToString(rule))
	}

	database.GetDb().Model(&model.Rule{}).Where("id in ?", ids).Find(&rules)

	var rulesArray []string

	for _, rule := range rules {

		rulesArray = append(rulesArray, rule.Rule)
	}

	var temp []*model.Menu

	for _, m := range *list {

		if m.Pid == 0 {

			temp = append(temp, m)

			continue
		}

		if tools.InArray(rulesArray, m.Rule) {

			temp = append(temp, m)
		}

	}

	return response.Resp().Api(1, "success", temp)
}

func RoleList(c *contextPlus.Context) *response.Response {

	roles := make([]*model.Role, 0)

	database.GetDb().Model(&model.Role{}).Find(&roles)

	return response.Resp().Api(1, "success", roles)

}

func GetAllRule(c *contextPlus.Context) *response.Response {

	admin, _ := c.Session().Get("admin")

	id := admin.(map[string]interface{})["id"].(float64)

	var r model.Admin

	database.GetDb().Where("id = ?", id).Preload("RoleDetail.Role").First(&r)

	if r.RoleDetail.RoleId == 0 {

		return response.Resp().Api(1, "success", true)
	}

	rules := make([]*model.Rule, 0)

	ids := make([]string, len(r.RoleDetail.Role.Rules))

	for i, rule := range r.RoleDetail.Role.Rules {

		ids[i] = cast.ToString(rule)
	}

	database.GetDb().Where("id in ?", ids).Find(&rules)

	list := make([]string, len(r.RoleDetail.Role.Rules))

	for i, rule := range rules {

		list[i] = rule.Rule

	}

	return response.Resp().Api(1, "success", list)

}

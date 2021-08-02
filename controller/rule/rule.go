package rule

import (
	"fmt"
	"gin-web/contextPlus"
	"gin-web/database"
	"gin-web/model"
	"gin-web/response"
	"github.com/gin-gonic/gin"
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

		return response.Resp().Json(gin.H{"code": 2, "mgs": err.Error()})

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

	return response.Resp().Json(gin.H{"data": r, "msg": "测试消息", "code": 1})

}

func List(c *contextPlus.Context) *response.Response {

	rules := make([]*model.Rule, 0)

	database.GetDb().Model(&model.Rule{}).Find(&rules)

	return response.Resp().Json(gin.H{"data": rules})

}

func Detail(c *contextPlus.Context) *response.Response {

	fmt.Println(c.Param("id"))

	type Form struct {
		Id int `json:"id" uri:"id"`
	}

	var form Form

	err := c.ShouldBindPlus(&form)

	if err != nil {

		return response.Resp().Json(gin.H{"code": 2, "mgs": err.Error()})

	}

	var r model.Rule

	database.GetDb().Where("id = ?", form.Id).First(&r)

	return response.Resp().Json(gin.H{"data": r})

}

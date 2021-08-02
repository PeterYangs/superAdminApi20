package role

import (
	"gin-web/contextPlus"
	"gin-web/database"
	"gin-web/model"
	"gin-web/response"
)

func GetRuleList(c *contextPlus.Context) *response.Response {

	rules := make([]*model.Rule, 0)

	database.GetDb().Find(&rules)

	return response.Resp().Api(1, "success", rules)
}

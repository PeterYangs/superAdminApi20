package access

import (
	"gin-web/common"
	"gin-web/contextPlus"
	"gin-web/database"
	"gin-web/model"
	"gin-web/response"
	"github.com/spf13/cast"
)

func List(c *contextPlus.Context) *response.Response {

	roles := make([]*model.Access, 0)

	tx := database.GetDb().Model(&model.Access{}).Order("id desc")

	data := common.Paginate(tx, &roles, cast.ToInt(c.DefaultQuery("p", "1")), 10)

	return response.Resp().Api(1, "success", data)

}

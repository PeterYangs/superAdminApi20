package access

import (
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/database"
	"github.com/PeterYangs/superAdminCore/response"
	"github.com/spf13/cast"
	"superadmin/common"
	"superadmin/model"
)

func List(c *contextPlus.Context) *response.Response {

	roles := make([]*model.Access, 0)

	tx := database.GetDb().Model(&model.Access{}).Order("id desc")

	data := common.Paginate(tx, &roles, cast.ToInt(c.DefaultQuery("p", "1")), 10)

	return response.Resp().Api(1, "success", data)

}

package authCheck

import (
	"github.com/PeterYangs/superAdminCore/v2/contextPlus"
	"github.com/PeterYangs/superAdminCore/v2/database"
	"github.com/PeterYangs/superAdminCore/v2/response"
	"github.com/spf13/cast"
	"superadmin/model"
)

func AuthCheck(c *contextPlus.Context) {

	//fmt.Println(c.Handler.Tag,"----------")

	//跳过权限检查
	if c.Handler.Tag == "skip_auth" {

		return
	}

	admin_, err := c.Session().Get("admin")

	if err != nil {

		c.JSON(200, response.Resp().Api(11, "session异常", "").GetData())

		c.Abort()

		return
	}

	admin := admin_.(map[string]interface{})

	adminId := admin["id"]

	var info model.Admin

	database.GetDb().Model(&model.Admin{}).Where("id = ?", adminId).Preload("RoleDetail.Role").First(&info)

	//超级管理员
	if info.RoleDetail.RoleId == 0 {

		return
	}

	rules := make([]*model.Rule, 0)

	var str []string

	for _, rule := range info.RoleDetail.Role.Rules {

		str = append(str, cast.ToString(rule))
	}

	database.GetDb().Model(&model.Rule{}).Where("id in ? ", str).Find(&rules)

	url := c.FullPath()

	for _, rule := range rules {

		if url == rule.Rule {

			return
		}
	}

	c.JSON(200, response.Resp().Api(51, "你没有这个权限", "").GetData())

	c.Abort()

}

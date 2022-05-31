package loginCheck

import (
	"github.com/PeterYangs/superAdminCore/v2/contextPlus"
	"github.com/PeterYangs/superAdminCore/v2/response"
)

func LoginCheck(c *contextPlus.Context) {

	if !c.Session().Exist("admin") {

		c.JSON(200, response.Resp().Api(10, "请登录", "").GetData())

		c.Abort()

		return
	}

	//fmt.Println(c.FullPath())

}

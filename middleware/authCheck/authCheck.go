package authCheck

import (
	"gin-web/contextPlus"
	"gin-web/response"
)

func AuthCheck(c *contextPlus.Context) {

	admin_, err := c.Session().Get("admin")

	if err != nil {

		c.JSON(200, response.Resp().Api(11, "session异常", "").GetData())

		c.Abort()

		return
	}

	admin := admin_.(map[string]interface{})

	adminId := admin["id"]

	//fmt.Println(adminId)
	_ = adminId
	//fmt.Println(admin)

}

package controller

import (
	"github.com/PeterYangs/superAdminCore/v2/contextPlus"
	"github.com/PeterYangs/superAdminCore/v2/response"
)

func Index(c *contextPlus.Context) *response.Response {

	return response.Resp().Api(1, "success", "index")
}

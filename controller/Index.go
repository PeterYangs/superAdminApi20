package controller

import (
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/response"
)

func Index(c *contextPlus.Context) *response.Response {

	return response.Resp().Api(1, "success", "index")
}

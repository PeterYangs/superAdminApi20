package test

import (
	"fmt"
	"github.com/PeterYangs/superAdminCore/cache"
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/response"
)

func Cache(c *contextPlus.Context) *response.Response {

	cache.Cache().Put("nice", "0", 0)

	fmt.Println(cache.Cache().Get("nice"))

	return response.Resp().Api(1, "success", "")
}

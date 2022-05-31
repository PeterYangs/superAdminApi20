package test

import (
	"fmt"
	"github.com/PeterYangs/superAdminCore/v2/cache"
	"github.com/PeterYangs/superAdminCore/v2/contextPlus"
	"github.com/PeterYangs/superAdminCore/v2/response"
)

func Cache(c *contextPlus.Context) *response.Response {

	cache.Cache().Put("nice", "0", 0)

	fmt.Println(cache.Cache().Get("nice"))

	return response.Resp().Api(1, "success", "")
}

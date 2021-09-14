package accessLog

import (
	"bytes"
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/queue"
	"superadmin/task/access"

	"io/ioutil"
)

func AccessLog(c *contextPlus.Context) {

	b, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	var id float64

	if c.Session().Exist("admin") {

		admin, _ := c.Session().Get("admin")

		id = admin.(map[string]interface{})["id"].(float64)

	}

	queue.Dispatch(access.NewAccessTask(c.ClientIP(), c.Request.URL.String(), string(b), id)).Run()

}

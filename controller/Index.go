package controller

import (
	"gin-web/contextPlus"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Index 主页
func Index(c *contextPlus.Context) interface{} {

	//fmt.Println(c.Session().Get("test3"))

	c.Session().Set("test3", "123")

	return gin.H{"code": 1, "msg": "hello world"}
}

// SessionSet 并发写入demo
func SessionSet(c *contextPlus.Context) interface{} {

	s := c.Session()

	for i := 0; i < 100; i++ {

		go func(ii int, ss *contextPlus.Session) {

			ss.Set("test"+strconv.Itoa(ii), ii)

			//fmt.Println(e)
		}(i, s)
	}

	return gin.H{"code": 1, "msg": "hello world"}
}

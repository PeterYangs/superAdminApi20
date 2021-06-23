package controller

import (
	"gin-web/contextPlus"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Index 主页
func Index(c *contextPlus.Context) interface{} {

	//fmt.Println(c.Jj)

	//fmt.Println(c.Get("test"))

	//panic("123")

	return gin.H{"code": 1, "msg": "hello world"}
}

// SessionSet 并发写入demo
func SessionSet(c *contextPlus.Context) interface{} {

	s := c.Session()

	for i := 0; i < 100; i++ {

		go func(ii int, ss *contextPlus.Session) {

			ss.Set("test"+strconv.Itoa(ii), ii)

		}(i, s)
	}

	return gin.H{"code": 1, "msg": "hello world"}
}

func Captcha(c *contextPlus.Context) interface{} {

	b := c.GetCaptcha()

	c.Header("content-type", "image/png")

	return b
}

func CheckCaptcha(c *contextPlus.Context) interface{} {

	code := c.Query("code")

	return gin.H{"data": c.CheckCaptcha(code)}

}

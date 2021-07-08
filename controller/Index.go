package controller

import (
	"gin-web/component/logs"
	"gin-web/contextPlus"
	"gin-web/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Index 主页
func Index(c *contextPlus.Context) *response.Response {

	//panic("自定义错误")

	logs.NewLogs().Debug("yy").Stdout()

	return response.Resp().Json(gin.H{"data": "hello world"})
}

func Index2(c *contextPlus.Context) *response.Response {

	return response.Resp().Json(gin.H{"code": 1, "msg": "hello world"})
}

// SessionSet 并发写入demo
func SessionSet(c *contextPlus.Context) *response.Response {

	s := c.Session()

	for i := 0; i < 100; i++ {

		go func(ii int, ss *contextPlus.Session) {

			ss.Set("test"+strconv.Itoa(ii), ii)

		}(i, s)
	}

	return response.Resp().Json(gin.H{"code": 1, "msg": "hello world"})
}

func Captcha(c *contextPlus.Context) *response.Response {

	b := c.GetCaptcha()

	c.Header("content-type", "image/png")

	return response.Resp().Byte(b)
}

func CheckCaptcha(c *contextPlus.Context) *response.Response {

	code := c.Query("code")

	return response.Resp().Json(gin.H{"bool": code})

}

package captcha

import (
	"gin-web/contextPlus"
	"gin-web/response"
)

func Captcha(c *contextPlus.Context) *response.Response {

	b := c.GetCaptcha()

	c.Header("content-type", "image/png")

	return response.Resp().Byte(b)
}

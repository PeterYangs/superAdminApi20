package captcha

import (
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/response"
)

func Captcha(c *contextPlus.Context) *response.Response {

	b := c.GetCaptcha()

	c.Header("content-type", "image/png")

	return response.Resp().Byte(b)
}

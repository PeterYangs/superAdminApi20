package captcha

import (
	"github.com/PeterYangs/superAdminCore/v2/contextPlus"
	"github.com/PeterYangs/superAdminCore/v2/response"
)

func Captcha(c *contextPlus.Context) *response.Response {

	b := c.GetCaptcha()

	c.Header("content-type", "image/png")

	return response.Resp().Byte(b)
}

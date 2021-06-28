package captcha

import "gin-web/contextPlus"

func Captcha(c *contextPlus.Context) interface{} {

	b := c.GetCaptcha()

	c.Header("content-type", "image/png")

	return b
}

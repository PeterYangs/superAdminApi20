package loginLimiter

import (
	"gin-web/component/limiter"
	"gin-web/contextPlus"
	"golang.org/x/time/rate"
	"time"
)

func LoginLimiter(c *contextPlus.Context) {

	if !limiter.NewLimiter(rate.Every(1*time.Second), 10, c.ClientIP()).Allow() {

		c.String(500, "访问频率过高")

		c.Abort()

	}

}

package routes

import (
	"gin-web/controller/captcha"
	"gin-web/controller/login"
	"gin-web/middleware/loginLimiter"
)

func _init(_r group) {

	_r.Group("/login", func(g group) {

		g.Registered(POST, "/login", login.Login, loginLimiter.LoginLimiter).Bind()

		g.Registered(POST, "/registered", login.Registered).Bind()
		g.Registered(POST, "/logout", login.Logout).Bind()

	})

	_r.Group("/captcha", func(g group) {

		g.Registered(GET, "/captcha", captcha.Captcha).Bind()
	})

}

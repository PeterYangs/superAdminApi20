package routes

import (
	"gin-web/controller/captcha"
	"gin-web/controller/login"
)

func _init(_r group) {

	_r.Group("/login", func(g group) {

		g.Registered(POST, "/login", login.Login).Bind()

		g.Registered(POST, "/registered", login.Registered).Bind()

	})

	_r.Group("/captcha", func(g group) {

		g.Registered(GET, "/captcha", captcha.Captcha).Bind()
	})

}

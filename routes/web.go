package routes

import (
	"gin-web/controller"
	"gin-web/controller/file"
	"gin-web/controller/regex"
	"gin-web/middleware"
	"gin-web/middleware/loginLimiter"
)

func _init(_r group) {

	//路由组，支持嵌套
	_r.Group("/api", func(g group) {

		g.Registered(GET, "/", controller.Index).Bind()
		g.Registered(GET, "/gg", controller.Index).Bind()

		g.Group("/login", func(g2 group) {

			g2.Registered(GET, "/", controller.Index).Bind()
		})

	}, middleware.GoOn)

	//单路由
	_r.Registered(GET, "/", controller.Index, loginLimiter.LoginLimiter).Bind()
	_r.Registered(GET, "/2", controller.Index2).Bind()

	_r.Registered(GET, "/c", controller.Captcha).Bind()

	_r.Registered(GET, "/check", controller.CheckCaptcha).Bind()

	//文件上传
	_r.Registered(POST, "/file", file.File).Bind()

	_r.Registered(GET, "/regex/:name", regex.Regex).Bind()

}

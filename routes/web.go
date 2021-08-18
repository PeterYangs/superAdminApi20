package routes

import (
	"gin-web/controller"
	admin2 "gin-web/controller/admin"
	"gin-web/controller/captcha"
	category2 "gin-web/controller/category"
	"gin-web/controller/login"
	menu2 "gin-web/controller/menu"
	queue2 "gin-web/controller/queue"
	role2 "gin-web/controller/role"
	rule2 "gin-web/controller/rule"
	upload2 "gin-web/controller/upload"
	"gin-web/middleware/authCheck"
	"gin-web/middleware/loginCheck"
	"gin-web/middleware/loginLimiter"
)

func _init(_r group) {

	_r.Group("/login", func(g group) {

		g.Registered(POST, "/login", login.Login, loginLimiter.LoginLimiter).Bind()

		g.Registered(POST, "/logout", login.Logout).Bind()

	})

	_r.Group("/admin", func(admin group) {

		admin.Group("/rule", func(rule group) {

			rule.Registered(ANY, "/update", rule2.Update).Bind()
			rule.Registered(GET, "/list", rule2.List).Bind()
			rule.Registered(ANY, "/detail/:id", rule2.Detail).Bind()
			rule.Registered(ANY, "/destroy/:id", rule2.Destroy).Bind()

		})

		admin.Group("/role", func(role group) {

			role.Registered(GET, "/GetRuleList", role2.GetRuleList).SetTag("skip_auth").Bind()
			role.Registered(POST, "/update", role2.Update).Bind()
			role.Registered(GET, "/list", role2.List).Bind()
			role.Registered(GET, "/detail/:id", role2.Detail).Bind()
			role.Registered(POST, "/destroy/:id", role2.Destroy).Bind()

		})

		admin.Group("/admin", func(admins group) {

			admins.Registered(GET, "/getRoleList", admin2.GetRoleList).Bind()
			admins.Registered(POST, "/registered", login.Registered).Bind()
			admins.Registered(POST, "/list", admin2.List).Bind()
			admins.Registered(POST, "/detail/:id", admin2.Detail).Bind()
			admins.Registered(GET, "/info", admin2.Info).SetTag("skip_auth").Bind()
			admins.Registered(GET, "/SearchRule", admin2.SearchRule).SetTag("skip_auth").Bind()
			admins.Registered(POST, "/destroy/:id", admin2.Destroy).Bind()
			admins.Registered(POST, "/getMyMenu", admin2.GetMyMenu).SetTag("skip_auth").Bind()

		})

		admin.Group("/menu", func(menu group) {

			menu.Registered(GET, "/getFatherMenu", menu2.GetFatherMenu).SetTag("skip_auth").Bind()
			menu.Registered(POST, "/update", menu2.Update).Bind()
			menu.Registered(POST, "/list", menu2.List).Bind()
			menu.Registered(GET, "/detail/:id", menu2.Detail).Bind()

		})

		admin.Group("/category", func(category group) {

			category.Registered(GET, "/list", category2.List).Bind()
			category.Registered(POST, "/update", category2.Update).Bind()

		})

		admin.Group("/upload", func(upload group) {

			upload.Registered(POST, "/upload", upload2.Upload).Bind()
		})

		admin.Group("/queue", func(queue group) {

			queue.Registered(POST, "/list", queue2.List).Bind()
			queue.Registered(POST, "/delay_list", queue2.DelayList).Bind()

		})

	}, loginCheck.LoginCheck, authCheck.AuthCheck)

	_r.Group("/captcha", func(g group) {

		g.Registered(ANY, "/captcha", captcha.Captcha).Bind()
	})

	_r.Registered(GET, "/task", controller.Task).Bind()
	_r.Registered(GET, "/task2", controller.Task2).Bind()

}

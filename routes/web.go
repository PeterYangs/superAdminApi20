package routes

import (
	admin2 "gin-web/controller/admin"
	"gin-web/controller/captcha"
	"gin-web/controller/login"
	role2 "gin-web/controller/role"
	rule2 "gin-web/controller/rule"
	"gin-web/middleware/loginLimiter"
)

func _init(_r group) {

	_r.Group("/login", func(g group) {

		g.Registered(ANY, "/login", login.Login, loginLimiter.LoginLimiter).Bind()

		//g.Registered(ANY, "/registered", login.Registered).Bind()
		g.Registered(ANY, "/logout", login.Logout).Bind()

	})

	_r.Group("/admin", func(admin group) {

		admin.Group("/rule", func(rule group) {

			rule.Registered(ANY, "/update", rule2.Update).Bind()
			rule.Registered(GET, "/list", rule2.List).Bind()
			rule.Registered(ANY, "/detail/:id", rule2.Detail).Bind()
			rule.Registered(ANY, "/destroy/:id", rule2.Destroy).Bind()

		})

		admin.Group("/role", func(role group) {

			role.Registered(GET, "/GetRuleList", role2.GetRuleList).Bind()
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

		})

	})

	_r.Group("/captcha", func(g group) {

		g.Registered(ANY, "/captcha", captcha.Captcha).Bind()
	})

}

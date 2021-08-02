package routes

import (
	"gin-web/controller/captcha"
	"gin-web/controller/login"
	role2 "gin-web/controller/role"
	rule2 "gin-web/controller/rule"
	"gin-web/middleware/loginLimiter"
)

func _init(_r group) {

	_r.Group("/login", func(g group) {

		g.Registered(POST, "/login", login.Login, loginLimiter.LoginLimiter).Bind()

		g.Registered(POST, "/registered", login.Registered).Bind()
		g.Registered(POST, "/logout", login.Logout).Bind()

	})

	_r.Group("/admin", func(admin group) {

		admin.Group("/rule", func(rule group) {

			rule.Registered(GET, "/update", rule2.Update).Bind()
			rule.Registered(GET, "/list", rule2.List).Bind()
			rule.Registered(GET, "/detail/:id", rule2.Detail).Bind()
			rule.Registered(GET, "/destroy/:id", rule2.Destroy).Bind()

		})

		admin.Group("/role", func(role group) {

			role.Registered(GET, "/GetRuleList", role2.GetRuleList).Bind()

		})

	})

	_r.Group("/captcha", func(g group) {

		g.Registered(GET, "/captcha", captcha.Captcha).Bind()
	})

}

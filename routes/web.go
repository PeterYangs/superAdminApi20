package routes

import (
	"gin-web/contextPlus"
	access2 "gin-web/controller/access"
	admin2 "gin-web/controller/admin"
	"gin-web/controller/captcha"
	category2 "gin-web/controller/category"
	file2 "gin-web/controller/file"
	"gin-web/controller/login"
	menu2 "gin-web/controller/menu"
	queue2 "gin-web/controller/queue"
	role2 "gin-web/controller/role"
	rule2 "gin-web/controller/rule"
	upload2 "gin-web/controller/upload"
	"gin-web/kernel"
	"gin-web/middleware/authCheck"
	"gin-web/middleware/loginCheck"
	"gin-web/middleware/loginLimiter"
	"gin-web/response"
)

func _init(_r group) {

	_r.Group("/login", func(g group) {

		g.Registered(POST, "/login", login.Login, loginLimiter.LoginLimiter).Bind()

		g.Registered(POST, "/logout", login.Logout).Bind()

	})

	_r.Group("/admin", func(admin group) {

		admin.Group("/rule", func(rule group) {

			rule.Registered(POST, "/update", rule2.Update).Bind()
			rule.Registered(GET, "/list", rule2.List).Bind()
			rule.Registered(GET, "/detail/:id", rule2.Detail).Bind()
			rule.Registered(POST, "/destroy/:id", rule2.Destroy).Bind()

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
			admins.Registered(POST, "/roleList", admin2.RoleList).Bind()
			admins.Registered(GET, "/getAllRule", admin2.GetAllRule).SetTag("skip_auth").Bind()

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
			upload.Registered(ANY, "/big_file", upload2.BigFile).Bind()
		})

		admin.Group("/queue", func(queue group) {

			queue.Registered(POST, "/list", queue2.List).Bind()
			queue.Registered(POST, "/delay_list", queue2.DelayList).Bind()

		})

		admin.Group("/access", func(access group) {

			access.Registered(POST, "/list", access2.List).Bind()

		})

		admin.Group("/file", func(file group) {

			file.Registered(POST, "/update", file2.Update).Bind()
			file.Registered(POST, "/list", file2.List).Bind()
			file.Registered(POST, "/destroy/:id", file2.Destroy).Bind()

		})

	}, loginCheck.LoginCheck, authCheck.AuthCheck)

	_r.Group("/captcha", func(g group) {

		g.Registered(GET, "/captcha", captcha.Captcha).Bind()
	})

	//判断http服务已启动接口
	_r.Registered(GET, "/ping/:id", func(c *contextPlus.Context) *response.Response {

		id := c.Param("id")

		//判断服务id
		if id == kernel.Id {

			return response.Resp().String("success")
		}

		return response.Resp().String("fail")

	}).Bind()

	//
	//_r.Registered(GET, "/task", controller.Task).Bind()
	//_r.Registered(GET, "/task2", controller.Task2).Bind()

}
